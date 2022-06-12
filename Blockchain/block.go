package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Data     []byte
	Hash     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte(data), []byte{}, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}
func Genesis() *Block {
	return CreateBlock("Genesis Block", []byte{})
}
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(encoder)
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
