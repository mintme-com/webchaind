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

## GPU mining

The ethash algorithm is memory hard and in order to fit the DAG into memory, it needs 1-2GB of RAM on each GPU. If you get ``Error GPU mining. GPU memory fragmentation?`` you don’t have enough memory.

## Installing ethminer

To get ethminer, you need to install the ethminer binary package or build it from source. See https://github.com/ethereum-mining/ethminer/#build for the official ethminer build/install instructions. At the time of writing, ethminer only provides a binary for Microsoft Windows.

## Using ethminer with geth

First create an account to hold your block rewards.

`geth account new`

Follow the prompts and enter a good password. **DO NOT FORGET YOUR PASSWORD**. Also take note of the public Ethereum address which is printed at the end of the account creation process. In the following examples, we will use 0xC95767AC46EA2A9162F0734651d6cF17e5BfcF10 as the example address.

Now start geth and wait for it to sync the blockchain. This will take quite a while.

`geth --http --miner.etherbase 0xC95767AC46EA2A9162F0734651d6cF17e5BfcF10`

To monitor the syncing, in another terminal you can attach the geth JavaScript console to the running node like so:

`geth attach https://127.0.0.1:8545`

and then at the > prompt type

`eth.syncing`

You’ll see something like the example output below – it’s a two stage process as described in much more detail in our FAQ. In the first stage, the difference between the “currentBlock” and the “highestBlock” will decrease until they are almost equal. It will then look stuck and appear as never becoming equal. But you should see “pulledStates” rising to equal “knownStates.” When both are equal, you are synced.

Example output of first stage of block downloading:

```
{
  currentBlock: 10707814,
  highestBlock: 13252182,
  knownStates: 0,
  pulledStates: 0,
  startingBlock: 3809258 }
```

You will import up to the highestBlock and knownStates. Block importing will stop ~64 blocks behind head and finish importing states.

Once all states are downloaded, geth will switch into a full node and sync the remaining ~64 blocks fully, as well as new ones. In this context, eth.syncing returns false once synced.

Now we’re ready to start mining. In a new terminal session, run ethminer and connect it to geth:

## License

The core-geth library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The core-geth binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
