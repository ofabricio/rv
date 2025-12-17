package model

import (
	"bufio"
	"cmp"
	"encoding/json/v2"
	"flag"
	"fmt"
	"iter"
	"maps"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type State struct {
	Operacoes Operacoes

	Settings Settings
}

type Settings struct {
	// Determina a estratégia de Preço Médio usada na Bonificação.
	//
	// Se true, as Bonificações vão alterar seu preço médio.
	// É como se você tivesse pagado o valor que recebeu de bonificação.
	// Significa que na hora de vender você pagará menos impostos.
	//
	// Se false, as Bonificações NÃO vão alterar seu preço médio.
	// É como se você tivesse recebido as ações de bonificação de graça.
	// Significa que na hora de vender você pagará mais impostos.
	//
	// Não há consenso sobre qual estratégia é a correta.
	// Boatos afirmam que a maioria prefere a estratégia que reduz o imposto a pagar (true).
	// Para nunca ter problemas com a Receita Federal use (false), pois pagará mais impostos.
	AlterarPrecoMedioNaBonificacao bool

	// Mostra os valores exatos.
	MostrarValorExato bool

	AcaoSwingTradeIR decimal.Decimal
	AcaoDayTradeIR   decimal.Decimal
	AcaoLimiteIsento decimal.Decimal
}

func (s *State) Load(file string) {
	s.Operacoes = nil
	for line := range FileLines(file) {
		var o Operacao
		unmarshal(line, &o)
		if o.ID == 0 {
			var state State
			unmarshal(line, &state)
			s.Settings = state.Settings
		} else {
			s.Calculate(&o)
			s.Operacoes = append(s.Operacoes, o)
		}
	}
}

func (s *State) CommandLine() {

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("  Usage: rv <action>")
		fmt.Println()
		fmt.Println("  <action>:")
		fmt.Println("    show (default; pode ser omitido)")
		fmt.Println("        mostra os ativos")
		fmt.Println("    add")
		fmt.Println("        adiciona um ativo")
		fmt.Println()
		flag.PrintDefaults()
	}
	flag.BoolVar(&s.Settings.AlterarPrecoMedioNaBonificacao, "AlterarPrecoMedioNaBonificacao", false, "Determina a estratégia de Preço Médio usada na Bonificação.\nSe true, as Bonificações vão alterar seu preço médio.\nSe false, as Bonificações NÃO vão alterar seu preço médio.\nDefault: false")
	flag.Parse()
}

func (s *State) BensDireitos() []BensDireitos {

	var bens []BensDireitos
	var opcs []BensDireitos

	prevYear := map[string]Operacao{}
	for _, oprs := range s.Operacoes.PartitionByYear() {
		bem := BensDireitos{
			AnoAnterior: oprs[0].Data.Year() - 1,
			AnoCorrente: oprs[0].Data.Year(),
			Grupo:       "03",
			Codigo:      "01",
		}
		currYear := map[string]Operacao{}
		for ticker, oprs := range oprs.GroupByTicker() {
			currYear[ticker] = lo.LastOrEmpty(oprs) // Apenas a última operação é relevante, pois tem o saldo acumulado.
		}
		merge := lo.Assign(prevYear, currYear)
		for ticker := range sortKeysIter(merge) {
			if _, ok := currYear[ticker]; !ok && prevYear[ticker].Agg.ValorTotal.IsZero() && currYear[ticker].Agg.ValorTotal.IsZero() {
				continue // Situações zeradas no ano anterior e corrente cujo ano corrente não existe não são mostradas.
			}
			o := merge[ticker]
			bem.Tickers = append(bem.Tickers, BensDireitoTicker{
				Ticker:           ticker,
				SituacaoAnterior: prevYear[ticker].Agg.ValorTotal,
				SituacaoCorrente: o.Agg.ValorTotal,
				Discriminacao:    fmt.Sprintf("%s AÇÕES %s COM PREÇO MÉDIO DE R$ %s", o.Agg.Qtd.String(), ticker, s.formatDecimal(o.Agg.PrecoMedio)),
			})
		}
		prevYear = merge
		bens = append(bens, bem)

		if opc := s.BensDireitosOpcoes(oprs); len(opc) > 0 {
			opcs = append(opcs, BensDireitos{
				AnoAnterior: oprs[0].Data.Year() - 1,
				AnoCorrente: oprs[0].Data.Year(),
				Grupo:       "04",
				Codigo:      "04",
				Tickers:     opc,
			})
		}
	}
	if len(bens)%2 == 0 && len(bens) > 0 {
		bens = bens[1:] // Remove a primeira posição se o número de anos for par.
	}
	return append(bens, opcs...)
}

