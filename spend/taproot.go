package spend

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const TAPROOT_ADDRESS = "tb1p6s6mmwnjpea9x4ju6zhq4l4fslcmgqrxeuvjutvpfcltnnfp58lqdhlu6g"

/*
Send 0.0001 coin to "tb1p6s6mmwnjpea9x4ju6zhq4l4fslcmgqrxeuvjutvpfcltnnfp58lqdhlu6g" address
this is the TxIndex =  3f7d519759c210af65ee124ae665bddff2ad2b3ddce4c5de7a35d2b5c8a316eb
this is Vout = 1
*/

const TX_ID = "3f7d519759c210af65ee124ae665bddff2ad2b3ddce4c5de7a35d2b5c8a316eb"
const VOUT = 1

func SendCoinWithP2TR(client *rpcclient.Client, address string, amount uint) (*chainhash.Hash, error) {
	// calc commitTxHash this is the hash of the transaction send coin to Taproot Address
	// the for loop is revert the hash, because the NewHash() function
	txHash, _ := hex.DecodeString(TX_ID)
	for i, j := 0, len(txHash)-1; i < j; i, j = i+1, j-1 {
		txHash[i], txHash[j] = txHash[j], txHash[i]
	}
	commitTxHash, err := chainhash.NewHash(txHash)
	log.Println(commitTxHash)
	commitTx, err := client.GetRawTransaction(commitTxHash)
	if err != nil {
		log.Println("cannot get raw transaction")
		return nil, err
	}

	defaultAddress, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	if err != nil {
		log.Println("Cannot get default address!")
		return nil, err
	}
	log.Println(defaultAddress)

	wif, err := client.DumpPrivKey(defaultAddress)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	builder := txscript.NewScriptBuilder()
	builder.AddData(schnorr.SerializePubKey(wif.PrivKey.PubKey()))
	builder.AddOp(txscript.OP_CHECKSIG)
	hashLockScript, err := builder.Script()

	tapLeaf := txscript.NewBaseTapLeaf(hashLockScript)
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapLeaf)
	tapScriptRootHash := tapScriptTree.LeafMerkleProofs[0].RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(wif.PrivKey.PubKey(), tapScriptRootHash[:])
	log.Println(outputKey)

	tx := wire.NewMsgTx(2)
	tx.AddTxIn(&wire.TxIn{
		PreviousOutPoint: wire.OutPoint{
			Hash:  *commitTxHash,
			Index: VOUT,
		},
	})

	pkScript, _ := txscript.PayToAddrScript(defaultAddress)
	txOut := &wire.TxOut{
		Value: 500, PkScript: pkScript,
	}
	tx.AddTxOut(txOut)

	inputFetcher := txscript.NewCannedPrevOutputFetcher(
		commitTx.MsgTx().TxOut[VOUT].PkScript,
		commitTx.MsgTx().TxOut[VOUT].Value,
	)
	sigHashes := txscript.NewTxSigHashes(tx, inputFetcher)
	sig, err := txscript.RawTxInTapscriptSignature(
		tx, sigHashes, 0, txOut.Value,
		txOut.PkScript, tapLeaf, txscript.SigHashDefault,
		wif.PrivKey,
	)

	ctrlBlock := tapScriptTree.LeafMerkleProofs[0].ToControlBlock(wif.PrivKey.PubKey())
	ctrlBlockBytes, err := ctrlBlock.ToBytes()
	tx.TxIn[0].Witness = wire.TxWitness{
		sig, hashLockScript, ctrlBlockBytes,
	}
	var buf bytes.Buffer
	_ = tx.Serialize(&buf)
	fmt.Printf("Raw tx: %x\n", buf.Bytes())

	hashTx, err := client.SendRawTransaction(tx, true)
	if err != nil {
		log.Println("Cannot send raw transaction!")
		return nil, err
	}
	log.Println("Success!!! : ", hashTx)
	return hashTx, nil
}
