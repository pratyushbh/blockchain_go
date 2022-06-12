package blockchain

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}
type BlockchainIterator struct {
	currentHash []byte
	Database    *badger.DB
}

const (
	dbPath = "./tmp/blocks"
)

func InitBlockChain() *Blockchain {
	var lastHash []byte
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		err := os.MkdirAll(dbPath, os.ModeDir)
		Handle(err)
	}
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	db, err := badger.Open(opts)
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Print("No blockExist")
			genesis := Genesis()
			fmt.Print("Genesis Block created")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)
			Handle(err)
			lastHash = genesis.Hash
			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.Value()
			Handle(err)
			return err
		}
	})
	Handle(err)
	blockchain := Blockchain{lastHash, db}
	return &blockchain
}
func (chain *Blockchain) Addblock(data string) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()
		Handle(err)
		return err
	})
	newBlock := CreateBlock(data, lastHash)
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("ln"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
	// prevBlock := chain.Blocks[len(chain.Blocks)-1]
	// new := CreateBlock(data, prevBlock.Hash)
	// chain.Blocks = append(chain.Blocks, new)
}
func (chain *Blockchain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}
	return iter
}
func (chain *BlockchainIterator) Next() *Block {
	var block *Block
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(chain.currentHash)
		Handle(err)
		encodedBlock, err := item.Value()
		block = Deserialized(encodedBlock)
		return err
	})
	Handle(err)
	chain.currentHash = block.PrevHash
	return block
}
