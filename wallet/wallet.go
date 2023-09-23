package wallet

import (
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

const WALLET_PASSPHRASE = "161002"

func OpenWallet(client *rpcclient.Client) error {
	err := client.WalletPassphrase(WALLET_PASSPHRASE, 100)
	if err != nil {
		log.Println("Cannot open wallet!")
		return err
	}
	log.Println("Open wallet success!")
	return nil
}
