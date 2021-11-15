package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/guiwoo/gocoin/db"
	"github.com/guiwoo/gocoin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"defficulty"`
	Nonce        int    `json:"nonce"`
	TimeStamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.TimeStamp = int(time.Now().Unix())
		hash := utils.Hashing(b)
		fmt.Printf("\n\nTarget: %s\nHash: %s\nNonce: %d\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height int) *Block {
	block := Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty(BlockChain()),
		Nonce:      0,
	}
	block.mine()
	block.Transactions = Mempool.txToConfirm()
	block.persist()
	return &block
}

var ErrNotFound = errors.New("Block not Found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}
