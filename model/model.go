package model

import (
	"bufio"
	"bytes"
	"cmp"
	"encoding/json/v2"
	"fmt"
	"iter"
	"maps"
	"os"
	"slices"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type State struct {
	Operacoes Operacoes

	Config Config
}

type Config struct {
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

	// Formata valores usando esse separador de casas decimais.
	// Exemplo: "1234,56" (padrão) ou "1234.56".
	SeparadorDecimal string

	// Formata datas usando esse formato.
	// Exemplo: "02/01/2006" (padrão) ou "2006-01-02".
	FormatoData string

	AcaoSwingTradeIR decimal.Decimal
	AcaoDayTradeIR   decimal.Decimal
	AcaoLimiteIsento decimal.Decimal
}

func DefaultConfig() Config {
	return Config{
		AlterarPrecoMedioNaBonificacao: false,
		MostrarValorExato:              false,
		SeparadorDecimal:               ",",
		FormatoData:                    "02/01/2006",
		AcaoSwingTradeIR:               decimal.NewFromFloat(0.15),
		AcaoDayTradeIR:                 decimal.NewFromFloat(0.20),
		AcaoLimiteIsento:               decimal.NewFromInt(20000),
	}
}

func (s *State) Load(file string) {
	s.Operacoes = nil
	s.Config = DefaultConfig()
	aggr := map[string]Agregado{}
	for line := range FileLines(file) {
		var o Operacao
		unmarshal(line, &o)
		if o.Tipo == "Config" {
			unmarshal(line, &s.Config)
			continue
		}
		a := aggr[o.Ticker]
		s.Calculate(&o, &a)
		aggr[o.Ticker] = a
		s.Operacoes = append(s.Operacoes, o)
	}
}

func (s *State) BensDireitos() []BensDireitos {

	var bens []BensDireitos

	prevYear := map[string]Operacao{}
	for _, oprs := range s.Operacoes.PartitionByYear() {
		bem := BensDireitos{
			AnoAnterior: oprs[0].Data.Year() - 1,
			AnoCorrente: oprs[0].Data.Year(),
			Grupo:       "03",
			Codigo:      "01",
		}
		bonifics := map[string]decimal.Decimal{}
		currYear := map[string]Operacao{}
		for ticker, oprs := range oprs.GroupByTicker() {
			currYear[ticker] = lo.LastOrEmpty(oprs) // Apenas a última operação é relevante, pois tem o saldo acumulado.
			bonifics[ticker] = oprs.FilterByTipo(BONIFICACAO).QtdAcumulado()
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
				Discriminacao: fmt.Sprintf("%s AÇÕES %s COM PREÇO MÉDIO DE R$ %s%s", o.Agg.Qtd, ticker, s.formatDecimal(o.Agg.PrecoMedio),
					lo.Ternary(bonifics[ticker].IsPositive(), fmt.Sprintf(" ONDE %s AÇÕES SÃO PROVENIENTES DE BONIFICAÇÃO", bonifics[ticker]), ""),
				),
			})
		}
		prevYear = merge
		bens = append(bens, bem)
	}
	if len(bens)%2 == 0 && len(bens) > 0 {
		bens = bens[1:] // Remove a primeira posição se o número de anos for par.
	}
	return append(bens, s.BensDireitosOpcoes()...)
}

func (s *State) BensDireitosOpcoes() []BensDireitos {
	return s.bensDireitosOpcoes(func(o Operacao) bool { return o.Tipo == PUT_COMPRA || o.Tipo == CALL_COMPRA }, "COMPRADAS", "04", "04")
}

func (s *State) DividaOnusReais() []BensDireitos {
	return s.bensDireitosOpcoes(func(o Operacao) bool { return o.Tipo == PUT_VENDA || o.Tipo == CALL_VENDA }, "VENDIDAS", "", "16")
}

