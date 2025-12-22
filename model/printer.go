package model

import (
	"fmt"
	"io"
	"strings"

	"github.com/aquasecurity/table"
	"github.com/shopspring/decimal"
)

func (s *State) Print(w io.Writer) {
	s.PrintOperacoesAcoes(w)
	s.PrintBensDireitos(w)
	s.PrintDividaOnusReais(w)
	s.PrintRendimentosIsentosNaoTributaveis(w)
	s.PrintRendimentosIsentosNaoTributaveisAte20k(w)
	s.PrintRendimentosSujeitosTributacaoExclusiva(w)
	s.PrintOperacoesComunsDayTrade(w)
}

func (s *State) PrintOperacoesAcoes(w io.Writer) {
	t := table.New(w)
	t.SetRowLines(false)
	t.SetHeaderColSpans(0, 13)
	t.AddHeaders("OPERAÇÕES COM AÇÕES")
	t.AddHeaders(
		"ID",
		"Ticker",
		"Data",
		"Opr",
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
		table.AlignCenter,
		table.AlignLeft,
		table.AlignLeft,
		table.AlignLeft,
		table.AlignRight,
		table.AlignRight,
		table.AlignRight,
		table.AlignRight,
		table.AlignRight,
		table.AlignRight,
		table.AlignRight,
		table.AlignRight,
		table.AlignRight,
	)
	for i, v := range s.Operacoes {
		t.AddRow(
			fmt.Sprint(i+1),
			fmt.Sprintf("%s%12s", v.Ticker, v.Serie),
			s.formatColumnData(v),
			s.formatColumnTipo(v),
			s.fracOrInt(v.Qtd)+s.sprintFracao(v),
			s.formatColumnValorUnitario(v),
			s.formatColumnValorTotal(v),
			s.formatDecimal(v.Taxas),
			s.fracOrInt(v.Agg.Qtd),
			s.formatDecimal(v.Agg.ValorTotal),
			s.formatDecimal(v.Agg.PrecoMedio),
			s.formatDecimal(v.ValorCompra),
			s.formatColumnLucro(v),
		)
	}
	t.Render()
}

func (s *State) PrintBensDireitos(w io.Writer) {

	for _, bens := range s.BensDireitos() {
		t := table.New(w)
		t.SetColumnMaxWidth(100)
		t.SetHeaderColSpans(0, 4)
		t.SetHeaderAlignment(table.AlignCenter)
		t.AddHeaders(fmt.Sprintf("BENS E DIREITOS ── Grupo %s ── Código %s", bens.Grupo, bens.Codigo))
		t.SetAlignment(table.AlignLeft, table.AlignCenter, table.AlignCenter, table.AlignLeft)
		t.AddHeaders(
			"Ticker",
			fmt.Sprintf("Situação em %d", bens.AnoAnterior),
			fmt.Sprintf("Situação em %d", bens.AnoCorrente),
			"Discriminação",
		)
		t.SetAlignment(table.AlignLeft, table.AlignCenter, table.AlignCenter, table.AlignLeft)
		for _, ticker := range bens.Tickers {
			t.AddRow(
				ticker.Ticker,
				s.formatDecimal(ticker.SituacaoAnterior),
				s.formatDecimal(ticker.SituacaoCorrente),
				ticker.Discriminacao,
			)
		}
		t.Render()
	}
}

func (s *State) PrintDividaOnusReais(w io.Writer) {

	for _, bens := range s.DividaOnusReais() {
		t := table.New(w)
		t.SetColumnMaxWidth(100)
		t.SetHeaderColSpans(0, 4)
		t.SetHeaderAlignment(table.AlignCenter)
		t.AddHeaders(fmt.Sprintf("DÍVIDA E ÔNUS REAIS ── Código %s", bens.Codigo))
		t.SetAlignment(table.AlignLeft, table.AlignCenter, table.AlignCenter, table.AlignLeft)
		t.AddHeaders(
			"Ticker",
			fmt.Sprintf("Situação em %d", bens.AnoAnterior),
			fmt.Sprintf("Situação em %d", bens.AnoCorrente),
			"Discriminação",
		)
		t.SetAlignment(table.AlignLeft, table.AlignCenter, table.AlignCenter, table.AlignLeft)
		for _, ticker := range bens.Tickers {
			t.AddRow(
				ticker.Ticker,
				s.formatDecimal(ticker.SituacaoAnterior),
				s.formatDecimal(ticker.SituacaoCorrente),
				ticker.Discriminacao,
			)
		}
		t.Render()
	}
}

func (s *State) PrintRendimentosIsentosNaoTributaveis(w io.Writer) {
	for _, ano := range s.RendimentosIsentosNaoTributaveis() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 2)
		t.AddHeaders(fmt.Sprint(ano.Ano), "RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS")
		t.AddHeaders("Ticker", "Valor", "Código")
		t.SetAlignment(table.AlignLeft, table.AlignRight, table.AlignLeft)
		for _, r := range ano.Rendimentos {
			t.AddRow(r.Ticker, s.formatDecimal(r.Valor), r.Codigo)
		}
		t.Render()
	}
}

