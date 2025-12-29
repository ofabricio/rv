package data

import (
	"encoding/csv"
	"io"
)

type CSVPrinter struct {
	c *Carteira
}

func (p *CSVPrinter) PrintOperacoesComAcoes(w io.Writer) {
	rb := OperacoesRowBuilder{Param: p.c.Param}
	csv := csv.NewWriter(w)
	csv.Write(rb.Headers())
	for o := range p.c.Acoes.Iter() {
		csv.Write(rb.Build(&o))
	}
	csv.Flush()
}

func (p *CSVPrinter) PrintBensDireitos(w io.Writer) {}

func (p *CSVPrinter) PrintDividaOnusReais(w io.Writer) {}

func (p *CSVPrinter) PrintRendimentosIsentosNaoTributaveis(w io.Writer) {}

func (p *CSVPrinter) PrintRendimentosSujeitosTributacaoExclusiva(w io.Writer) {}

func (p *CSVPrinter) PrintOperacoesComunsDayTrade(w io.Writer) {}
