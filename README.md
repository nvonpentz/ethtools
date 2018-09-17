# ethtools
Ethereum development tools to be used with a local Ethereum JSON RPC endpoint `localhost:8545`.

```
NAME:
   ethtools - ethereum development tools

USAGE:
   ethtools command [arguments...]
   
COMMANDS:
   sendwei           Sends a transaction. e.g sendwei <from private key> <to address> <amount> <gas limit> <gas price>
   genkeys           Prints out a public and private key pair
```

# setup

* Install [Go](https://golang.org/doc/install)
* `git clone github.com/nvonpentz/ethtools.git`
* `go install`
