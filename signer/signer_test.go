package signer

import (
	"encoding/hex"
	"testing"

	mint "github.com/void616/gm.mint"
)

func TestSignerGeneration(t *testing.T) {

	gen1, _ := New()
	gen2, _ := New()

	p1 := gen1.PrivateKey()
	p2 := gen2.PrivateKey()

	if hex.EncodeToString(p1[:]) == hex.EncodeToString(p2[:]) {
		t.Fatal(hex.EncodeToString(p1[:]), "==", hex.EncodeToString(p2[:]))
	}
}

func TestSignerFromSumusPrivateKey(t *testing.T) {

	spvt, _ := mint.Unpack58("TBzyWv8Dga5aN4Hai2nFTwyTXvDJKkJhq8HMDPC9zqTWLSTLo4jFFKKnVS52a1kp7YJdm2b8HrR2Buk9PqyD1DwhxUzsJ")
	spub, _ := mint.Unpack58("2p6QCcwAMLSSXfFFVQT4vYCe8VPwm3rvK4zdNGAM7zeLBqrVLW")

	sig, _ := FromBytes(spvt)

	x := sig.PublicKey()
	if hex.EncodeToString(x[:]) != hex.EncodeToString(spub) {
		t.Fatal(hex.EncodeToString(x[:]), "!=", hex.EncodeToString(spub))
	}
}

func TestSignerFromSumusPrivateKey2(t *testing.T) {

	spvt, _ := mint.Unpack58("4CdzVBba43H7B12zNoSCE8dz8RM9ggUSagfxPdZ1kQ7hbrXLqNNUwGQiiV1VxU3xuEcj4ybxTZPnjq8BAhBUuJxzU8XxQ")
	spub, _ := mint.Unpack58("2PztA94iHZdeX8d5hPJbQfUGcN6WWUhfmU6G5ySJQ9cnUueiuk")

	sig, _ := FromBytes(spvt)

	x := sig.PublicKey()
	if hex.EncodeToString(x[:]) != hex.EncodeToString(spub) {
		t.Fatal(hex.EncodeToString(x[:]), "!=", hex.EncodeToString(spub))
	}
}
