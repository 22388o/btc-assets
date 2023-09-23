package address

import (
	"demo_btc_assets/client"
	"demo_btc_assets/wallet"
	"encoding/hex"
	"log"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

const DATA = "3f7d519759c210af65ee124ae665bddff2ad2b3ddce4c5de7a35d2b5c8a316eb"

func TestCreateTaprootAddress(t *testing.T) {
	client := client.New()
	wallet.OpenWallet(client)
	address, err := CreateTaprootAddress(client)
	if err != nil {
		log.Fatal("Cannot create taproot address")
	}

	log.Println(address)
}

func TestCreateTaprootAddressWithData(t *testing.T) {
	client := client.New()
	wallet.OpenWallet(client)
	defaultAddress, err := btcutil.DecodeAddress(ADDRESS_TEST, &chaincfg.TestNet3Params)
	if err != nil {
		log.Fatal("Can not decode address")
	}
	wif, err := client.DumpPrivKey(defaultAddress)
	if err != nil {
		log.Fatal("Can not dump privkey!")
	}

	data, _ := hex.DecodeString(DATA)

	address, err := CreateTaprootAddressWithData(wif.PrivKey.PubKey(), [32]byte(data[:]))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(address)
}

// address : tb1pyz23zhkq46wgwrkeun6gyxrp45jd5dhws22xxde06egcqfw9qjyqgy4ra3
// txId : f24b27a8061f7dc1f863c8ad12c1cf928c3efcf16073e67358895a2a3d6f1cc5
