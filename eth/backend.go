// Copyright 2014 The go-ethereum Authors
// Copyright 2018 Webchain project
// This file is part of Webchain.
//
// Webchain is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Webchain is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Webchain. If not, see <http://www.gnu.org/licenses/>.

// Package eth implements the Ethereum protocol.
package eth

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/webchain-network/cryptonight"
	"github.com/webchain-network/webchaind/accounts"
	"github.com/webchain-network/webchaind/common"
	"github.com/webchain-network/webchaind/common/compiler"
	"github.com/webchain-network/webchaind/common/httpclient"
	"github.com/webchain-network/webchaind/common/registrar/ethreg"
	"github.com/webchain-network/webchaind/core"
	"github.com/webchain-network/webchaind/core/types"
	"github.com/webchain-network/webchaind/eth/downloader"
	"github.com/webchain-network/webchaind/eth/filters"
	"github.com/webchain-network/webchaind/ethdb"
	"github.com/webchain-network/webchaind/event"
	"github.com/webchain-network/webchaind/logger"
	"github.com/webchain-network/webchaind/logger/glog"
	"github.com/webchain-network/webchaind/miner"
	"github.com/webchain-network/webchaind/node"
	"github.com/webchain-network/webchaind/p2p"
	"github.com/webchain-network/webchaind/rlp"
	"github.com/webchain-network/webchaind/rpc"
)

const (
	cryptonightRevision = 1
)

type Config struct {
	ChainConfig *core.ChainConfig // chain configuration

	NetworkId int // Network ID to use for selecting peers to connect to
	Genesis   *core.GenesisDump
	FastSync  bool // Enables the state download based fast synchronisation algorithm
	MaxPeers  int

	BlockChainVersion  int
	SkipBcVersionCheck bool // e.g. blockchain export
	DatabaseCache      int
	DatabaseHandles    int

	NatSpec   bool
	DocRoot   string
	PowTest   bool
	PowShared bool

	AccountManager *accounts.Manager
	Etherbase      common.Address
	GasPrice       *big.Int
	MinerThreads   int
	SolcPath       string

	UseAddrTxIndex bool

	GpoMinGasPrice          *big.Int
	GpoMaxGasPrice          *big.Int
	GpoFullBlockRatio       int
	GpobaseStepDown         int
	GpobaseStepUp           int
	GpobaseCorrectionFactor int

	TestGenesisBlock *types.Block   // Genesis block to seed the chain database with (testing only!)
	TestGenesisState ethdb.Database // Genesis state to seed the database with (testing only!)
}

type Ethereum struct {
	config      *Config
	chainConfig *core.ChainConfig
	// Channel for shutting down the ethereum
	shutdownChan chan bool

	// DB interfaces
	chainDb   ethdb.Database // Block chain database
	dappDb    ethdb.Database // Dapp database
	indexesDb ethdb.Database // Indexes database (optional -- eg. add-tx indexes)

	// Handlers
	txPool          *core.TxPool
	txMu            sync.Mutex
	blockchain      *core.BlockChain
	accountManager  *accounts.Manager
	pow             *cryptonight.Cryptonight
	protocolManager *ProtocolManager
	SolcPath        string
	solc            *compiler.Solidity
	gpo             *GasPriceOracle

	GpoMinGasPrice          *big.Int
	GpoMaxGasPrice          *big.Int
	GpoFullBlockRatio       int
	GpobaseStepDown         int
	GpobaseStepUp           int
	GpobaseCorrectionFactor int

	httpclient *httpclient.HTTPClient

	eventMux *event.TypeMux
	miner    *miner.Miner

	Mining        bool
	MinerThreads  int
	NatSpec       bool
	PowTest       bool
	etherbase     common.Address
	netVersionId  int
	netRPCService *PublicNetAPI
}

