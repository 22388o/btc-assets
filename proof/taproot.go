package proof

import (
	"demo_btc_assets/mssmt"

	"github.com/btcsuite/btcd/btcec/v2"
)

type TaprootProof struct {
	// OutputIndex is the index of the output for which the proof applies.
	OutputIndex uint32

	// InternalKey is the internal key of the taproot output at OutputIndex.
	InternalKey *btcec.PublicKey

	// CommitmentProof represents a commitment proof for an asset, proving
	// inclusion or exclusion of an asset within a Taproot Asset commitment.
	CommitmentProof mssmt.Proof

	// // TapscriptProof represents a taproot control block to prove that a
	// // taproot output is not committing to a Taproot Asset commitment.
	// //
	// // NOTE: This field will be set only if the output does NOT contain a
	// // valid Taproot Asset commitment.
	// TapscriptProof *TapscriptProof
}
