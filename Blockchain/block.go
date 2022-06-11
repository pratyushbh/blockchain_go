package blockchain

type Blockchain struct {
	Blocks []*Block
}
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
func (chain *Blockchain) Addblock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}
func Genesis() *Block {
	return CreateBlock("Genesis Block", []byte{})
}
func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}
