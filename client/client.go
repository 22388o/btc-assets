package client

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
)

func New() *rpcclient.Client {
	certHomeDir := btcutil.AppDataDir("btcwallet", false)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18332",
		Endpoint:     "ws",
		User:         "eeyHzXbqKIxOrZquE5c9POR/8mI=",
		Pass:         "q7pfGKtXrHE95qpS3oFGdAiAAXM=",
		Certificates: certs,
		Params:       "testnet3",
	}
	client, err := rpcclient.New(connCfg, nil)
	return client
}
