package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guiwoo/gocoin/blockchain"
	"github.com/guiwoo/gocoin/p2p"
	"github.com/guiwoo/gocoin/utils"
	"github.com/guiwoo/gocoin/wallet"
)

type uRLaddress string

var port string

func (u uRLaddress) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type uRL struct {
	URL         uRLaddress `json:"url"`
	Method      string     `json:"method"`
	Description string     `json:"description"`
	Payload     string     `json:"payload,omitempty"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type errResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addPeerPayload struct {
	Address string
	Port    string
}

func (u uRL) String() string {
	return "Hello i am url description"
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []uRL{
		{
			URL:         uRLaddress("/"),
			Method:      "GET",
			Description: "See Documentation",
		}, {
			URL:         uRLaddress("/status"),
			Method:      "GET",
			Description: "See the status of the block",
		},
		{
			URL:         uRLaddress("/blocks"),
			Method:      "GET",
			Description: "See all blocks",
		},
		{
			URL:         uRLaddress("/blocks"),
			Method:      "Post",
			Description: "Add a block",
			Payload:     "data:string",
		},
		{
			URL:         uRLaddress("/blocks/{hash}"),
			Method:      "Post",
			Description: "See a block",
		},
		{
			URL:         uRLaddress("/balance/{address}"),
			Method:      "GET",
			Description: "GET Txouts for an Address",
		},
		{
			URL:         uRLaddress("/ws}"),
			Method:      "GET",
			Description: "Upgrade to Web Sockets",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
	case "POST":
		blockchain.BlockChain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err != nil {
		encoder.Encode(errResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//Do something
		rw.Header().Add("Content-Type", "application/json")
		log.Println(r.RequestURI)
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BlockChain())
}
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.BlockChain())
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.BlockChain())))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errResponse{err.Error()})
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(struct{ Address string }{Address: address})
}

func peersHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port)
		rw.WriteHeader(http.StatusOK)

	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peersHandler).Methods("GET", "POST")
	fmt.Printf("✅ http://localhost%s Connected\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
