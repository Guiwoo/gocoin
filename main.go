package main

import (
	"github.com/guiwoo/gocoin/cli"
	"github.com/guiwoo/gocoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
