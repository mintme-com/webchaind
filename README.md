## Webchaind

Official Go language implementation of the Webchain daemon.

## Install

### :rocket: From a release binary
The simplest way to get started running a node is to visit our [Releases page](https://github.com/webchain-network/webchaind/releases) and download a zipped executable binary (matching your operating system, of course), then moving the unzipped file `webchaind` to somewhere in your `$PATH`. Now you should be able to open a terminal and run `$ webchaind help` to make sure it's working.

### :hammer: Building the source

If your heart is set on the bleeding edge, install from source. However, please be advised that you may encounter some strange things, and we can't prioritize support beyond the release versions. Recommended for developers only.

#### Dependencies
Building webchaind requires both Go >=1.9 and a C compiler.

#### Get source and dependencies
`$ go get -v github.com/webchain-network/webchaind/...`

#### Installing command executables

To install...

- the full suite of utilities: `$ go install github.com/webchain-network/webchaind/cmd/...`
- just __webchaind__: `$ go install github.com/webchain-network/webchaind/cmd/webchaind`

Executables built from source will, by default, be installed in `$GOPATH/bin/`.

#### Building specific release
All the above commands results with building binaries from `HEAD`. To use a specific release/tag, use the following:
```
$ go get -d github.com/webchain-network/webchaind/...
$ cd $GOPATH/src/github.com/webchain-network/webchaind
$ git checkout <TAG OR REVISION>
$ go install -ldflags "-X main.Version="`git describe --tags` ./cmd/...
```

#### Using release source code tarball
Because of strict Go directory structure, the tarball needs to be extracted into the proper subdirectory under `$GOPATH`.
The following commands are an example of building the v4.1.1 release:
```
$ mkdir -p $GOPATH/src/github.com/webchain-network
$ cd $GOPATH/src/github.com/webchain-network
$ tar xzf /path/to/v0.1.0.tar.gz
$ mv v0.1.0 webchaind
$ cd webchaind
$ go install -ldflags "-X main.Version=v0.1.0" ./cmd/...
```

#### Building with [SputnikVM](https://github.com/ethereumproject/sputnikvm)
Have Rust (>= 1.21) and Golang (>= 1.9) installed.

> For __Linux__ and __macOS__:

```
cd $GOPATH/src/github.com/ethereumproject
git clone https://github.com/ethereumproject/sputnikvm-ffi
cd sputnikvm-ffi/c/ffi
cargo build --release
cp $GOPATH/src/github.com/ethereumproject/sputnikvm-ffi/c/ffi/target/release/libsputnikvm_ffi.a $GOPATH/src/github.com/ethereumproject/sputnikvm-ffi/c/libsputnikvm.a
```
And then build webchaind with CGO_LDFLAGS:

- In Linux:

```
cd $GOPATH/src/github.com/webchain-network/webchaind/cmd/webchaind
CGO_LDFLAGS="$GOPATH/src/github.com/ethereumproject/sputnikvm-ffi/c/libsputnikvm.a -ldl" go build -tags=sputnikvm .
```

- In macOS:

```
cd $GOPATH/src/github.com/webchain-network/webchaind/cmd/webchaind
CGO_LDFLAGS="$GOPATH/src/github.com/ethereumproject/sputnikvm-ffi/c/libsputnikvm.a -ldl -lresolv" go build -tags=sputnikvm .
```

> For __Windows__:

```
cd %GOPATH%\src\github.com\ethereumproject
git clone https://github.com/ethereumproject/sputnikvm-ffi
cd sputnikvm-ffi\c\ffi
cargo build --release
copy %GOPATH%\src\github.com\ethereumproject\sputnikvm-ffi\c\ffi\target\release\sputnikvm.lib %GOPATH%\src\github.com\ethereumproject\sputnikvm-ffi\c\sputnikvm.lib
```
And then build webchaind with CGO_LDFLAGS:
```
cd %GOPATH%\src\github.com\webchain-network\webchaind\cmd\webchaind
set CGO_LDFLAGS=-Wl,--allow-multiple-definition %GOPATH%\src\github.com\ethereumproject\sputnikvm-ffi\c\sputnikvm.lib -lws2_32 -luserenv
go build -tags=sputnikvm .
```

## Executables

This repository includes several wrappers/executables found in the `cmd` directory.

| Command    | Description |
|:----------:|-------------|
| **`webchaind`** | The main Webchain CLI client. It is the entry point into the Webchain network (main-, test-, or private net), capable of running as a full node (default) archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the Webchain network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports.|
| `abigen` | Source code generator to convert Webchain contract definitions into easy to use, compile-time type-safe Go packages. It operates on plain [Ethereum contract ABIs](https://github.com/ethereumproject/wiki/wiki/Ethereum-Contract-ABI) with expanded functionality if the contract bytecode is also available. However it also accepts Solidity source files, making development much more streamlined. Please see our [Native DApps](https://github.com/webchain-network/webchaind/wiki/Native-DApps-in-Go) wiki page for details. |
| `bootnode` | Stripped down version of our Webchain client implementation that only takes part in the network node discovery protocol, but does not run any of the higher level application protocols. It can be used as a lightweight bootstrap node to aid in finding peers in private networks. |
| `disasm` | Bytecode disassembler to convert EVM (Ethereum Virtual Machine) bytecode into more user friendly assembly-like opcodes (e.g. `echo "6001" | disasm`). For details on the individual opcodes, please see pages 22-30 of the [Ethereum Yellow Paper](http://gavwood.com/paper.pdf). |
| `evm` | Developer utility version of the EVM (Ethereum Virtual Machine) that is capable of running bytecode snippets within a configurable environment and execution mode. Its purpose is to allow insolated, fine graned debugging of EVM opcodes (e.g. `evm --code 60ff60ff --debug`). |
| `rlpdump` | Developer utility tool to convert binary RLP ([Recursive Length Prefix](https://github.com/ethereumproject/wiki/wiki/RLP)) dumps (data encoding used by the Webchain protocol both network as well as consensus wise) to user friendlier hierarchical representation (e.g. `rlpdump --hex CE0183FFFFFFC4C304050583616263`). |

## :green_book: Webchaind: the basics

### Data directory
By default, webchaind will store all node and blockchain data in a __parent directory__ depending on your OS:
- Linux: `$HOME/.webchain/`
- Mac: `$HOME/Library/Webchain/`
- Windows: `$HOME/AppData/Roaming/Webchain/`

__You can specify this directory__ with `--data-dir=$HOME/id/rather/put/it/here`.

Within this parent directory, webchaind will use a __/subdirectory__ to hold data for each network you run. The defaults are:
 - `/mainnet` for the Mainnet
 - `/morden` for the Morden Testnet

__You can specify this subdirectory__ with `--chain=mycustomnet`.

### Full node on the main Webchain network

```
$ webchaind
```

It's that easy! This will establish an WEB blockchain node and download ("sync") the full blocks for the entirety of the WEB blockchain. __However__, before you go ahead with plain ol' `webchaind`, we would encourage reviewing the following section...

#### :speedboat: `--fast`

The most common scenario is users wanting to simply interact with the Webchain network: create accounts; transfer funds; deploy and interact with contracts, and mine. For this particular use-case the user doesn't care about years-old historical data, so we can _fast-sync_ to the current state of the network. To do so:

```
$ webchaind --fast
```

Using webchaind in fast sync mode causes it to download only block _state_ data -- leaving out bulky transaction records -- which avoids a lot of CPU and memory intensive processing.

Fast sync will be automatically __disabled__ (and full sync enabled) when:
- your chain database contains *any* full blocks
- your node has synced up to the current head of the network blockchain

In case of using `--mine` together with `--fast`, webchaind will operate as described; syncing in fast mode up to the head, and then begin mining once it has synced its first full block at the head of the chain.

*Note:* To further increase webchaind performace, you can use a `--cache=512` flag to bump the memory allowance of the database (e.g. 512MB) which can significantly improve sync times, especially for HDD users. This flag is optional and you can set it as high or as low as you'd like, though we'd recommend the 512MB - 2GB range.

### Create or manage account(s)

Webchaind is able to create, import, update, unlock, and otherwise manage your private (encrypted) key files. Key files are in JSON format and, by default, stored in the respective chain folder's `/keystore` directory; you can specify a custom location with the `--keystore` flag.

```
$ webchaind account new
```

This command will create a new account and prompt you to enter a passphrase to protect your account.

Other `account` subcommands include:
```
SUBCOMMANDS:

        list    print account addresses
        new     create a new account
        update  update an existing account
        import  import a private key into a new account

```

Learn more at the [Accounts Wiki Page](https://github.com/webchain-network/webchaind/wiki/Managing-Accounts). If you're interested in using webchaind to manage a lot (~100,000+) of accounts, please visit the [Indexing Accounts Wiki page](https://github.com/webchain-network/webchaind/wiki/Indexing-Accounts).


### Interact with the Javascript console
```
$ webchaind console
```

This command will start up webchaind built-in interactive [JavaScript console](https://github.com/webchain-network/webchaind/wiki/JavaScript-Console), through which you can invoke all official [`web3` methods](https://github.com/ethereumproject/wiki/wiki/JavaScript-API) as well as webchaind own [management APIs](https://github.com/webchain-network/webchaind/wiki/Management-APIs). This too is optional and if you leave it out you can always attach to an already running webchaind instance with `webchaind attach`.

Learn more at the [Javascript Console Wiki page](https://github.com/webchain-network/webchaind/wiki/JavaScript-Console).


### And so much more!

For a comprehensive list of command line options, please consult our [CLI Wiki page](https://github.com/webchain-network/webchaind/wiki/Command-Line-Options).

## :orange_book: Webchaind: developing and advanced useage

### Morden Testnet
If you'd like to play around with creating Webchain contracts, you
almost certainly would like to do that without any real money involved until you get the hang of the entire system. In other words, instead of attaching to the main network, you want to join the **test** network with your node, which is fully equivalent to the main network, but with play-Ether only.

```
$ webchaind --chain=morden --fast --cache=512 console
```

The `--fast`, `--cache` flags and `console` subcommand have the exact same meaning as above and they are equally useful on the testnet too. Please see above for their explanations if you've skipped to here.

Specifying the `--chain=morden` flag will reconfigure your webchaind instance a bit:

 -  As mentioned above, webchaind will host its testnet data in a `morden` subfolder (`~/.webchain/morden`).
 - Instead of connecting the main Webchain network, the client will connect to the test network, which uses different P2P bootnodes, different network IDs and genesis states.

You may also optionally use `--testnet` or `--chain=testnet` to enable this configuration.

> *Note: Although there are some internal protective measures to prevent transactions from crossing over between the main network and test network (different starting nonces), you should make sure to always use separate accounts for play-money and real-money. Unless you manually move accounts, webchaind
will by default correctly separate the two networks and will not make any accounts available between them.*

### Programatically interfacing webchaind nodes

As a developer, sooner rather than later you'll want to start interacting with webchaind and the Webchain network via your own programs and not manually through the console. To aid this, webchaind has built in support for a JSON-RPC based APIs ([standard APIs](https://github.com/ethereumproject/wiki/wiki/JSON-RPC) and
[Webchaind specific APIs](https://github.com/webchain-network/webchaind/wiki/Management-APIs)). These can be exposed via HTTP, WebSockets and IPC (unix sockets on unix based platroms, and named pipes on Windows).

The IPC interface is enabled by default and exposes all the APIs supported by webchaind, whereas the HTTP and WS interfaces need to manually be enabled and only expose a subset of APIs due to security reasons. These can be turned on/off and configured as you'd expect.

HTTP based JSON-RPC API options:

  * `--rpc` Enable the HTTP-RPC server
  * `--rpc-addr` HTTP-RPC server listening interface (default: "localhost")
  * `--rpc-port` HTTP-RPC server listening port (default: 8545)
  * `--rpc-api` API's offered over the HTTP-RPC interface (default: "eth,net,web3")
  * `--rpc-cors-domain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
  * `--ws` Enable the WS-RPC server
  * `--ws-addr` WS-RPC server listening interface (default: "localhost")
  * `--ws-port` WS-RPC server listening port (default: 8546)
  * `--ws-api` API's offered over the WS-RPC interface (default: "eth,net,web3")
  * `--ws-origins` Origins from which to accept websockets requests
  * `--ipc-disable` Disable the IPC-RPC server
  * `--ipc-api` API's offered over the IPC-RPC interface (default: "admin,debug,eth,miner,net,personal,shh,txpool,web3")
  * `--ipc-path` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to connect via HTTP, WS or IPC to a webchaind node configured with the above flags and you'll need to speak [JSON-RPC](http://www.jsonrpc.org/specification) on all transports. You can reuse the same connection for multiple requests!

> Note: Please understand the security implications of opening up an HTTP/WS based transport before doing so! Hackers on the internet are actively trying to subvert Webchain nodes with exposed APIs! Further, all browser tabs can access locally running webservers, so malicious webpages could try to subvert locally available APIs!*

## Contribution

Thank you for considering to help out with the source code!

The core values of democratic engagement, transparency, and integrity run deep with us. We welcome contributions from everyone, and are grateful for even the smallest of fixes.  :clap:

If you'd like to contribute to webchaind, please fork, fix, commit and send a pull request for the maintainers to review and merge into the main code base.

Please see the [Wiki](https://github.com/webchain-network/webchaind/wiki) for more details on configuring your environment, managing project dependencies, and testing procedures.

## License

The webchaind library (i.e. all code outside of the `cmd` directory) is licensed under the [GNU Lesser General Public License v3.0](http://www.gnu.org/licenses/lgpl-3.0.en.html), also included in our repository in the `COPYING.LESSER` file.

The webchaind binaries (i.e. all code inside of the `cmd` directory) is licensed under the [GNU General Public License v3.0](http://www.gnu.org/licenses/gpl-3.0.en.html), also included in our repository in the `COPYING` file.
