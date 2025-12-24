package data

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aquasecurity/table"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type TablePrinter struct {
	c *Carteira
}

func (p *TablePrinter) PrintOperacoesComAcoes(w io.Writer) {
	t := table.New(w)
	t.SetRowLines(false)
	t.SetHeaderColSpans(0, 13)
	t.AddHeaders("OPERAÇÕES COM AÇÕES")
	t.AddHeaders(
		"ID",
		"Ticker",
		"Data",
		"Operação",
		"Qtd",
		"V. Unit.",
		"V. Total",
		"Taxas",
		"Qtd Ac.",
		"V. Total Ac.",
		"PM",
		"V. Compra",
		"Lucro",
	)
	t.SetAlignment(
		table.AlignLeft,  // ID
		table.AlignLeft,  // Data
		table.AlignLeft,  // Ticker
		table.AlignLeft,  // Operação
		table.AlignRight, // Qtd
		table.AlignRight, // V. Unit.
		table.AlignRight, // V. Total
		table.AlignRight, // Taxas
		table.AlignRight, // Qtd Ac.
		table.AlignRight, // V. Total Ac.
		table.AlignRight, // PM
		table.AlignRight, // V. Compra
		table.AlignRight, // Lucro
	)
	for i, o := range p.c.Acoes.IterI() {
		t.AddRow(
			fmt.Sprint(i+1),
			lo.Ternary(o.Serie != "", fmt.Sprintf("%s %s", o.Ticker, o.Serie), o.Ticker),
			lo.Ternary(!o.Vencimento.IsZero(), fmt.Sprintf("%s V %s", o.Data.Format(time.DateOnly), o.Vencimento.Format(time.DateOnly)), o.Data.Format(time.DateOnly)),
			string(o.Tipo),
			lo.Ternary(o.Fracao.IsPositive(), fmt.Sprintf("%s (%s)", o.Qtd, p.formatDecimal(o.Fracao)), o.Qtd.String()),
			lo.Ternary(o.ValorExercicio.IsPositive(), fmt.Sprintf("(E %s) %s", p.formatDecimal(o.ValorExercicio), p.formatDecimal(o.ValorUnitario)), p.formatDecimal(o.ValorUnitario)),
			lo.Ternary(o.ValorExercicio.IsPositive(), fmt.Sprintf("(E %s) %s", p.formatDecimal(o.ValorExercicio.Mul(o.Qtd)), p.formatDecimal(o.ValorTotal)), p.formatDecimal(o.ValorTotal)),
			p.formatDecimal(o.Taxas),
			o.Agg.Qtd.String(),
			p.formatDecimal(o.Agg.ValorTotal),
			p.formatDecimal(o.Agg.PrecoMedio),
			p.formatDecimal(o.ValorCompra),
			lo.Ternary(o.Premio.IsPositive(), fmt.Sprintf("(P %s) %s", p.formatDecimal(o.Premio), p.formatDecimal(o.Lucro)), p.formatDecimal(o.Lucro)),
		)
	}
	t.Render()
}

func (p *TablePrinter) PrintBensDireitos(w io.Writer) {

	for _, bens := range p.c.BensDireitos() {
		t := table.New(w)
		t.SetColumnMaxWidth(100)
		t.SetHeaderColSpans(0, 5)
		t.AddHeaders("BENS E DIREITOS")
		t.AddHeaders(
			"Ticker",
			fmt.Sprintf("Situação em %d", bens.AnoAnterior),
			fmt.Sprintf("Situação em %d", bens.AnoCorrente),
			"Grp Cód",
			"Discriminação",
		)
		t.SetAlignment(table.AlignLeft, table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignLeft)
		for _, ticker := range bens.Tickers {
			t.AddRow(
				ticker.Ticker,
				p.formatDecimal(ticker.SituacaoAnterior),
				p.formatDecimal(ticker.SituacaoCorrente),
				fmt.Sprintf("%s %s", ticker.Grupo, ticker.Codigo),
				ticker.Discriminacao,
			)
		}
		t.Render()
	}
}

