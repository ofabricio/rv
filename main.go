package main

import (
	"io"
	"os"

	"github.com/ofabricio/rv/data"
)

func main() {

	var r io.Reader

	if isPipe() {
		r = os.Stdin
	}

	c := data.NewCarteira()
	if err := c.CommandLine(r, os.Stdout); err != nil {
		panic(err)
	}
}

// Verifica se há dados sendo passados via pipe.
func isPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
