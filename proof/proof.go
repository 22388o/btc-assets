package proof

import (
	"github.com/btcsuite/btcd/wire"
)

// Proof encodes all of the data necessary to prove a valid state transition for
// an asset has occurred within an on-chain transaction.
type Proof struct {
	// PrevOut is the previous on-chain outpoint of the asset.
	PrevOut wire.OutPoint

	// BlockHeader is the current block header committing to the on-chain
	// transaction attempting an asset state transition.
	BlockHeader wire.BlockHeader

	// BlockHeight is the height of the current block committing to the
	// on-chain transaction attempting an asset state transition.
	BlockHeight uint32

	// AnchorTx is the on-chain transaction attempting the asset state
	// transition.
	AnchorTx wire.MsgTx

	// TxMerkleProof is the merkle proof for AnchorTx used to prove its
	// inclusion within BlockHeader.
	//
	// TODO(roasbeef): also store height+index information?
	// TxMerkleProof TxMerkleProof

	// // Asset is the resulting asset after its state transition.
	// Asset Asset

	// InclusionProof is the TaprootProof proving the new inclusion of the
	// resulting asset within AnchorTx.
	InclusionProof TaprootProof

	// ExclusionProofs is the set of TaprootProofs proving the exclusion of
	// the resulting asset from all other Taproot outputs within AnchorTx.
	ExclusionProofs []TaprootProof

	// SplitRootProof is an optional TaprootProof needed if this asset is
	// the result of a split. SplitRootProof proves inclusion of the root
	// asset of the split.
	SplitRootProof *TaprootProof

	// // MetaReveal is the set of bytes that were revealed to prove the
	// // derivation of the meta data hash contained in the genesis asset.
	// //
	// // TODO(roasbeef): use even/odd framing here?
	// //
	// // NOTE: This field is optional, and can only be specified if the asset
	// // above is a genesis asset. If specified, then verifiers _should_ also
	// // verify the hashes match up.
	// MetaReveal *MetaReveal

	// // AdditionalInputs is a nested full proof for any additional inputs
	// // found within the resulting asset.
	// AdditionalInputs []File

	// // ChallengeWitness is an optional virtual transaction witness that
	// // serves as an ownership proof for the asset. If this is non-nil, then
	// // it is a valid transfer witness for a 1-input, 1-output virtual
	// // transaction that spends the asset in this proof and sends it to the
	// // NUMS key, to prove that the creator of the proof is able to produce
	// // a valid signature to spend the asset.
	// ChallengeWitness wire.TxWitness
}
