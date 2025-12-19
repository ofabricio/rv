package main

import (
	"os"

	"github.com/ofabricio/rv/model"
)

func main() {

	var s model.State

	s.Load("db.ndjson")
	s.Print(os.Stdout)
}
