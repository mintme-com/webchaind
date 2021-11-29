## CoreGeth: An Ethereum Protocol Provider

> An [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) downstream effort to make the Ethereum Protocol accessible and extensible for a diverse ecosystem.

Priority is given to reducing opinions around chain configuration, IP-based feature implementations, and API predictability.
Upstream development from [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) is merged to this repository regularly,
 usually at every upstream tagged release. Every effort is made to maintain seamless compatibility with upstream source, including compatible RPC, JS, and CLI
 APIs, data storage locations and schemas, and, of course, interoperable node protocols. Applicable bug reports, bug fixes, features, and proposals should be
 made upstream whenever possible.

[![OpenRPC](https://img.shields.io/static/v1.svg?label=OpenRPC&message=1.14.0&color=blue)](#openrpc-discovery)
[![API Reference](https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667)](https://godoc.org/github.com/etclabscore/core-geth)
[![Go Report Card](https://goreportcard.com/badge/github.com/etclabscore/core-geth)](https://goreportcard.com/report/github.com/etclabscore/core-geth)
[![Travis](https://travis-ci.org/etclabscore/core-geth.svg?branch=master)](https://travis-ci.org/etclabscore/core-geth)
[![Gitter](https://badges.gitter.im/core-geth/community.svg)](https://gitter.im/core-geth/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Network/provider comparison

Networks supported by the respective go-ethereum packaged `geth` program.

| Ticker | Consensus         | Network                               | core-geth                                                | ethereum/go-ethereum |
| ---    | ---               | ---                                   | ---                                                      | ---                  |
| ETC    | :zap:             | Ethereum Classic                      | :heavy_check_mark:                                       |                      |
| ETH    | :zap:             | Ethereum (Foundation)                 | :heavy_check_mark:                                       | :heavy_check_mark:   |
| -      | :zap: :handshake: | Private chains                        | :heavy_check_mark:                                       | :heavy_check_mark:   |
|        | :zap:             | Mordor (Geth+Parity ETH PoW Testnet)  | :heavy_check_mark:                                       |                      |
|        | :zap:             | Morden (Geth+Parity ETH PoW Testnet)  |                                                          |                      |
|        | :zap:             | Ropsten (Geth+Parity ETH PoW Testnet) | :heavy_check_mark:                                       | :heavy_check_mark:   |
|        | :handshake:       | Rinkeby (Geth-only ETH PoA Testnet)   | :heavy_check_mark:                                       | :heavy_check_mark:   |
|        | :handshake:       | Goerli (Geth+Parity ETH PoA Testnet)  | :heavy_check_mark:                                       | :heavy_check_mark:   |
|        | :handshake:       | Kotti (Geth+Parity ETC PoA Testnet)   | :heavy_check_mark:                                       |                      |
|        | :handshake:       | Kovan (Parity-only ETH PoA Testnet)   |                                                          |                      |
|        |                   | Tobalaba (EWF Testnet)                |                                                          |                      |
|        |                   | Ephemeral development PoA network     | :heavy_check_mark:                                       | :heavy_check_mark:   |
| MINTME | :zap:             | MintMe.com Coin                       | :heavy_check_mark:                                       |                      |

- :zap: = __Proof of Work__
- :handshake: = __Proof of Authority__

<a name="ellaism-footnote">1</a>: This is originally an [Ellaism
Project](https://github.com/ellaism). However, A [recent hard
fork](https://github.com/ellaism/specs/blob/master/specs/2018-0003-wasm-hardfork.md)
makes Ellaism not feasible to support with go-ethereum any more. Existing
Ellaism users are asked to switch to
[Parity](https://github.com/paritytech/parity).

<a name="configuration-capable">2</a>: Network not supported by default, but network configuration is possible. Make a PR!

## Documentation

- CoreGeth documentation is available [here](https://etclabscore.github.io/core-geth).
  + Getting Started [Installation](https://etclabscore.github.io/core-geth/getting-started/installation) and [CLI](https://etclabscore.github.io/core-geth/getting-started/run-cli)
  + [JSONRPC API](https://etclabscore.github.io/core-geth/apis/jsonrpc-apis)
  + [Developers](https://etclabscore.github.io/core-geth/developers/build-from-source)
  + [Tutorials](https://etclabscore.github.io/core-geth/tutorials/private-network)
- Further [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) documentation about can be found [here](https://geth.ethereum.org/docs/).
- Documentation about documentation lives [here](./docs/developers/documentation.md).

## Contribution

Thank you for considering to help out with the source code! We welcome contributions
from anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to core-geth, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base. If you wish to submit
more complex changes though, please check up with the core devs first on [our gitter channel](https://gitter.im/etclabscore/core-geth)
to ensure those changes are in line with the general philosophy of the project and/or get
some early feedback which can make both your efforts much lighter as well as our review
and merge procedures quick and simple.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting)
   guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary)
   guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 * Commit messages should be prefixed with the package(s) they modify.
   * E.g. "eth, rpc: make trace configs optional"

Please see the [Developers' Guide](https://github.com/ethereum/go-ethereum/wiki/Developers'-Guide)
for more details on configuring your environment, managing project dependencies, and
testing procedures.

## Mining

Mining is the process through which new blocks are created. Geth actually creates new blocks all the time, but these blocks need to be secured through proof-of-work so they will be accepted by other nodes. Mining is all about creating these proof-of-work values.

The proof-of-work computation can be performed in multiple ways. Geth includes a CPU miner, which does mining within the geth process. We discourage using the CPU miner with the Ethereum mainnet. If you want to mine real ether, use GPU mining. Your best option for doing that is the ethminer software.

Always ensure your blockchain is fully synchronised with the chain before starting to mine, otherwise you will not be mining on the correct chain and your block rewards will not be valueable.

## Solo mining

When you start up your ethereum node with geth it is not mining by default. To start it in mining mode, you use the --mine command-line flag. The --miner.threads parameter can be used to set the number parallel mining threads (defaulting to the total number of processor cores).

`geth --mine --miner.threads=4`

You can also start and stop CPU mining at runtime using the console. miner.start takes an optional parameter for the number of miner threads.
```
> miner.start(8)
true
> miner.stop()
true
```

Note that mining for real ether only makes sense if you are in sync with the network (since you mine on top of the consensus block). Therefore the eth blockchain downloader/synchroniser will delay mining until syncing is complete, and after that mining automatically starts unless you cancel your intention with miner.stop().

In order to earn ether you must have your etherbase (or coinbase) address set. This etherbase defaults to your primary account. If you don’t have an etherbase address, then geth --mine will not start up.

You can set your etherbase on the command line:

`geth --miner.etherbase '0xC95767AC46EA2A9162F0734651d6cF17e5BfcF10' --mine 2>> geth.log`

You can reset your etherbase on the console too:

`> miner.setEtherbase(eth.accounts[2])`

Note that your etherbase does not need to be an address of a local account, just an existing one.

There is an option to add extra data (32 bytes only) to your mined blocks. By convention this is interpreted as a unicode string, so you can set your short vanity tag.

`> miner.setExtra("ΞTHΞЯSPHΞЯΞ")`

You can check your hashrate with miner.hashrate, the result is in H/s (Hash operations per second).

```
> eth.hashrate
712000
```

After you successfully mined some blocks, you can check the ether balance of your etherbase account. Now assuming your etherbase is a local account:

```
> eth.getBalance(eth.coinbase).toNumber();
'34698870000000'
```

You can check which blocks are mined by a particular miner (address) with the following code snippet on the console:
```
> function minedBlocks(lastn, addr) {
    addrs = [];
    if (!addr) {
        addr = eth.coinbase
    }
    limit = eth.blockNumber - lastn
    for (i = eth.blockNumber; i >= limit; i--) {
        if (eth.getBlock(i).miner == addr) {
            addrs.push(i)
        }
    }
    return addrs
}
// scans the last 1000 blocks and returns the blocknumbers of blocks mined by your coinbase
// (more precisely blocks the mining reward for which is sent to your coinbase).
> minedBlocks(1000, eth.coinbase)
[352708, 352655, 352559]
```
Note that it will happen often that you find a block yet it never makes it to the canonical chain. This means when you locally include your mined block, the current state will show the mining reward credited to your account, however, after a while, the better chain is discovered and we switch to a chain in which your block is not included and therefore no mining reward is credited. Therefore it is quite possible that as a miner monitoring their coinbase balance will find that it may fluctuate quite a bit.

The logs show locally mined blocks confirmed after 5 blocks. At the moment you may find it easier and faster to generate the list of your mined blocks from these logs.

## Pool mining

We download miner binary for your system from here https://github.com/webchain-network/webchain-miner/releases

Then unpack it and in dir of unpacked miner we see file config.json , which is our miner configuration file, which contains main mining setting, like wallet number, threads , max-cpu-usage , etc.

The main thing we need to do is to set your wallet number where we will get reward to config.json

So we put our wallet number created by webchaind to "user" field of config.json, it shoud be set like this:

"user": "0x918e173c8426593bd37d5bc7d03f17dcc154cd5b",

Now we can save config.json , start mining and get rewards

We ready to start miner, it done by command: ./webchain-miner (Linux) or webchain-miner.exe (Windows)

Now we will see screen like that:

`./webchain-miner`

    VERSIONS: webchain-miner/2.6.2 libuv/1.20.3-dev gcc/6.3.0
    CPU: Intel(R) Xeon(R) CPU D-1521 @ 2.40GHz (1) x64 AES-NI
    CPU L2/L3: 1.0 MB/6.0 MB
    THREADS: 5, cryptonight-webchain, av=1, donate=5%
    POOL #1: pool2.webchain.network:2222
    COMMANDS: hashrate, pause, resume [2018-05-10 16:54:36] use pool pool2.webchain.network:2222 212.32.255.73 [2018-05-10 16:54:36] new job from pool2.webchain.network:2222 diff 5000 algo cn-web/1 [2018-05-10 16:54:36] READY (CPU) threads 5(5) huge pages 0/5 0% memory 10.0 MB [2018-05-10 23:26:02] speed 2.5s/60s/15m 14.6 14.8 13.6 H/s max: 20.3 H/s [2018-05-10 23:26:23] accepted (1/0) diff 5000 (170 ms)

Short description:

Pool is service which devide block finding task from webchain network between all miners in pool and also divide reward based on miners hashrate

Job is mining block finding task from network H/s is how much hashes in second your CPU can generate to find correct solution of block task accepted means your CPU found correct solution of task , send it to pool and it can be one of multiply correct block solution for what network will reward pool by 50 WEB, which will divide between miners based on their hasrate

So all you need is to run miner and wait some time (maybe 5 min, maybe 1 hour) when your miner find accepted shares and pool will give some part of reward coins to your account.

Step 4. Pool and payments

Currently we have official pool - pool.webchain.network:3333

So we start mining and then open https://pool.webchain.network/ in search string below "Your Stats & Payment History" we put our wallet number: `0x918e173c8426593bd37d5bc7d03f17dcc154cd5b`

and after miner send his first accepted share we will see realtime stats on top of page

First of all you will get coins in Immature Balance, it means you got coins from network, but they need confirmation. Pending Balance is your money already which wait pool payment schedule to pay to your wallet Total Paid: is how much money already sent from pool to your wallet Last Share Submitted means how much time ago your miner(s) submit accepted share to pool Workers Online means how much computers mine with the same wallet Hashrate (30m) and (3h) shows your common average from all your computers where you mining with this wallet number Blocks Found means how much block was found by your share solution Total payments means how much payment schedules was completed from pool with payment to your wallet Your round share is how much job of block solution your miner(s) do in comparison with all miners in pool

So when you will see some amount in Total paid, you can go back to your webchaind console and write command again:

`web3.fromWei(webchain.getBalance(webchain.coinbase), "ether")`

`5`

You will see other amount then 5 but it will be different from 0

If not, there can be some reason. First - your webchaind not synchronizing with network. So you need some time to wait (maybe about 1 hour) Also you can try to run

`webchaind --fast console`

or you can download blockchain file here: https://webchain.network/blockchain.raw and import it by executing:

`$ webchaind --fakepow import <path where you downloaded the blockchain>/blockchain.raw`

and retry to get balance, or wait some time and try again

Second - your local time is not correct Check please your OS for how to resync your clock (example `sudo ntpdate -s time.nist.gov`) because even 12 seconds too fast can lead to 0 peers.

Also you can see your real time balance of coins at https://explorer.webchain.network, enter you wallet number at the right top of page and you will how much coins you have for now

How to transer funds you can see at https://github.com/mintme-com/webchaind/wiki/Transfers

## License

The core-geth library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The core-geth binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
