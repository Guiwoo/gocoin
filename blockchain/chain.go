package blockchain

import (
	"fmt"
	"sync"

	"github.com/guiwoo/gocoin/db"
	"github.com/guiwoo/gocoin/utils"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentdifficulty"`
}

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

var b *blockchain
var once sync.Once

func (b *blockchain) recaulculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculateBlock := allBlocks[difficultyInterval-1]
	fmt.Println("✅", allBlocks)
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculateBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime < expectedTime-allowedRange {
		return b.CurrentDifficulty + 1
	} else if actualTime > expectedTime+allowedRange {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.recaulculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) restore(data []byte) {
	utils.RestoreFromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := CreateBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
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

func (b *blockchain) txOut() []*TxOut {
	var txOuts []*TxOut
	blocks := b.Blocks()
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
}

func (b *blockchain) TxOutsByAddress(address string) []*TxOut {
	var ownTxOuts []*TxOut
	txOuts := b.txOut()
	for _, txOut := range txOuts {
		if txOut.Owner == address {
			ownTxOuts = append(ownTxOuts, txOut)
		}
	}
	return ownTxOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	txOuts := b.TxOutsByAddress(address)
	var amount int
	for _, v := range txOuts {
		amount += v.Amount
	}
	return amount
}

func BlockChain() *blockchain {
	if b == nil {
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
	}
	return b
}
