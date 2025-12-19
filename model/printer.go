package model

import (
	"fmt"
	"io"

	"github.com/aquasecurity/table"
	"github.com/shopspring/decimal"
)

func (s *State) Print(w io.Writer) {
	s.PrintOperacoesAcoes(w)
	s.PrintBensDireitos(w)
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
	for _, v := range s.Operacoes {
		t.AddRow(
			fmt.Sprint(v.ID),
			v.Ticker,
			v.Data.Format("2006-01-02"),
			string(v.Tipo),
			s.fracOrInt(v.Qtd)+s.sprintFracao(v),
			s.formatDecimal(v.ValorUnitario),
			s.formatDecimal(v.ValorTotal),
			s.formatDecimal(v.Taxas),
			s.fracOrInt(v.Agg.Qtd),
			s.formatDecimal(v.Agg.ValorTotal),
			s.formatDecimal(v.Agg.PrecoMedio),
			s.formatDecimal(v.ValorCompra),
			s.formatDecimal(v.Lucro),
		)
	}
	t.Render()
}

func (s *State) PrintBensDireitos(w io.Writer) {

	for _, bens := range s.BensDireitos() {
		t := table.New(w)
		t.SetHeaderColSpans(0, 4)
		t.SetHeaderAlignment(table.AlignCenter)
		t.AddHeaders("BENS E DIREITOS ── Grupo 03 ── Código 01")
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
				ticker.Discriminacao(),
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
			t.AddRow(translateMonth(mes.Mes), s.formatDecimal(mes.Valor))
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
		t.SetHeaderColSpans(0, 1, 3)
		t.AddHeaders(fmt.Sprintf("%s %d", r.Ticker, r.Ano), "OPERAÇÕES COMUNS/DAY-TRADE")
		t.AddHeaders("Mês", "Lucro", "Lucro Ac.", "IR (15%)")
		t.SetAlignment(table.AlignCenter, table.AlignRight, table.AlignRight, table.AlignRight)
		for _, v := range r.Meses {
			t.AddRow(translateMonth(v.Mes), s.formatDecimal(v.Lucro), s.formatDecimal(v.LucroAc), s.formatDecimal(v.IR))
		}
		t.SetFooterAlignment(table.AlignCenter, table.AlignRight, table.AlignRight, table.AlignRight)
		t.AddFooters("Total", s.formatDecimal(r.Total), s.formatDecimal(r.TotalAc), s.formatDecimal(r.TotalIR))
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
	if s.Settings.MostrarValorExato {
		return v.String()
	}
	return v.StringFixed(2)
}

func translateMonth(monthNumber string) string {
	return map[string]string{
		"01": "JAN",
		"02": "FEV",
		"03": "MAR",
		"04": "ABR",
		"05": "MAI",
		"06": "JUN",
		"07": "JUL",
		"08": "AGO",
		"09": "SET",
		"10": "OUT",
		"11": "NOV",
		"12": "DEZ",
	}[monthNumber]
}