func (s *State) bensDireitosOpcoes(cond func(Operacao) bool, tipo, grupo, cod string) []BensDireitos {

	var bens []BensDireitos

	oprs := lo.Filter(s.Operacoes, func(o Operacao, _ int) bool { return cond(o) && o.Vencimento.Year() > o.Data.Year() })

	prevYear := map[string]Operacao{}
	for _, oprs := range oprs.PartitionByYear() {
		bem := BensDireitos{
			AnoAnterior: oprs[0].Data.Year() - 1,
			AnoCorrente: oprs[0].Data.Year(),
			Grupo:       grupo,
			Codigo:      cod,
		}
		currYear := map[string]Operacao{}
		for _, op := range oprs {
			currYear[op.Serie] = op
		}
		merge := lo.Assign(prevYear, currYear)
		for ticker := range sortKeysIter(merge) {
			if merge[ticker].Vencimento.Year() <= bem.AnoCorrente {
				continue // Se já venceu não mostra.
			}
			prev := prevYear[ticker]
			curr := merge[ticker]
			bem.Tickers = append(bem.Tickers, BensDireitoTicker{
				Ticker:           ticker,
				SituacaoAnterior: prev.Premio.Mul(prev.Qtd),
				SituacaoCorrente: curr.Premio.Mul(curr.Qtd),
				Discriminacao:    fmt.Sprintf("%s OPÇÕES %s SÉRIE %s VENCIMENTO %s", curr.Qtd, tipo, curr.Serie, curr.Vencimento.Format("02/01/2006")),
			})
		}
		prevYear = merge
		bens = append(bens, bem)
	}
	if len(bens)%2 == 0 && len(bens) > 0 {
		bens = bens[1:] // Remove a primeira posição se o número de anos for par.
	}
	return bens
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
	return s.rendimentos(func(o Operacao) bool { return tipoOprConfig[o.Tipo].IsRendimentoIsentoNaoTributavel })
}

func (s *State) RendimentosSujeitosTributacaoExclusiva() []RendimentosIsentosNaoTributaveis {
	return s.rendimentos(func(o Operacao) bool { return tipoOprConfig[o.Tipo].IsRendimentoSujeitoTributacaoExclusiva })
}