func (s *State) BensDireitosOpcoes(oprs Operacoes) []BensDireitoTicker {
	var opc []BensDireitoTicker
	for _, o := range lo.Filter(oprs, func(o Operacao, _ int) bool { return o.CID == 0 && (o.Tipo == PUT_COMPRA || o.Tipo == CALL_COMPRA) }) {
		opc = append(opc, BensDireitoTicker{
			Ticker:        o.Opcao,
			Discriminacao: fmt.Sprintf("%s OPÇÕES COMPRADAS %s DE STRIKE R$ %s", o.Qtd.String(), o.Opcao, s.formatDecimal(o.ValorUnitario)),
		})
	}
	return opc
}

type BensDireitos struct {
	AnoAnterior int
	AnoCorrente int
	Grupo       string
	Codigo      string
	Tickers     []BensDireitoTicker
}

type BensDireitoTicker struct {
	Ticker           string
	SituacaoAnterior decimal.Decimal
	SituacaoCorrente decimal.Decimal
	Discriminacao    string
}

func (s *State) RendimentosIsentosNaoTributaveis() []RendimentosIsentosNaoTributaveis {

	var rs []RendimentosIsentosNaoTributaveis

	oprs := lo.Filter(s.Operacoes, func(o Operacao, _ int) bool {
		return tipoOprConfig[o.Tipo].IsRendimentoIsentoNaoTributavel && o.Lucro.IsPositive()
	})

	for year, oprs := range oprs.GroupByYear() {
		var r RendimentosIsentosNaoTributaveis
		r.Ano = year
		for ticker, oprs := range oprs.GroupByTicker() {
			g := lo.GroupBy(oprs, func(o Operacao) string { return tipoOprConfig[o.Tipo].Codigo })
			for cod, oprs := range sortKeysIter(g) {
				r.Rendimentos = append(r.Rendimentos, RendimentoIsentoNaoTributavel{
					Ticker: ticker,
					Valor:  oprs.LucroAcumulado(),
					Codigo: cod,
				})
			}
		}
		rs = append(rs, r)
	}

	return rs
}

func (s *State) RendimentosSujeitosTributacaoExclusiva() []RendimentosIsentosNaoTributaveis {

	var rs []RendimentosIsentosNaoTributaveis

	oprs := lo.Filter(s.Operacoes, func(o Operacao, _ int) bool {
		return tipoOprConfig[o.Tipo].IsRendimentoSujeitoTributacaoExclusiva && o.Lucro.IsPositive()
	})

	for year, oprs := range oprs.GroupByYear() {
		var r RendimentosIsentosNaoTributaveis
		r.Ano = year
		for ticker, oprs := range oprs.GroupByTicker() {
			g := lo.GroupBy(oprs, func(o Operacao) string { return tipoOprConfig[o.Tipo].Codigo })
			for cod, oprs := range sortKeysIter(g) {
				r.Rendimentos = append(r.Rendimentos, RendimentoIsentoNaoTributavel{
					Ticker: ticker,
					Valor:  oprs.LucroAcumulado(),
					Codigo: cod,
				})
			}
		}
		rs = append(rs, r)
	}

	return rs
}

type RendimentosIsentosNaoTributaveis struct {
	Ano         int
	Rendimentos []RendimentoIsentoNaoTributavel
}

type RendimentoIsentoNaoTributavel struct {
	Ticker string
	Valor  decimal.Decimal
	Codigo string
}

