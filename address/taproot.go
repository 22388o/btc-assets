package address

import (
	"log"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

const ADDRESS_TEST = "mgUE4a5G1tnKLxmRTPCzQSURueDwzPnMSm"

func CreateTaprootAddress(client *rpcclient.Client) (*btcutil.AddressTaproot, error) {
	defaultAddress, err := btcutil.DecodeAddress(ADDRESS_TEST, &chaincfg.TestNet3Params)
	if err != nil {
		log.Println("Cannot decode Address!")
		return nil, err
	}
	wif, err := client.DumpPrivKey(defaultAddress)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	builder := txscript.NewScriptBuilder()
	builder.AddData(schnorr.SerializePubKey(wif.PrivKey.PubKey()))
	builder.AddOp(txscript.OP_CHECKSIG)
	script, err := builder.Script()
	if err != nil {
		log.Println("Cannot get script in builder")
		return nil, err
	}

	tapleaf := txscript.NewBaseTapLeaf(script)
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapleaf)
	tapScriptRootHash := tapScriptTree.LeafMerkleProofs[0].RootNode.TapHash()

	outputKey := txscript.ComputeTaprootOutputKey(wif.PrivKey.PubKey(), tapScriptRootHash[:])
	address, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(outputKey), &chaincfg.TestNet3Params)
	if err != nil {
		log.Println("Cannot creaate taproot address!")
		return nil, err
	}
	return address, nil
}

//	addr := &AddressTaproot{
//		AddressSegWit{
//			hrp:            strings.ToLower(hrp),
//			witnessVersion: 0x01,
//			witnessProgram: witnessProg, // this is x
//		},
//	}
func CreateTaprootAddressWithData(pubkey *secp256k1.PublicKey, data [32]byte) (*btcutil.AddressTaproot, error) {
	tapleaf := txscript.NewBaseTapLeaf(data[:])
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapleaf)
	tapScriptRootHash := tapScriptTree.LeafMerkleProofs[0].RootNode.TapHash()

	// taprootKey = internalKey + (h_tapTweak(internalKey || merkleRoot)*G)
	outputKey := txscript.ComputeTaprootOutputKey(pubkey, tapScriptRootHash[:])

	address, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(outputKey), &chaincfg.TestNet3Params)
	if err != nil {
		log.Println("Cannot create taproot address!")
		return nil, err
	}

	// Generate Proof and store ctrBlock
	ctrBlock := tapScriptTree.LeafMerkleProofs[0].ToControlBlock(pubkey)

	//// using ctrBlock to verify the data in this address
	witnessProgram := address.WitnessProgram() // x in address

	//verify
	err = txscript.VerifyTaprootLeafCommitment(&ctrBlock, witnessProgram, data[:])
	if err != nil {
		log.Println("Verify failse")
	}
	log.Println("Verify success")
	return address, nil
}
