package data

import (
	"fmt"
	"io"

	"github.com/aquasecurity/table"
	"github.com/shopspring/decimal"
)

type PrinterTable struct {
	c *Carteira
}

func (p *PrinterTable) PrintOperacoesComAcoes(w io.Writer) {
	rb := OperacoesRowBuilder{Param: p.c.Param}
	t := table.New(w)
	t.SetRowLines(false)
	t.SetHeaderColSpans(0, 13)
	t.AddHeaders("OPERAÇÕES COM AÇÕES")
	t.AddHeaders(rb.Headers()...)
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
	for o := range p.c.Acoes.Iter() {
		t.AddRow(rb.Build(&o)...)
	}
	t.Render()
}

func (p *PrinterTable) PrintBensDireitos(w io.Writer) {

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
				p.c.Param.FormatDecimal(ticker.SituacaoAnterior),
				p.c.Param.FormatDecimal(ticker.SituacaoCorrente),
				fmt.Sprintf("%s %s", ticker.Grupo, ticker.Codigo),
				ticker.Discriminacao,
			)
		}
		t.Render()
	}
}

func (p *PrinterTable) PrintDividaOnusReais(w io.Writer) {

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
				p.c.Param.FormatDecimal(ticker.SituacaoAnterior),
				p.c.Param.FormatDecimal(ticker.SituacaoCorrente),
				ticker.Codigo,
				ticker.Discriminacao,
			)
		}
		t.Render()
	}
}

func (p *PrinterTable) PrintRendimentosIsentosNaoTributaveis(w io.Writer) {
	for _, ano := range p.c.RendimentosIsentosNaoTributaveis() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 2)
		t.AddHeaders(fmt.Sprint(ano.Ano), "RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS")
		t.AddHeaders("Ticker", "Valor", "Código")
		t.SetAlignment(table.AlignLeft, table.AlignRight, table.AlignLeft)
		for _, r := range ano.Rendimentos {
			t.AddRow(r.Ticker, p.c.Param.FormatDecimal(r.Valor), r.Codigo)
		}
		t.Render()
	}
}

func (p *PrinterTable) PrintRendimentosSujeitosTributacaoExclusiva(w io.Writer) {
	for _, ano := range p.c.RendimentosSujeitosTributacaoExclusiva() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 2)
		t.AddHeaders(fmt.Sprint(ano.Ano), "RENDIMENTOS SUJEITOS À TRIBUTAÇÃO EXCLUSIVA")
		t.AddHeaders("Ticker", "Valor", "Código")
		t.SetAlignment(table.AlignLeft, table.AlignRight, table.AlignLeft)
		for _, r := range ano.Rendimentos {
			t.AddRow(r.Ticker, p.c.Param.FormatDecimal(r.Valor), r.Codigo)
		}
		t.Render()
	}
}

func (p *PrinterTable) PrintOperacoesComunsDayTrade(w io.Writer) {
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
				p.c.Param.FormatDecimal(v.Lucro),
				p.c.Param.FormatDecimal(v.LucroOp),
				p.c.Param.FormatDecimal(v.LucroAc),
				p.c.Param.FormatDecimal(v.IR),
			)
		}
		t.SetFooterAlignment(table.AlignCenter, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight)
		t.AddFooters("Total", p.c.Param.FormatDecimal(r.TotalAcoes), p.c.Param.FormatDecimal(r.TotalOpcao), p.c.Param.FormatDecimal(r.TotalAc), p.c.Param.FormatDecimal(r.TotalIR))
		t.Render()
	}
}