func (s *State) PrintRendimentosSujeitosTributacaoExclusiva(w io.Writer) {
	for _, ano := range s.RendimentosSujeitosTributacaoExclusiva() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 2)
		t.AddHeaders(fmt.Sprint(ano.Ano), "RENDIMENTOS SUJEITOS À TRIBUTAÇÃO EXCLUSIVA")
		t.AddHeaders("Ticker", "Valor", "Código")
		t.SetAlignment(table.AlignLeft, table.AlignRight, table.AlignLeft)
		for _, r := range ano.Rendimentos {
			t.AddRow(r.Ticker, s.formatDecimal(r.Valor), r.Codigo)
		}
		t.Render()
	}
}

func (s *State) PrintRendimentosIsentosNaoTributaveisAte20k(w io.Writer) {
	for _, r := range s.RendimentosIsentosNaoTributaveisAte20k() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderAlignment(table.AlignCenter, table.AlignCenter)
		t.SetHeaders(fmt.Sprint(r.Ano), "RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS")
		t.SetAlignment(table.AlignCenter, table.AlignRight)
		for _, mes := range r.Meses {
			t.AddRow(translateMonth[mes.Mes], s.formatDecimal(mes.Valor))
		}
		t.SetFooterAlignment(table.AlignCenter, table.AlignRight)
		t.SetFooters("Códig 20", "Isenção até R$ 20000"+" │ "+s.formatDecimal(r.Total))
		t.Render()
	}
}

func (s *State) PrintOperacoesComunsDayTrade(w io.Writer) {
	for _, r := range s.OperacoesComunsDayTrade() {
		t := table.New(w)
		t.SetRowLines(false)
		t.SetHeaderColSpans(0, 1, 4)
		t.AddHeaders(fmt.Sprintf("%s %d", r.Ticker, r.Ano), "OPERAÇÕES COMUNS/DAY-TRADE")
		t.AddHeaders("Mês", "Ações", "Opções", "Acumulado", "IR (15%)")
		t.SetAlignment(table.AlignCenter, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight)
		for _, v := range r.Meses {
			t.AddRow(
				translateMonth[v.Mes],
				s.formatDecimal(v.Lucro),
				s.formatDecimal(v.LucroOp),
				s.formatDecimal(v.LucroAc),
				s.formatDecimal(v.IR),
			)
		}
		t.SetFooterAlignment(table.AlignCenter, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight)
		t.AddFooters("Total", s.formatDecimal(r.Total), s.formatDecimal(r.TotalOp), s.formatDecimal(r.TotalAc), s.formatDecimal(r.TotalIR))
		t.Render()
	}
}

func (s *State) sprintFracao(o Operacao) string {
	if o.Fracao.IsPositive() {
		return fmt.Sprintf(" (%s)", s.formatDecimal(o.Fracao))
	}
	return ""
}

func (s *State) fracOrInt(d decimal.Decimal) string {
	if d.Truncate(0).Equal(d) {
		return d.Truncate(0).String()
	}
	return s.formatDecimal(d)
}

func (s *State) formatDecimal(v decimal.Decimal) string {
	var r string
	if s.Config.MostrarValorExato {
		r = v.String()
	} else {
		r = v.StringFixed(2)
	}
	return strings.Replace(r, ".", s.Config.SeparadorDecimal, 1)
}

func (s *State) formatColumnLucro(o Operacao) string {
	if !o.Premio.IsZero() {
		return fmt.Sprintf("(P %s) %s", s.formatDecimal(o.Premio), s.formatDecimal(o.Lucro))
	}
	return s.formatDecimal(o.Lucro)
}

func (s *State) formatColumnValorUnitario(o Operacao) string {
	if !o.ValorExercicio.IsZero() {
		return fmt.Sprintf("(E %s) %s", s.formatDecimal(o.ValorExercicio), s.formatDecimal(o.ValorUnitario))
	}
	return s.formatDecimal(o.ValorUnitario)
}

func (s *State) formatColumnValorTotal(o Operacao) string {
	if !o.ValorExercicio.IsZero() {
		return fmt.Sprintf("(E %s) %s", s.formatDecimal(o.ValorExercicio.Mul(o.Qtd)), s.formatDecimal(o.ValorTotal))
	}
	return s.formatDecimal(o.ValorTotal)
}

func (s *State) formatColumnTipo(o Operacao) string {
	return string(o.Tipo)
}

func (s *State) formatColumnData(o Operacao) string {
	if !o.Vencimento.IsZero() {
		return fmt.Sprintf("%s V %s", o.Data.Format(s.Config.FormatoData), o.Vencimento.Format(s.Config.FormatoData))
	}
	return o.Data.Format(s.Config.FormatoData)
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
