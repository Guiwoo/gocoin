package p2p

import (
	"encoding/json"

	"github.com/guiwoo/gocoin/blockchain"
	"github.com/guiwoo/gocoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota + 200
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func (m *Message) addPayload(p interface{}) {
	b, err := json.Marshal(p)
	utils.HandleErr(err)
	m.Payload = b
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind: kind,
	}
	m.addPayload(payload)
	mJson, err := json.Marshal(m)
	utils.HandleErr(err)
	return mJson
}

func sendNewestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}
