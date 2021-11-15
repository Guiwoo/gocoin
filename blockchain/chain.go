package blockchain

import (
	"sync"

	"github.com/guiwoo/gocoin/db"
	"github.com/guiwoo/gocoin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `josn:"currentDifficulty"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func persistBlockchain(b *blockchain) {
	db.SaveBlockChain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func Blocks(b *blockchain) []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalculateBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.TimeStamp / 60) - (lastRecalculateBlock.TimeStamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime < (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime > (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func difficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		//recalculate the difficulty
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func UTxOutByaddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)

	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutByaddress(address, b)
	var amout int
	for _, txOut := range txOuts {
		amout += txOut.Amount
	}
	return amout
}

func BlockChain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkPoint := db.CheckPoint()
		if checkPoint == nil {
			b.AddBlock()
		} else {
			b.restore(checkPoint)
		}
	})
	return b
}
