package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"runtime"

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
	genesisData = "First Transaction from Genesis"
	dbFile      = "./tmp/blocks/MANIFEST"
	dbPath      = "./tmp/blocks"
)

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func InitBlockChain(address string) *Blockchain {
	var lastHash []byte
	if DBExists() {
		fmt.Println("BlockChain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	db, err := badger.Open(opts)
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error {
		// if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
		// 	fmt.Print("No blockExist")
		// 	genesis := Genesis()
		// 	fmt.Print("Genesis Block created")
		// 	err = txn.Set(genesis.Hash, genesis.Serialize())
		// 	Handle(err)
		// 	err = txn.Set([]byte("lh"), genesis.Hash)
		// 	Handle(err)
		// 	lastHash = genesis.Hash
		// 	return err
		// } else {
		// 	item, err := txn.Get([]byte("lh"))
		// 	Handle(err)
		// 	lastHash, err = item.Value()
		// 	Handle(err)
		// 	return err
		// }
		cbtx := CoinBaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Print("Genesis Block created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		Handle(err)
		lastHash = genesis.Hash
		return err
	})
	Handle(err)
	blockchain := Blockchain{lastHash, db}
	return &blockchain
}
func ContinueBlockChain(address string) *Blockchain {
	if DBExists() == false {
		fmt.Println("No existing blockchain fount, create one!")
		runtime.Goexit()
	}
	var lastHash []byte
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	db, err := badger.Open(opts)
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()
		Handle(err)
		return err
	})
	Handle(err)
	blockchain := Blockchain{lastHash, db}
	return &blockchain
}
func (chain *Blockchain) Addblock(transactions []*Transaction) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()
		Handle(err)
		return err
	})
	Handle(err)
	newBlock := CreateBlock(transactions, lastHash)
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
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
func (chain *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTxs []Transaction
	spentTXOs := make(map[string][]int)
	iter := chain.Iterator()
	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txId] != nil {
					for _, spentOut := range spentTXOs[txId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if tx.isCoinBase() == false {
				for _, in := range tx.Inputs {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

func (chain *Blockchain) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXO []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(pubKeyHash)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXO = append(UTXO, out)
			}
		}
	}
	return UTXO
}

func (chain *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOuts
}
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()
	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.PrevHash) == 0 {
			break

		}
	}
	return Transaction{}, errors.New("Transaction does not exist")
}
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey, prevTXs)
}
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs)
}
