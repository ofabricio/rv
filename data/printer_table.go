package data

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"

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
		"Ticker",
		"Data",
		"Operação",
		"Qtd",
		"V. Unit.",
		"V. Total",
		"Taxas",
		"IRRF",
		"Qtd Ac.",
		"V. Total Ac.",
		"PM",
		"V. Compra",
		"Lucro",
	)
	t.SetAlignment(
		table.AlignLeft,  // Data
		table.AlignLeft,  // Ticker
		table.AlignLeft,  // Operação
		table.AlignRight, // Qtd
		table.AlignRight, // V. Unit.
		table.AlignRight, // V. Total
		table.AlignRight, // Taxas
		table.AlignRight, // IRRF
		table.AlignRight, // Qtd Ac.
		table.AlignRight, // V. Total Ac.
		table.AlignRight, // PM
		table.AlignRight, // V. Compra
		table.AlignRight, // Lucro
	)
	for o := range p.c.Acoes.Iter() {
		if !p.c.Param.FilterOperacao(o) {
			continue
		}
		t.AddRow(
			// Ticker
			lo.Ternary(o.Serie != "", fmt.Sprintf("%s %s", o.Ticker, o.Serie), o.Ticker),
			// Data
			lo.Ternary(!o.Vencimento.IsZero(), fmt.Sprintf("%s V %s", p.c.Param.FormatDate(o.Data), p.c.Param.FormatDate(o.Vencimento)), p.c.Param.FormatDate(o.Data)),
			// Operação
			lo.Ternary(o.Fator.IsPositive(), fmt.Sprintf("%s (%s)", o.Tipo, p.c.Param.FormatDecimal(o.Fator)), string(o.Tipo)),
			// Qtd
			lo.Ternary(o.Qtd.IsZero() && o.Tipo != VENDA, "-",
				lo.Ternary(o.QtdFracao().IsPositive(), fmt.Sprintf("%s (%s)", o.QtdInt(), p.c.Param.FormatDecimal(o.QtdFracao())), o.QtdInt().String()),
			),
			// V. Unit.
			lo.Ternary(o.ValorUnitario.IsZero(), "-",
				lo.Ternary(o.ValorExercicio.IsPositive(), fmt.Sprintf("(E %s) %s", p.c.Param.FormatDecimal(o.ValorExercicio), p.c.Param.FormatDecimal(o.ValorUnitario)), p.c.Param.FormatDecimal(o.ValorUnitario)),
			),
			// V. Total
			lo.Ternary(o.ValorTotal.IsZero(), "-",
				lo.Ternary(o.ValorExercicio.IsPositive(), fmt.Sprintf("(E %s) %s", p.c.Param.FormatDecimal(o.ValorExercicio.Mul(o.Qtd)), p.c.Param.FormatDecimal(o.ValorTotal)), p.c.Param.FormatDecimal(o.ValorTotal)),
			),
			// Taxas
			lo.Ternary(o.Taxas.IsZero(), "-", p.c.Param.FormatDecimal(o.Taxas)),
			// IRRF
			lo.Ternary(o.IRRF.IsPositive(), p.c.Param.FormatDecimal(o.IRRF.Truncate(2)), "-"),
			// Qtd Ac.
			o.Agg.Qtd.String(),
			// V. Total Ac.
			p.c.Param.FormatDecimal(o.Agg.ValorTotal),
			// PM
			p.c.Param.FormatDecimal(o.Agg.PrecoMedio),
			// V. Compra
			lo.Ternary(o.ValorCompra.IsZero(), "-", p.c.Param.FormatDecimal(o.ValorCompra)),
			// Lucro
			lo.Ternary(o.Premio.Add(o.Lucro).IsZero(), "-",
				lo.Ternary(o.Premio.IsPositive(), fmt.Sprintf("(P %s) %s", p.c.Param.FormatDecimal(o.Premio), p.c.Param.FormatDecimal(o.Lucro)), p.c.Param.FormatDecimal(o.Lucro)),
			),
		)
	}
	t.Render()
}

