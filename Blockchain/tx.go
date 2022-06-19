package blockchain

type TxOutput struct {
	Value  int
	Pubkey string
}
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}
