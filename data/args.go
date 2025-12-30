package data

import (
	"flag"
	"fmt"
	"os"
)

func (c *Carteira) CommandLine() {

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("  Usage: rv")
		fmt.Println()
		flag.PrintDefaults()
	}

	file := flag.String("file", "db.ndjson", "arquivo de operações")
	frmt := flag.String("format", "table", "mostra resultado no formato especificado (table, csv)")
	flag.StringVar(&c.Param.FormatoData, "date-format", "2006-01-02", "formato de data (ex. 02/01/2006)")
	flag.BoolVar(&c.Param.MostrarValorExato, "exact-value", false, "mostra valores exatos (ex. 0,14952765 em vez de 0,15)")

	flag.Parse()

	c.Load(*file)
	c.Print(*frmt, os.Stdout)
}
