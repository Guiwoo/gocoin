package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/guiwoo/gocoin/explorer"
	"github.com/guiwoo/gocoin/rest"
)

func usage() {
	fmt.Printf("Welcome to 귀우 코인\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port:		Set the Port on the server\n")
	fmt.Printf("-mode:		Choose between 'html and rest'\n")
	runtime.Goexit()
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")

	mode := flag.String("mode", "rest", "Choose between html and rest")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}
