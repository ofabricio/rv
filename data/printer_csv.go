package data

import (
	"encoding/csv"
	"io"
)

type PrinterCSV struct {
	c *Carteira
}

func (p *PrinterCSV) PrintOperacoesComAcoes(w io.Writer) {
	rb := OperacoesRowBuilder{Param: p.c.Param}
	csv := csv.NewWriter(w)
	csv.Write(rb.Headers())
	for o := range p.c.Acoes.Iter() {
		csv.Write(rb.Build(&o))
	}
	csv.Flush()
}

func (p *PrinterCSV) PrintBensDireitos(w io.Writer) {}

func (p *PrinterCSV) PrintDividaOnusReais(w io.Writer) {}

func (p *PrinterCSV) PrintRendimentosIsentosNaoTributaveis(w io.Writer) {}

func (p *PrinterCSV) PrintRendimentosSujeitosTributacaoExclusiva(w io.Writer) {}

func (p *PrinterCSV) PrintOperacoesComunsDayTrade(w io.Writer) {}
