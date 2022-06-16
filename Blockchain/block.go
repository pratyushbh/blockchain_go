package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	//Data string
	Transactions []*Transaction
	Hash         []byte
	PrevHash     []byte
	Nonce        int
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{txs, []byte{}, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	Handle(err)
	return res.Bytes()
}
func Deserialized(data []byte) *Block {
	var Block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&Block)
	Handle(err)
	return &Block
}
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