func (p *TablePrinter) PrintDividaOnusReais(w io.Writer) {

	for _, bens := range p.c.DividaOnusReais() {
		t := table.New(w)
		t.SetColumnMaxWidth(100)
		t.SetHeaderColSpans(0, 5)
		t.AddHeaders("DÍVIDA E ÔNUS REAIS")
		t.AddHeaders(
			"Ticker",
			fmt.Sprintf("Situação em %d", bens.AnoAnterior),
			fmt.Sprintf("Situação em %d", bens.AnoCorrente),
			"Código",
			"Discriminação",
		)
		t.SetAlignment(table.AlignLeft, table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignLeft)
		for _, ticker := range bens.Tickers {
			t.AddRow(
				ticker.Ticker,
				p.formatDecimal(ticker.SituacaoAnterior),
				p.formatDecimal(ticker.SituacaoCorrente),
				ticker.Codigo,
				ticker.Discriminacao,
			)
		}
		t.Render()
	}
}

func (p *TablePrinter) PrintRendimentosIsentosNaoTributaveis(w io.Writer) {
	for _, ano := range p.c.RendimentosIsentosNaoTributaveis() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 2)
		t.AddHeaders(fmt.Sprint(ano.Ano), "RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS")
		t.AddHeaders("Ticker", "Valor", "Código")
		t.SetAlignment(table.AlignLeft, table.AlignRight, table.AlignLeft)
		for _, r := range ano.Rendimentos {
			t.AddRow(r.Ticker, p.formatDecimal(r.Valor), r.Codigo)
		}
		t.Render()
	}
}

func (p *TablePrinter) PrintRendimentosSujeitosTributacaoExclusiva(w io.Writer) {
	for _, ano := range p.c.RendimentosSujeitosTributacaoExclusiva() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 2)
		t.AddHeaders(fmt.Sprint(ano.Ano), "RENDIMENTOS SUJEITOS À TRIBUTAÇÃO EXCLUSIVA")
		t.AddHeaders("Ticker", "Valor", "Código")
		t.SetAlignment(table.AlignLeft, table.AlignRight, table.AlignLeft)
		for _, r := range ano.Rendimentos {
			t.AddRow(r.Ticker, p.formatDecimal(r.Valor), r.Codigo)
		}
		t.Render()
	}
}

func (p *TablePrinter) PrintOperacoesComunsDayTrade(w io.Writer) {
	for _, r := range p.c.OperacoesComunsDayTrade() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 4)
		t.AddHeaders(fmt.Sprintf("%s %d", r.Ticker, r.Ano), "OPERAÇÕES COMUNS/DAY-TRADE")
		t.AddHeaders("Mês", "Ações", "Opções", "Acumulado", fmt.Sprintf("IR (%s%%)", r.SwingTradeIR.Mul(decimal.NewFromInt(100))))
		t.SetAlignment(table.AlignCenter, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight)
		for _, v := range r.Meses {
			t.AddRow(
				translateMonth[v.Mes],
				p.formatDecimal(v.Lucro),
				p.formatDecimal(v.LucroOp),
				p.formatDecimal(v.LucroAc),
				p.formatDecimal(v.IR),
			)
		}
		t.SetFooterAlignment(table.AlignCenter, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight)
		t.AddFooters("Total", p.formatDecimal(r.TotalAcoes), p.formatDecimal(r.TotalOpcao), p.formatDecimal(r.TotalAc), p.formatDecimal(r.TotalIR))
		t.Render()
	}
}

func (p *TablePrinter) formatDecimal(d decimal.Decimal) string {
	return strings.Replace(d.StringFixed(2), ".", ",", 1)
}

var translateMonth = map[int]string{
	1:  "JAN",
	2:  "FEV",
	3:  "MAR",
	4:  "ABR",
	5:  "MAI",
	6:  "JUN",
	7:  "JUL",
	8:  "AGO",
	9:  "SET",
	10: "OUT",
	11: "NOV",
	12: "DEZ",
}
