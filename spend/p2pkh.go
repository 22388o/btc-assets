package spend

import (
	"bytes"
	"fmt"
	"log"
	"math"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const FEE btcutil.Amount = 500

func SendCoinWithP2PKH(client *rpcclient.Client, address string, amount uint) (*chainhash.Hash, error) {
	btcAmount := btcutil.Amount(amount)
	balance, err := client.GetBalance("*") // * mean all accounts
	if err != nil {
		log.Println("Can not get the balance of the wallet!")
		return nil, err
	}

	if balance-FEE < btcAmount {
		log.Println("You don't have enough coin to send!")
		return nil, err
	}

	unspents, err := client.ListUnspent()
	if err != nil {
		log.Println("Can not get the list unspent!")
	}

	rawTx := wire.NewMsgTx(2)

	// inputs := make([]btcjson.TransactionInput, 0, len(unspents))
	var inputAmount btcutil.Amount = 0
	for _, value := range unspents {
		txHash, err := chainhash.NewHashFromStr(value.TxID)
		if err != nil {
			log.Println("Can not convert TxID to type chainhash.Hash")
			return nil, err
		}

		txIn := wire.NewTxIn(&wire.OutPoint{Hash: *txHash, Index: value.Vout}, nil, nil)
		rawTx.AddTxIn(txIn)
		inputAmount = inputAmount + btcutil.Amount(value.Amount*math.Pow10(8))

		log.Printf("Inputamount: %d", inputAmount)
		if inputAmount-FEE > btcAmount {
			break
		}
	}

	destinationAddress, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	if err != nil {
		log.Println("Your address is invalid!")
		return nil, err
	}

	changeAddress, err := client.GetNewAddress("default")
	if err != nil {
		log.Println("Can not get the change address!")
	}

	// add output for the receiver
	// builder := txscript.NewScriptBuilder()
	// builder.AddData(destinationAddress.ScriptAddress())
	// builder.AddOp(txscript.OP_CHECKSIG)
	// hashLockScript, err := builder.Script()
	hashLockScript, _ := txscript.PayToAddrScript(destinationAddress)
	txOut := wire.NewTxOut(int64(btcAmount), hashLockScript)
	rawTx.AddTxOut(txOut)

	// add output for chagne
	changeLockScript, _ := txscript.PayToAddrScript(changeAddress)
	txOut = wire.NewTxOut(int64(inputAmount-btcAmount-FEE), changeLockScript)
	rawTx.AddTxOut(txOut)

	finalTx, isSign, err := client.SignRawTransaction(rawTx)
	if err != nil {
		log.Println("Cannot sign rawTx!")
		return nil, err
	}

	if !isSign {
		var buff bytes.Buffer
		finalTx.Serialize(&buff)
		fmt.Printf("Raw tx: %x\n", buff.Bytes())
		log.Println("RawTx is n't Sign!")
		return nil, nil
	}

	commitTxHash, err := client.SendRawTransaction(finalTx, false)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("transaction hash", commitTxHash)
	return commitTxHash, nil
}
