package main

import (
  "os"
  "fmt"
  "strconv"
  "encoding/hex"
  "crypto/ecdsa"
  "math/big"
  "context"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/crypto"
  "github.com/ethereum/go-ethereum/ethclient"
)

const EthereumNodeURL = "http://127.0.0.1:8545"

func main() {
  if os.Args[1] == "genkeys" {
    genKeys()
  } else if os.Args[1] == "sendwei" { 
    fromPrivateKeyHex := os.Args[2]
    toAddressHex := os.Args[3]

    amountWei := new(big.Int)
    amountWei, ok := amountWei.SetString(os.Args[4], 10)
    if !ok {
      fmt.Println("Error converting amountWei to big Int")
    }

    gasLimit, err := strconv.ParseUint(os.Args[5], 10, 64)
    if err != nil {
      fmt.Println(err)
    }    

    gasPrice := new(big.Int)
    gasPrice, ok = gasPrice.SetString(os.Args[6], 10)
    if !ok {
      fmt.Println("Error converting gasPrice to big Int")
    }

    sendWei(fromPrivateKeyHex, toAddressHex, amountWei, gasLimit, gasPrice)    
  } else if os.Args[1] == "help"{
    printHelp()
  } else {
    fmt.Printf("'%v' is not a valid command. Run 'ethtool help' for a list of available commands.\n", os.Args[1])
  }
}

func genKeys(){
  // Create an account
  key, err := crypto.GenerateKey()
  if err != nil {
    fmt.Println(err)
  }

  // Get the address
  address := crypto.PubkeyToAddress(key.PublicKey).Hex()

  // Get the private key
  privateKey := hex.EncodeToString(key.D.Bytes())

  fmt.Println("Public key:", address)
  fmt.Println("Private key:", privateKey)
}


func sendWei(fromPrivateKeyHex string, toAddressHex string, amountWei *big.Int, gasLimit uint64, gasPrice *big.Int) {
  // Connect to geth
  client, err := ethclient.Dial(EthereumNodeURL)
  if err != nil {
    fmt.Println(err)
  }

  // Convert private key to ECDSA format
  privateKey, err := crypto.HexToECDSA(fromPrivateKeyHex)
  if err != nil {
    fmt.Println(err)
  }

  // Get the public associated public key
  publicKey := privateKey.Public()

  // Convert public key to ECDSA format
  publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
  if !ok {
    fmt.Println("error casting public key to ECDSA")
  }

  // Get the address from the public key
  fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

  // Get the nonce for this transaction from the eth client
  nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
  if err != nil {
    fmt.Println(err)
  }

  // ToDo: Default to suggested gas price if none supplied
  // gasPrice, err = client.SuggestGasPrice(context.Background())
  // if err != nil {
  //   fmt.Println(err)
  // }

  var data []byte // to fill in with message saying this was processed with gas
  toAddress := common.HexToAddress(toAddressHex)

  transaction := types.NewTransaction(nonce, toAddress, amountWei, gasLimit, gasPrice, data)
  signedTransaction, err := types.SignTx(transaction, types.HomesteadSigner{}, privateKey)
  if err != nil {
    fmt.Println(err)
  }

  // Send the transaction to geth to gossip
  err = client.SendTransaction(context.Background(), signedTransaction)
  if err != nil {
    fmt.Println(err)
  }
}

func printHelp(){
  fmt.Print(`NAME:
   ethtool - ethereum development tools

USAGE:
   ethtool command [arguments...]
   
COMMANDS:
   sendwei           Sends a transaction. e.g sendwei <from private key> <to address> <amount> <gas limit> <gas price>
   genkeys           Generates an public and private key pair.`)
}