func New(ctx *node.ServiceContext, config *Config) (*Ethereum, error) {
	// Open the chain database and perform any upgrades needed
	chainDb, err := ctx.OpenDatabase("chaindata", config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if err := upgradeChainDatabase(chainDb); err != nil {
		return nil, err
	}
	if err := addMipmapBloomBins(chainDb); err != nil {
		return nil, err
	}

	dappDb, err := ctx.OpenDatabase("dapp", config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}

	glog.V(logger.Info).Infof("Protocol Versions: %v, Network Id: %v, Chain Id: %v", ProtocolVersions, config.NetworkId, config.ChainConfig.GetChainID())
	glog.D(logger.Warn).Infof("Protocol Versions: %v, Network Id: %v, Chain Id: %v", logger.ColorGreen(fmt.Sprintf("%v", ProtocolVersions)), logger.ColorGreen(strconv.Itoa(config.NetworkId)), logger.ColorGreen(func() string {
		cid := config.ChainConfig.GetChainID().String()
		if cid == "0" {
			cid = "not set"
		}
		return cid
	}()))

	// Load up any custom genesis block if requested
	if config.Genesis != nil {
		_, err := core.WriteGenesisBlock(chainDb, config.Genesis)
		if err != nil {
			return nil, err
		}
	}

	// Load up a test setup if directly injected
	if config.TestGenesisState != nil {
		chainDb = config.TestGenesisState
	}
	if config.TestGenesisBlock != nil {
		core.WriteTd(chainDb, config.TestGenesisBlock.Hash(), config.TestGenesisBlock.Difficulty())
		core.WriteBlock(chainDb, config.TestGenesisBlock)
		core.WriteCanonicalHash(chainDb, config.TestGenesisBlock.Hash(), config.TestGenesisBlock.NumberU64())
		core.WriteHeadBlockHash(chainDb, config.TestGenesisBlock.Hash())
	}

	if !config.SkipBcVersionCheck {
		bcVersion := core.GetBlockChainVersion(chainDb)
		if bcVersion != config.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run webchaind upgradedb.\n", bcVersion, config.BlockChainVersion)
		}
		core.WriteBlockChainVersion(chainDb, config.BlockChainVersion)
	}
	glog.V(logger.Info).Infof("Blockchain DB Version: %d", config.BlockChainVersion)

	eth := &Ethereum{
		config:                  config,
		shutdownChan:            make(chan bool),
		chainDb:                 chainDb,
		dappDb:                  dappDb,
		eventMux:                ctx.EventMux,
		accountManager:          config.AccountManager,
		etherbase:               config.Etherbase,
		netVersionId:            config.NetworkId,
		NatSpec:                 config.NatSpec,
		MinerThreads:            config.MinerThreads,
		SolcPath:                config.SolcPath,
		PowTest:                 config.PowTest,
		GpoMinGasPrice:          config.GpoMinGasPrice,
		GpoMaxGasPrice:          config.GpoMaxGasPrice,
		GpoFullBlockRatio:       config.GpoFullBlockRatio,
		GpobaseStepDown:         config.GpobaseStepDown,
		GpobaseStepUp:           config.GpobaseStepUp,
		GpobaseCorrectionFactor: config.GpobaseCorrectionFactor,
		httpclient:              httpclient.New(config.DocRoot),
	}
	switch {
	case config.PowTest:
		glog.V(logger.Info).Infof("Consensus: cryptonight used in test mode")
		eth.pow, err = cryptonight.NewForTesting()
		if err != nil {
			return nil, err
		}
	case config.PowShared:
		glog.V(logger.Info).Infof("Consensus: cryptonight used in shared mode")
		eth.pow = cryptonight.NewShared()

	default:
		eth.pow = cryptonight.New()
	}

	// Initialize indexes db if enabled
	// Blockchain will be assigned the db and atx enabled after blockchain is initialized below.
	var indexesDb ethdb.Database
	if config.UseAddrTxIndex {
		// TODO: these are arbitrary numbers I just made up. Optimize?
		// The reason these numbers are different than the atxi-build command is because for "appending" (vs. building)
		// the atxi database should require far fewer resources since application performance is limited primarily by block import (chaindata db).
		ethdb.SetCacheRatio("chaindata", 0.95)
		ethdb.SetHandleRatio("chaindata", 0.95)
		ethdb.SetCacheRatio("indexes", 0.05)
		ethdb.SetHandleRatio("indexes", 0.05)
		indexesDb, err = ctx.OpenDatabase("indexes", config.DatabaseCache, config.DatabaseCache)
		if err != nil {
			return nil, err
		}
		eth.indexesDb = indexesDb
	}

	// load the genesis block or write a new one if no genesis
	// block is present in the database.
	genesis := core.GetBlock(chainDb, core.GetCanonicalHash(chainDb, 0))
	if genesis == nil {
		genesis, err = core.WriteGenesisBlock(chainDb, core.DefaultConfigMainnet.Genesis)
		if err != nil {
			return nil, err
		}
		glog.V(logger.Info).Infof("Successfully wrote default ethereum mainnet genesis block: %s", logger.ColorGreen(genesis.Hash().Hex()))
		glog.D(logger.Warn).Infof("Wrote mainnet genesis block: %s", logger.ColorGreen(genesis.Hash().Hex()))
	}

	// Log genesis block information.
	var genName string
	if fmt.Sprintf("%x", genesis.Hash()) == "0cd786a2425d16f152c658316c423e6ce1181e15c3295826d7c9904cba9ce303" {
		genName = "morden testnet"
	} else if fmt.Sprintf("%x", genesis.Hash()) == "78e066e78f30695e6c4218db4af16a670085b58e592378e864b40156b87a4c19" {
		genName = "mainnet"
	} else {
		genName = "custom"
	}
	glog.V(logger.Info).Infof("Successfully established %s genesis block: %s", genName, genesis.Hash().Hex())
	glog.D(logger.Warn).Infof("Genesis block: %s (%s)", logger.ColorGreen(genesis.Hash().Hex()), genName)

	if config.ChainConfig == nil {
		return nil, errors.New("missing chain config")
	}

	eth.chainConfig = config.ChainConfig

	eth.blockchain, err = core.NewBlockChain(chainDb, eth.chainConfig, eth.pow, eth.EventMux())
	if err != nil {
		if err == core.ErrNoGenesis {
			return nil, fmt.Errorf(`No chain found. Please initialise a new chain using the "init" subcommand.`)
		}
		return nil, err
	}
	// Configure enabled atxi for blockchain
	if config.UseAddrTxIndex {
		eth.blockchain.SetAtxi(&core.AtxiT{
			Db: eth.indexesDb,
		})
	}

	eth.gpo = NewGasPriceOracle(eth)

	newPool := core.NewTxPool(eth.chainConfig, eth.EventMux(), eth.blockchain.State, eth.blockchain.GasLimit)
	eth.txPool = newPool

	m := downloader.FullSync
	if config.FastSync {
		m = downloader.FastSync
	}
	if eth.protocolManager, err = NewProtocolManager(eth.chainConfig, m, uint64(config.NetworkId), eth.eventMux, eth.txPool, eth.pow, eth.blockchain, chainDb); err != nil {
		return nil, err
	}
	eth.miner = miner.New(eth, eth.chainConfig, eth.EventMux(), eth.pow)
	if err = eth.miner.SetGasPrice(config.GasPrice); err != nil {
		return nil, err
	}

	return eth, nil
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Ethereum) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEthereumAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicAccountAPI(s.accountManager),
			Public:    true,
		}, {
			Namespace: "personal",
			Version:   "1.0",
			Service:   NewPrivateAccountAPI(s),
			Public:    false,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicBlockChainAPI(s.chainConfig, s.blockchain, s.miner, s.chainDb, s.gpo, s.eventMux, s.accountManager),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicTransactionPoolAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "txpool",
			Version:   "1.0",
			Service:   NewPublicTxPoolAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.chainDb, s.eventMux),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   ethreg.NewPrivateRegistarAPI(s.chainConfig, s.blockchain, s.chainDb, s.txPool, s.accountManager),
		}, {
			Namespace: "geth",
			Version:   "1.0",
			Service:   NewPublicGethAPI(s),
			Public:    true,
		},
	}
}