func (s *State) RendimentosIsentosNaoTributaveisAte20k() []RendimentosIsentosAte20k {

	var ris []RendimentosIsentosAte20k

	oprs := lo.Filter(s.Operacoes, func(o Operacao, _ int) bool { return tipoOprConfig[o.Tipo].IsRendimentoIsentoAte20k })

	for year, oprs := range oprs.GroupByYear() {

		var ri RendimentosIsentosAte20k
		ri.Ano = year

		for month, oprs := range oprs.GroupByMonth() {

			vendaNoMes := oprs.ValorTotalAcumulado()
			lucroNoMes := oprs.LucroAcumulado()

			if vendaNoMes.LessThanOrEqual(s.Settings.AcaoLimiteIsento) && lucroNoMes.IsPositive() {
				ri.Meses = append(ri.Meses, RendimentoIsentosAte20kMensal{Mes: month, Valor: lucroNoMes})
				ri.Total = ri.Total.Add(lucroNoMes)
			}
		}
		if len(ri.Meses) > 0 {
			ris = append(ris, ri)
		}
	}

	return ris
}

type RendimentosIsentosAte20k struct {
	Ticker string
	Ano    int
	Meses  []RendimentoIsentosAte20kMensal
	Total  decimal.Decimal
}

type RendimentoIsentosAte20kMensal struct {
	Mes   string
	Valor decimal.Decimal
}

func (s *State) OperacoesComunsDayTrade() []RendimentosTributaveis {

	var rts []RendimentosTributaveis

	lucroAcumuladoPelosAnos := decimal.Zero
	for year, oprs := range s.Operacoes.GroupByYear() {

		var rt RendimentosTributaveis
		rt.Ano = year

		for month, oprs := range oprs.GroupByMonth() {

			lucroNoMes := decimal.Zero
			lucroNoMesOp := decimal.Zero

			for _, p := range oprs.PartitionByAcaoOpcao() {
				vendaNoMes := decimal.Zero
				var lucro *decimal.Decimal
				if p[0].IsOpcao() {
					lucro = &lucroNoMesOp
				} else {
					vendaNoMes = p.TotalVendas()
					lucro = &lucroNoMes
				}
				for _, o := range p {
					if tipoOprConfig[o.Tipo].IsRendimentoTributavelApos20k && vendaNoMes.GreaterThan(s.Settings.AcaoLimiteIsento) && o.Lucro.IsPositive() {
						*lucro = lucro.Add(o.Lucro)
					}
					if tipoOprConfig[o.Tipo].IsLucroTributavel && o.Lucro.IsPositive() {
						*lucro = lucro.Add(o.Lucro)
					}
					if tipoOprConfig[o.Tipo].IsPrejuizoAbativel && o.Lucro.IsNegative() {
						*lucro = lucro.Add(o.Lucro)
					}
				}
			}

			if lucroNoMes.IsZero() && lucroNoMesOp.IsZero() {
				continue
			}

			lucroAcumuladoPelosAnos = lucroAcumuladoPelosAnos.Add(lucroNoMes).Add(lucroNoMesOp)
			rtm := RendimentoTributavelMensal{Mes: month, Lucro: lucroNoMes, LucroOp: lucroNoMesOp, LucroAc: lucroAcumuladoPelosAnos}
			if lucroAcumuladoPelosAnos.IsPositive() {
				rtm.IR = lucroAcumuladoPelosAnos.Mul(s.Settings.AcaoSwingTradeIR)
				lucroAcumuladoPelosAnos = decimal.Zero
			}
			rt.Meses = append(rt.Meses, rtm)
			rt.Total = rt.Total.Add(lucroNoMes)
			rt.TotalOp = rt.TotalOp.Add(lucroNoMesOp)
			rt.TotalAc = lucroAcumuladoPelosAnos
			rt.TotalIR = rt.TotalIR.Add(rtm.IR)
		}
		if len(rt.Meses) > 0 {
			rts = append(rts, rt)
		}
	}

	return rts
}

type RendimentosTributaveis struct {
	Ticker  string
	Ano     int
	Meses   []RendimentoTributavelMensal
	Total   decimal.Decimal
	TotalOp decimal.Decimal
	TotalAc decimal.Decimal
	TotalIR decimal.Decimal
}

