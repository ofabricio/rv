package main

import (
	"github.com/ofabricio/rv/data"
)

func main() {

	c := data.NewCarteira()
	if err := c.CommandLine(); err != nil {
		panic(err)
	}
}