func (s *Ethereum) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Ethereum) Etherbase() (eb common.Address, err error) {
	eb = s.etherbase
	if eb.IsEmpty() {
		firstAccount, err := s.AccountManager().AccountByIndex(0)
		eb = firstAccount.Address
		if err != nil {
			return eb, fmt.Errorf("etherbase address must be explicitly specified")
		}
	}
	return eb, nil
}

// set in js console via admin interface or wrapper from cli flags
func (self *Ethereum) SetEtherbase(etherbase common.Address) {
	self.etherbase = etherbase
	self.miner.SetEtherbase(etherbase)
}

func (s *Ethereum) StopMining()         { s.miner.Stop() }
func (s *Ethereum) IsMining() bool      { return s.miner.Mining() }
func (s *Ethereum) Miner() *miner.Miner { return s.miner }

func (s *Ethereum) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Ethereum) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Ethereum) TxPool() *core.TxPool               { return s.txPool }
func (s *Ethereum) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Ethereum) ChainDb() ethdb.Database            { return s.chainDb }
func (s *Ethereum) DappDb() ethdb.Database             { return s.dappDb }
func (s *Ethereum) IsListening() bool                  { return true } // Always listening
func (s *Ethereum) EthVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Ethereum) NetVersion() int                    { return s.netVersionId }
func (s *Ethereum) ChainConfig() *core.ChainConfig     { return s.chainConfig }
func (s *Ethereum) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Ethereum) Protocols() []p2p.Protocol {
	return s.protocolManager.SubProtocols
}

