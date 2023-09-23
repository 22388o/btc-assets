package main

import (
	"context"
	"crypto/sha256"
	"demo_btc_assets/address"
	"demo_btc_assets/mssmt"
	"demo_btc_assets/spend"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
)

type Server interface {
	CreateWallet(name, passphrase string) error
	MintAsset(name string, amount uint64) (*Asset, error)
}

type server struct {
	client *rpcclient.Client
}

var _ Server = (*server)(nil)

func NewServer() Server {
	certPath := filepath.Join(btcutil.AppDataDir("btcwallet", false), "rpc.cert")
	cert, err := os.ReadFile(certPath)
	if err != nil {
		panic(err)
	}

	client, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         "localhost:18332",
		Params:       "testnet3",
		Endpoint:     "ws",
		User:         "admin",
		Pass:         "admin123",
		Certificates: cert,
	}, nil)
	if err != nil {
		panic(err)
	}

	err = client.WalletPassphrase("admin", 100)
	if err != nil {
		fmt.Println("Cannot unlock wallet with passphrase")
	}

	return &server{
		client: client,
	}
}

func (s *server) CreateWallet(name, passphrase string) error {
	_, err := s.client.CreateWallet(name, rpcclient.WithCreateWalletPassphrase(passphrase))
	if err != nil {
		return err
	}

	return nil
}

type Asset struct {
	Name         string
	Amount       uint64
	ScriptPubkey *btcec.PublicKey
}

func (s *server) ComputeAssetDataByte(data *Asset) ([]byte, [32]byte) {
	h := sha256.New()
	_, _ = h.Write([]byte(data.Name))
	_, _ = h.Write(schnorr.SerializePubKey(data.ScriptPubkey))
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, [32]byte{}
	}

	return rawData, *(*[32]byte)(h.Sum(nil))
}

func (s *server) MintAsset(name string, amount uint64) (*Asset, error) {
	addr, err := s.client.GetAccountAddress("default")
	if err != nil {
		fmt.Println(err)
	}

	wif, err := s.client.DumpPrivKey(addr)
	if err != nil {
		fmt.Println(err)
	}

	asset := &Asset{Name: name, Amount: amount, ScriptPubkey: wif.PrivKey.PubKey()}
	tree := mssmt.NewCompactedTree(mssmt.NewDefaultStore())
	rawData, key := s.ComputeAssetDataByte(asset)

	leafNode := mssmt.NewLeafNode(rawData, asset.Amount)
	tree.Insert(context.TODO(), key, leafNode)

	root, err := tree.Root(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	tapAddr, err := address.CreateTaprootAddressWithData(asset.ScriptPubkey, root.NodeHash())
	if err != nil {
		fmt.Println(err)
	}

	chainHash, err := spend.SendCoinWithP2PKH(s.client, tapAddr.String(), 1000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(chainHash)

	return asset, nil
}