func (p *PrinterTable) PrintBensDireitos(w io.Writer) {

	for _, bens := range p.c.BensDireitos() {
		if !p.c.Param.FilterYear(bens.AnoCorrente) {
			continue
		}
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
			if !p.c.Param.FilterTicker(ticker.Ticker) {
				continue
			}
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
		if !p.c.Param.FilterYear(bens.AnoCorrente) {
			continue
		}
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
			if !p.c.Param.FilterTicker(ticker.Ticker) {
				continue
			}
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
		if !p.c.Param.FilterYear(ano.Ano) {
			continue
		}
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 3)
		t.AddHeaders(fmt.Sprint(ano.Ano), "RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS")
		t.AddHeaders("Ticker", "Valor", "Código", "Descrição")
		t.SetAlignment(table.AlignLeft, table.AlignRight, table.AlignCenter, table.AlignLeft)
		for _, r := range ano.Rendimentos {
			if !p.c.Param.FilterTicker(r.Ticker) {
				continue
			}
			t.AddRow(r.Ticker, p.c.Param.FormatDecimal(r.Valor), r.Codigo, r.Descr)
		}
		t.Render()
	}
}

func (p *PrinterTable) PrintRendimentosSujeitosTributacaoExclusiva(w io.Writer) {
	for _, ano := range p.c.RendimentosSujeitosTributacaoExclusiva() {
		if !p.c.Param.FilterYear(ano.Ano) {
			continue
		}
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 3)
		t.AddHeaders(fmt.Sprint(ano.Ano), "RENDIMENTOS SUJEITOS À TRIBUTAÇÃO EXCLUSIVA")
		t.AddHeaders("Ticker", "Valor", "Código", "Descrição")
		t.SetAlignment(table.AlignLeft, table.AlignRight, table.AlignCenter, table.AlignLeft)
		for _, r := range ano.Rendimentos {
			if !p.c.Param.FilterTicker(r.Ticker) {
				continue
			}
			t.AddRow(r.Ticker, p.c.Param.FormatDecimal(r.Valor), r.Codigo, r.Descr)
		}
		t.Render()
	}
}

func (p *PrinterTable) PrintOperacoesComunsDayTrade(w io.Writer) {

	print := func(w io.Writer, rs []RendimentosTributaveis, title string) {
		for _, r := range rs {
			if !p.c.Param.FilterYear(r.Ano) {
				continue
			}
			t := table.New(w)
			t.SetRowLines(false)
			t.SetHeaderColSpans(0, 1, 6)
			t.AddHeaders(fmt.Sprint(r.Ano), title)
			t.AddHeaders(
				"Mês",
				"Ações",
				"Opções",
				"Acumulado",
				fmt.Sprintf("IR (%s%%)", r.IR.Mul(decimal.NewFromInt(100))),
				fmt.Sprintf("IRRF (%s%%)", r.IRRF.Mul(decimal.NewFromInt(100))),
				"DARF",
			)
			t.SetAlignment(table.AlignCenter, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight)
			for _, v := range r.Meses {
				t.AddRow(
					translateMonth[v.Mes],
					p.c.Param.FormatDecimal(v.Lucro),
					p.c.Param.FormatDecimal(v.LucroOp),
					p.c.Param.FormatDecimal(v.LucroAc),
					p.c.Param.FormatDecimal(v.IR),
					p.c.Param.FormatDecimal(v.IRRF),
					p.c.Param.FormatDecimal(v.DARF),
				)
			}
			t.SetFooterAlignment(table.AlignCenter, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight)
			t.AddFooters(
				"Total",
				p.c.Param.FormatDecimal(r.TotalAcoes),
				p.c.Param.FormatDecimal(r.TotalOpcao),
				p.c.Param.FormatDecimal(r.TotalAc),
				p.c.Param.FormatDecimal(r.TotalIR),
				p.c.Param.FormatDecimal(r.TotalIRRF),
				p.c.Param.FormatDecimal(r.TotalDARF),
			)
			t.Render()
		}
	}

	var a, b bytes.Buffer
	print(&a, p.c.OperacoesComunsDayTradeSwingTrade(), "OPERAÇÕES COMUNS/DAY-TRADE [SWING TRADE]")
	print(&b, p.c.OperacoesComunsDayTradeDayTrade(), "OPERAÇÕES COMUNS/DAY-TRADE [DAY TRADE]")

	// Mescla as duas tabelas lado a lado.
	sa := bytes.Split(a.Bytes(), []byte("\n"))
	sb := bytes.Split(b.Bytes(), []byte("\n"))
	pd := bytes.Repeat([]byte(" "), utf8.RuneCount(sa[0]))
	for _, v := range lo.Zip2(sa, sb) {
		w.Write(v.A)
		if len(v.A) == 0 {
			w.Write(pd)
		}
		if len(v.B) != 0 {
			w.Write([]byte(" "))
			w.Write(v.B)
		}
		w.Write([]byte("\n"))
	}
}
