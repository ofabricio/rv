package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/samber/lo"
)

type PrinterCSV struct {
	c *Carteira
}

func (p *PrinterCSV) PrintOperacoesComAcoes(w io.Writer) {
	csv := csv.NewWriter(w)
	csv.Write([]string{
		"ID",
		"Ticker",
		"Série",
		"Data",
		"Vencimento",
		"Operação",
		"Qtd",
		"Fração",
		"V. Unit.",
		"V. Total",
		"Taxas",
		"Qtd Ac.",
		"V. Total Ac.",
		"PM",
		"V. Compra",
		"Lucro",
		"Prêmio",
		"V. Exercício",
	})
	i := 0
	for o := range p.c.Acoes.Iter() {
		i++
		csv.Write([]string{
			// ID
			fmt.Sprint(i),
			// Ticker
			o.Ticker,
			// Série
			o.Serie,
			// Data
			o.Data.Format(time.DateOnly),
			// Vencimento
			lo.Ternary(o.Vencimento.IsZero(), "", o.Vencimento.Format(time.DateOnly)),
			// Operação
			string(o.Tipo),
			// Qtd
			o.Qtd.String(),
			// Fração
			p.c.Param.FormatDecimal(o.Fracao),
			// V. Unit.
			p.c.Param.FormatDecimal(o.ValorUnitario),
			// V. Total
			p.c.Param.FormatDecimal(o.ValorTotal),
			// Taxas
			p.c.Param.FormatDecimal(o.Taxas),
			// Qtd Ac.
			o.Agg.Qtd.String(),
			// V. Total Ac.
			p.c.Param.FormatDecimal(o.Agg.ValorTotal),
			// PM
			p.c.Param.FormatDecimal(o.Agg.PrecoMedio),
			// V. Compra
			p.c.Param.FormatDecimal(o.ValorCompra),
			// Lucro
			p.c.Param.FormatDecimal(o.Lucro),
			// Premio
			p.c.Param.FormatDecimal(o.Premio),
			// ValorExercicio
			p.c.Param.FormatDecimal(o.ValorExercicio),
		})
	}
	csv.Flush()
}

func (p *PrinterCSV) PrintBensDireitos(w io.Writer) {}

func (p *PrinterCSV) PrintDividaOnusReais(w io.Writer) {}

func (p *PrinterCSV) PrintRendimentosIsentosNaoTributaveis(w io.Writer) {}

func (p *PrinterCSV) PrintRendimentosSujeitosTributacaoExclusiva(w io.Writer) {}

func (p *PrinterCSV) PrintOperacoesComunsDayTrade(w io.Writer) {}