type RendimentoTributavelMensal struct {
	Mes     string
	Lucro   decimal.Decimal
	LucroOp decimal.Decimal
	LucroAc decimal.Decimal
	IR      decimal.Decimal
}

func (s *State) Calculate(o *Operacao) {

	p := s.Operacoes.GetID(o.PID)

	o.Inherit(p)

	switch o.Tipo {
	case COMPRA:
		o.CalcCompra(p)
	case VENDA:
		o.CalcVenda(p)
	case BONIFICACAO:
		o.CalcBonificacao(p, s.Settings.AlterarPrecoMedioNaBonificacao)
	case DESDOBRAMENTO:
		o.CalcDesdobramento(p)
	case GRUPAMENTO:
		o.CalcGrupamento(p)
	case LEILAO_FRACAO:
		o.CalcLeilaoFracao(p)
	case DIVIDENDOS:
		o.CalcDividendos(p)
	case JSCP:
		o.CalcJSCP(p)
	case REND_TRIB:
		o.CalcRendTrib(p)
	case REDUCAO_CAPITAL:
		o.CalcReducaoCapital(p)
	case SUBSCRICAO_COMPRA:
		o.CalcSubscricaoCompra(p)
	case SUBSCRICAO_VENDA:
		o.CalcSubscricaoVenda(p)
	case SUBSCRICAO_EXERCICIO:
		o.CalcSubscricaoExercicio(p)
	case PUT_VENDA:
		o.CalcVendaPut(p)
	case PUT_VENDA_EX:
		o.CalcVendaPutExercido(p)
	case PUT_VENDA_NE:
		o.CalcVendaPutNaoExercido(p)
	case PUT_COMPRA:
		o.CalcCompraPut(p)
	case PUT_COMPRA_EX:
		o.CalcCompraPutExercido(p)
	case PUT_COMPRA_NE:
		o.CalcCompraPutNaoExercido(p)
	case CALL_COMPRA:
		o.CalcCompraCall(p)
	case CALL_COMPRA_EX:
		o.CalcCompraCallExercido(p)
	case CALL_COMPRA_NE:
		o.CalcCompraCallNaoExercido(p)
	case CALL_VENDA:
		o.CalcVendaCall(p)
	case CALL_VENDA_EX:
		o.CalcVendaCallExercido(p)
	case CALL_VENDA_NE:
		o.CalcVendaCallNaoExercido(p)
	}
}

type Operacoes []Operacao

func (o Operacoes) GroupByTicker() iter.Seq2[string, Operacoes] {
	return func(yield func(string, Operacoes) bool) {
		g := lo.GroupByMap(o, func(o Operacao) (string, Operacao) { return o.Ticker, o })
		for _, k := range sortKeys(g) {
			if !yield(k, g[k]) {
				return
			}
		}
	}
}

func (o Operacoes) GroupByYear() iter.Seq2[int, Operacoes] {
	return func(yield func(int, Operacoes) bool) {
		g := lo.GroupByMap(o, func(o Operacao) (int, Operacao) { return o.Data.Year(), o })
		for _, k := range sortKeys(g) {
			if !yield(k, g[k]) {
				return
			}
		}
	}
}

func (o Operacoes) GroupByMonth() iter.Seq2[string, Operacoes] {
	return func(yield func(string, Operacoes) bool) {
		g := lo.GroupByMap(o, func(o Operacao) (string, Operacao) { return o.Data.Format("01"), o })
		for _, k := range sortKeys(g) {
			if !yield(k, g[k]) {
				return
			}
		}
	}
}

func (o Operacoes) ValorTotalAcumulado() decimal.Decimal {
	return lo.Reduce(o, func(agg decimal.Decimal, o Operacao, _ int) decimal.Decimal {
		return agg.Add(o.ValorTotal)
	}, decimal.Zero)
}

func (o Operacoes) LucroAcumulado() decimal.Decimal {
	return lo.Reduce(o, func(agg decimal.Decimal, o Operacao, _ int) decimal.Decimal {
		return agg.Add(o.Lucro)
	}, decimal.Zero)
}

