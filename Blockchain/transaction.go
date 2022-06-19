package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}
type TxOutput struct {
	Value  int
	Pubkey string
}
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func CoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}
	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetId()
	return &tx
}
func (tx *Transaction) SetId() {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}
func NewTransaction(from, to string, amount int, chain *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	acc, validOutputs := chain.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("Error: not enough funds")
	}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetId()
	return &tx
}
func (tx *Transaction) isCoinBase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.Pubkey == data
}
