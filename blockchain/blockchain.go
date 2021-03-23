package blockchain

type Blockchain struct { // blockchain struct
	blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) { // 블록체인에 블록 추가
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)

}

func NewBlockchain() *Blockchain { // 새로운 블록체인 생성
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) GetBlocks() []*Block { // 블록체인에 블록들 가져오기
	return bc.blocks
}
