package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/guiwoo/gocoin/explorer"
	"github.com/guiwoo/gocoin/rest"
)

func usage() {
	fmt.Printf("Welcome to 고코인\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port=4000:		Set the port of the server\n")
	fmt.Printf("-mode=rest:		✅Choose between 'html' and `rest\n")
	os.Exit(0)
}

func Start() {
	port := flag.Int("port", 4000, "✅Set port of the server")
	extraPort := flag.Int("extraPort", *port+1, "⭐️Running boss of them")
	mode := flag.String("mode", "run", "✅Choose between 'html' and `rest` or `run` for both")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	case "run":
		go rest.Start(*port)
		explorer.Start(*extraPort)
	default:
		usage()
	}
}