func (o Operacoes) TotalVendas() decimal.Decimal {
	return lo.Reduce(o, func(agg decimal.Decimal, o Operacao, _ int) decimal.Decimal {
		if o.Tipo == VENDA {
			return agg.Add(o.ValorTotal)
		}
		return agg
	}, decimal.Zero)
}

func (o Operacoes) PartitionByYear() []Operacoes {
	return lo.PartitionBy(o, func(o Operacao) int { return o.Data.Year() })
}

func (o Operacoes) PartitionByAcaoOpcao() []Operacoes {
	return lo.PartitionBy(o, func(o Operacao) bool { return o.IsOpcao() })
}

func (o Operacoes) GetID(id int64) Operacao {
	if id == 0 {
		return Operacao{}
	}
	return o[id-1]
}

type Agregado struct {
	Qtd        decimal.Decimal
	ValorTotal decimal.Decimal
	PrecoMedio decimal.Decimal
}

func (a *Agregado) CalcPrecoMedio() {
	if a.Qtd.IsZero() {
		a.PrecoMedio = decimal.Zero
	} else {
		a.PrecoMedio = a.ValorTotal.Div(a.Qtd)
	}
}

type Operacao struct {
	ID            int64
	PID           int64
	Ticker        string
	Tipo          TipoOpr
	Data          time.Time `json:"Data,format:DateOnly"`
	Qtd           decimal.Decimal
	ValorUnitario decimal.Decimal
	Taxas         decimal.Decimal
	ValorTotal    decimal.Decimal
	ValorCompra   decimal.Decimal
	Lucro         decimal.Decimal // Lucro ou prejuízo da operação de Venda, Bonificação, Grupamento, Subscrição Compra, Redução de Capital, Opções.
	Fracao        decimal.Decimal // Parte fracionária resultante de Bonificação, Grupamento ou Desdobramento.
	Fator         decimal.Decimal // Fator de Bonificação, Grupamento ou Desdobramento e Redução de Capital.
	Agg           Agregado

	// Opções.
	CID   int64
	Opcao string
}

func (o *Operacao) IsOpcao() bool {
	return o.Opcao != ""
}

