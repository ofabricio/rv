package data

import (
	"flag"
	"fmt"
	"io"
	"time"
)

func (c *Carteira) CommandLine(r io.Reader, w io.Writer) error {

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("  Usage: rv [flags]")
		fmt.Println()
		flag.PrintDefaults()
	}

	file := flag.String("file", "db.ndjson", "arquivo de operações")
	frmt := flag.String("format", "table", "mostra resultado no formato especificado (table, csv)")
	nota := flag.Bool("notas", false, "importa notas de corretagem de arquivos pdf que estiverem no diretório corrente")
	flag.StringVar(&c.Param.FormatoData, "date-format", "2006-01-02", "formato de data (ex. 02/01/2006)")
	flag.BoolVar(&c.Param.MostrarValorExato, "exact-value", false, "mostra valores exatos (ex. 0,14952765 em vez de 0,15) (default false)")
	flag.IntVar(&c.Param.FiltrarAno, "filter-year", 0, "filtra operações por ano (YYYY para mostrar o ano em questão; 0 para mostrar todos os anos; -1 para mostrar apenas o ano atual) (default 0)")
	flag.StringVar(&c.Param.FiltrarTicker, "filter-ticker", "", "filtra operações por ticker ")

	flag.Parse()

	if c.Param.FiltrarAno < 0 {
		c.Param.FiltrarAno = time.Now().Local().Year() + (c.Param.FiltrarAno + 1)
	}

	if *nota {
		return ProcessarNotasPDF(".", w)
	}

	if r != nil {
		if err := c.Read(r); err != nil {
			return err
		}
	} else {
		if err := c.Load(*file); err != nil {
			return err
		}
	}

	c.Print(*frmt, w)
	return nil
}