// Start implements node.Service, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *Ethereum) Start(srvr *p2p.Server) error {
	s.protocolManager.Start(s.config.MaxPeers)
	s.netRPCService = NewPublicNetAPI(srvr, s.NetVersion())
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *Ethereum) Stop() error {
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	s.dappDb.Close()
	close(s.shutdownChan)

	return nil
}

// This function will wait for a shutdown and resumes main thread execution
func (s *Ethereum) WaitForShutdown() {
	<-s.shutdownChan
}

// HTTPClient returns the light http client used for fetching offchain docs
// (natspec, source for verification)
func (self *Ethereum) HTTPClient() *httpclient.HTTPClient {
	return self.httpclient
}

func (self *Ethereum) Solc() (*compiler.Solidity, error) {
	var err error
	if self.solc == nil {
		self.solc, err = compiler.New(self.SolcPath)
	}
	return self.solc, err
}

// set in js console via admin interface or wrapper from cli flags
func (self *Ethereum) SetSolc(solcPath string) (*compiler.Solidity, error) {
	self.SolcPath = solcPath
	self.solc = nil
	return self.Solc()
}

// upgradeChainDatabase ensures that the chain database stores block split into
// separate header and body entries.
func upgradeChainDatabase(db ethdb.Database) error {
	// Short circuit if the head block is stored already as separate header and body
	data, err := db.Get([]byte("LastBlock"))
	if err != nil {
		return nil
	}
	head := common.BytesToHash(data)

	if block := core.GetBlockByHashOld(db, head); block == nil {
		return nil
	}
	// At least some of the database is still the old format, upgrade (skip the head block!)
	glog.V(logger.Info).Info("Old database detected, upgrading...")

	if db, ok := db.(*ethdb.LDBDatabase); ok {
		blockPrefix := []byte("block-hash-")
		for it := db.NewIterator(); it.Next(); {
			// Skip anything other than a combined block
			if !bytes.HasPrefix(it.Key(), blockPrefix) {
				continue
			}
			// Skip the head block (merge last to signal upgrade completion)
			if bytes.HasSuffix(it.Key(), head.Bytes()) {
				continue
			}
			// Load the block, split and serialize (order!)
			block := core.GetBlockByHashOld(db, common.BytesToHash(bytes.TrimPrefix(it.Key(), blockPrefix)))

			if err := core.WriteTd(db, block.Hash(), block.DeprecatedTd()); err != nil {
				return err
			}
			if err := core.WriteBody(db, block.Hash(), block.Body()); err != nil {
				return err
			}
			if err := core.WriteHeader(db, block.Header()); err != nil {
				return err
			}
			if err := db.Delete(it.Key()); err != nil {
				return err
			}
		}
		// Lastly, upgrade the head block, disabling the upgrade mechanism
		current := core.GetBlockByHashOld(db, head)

		if err := core.WriteTd(db, current.Hash(), current.DeprecatedTd()); err != nil {
			return err
		}
		if err := core.WriteBody(db, current.Hash(), current.Body()); err != nil {
			return err
		}
		if err := core.WriteHeader(db, current.Header()); err != nil {
			return err
		}
	}
	return nil
}

func addMipmapBloomBins(db ethdb.Database) (err error) {
	const mipmapVersion uint = 2

	// check if the version is set. We ignore data for now since there's
	// only one version so we can easily ignore it for now
	var data []byte
	data, _ = db.Get([]byte("setting-mipmap-version"))
	if len(data) > 0 {
		var version uint
		if err := rlp.DecodeBytes(data, &version); err == nil && version == mipmapVersion {
			return nil
		}
	}

	defer func() {
		if err == nil {
			var val []byte
			val, err = rlp.EncodeToBytes(mipmapVersion)
			if err == nil {
				err = db.Put([]byte("setting-mipmap-version"), val)
			}
			return
		}
	}()
	latestBlock := core.GetBlock(db, core.GetHeadBlockHash(db))
	if latestBlock == nil { // clean database
		return
	}

	tstart := time.Now()
	glog.V(logger.Info).Infoln("upgrading db log bloom bins")
	for i := uint64(0); i <= latestBlock.NumberU64(); i++ {
		hash := core.GetCanonicalHash(db, i)
		if (hash == common.Hash{}) {
			return fmt.Errorf("chain db corrupted. Could not find block %d.", i)
		}
		err := core.WriteMipmapBloom(db, i, core.GetBlockReceipts(db, hash))
		if err != nil {
			return err
		}
	}
	glog.V(logger.Info).Infoln("upgrade completed in", time.Since(tstart))
	return nil
}