func (o *Operacao) CalcCompra(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Agg.Qtd = p.Agg.Qtd.Add(o.Qtd)
	o.Agg.ValorTotal = p.Agg.ValorTotal.Add(o.ValorTotal)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcVenda(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Sub(o.Taxas)
	o.ValorCompra = p.Agg.PrecoMedio.Mul(o.Qtd)
	o.Agg.Qtd = p.Agg.Qtd.Sub(o.Qtd)
	o.Agg.ValorTotal = p.Agg.ValorTotal.Sub(o.ValorCompra)
	o.Lucro = o.ValorTotal.Sub(o.ValorCompra)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcBonificacao(p Operacao, alterarPrecoMedio bool) {
	if !o.Fator.IsZero() {
		o.Fracao = p.Agg.Qtd.Mul(o.Fator)
		o.Qtd = o.Fracao.Truncate(0)
		o.Fracao = o.Fracao.Sub(o.Qtd).Abs()
	}
	o.Agg.Qtd = p.Agg.Qtd.Add(o.Qtd)
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Agg.ValorTotal = o.ValorTotal.Add(p.Agg.ValorTotal)
	o.Lucro = o.ValorTotal
	if alterarPrecoMedio {
		o.Agg.CalcPrecoMedio()
	}
}

func (o *Operacao) CalcDesdobramento(p Operacao) {
	o.Agg.Qtd = p.Agg.Qtd.Mul(o.Fator)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcGrupamento(p Operacao) {
	qtdAcFra := p.Agg.Qtd.Div(o.Fator)
	qtdAcInt := qtdAcFra.Truncate(0)
	o.Fracao = qtdAcFra.Sub(qtdAcInt)
	o.Agg.Qtd = qtdAcInt
	o.Agg.ValorTotal = qtdAcInt.Mul(p.Agg.ValorTotal.Div(qtdAcFra))
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcLeilaoFracao(p Operacao) {
	o.Lucro = o.ValorUnitario.Mul(o.Qtd).Sub(o.Agg.PrecoMedio.Mul(o.Qtd))
}

func (o *Operacao) CalcDividendos(p Operacao) {
}

func (o *Operacao) CalcJSCP(p Operacao) {
}

func (o *Operacao) CalcRendTrib(p Operacao) {
}

func (o *Operacao) CalcSubscricaoCompra(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Agg.ValorTotal = p.Agg.ValorTotal.Add(o.ValorTotal)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcSubscricaoVenda(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Sub(o.Taxas)
	o.Lucro = o.ValorTotal
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcSubscricaoExercicio(p Operacao) {
	o.Agg.Qtd = p.Agg.Qtd.Add(o.Qtd)
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Agg.ValorTotal = p.Agg.ValorTotal.Add(o.ValorTotal)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcReducaoCapital(p Operacao) {
	if o.Fator.IsPositive() {
		redcap := p.Agg.ValorTotal.Mul(o.Fator)
		restituicao := p.Agg.Qtd.Mul(o.ValorUnitario)
		o.ValorTotal = restituicao
		o.ValorCompra = redcap
		o.Lucro = restituicao.Sub(redcap)
		o.Agg.ValorTotal = p.Agg.ValorTotal.Sub(redcap)
	} else {
		o.Agg.ValorTotal = p.Agg.ValorTotal.Sub(o.ValorTotal)
	}
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcVendaPut(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
}

func (o *Operacao) CalcVendaPutExercido(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(p.Qtd)
	o.Lucro = o.ValorUnitario.Sub(p.ValorUnitario).Mul(p.Qtd).Add(p.Lucro.Mul(p.Qtd)) // (VA - Strike) + Prêmio
	o.Agg.Qtd = p.Agg.Qtd.Add(p.Qtd)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcVendaPutNaoExercido(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(p.Qtd)
	o.Lucro = p.Lucro.Mul(p.Qtd)
}

func (o *Operacao) CalcCompraPut(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
}

func (o *Operacao) CalcCompraPutExercido(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(p.Qtd)
	o.Lucro = p.ValorUnitario.Sub(o.ValorUnitario.Add(p.Lucro)).Mul(p.Qtd) // Strike - (VA + Custo)
	o.Agg.Qtd = p.Agg.Qtd.Sub(p.Qtd)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcCompraPutNaoExercido(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(p.Qtd)
	o.Lucro = p.Lucro.Neg().Mul(p.Qtd)
}

func (o *Operacao) CalcCompraCall(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
}

func (o *Operacao) CalcCompraCallExercido(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(p.Qtd)
	o.Lucro = o.ValorUnitario.Sub(p.ValorUnitario.Add(p.Lucro)).Mul(p.Qtd) // VA - (Strike + Prêmio)
	o.Agg.Qtd = p.Agg.Qtd.Add(p.Qtd)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcCompraCallNaoExercido(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(p.Qtd)
	o.Lucro = p.Lucro.Neg().Mul(p.Qtd)
}

func (o *Operacao) CalcVendaCall(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
}

func (o *Operacao) CalcVendaCallExercido(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(p.Qtd)
	o.Lucro = p.ValorUnitario.Add(p.Lucro).Sub(o.ValorUnitario).Mul(p.Qtd) // (Strike + Prêmio) - VA
	o.Agg.Qtd = p.Agg.Qtd.Sub(p.Qtd)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcVendaCallNaoExercido(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(p.Qtd)
	o.Lucro = p.Lucro.Mul(p.Qtd)
}

func (o *Operacao) Inherit(p Operacao) {
	o.Agg = p.Agg
}

func unmarshal(line []byte, v any) {
	if err := json.Unmarshal(line, v); err != nil {
		panic(err)
	}
}

func FileLines(file string) func(yield func([]byte) bool) {
	return func(yield func([]byte) bool) {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		lines := bufio.NewScanner(f)
		for lines.Scan() {
			if lines.Text() == "" || strings.HasPrefix(lines.Text(), "//") {
				continue
			}
			if !yield(lines.Bytes()) {
				return
			}
		}
		if lines.Err() != nil {
			panic(lines.Err())
		}
	}
}

func sortKeysIter[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, v := range sortKeys(m) {
			if !yield(v, m[v]) {
				return
			}
		}
	}
}

func sortKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)
	return keys
}
