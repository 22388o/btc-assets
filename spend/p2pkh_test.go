package spend

import (
	"demo_btc_assets/client"
	"demo_btc_assets/wallet"
	"log"
	"testing"
)

const ADDRESS_TEST = "mgUE4a5G1tnKLxmRTPCzQSURueDwzPnMSm"
const AMOUNT_TEST = 100000

func TestP2PKH(t *testing.T) {

	client := client.New()
	wallet.OpenWallet(client)

	startBalance, err := client.GetBalance("default")
	if err != nil {
		log.Fatal("Cannot get start Balance")
	}

	SendCoinWithP2PKH(client, ADDRESS_TEST, AMOUNT_TEST)

	endBalance, err := client.GetBalance("default")
	if err != nil {
		log.Fatal("Cannot get end Balance")
	}

	if startBalance-AMOUNT_TEST-FEE != endBalance {
		log.Fatal("balance isn't decrease")
	}

}
