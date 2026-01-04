package data

import (
	"fmt"
	"io"

	"github.com/aquasecurity/table"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type PrinterTable struct {
	c *Carteira
}

func (p *PrinterTable) PrintOperacoesComAcoes(w io.Writer) {
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
	i := 0
	for o := range p.c.Acoes.Iter() {
		i++
		t.AddRow(
			// ID
			fmt.Sprint(i),
			// Ticker
			lo.Ternary(o.Serie != "", fmt.Sprintf("%s %s", o.Ticker, o.Serie), o.Ticker),
			// Data
			lo.Ternary(!o.Vencimento.IsZero(), fmt.Sprintf("%s V %s", p.c.Param.FormatDate(o.Data), p.c.Param.FormatDate(o.Vencimento)), p.c.Param.FormatDate(o.Data)),
			// Operação
			lo.Ternary(o.Fator.IsPositive(), fmt.Sprintf("%s (%s)", o.Tipo, p.c.Param.FormatDecimal(o.Fator)), string(o.Tipo)),
			// Qtd
			lo.Ternary(o.QtdFracao().IsPositive(), fmt.Sprintf("%s (%s)", o.QtdInt(), p.c.Param.FormatDecimal(o.QtdFracao())), o.QtdInt().String()),
			// V. Unit.
			lo.Ternary(o.ValorExercicio.IsPositive(), fmt.Sprintf("(E %s) %s", p.c.Param.FormatDecimal(o.ValorExercicio), p.c.Param.FormatDecimal(o.ValorUnitario)), p.c.Param.FormatDecimal(o.ValorUnitario)),
			// V. Total
			lo.Ternary(o.ValorExercicio.IsPositive(), fmt.Sprintf("(E %s) %s", p.c.Param.FormatDecimal(o.ValorExercicio.Mul(o.Qtd)), p.c.Param.FormatDecimal(o.ValorTotal)), p.c.Param.FormatDecimal(o.ValorTotal)),
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
			lo.Ternary(o.Premio.IsPositive(), fmt.Sprintf("(P %s) %s", p.c.Param.FormatDecimal(o.Premio), p.c.Param.FormatDecimal(o.Lucro)), p.c.Param.FormatDecimal(o.Lucro)),
		)
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
		t.SetFooterAlignment(table.AlignLeft, table.AlignRight, table.AlignLeft)
		for _, r := range ano.Totais {
			t.AddFooters(r.Ticker, p.c.Param.FormatDecimal(r.Valor), r.Codigo)
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
		t.AddHeaders("Mês", "Ações", "Opções", "Acumulado", fmt.Sprintf("IR (%s%%)", r.IR.Mul(decimal.NewFromInt(100))))
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
	for _, r := range p.c.OperacoesComunsDayTradeDayTrade() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 4)
		t.AddHeaders(fmt.Sprintf("%s %d", r.Ticker, r.Ano), "OPERAÇÕES COMUNS/DAY-TRADE [DAY TRADE]")
		t.AddHeaders("Mês", "Ações", "Opções", "Acumulado", fmt.Sprintf("IR (%s%%)", r.IR.Mul(decimal.NewFromInt(100))))
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
