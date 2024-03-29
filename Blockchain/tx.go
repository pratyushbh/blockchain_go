package blockchain

import (
	"bytes"

	"github.com/pratyushbh/blockchain_go/wallet"
)

type TxOutput struct {
	Value      int
	PubkeyHash []byte
}
type TxInput struct {
	ID     []byte
	Out    int
	Sig    []byte
	Pubkey []byte
}

func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.Pubkey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubkeyHash = pubKeyHash
}
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubkeyHash, pubKeyHash) == 0
}
