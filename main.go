package main

import (
	"os"

	"github.com/ofabricio/rv/data"
)

func main() {

	var c data.Carteira
	c.Load("db.ndjson")
	c.Print(os.Stdout)
}
