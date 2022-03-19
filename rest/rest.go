package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guiwoo/gocoin/blockchain"
	"github.com/guiwoo/gocoin/utils"
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
		}, {
			URL:         uRLaddress("/balance/{address}"),
			Method:      "GET",
			Description: "GET Txouts for an Address",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.BlockChain().Blocks())
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

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BlockChain())
}
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BlockChain().BalanceByAddress(address)
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.BlockChain().TxOutsByAddress(address)))
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
		json.NewEncoder(rw).Encode(errResponse{"not enough funds"})
	}
	rw.WriteHeader(http.StatusCreated)
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance)
	router.HandleFunc("/mempool", mempool)
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("✅ http://localhost%s Connected\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
