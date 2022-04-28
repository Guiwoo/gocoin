package p2p

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
