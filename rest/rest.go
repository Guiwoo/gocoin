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

type uRL struct {
	URL         uRLaddress `json:"url"`
	Method      string     `json:"method"`
	Description string     `json:"description"`
	Payload     string     `json:"payload,omitempty"`
}

type addBlockMessage struct {
	Message string
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
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.BlockChain().Blocks())
	case "POST":
		var addBlockMessage addBlockMessage
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockMessage))
		blockchain.BlockChain().AddBlock(addBlockMessage.Message)
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

func Start(aPort int) {
	rotuer := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	rotuer.Use(jsonContentTypeMiddleware)
	rotuer.HandleFunc("/", documentation).Methods("GET")
	rotuer.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	rotuer.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("✅ http://localhost%s Connected\n", port)
	log.Fatal(http.ListenAndServe(port, rotuer))
}
