package spend

import (
	"demo_btc_assets/client"
	"demo_btc_assets/wallet"
	"log"
	"testing"
)

func TestSendCoinWithP2TR(t *testing.T) {
	client := client.New()
	wallet.OpenWallet(client)
	_, err := SendCoinWithP2TR(client, ADDRESS_TEST, 500)
	if err != nil {
		log.Fatal(err)
	}
}