func (s *State) rendimentos(cond func(Operacao) bool) []RendimentosIsentosNaoTributaveis {

	var rs []RendimentosIsentosNaoTributaveis

	oprs := lo.Filter(s.Operacoes, func(o Operacao, _ int) bool { return cond(o) && o.Lucro.IsPositive() })

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

			vendaNoMes := oprs.VendasAcumulado()
			lucroNoMes := oprs.LucroAcumulado()

			if vendaNoMes.LessThanOrEqual(s.Config.AcaoLimiteIsento) && lucroNoMes.IsPositive() {
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
					vendaNoMes = p.VendasAcumulado()
					lucro = &lucroNoMes
				}
				for _, o := range p {
					if tipoOprConfig[o.Tipo].IsRendimentoTributavelApos20k && vendaNoMes.GreaterThan(s.Config.AcaoLimiteIsento) && o.Lucro.IsPositive() {
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
				rtm.IR = lucroAcumuladoPelosAnos.Mul(s.Config.AcaoSwingTradeIR)
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

func (s *State) Calculate(o *Operacao, a *Agregado) {
	o.Agg = *a
	switch o.Tipo {
	case COMPRA:
		o.CalcCompra()
	case VENDA:
		o.CalcVenda()
	case BONIFICACAO:
		o.CalcBonificacao(s.Config.AlterarPrecoMedioNaBonificacao)
	case DESDOBRAMENTO:
		o.CalcDesdobramento()
	case GRUPAMENTO:
		o.CalcGrupamento()
	case LEILAO_FRACAO:
		o.CalcLeilaoFracao()
	case DIVIDENDOS:
		o.CalcDividendos()
	case JSCP:
		o.CalcJSCP()
	case REND_TRIB:
		o.CalcRendTrib()
	case REDUCAO_CAPITAL:
		o.CalcReducaoCapital()
	case SUBSCRICAO_COMPRA:
		o.CalcSubscricaoCompra()
	case SUBSCRICAO_VENDA:
		o.CalcSubscricaoVenda()
	case SUBSCRICAO_COMPRA_EX:
		o.CalcSubscricaoExercicio()
	case PUT_VENDA:
		o.CalcVendaPut()
	case PUT_VENDA_EX:
		o.CalcVendaPutExercida()
	case PUT_VENDA_NE:
		o.CalcVendaPutNaoExercida()
	case PUT_COMPRA:
		o.CalcCompraPut()
	case PUT_COMPRA_EX:
		o.CalcCompraPutExercida()
	case PUT_COMPRA_NE:
		o.CalcCompraPutNaoExercida()
	case CALL_COMPRA:
		o.CalcCompraCall()
	case CALL_COMPRA_EX:
		o.CalcCompraCallExercida()
	case CALL_COMPRA_NE:
		o.CalcCompraCallNaoExercida()
	case CALL_VENDA:
		o.CalcVendaCall()
	case CALL_VENDA_EX:
		o.CalcVendaCallExercida()
	case CALL_VENDA_NE:
		o.CalcVendaCallNaoExercida()
	}
	*a = o.Agg
}

type Operacoes []Operacao

func (o Operacoes) FilterByTipo(tipo TipoOpr) Operacoes {
	return lo.Filter(o, func(v Operacao, _ int) bool { return v.Tipo == tipo })
}

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

func (o Operacoes) QtdAcumulado() decimal.Decimal {
	return lo.Reduce(o, func(agg decimal.Decimal, o Operacao, _ int) decimal.Decimal { return agg.Add(o.Qtd) }, decimal.Zero)
}

func (o Operacoes) VendasAcumulado() decimal.Decimal {
	return lo.Reduce(o, func(agg decimal.Decimal, o Operacao, _ int) decimal.Decimal {
		if o.Tipo == VENDA {
			return agg.Add(o.ValorTotal)
		}
		return agg
	}, decimal.Zero)
}

func (o Operacoes) LucroAcumulado() decimal.Decimal {
	return lo.Reduce(o, func(agg decimal.Decimal, o Operacao, _ int) decimal.Decimal {
		return agg.Add(o.Lucro)
	}, decimal.Zero)
}

func (o Operacoes) PartitionByYear() []Operacoes {
	return lo.PartitionBy(o, func(o Operacao) int { return o.Data.Year() })
}

func (o Operacoes) PartitionByAcaoOpcao() []Operacoes {
	return lo.PartitionBy(o, func(o Operacao) bool { return o.IsOpcao() })
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
	Ticker        string
	Tipo          TipoOpr
	Data          time.Time `json:",format:DateOnly"` // Opções: Data de Compra/Venda da opção ou Data de Vencimento ou Data de Encerramento.
	Qtd           decimal.Decimal
	ValorUnitario decimal.Decimal // Opções: Strike.
	Taxas         decimal.Decimal
	ValorTotal    decimal.Decimal
	ValorCompra   decimal.Decimal
	Lucro         decimal.Decimal // Lucro ou prejuízo da operação de Venda, Bonificação, Grupamento, Subscrição Compra, Redução de Capital, Opções.
	Fracao        decimal.Decimal // Parte fracionária resultante de Bonificação, Grupamento ou Desdobramento.
	Fator         decimal.Decimal // Fator de Bonificação, Grupamento ou Desdobramento e Redução de Capital.
	Agg           Agregado

	// Opções.
	Serie          string
	ValorExercicio decimal.Decimal // Valor da ação no dia do exercício da opção.
	Premio         decimal.Decimal
	Vencimento     time.Time `json:",format:DateOnly"`
}

func (o *Operacao) IsOpcao() bool {
	return o.Serie != ""
}

func (o *Operacao) CalcCompra() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Lucro = o.Premio.Mul(o.Qtd)
	o.Agg.Qtd = o.Agg.Qtd.Add(o.Qtd)
	o.Agg.ValorTotal = o.Agg.ValorTotal.Add(o.ValorTotal).Sub(o.Premio.Mul(o.Qtd))
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcVenda() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Sub(o.Taxas)
	o.ValorCompra = o.Agg.PrecoMedio.Mul(o.Qtd)
	o.Lucro = o.ValorTotal.Add(o.Premio.Mul(o.Qtd)).Sub(o.ValorCompra)
	o.Agg.Qtd = o.Agg.Qtd.Sub(o.Qtd)
	o.Agg.ValorTotal = o.Agg.ValorTotal.Sub(o.ValorCompra)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcBonificacao(alterarPrecoMedio bool) {
	if !o.Fator.IsZero() {
		o.Fracao = o.Agg.Qtd.Mul(o.Fator)
		o.Qtd = o.Fracao.Truncate(0)
		o.Fracao = o.Fracao.Sub(o.Qtd).Abs()
	}
	o.Agg.Qtd = o.Agg.Qtd.Add(o.Qtd)
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Agg.ValorTotal = o.ValorTotal.Add(o.Agg.ValorTotal)
	o.Lucro = o.ValorTotal
	if alterarPrecoMedio {
		o.Agg.CalcPrecoMedio()
	}
}

func (o *Operacao) CalcDesdobramento() {
	o.Agg.Qtd = o.Agg.Qtd.Mul(o.Fator)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcGrupamento() {
	qtdAcFra := o.Agg.Qtd.Div(o.Fator)
	qtdAcInt := qtdAcFra.Truncate(0)
	o.Fracao = qtdAcFra.Sub(qtdAcInt)
	o.Agg.Qtd = qtdAcInt
	o.Agg.ValorTotal = qtdAcInt.Mul(o.Agg.ValorTotal.Div(qtdAcFra))
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcLeilaoFracao() {
	o.Lucro = o.ValorUnitario.Mul(o.Fracao).Sub(o.Agg.PrecoMedio.Mul(o.Fracao))
}

func (o *Operacao) CalcDividendos() {
}

func (o *Operacao) CalcJSCP() {
}

func (o *Operacao) CalcRendTrib() {
}

func (o *Operacao) CalcSubscricaoCompra() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Agg.ValorTotal = o.Agg.ValorTotal.Add(o.ValorTotal)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcSubscricaoVenda() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Sub(o.Taxas)
	o.Lucro = o.ValorTotal
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcSubscricaoExercicio() {
	o.Agg.Qtd = o.Agg.Qtd.Add(o.Qtd)
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Agg.ValorTotal = o.Agg.ValorTotal.Add(o.ValorTotal)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcReducaoCapital() {
	if o.Fator.IsPositive() {
		redcap := o.Agg.ValorTotal.Mul(o.Fator)
		restituicao := o.Agg.Qtd.Mul(o.ValorUnitario)
		o.ValorTotal = restituicao
		o.ValorCompra = redcap
		o.Lucro = restituicao.Sub(redcap)
		o.Agg.ValorTotal = o.Agg.ValorTotal.Sub(redcap)
	} else {
		o.Agg.ValorTotal = o.Agg.ValorTotal.Sub(o.ValorTotal)
	}
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcVendaPut() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Lucro = o.Premio.Mul(o.Qtd)
}

func (o *Operacao) CalcVendaPutExercida() {
	o.CalcCompra()
}

func (o *Operacao) CalcVendaPutNaoExercida() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Lucro = o.Premio.Mul(o.Qtd)
}

func (o *Operacao) CalcCompraPut() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
}

func (o *Operacao) CalcCompraPutExercida() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.ValorUnitario.Sub(o.ValorExercicio).Sub(o.Premio).Mul(o.Qtd)
	o.Agg.Qtd = o.Agg.Qtd.Sub(o.Qtd)
	o.Agg.ValorTotal = o.Agg.ValorTotal.Sub(o.ValorTotal)
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcCompraPutNaoExercida() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
}

func (o *Operacao) CalcCompraCall() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
}

func (o *Operacao) CalcCompraCallExercida() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
	o.Agg.Qtd = o.Agg.Qtd.Add(o.Qtd)
	o.Agg.ValorTotal = o.Agg.ValorTotal.Add(o.ValorTotal).Add(o.Premio.Mul(o.Qtd))
	o.Agg.CalcPrecoMedio()
}

func (o *Operacao) CalcCompraCallNaoExercida() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
}

func (o *Operacao) CalcVendaCall() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Mul(o.Qtd)
}

func (o *Operacao) CalcVendaCallExercida() {
	o.CalcVenda()
}

func (o *Operacao) CalcVendaCallNaoExercida() {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Lucro = o.Premio.Mul(o.Qtd)
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
			line := bytes.TrimSpace(lines.Bytes())
			if len(line) == 0 || bytes.HasPrefix(line, []byte("//")) {
				continue
			}
			if !yield(line) {
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
