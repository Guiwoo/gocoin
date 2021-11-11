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

var port string

type uRL string
type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}
type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

func (u uRL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type uRLDescription struct {
	URL         uRL    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []uRLDescription{
		{
			URL:         uRL("/"),
			Method:      "GET",
			Description: "See Documentation",
		}, {
			URL:         uRL("/status"),
			Method:      "GET",
			Description: "See the Status of the BlockChain",
		},
		{
			URL:         uRL("/blocks"),
			Method:      "Post",
			Description: "Add A Block",
			Payload:     "data: string",
		}, {
			URL:         uRL("/blocks{hash}"),
			Method:      "GET",
			Description: "See A Block",
		},
		{
			URL:         uRL("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
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

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BlockChain())
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
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
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.BlockChain().TxOutByaddress(address)))
	}
}

func jsonContentTypeMiddlewar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddlewar)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	fmt.Printf("✅Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
