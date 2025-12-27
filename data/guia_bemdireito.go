package data

import (
	"fmt"
	"iter"
	"slices"
	"time"

	"github.com/samber/lo"
	"github.com/samber/lo/it"
	"github.com/shopspring/decimal"
)

func (c *Carteira) bensDireitos(v BensDireitosStrategy) []BemDireito {

	var bens []BemDireito

	prev := map[string]BensDireitoTicker{}
	for _, ano := range c.Acoes {

		bem := BemDireito{
			AnoAnterior: ano.Ano - 1,
			AnoCorrente: ano.Ano,
		}

		curr := v.SituacaoAnoCorrente(ano.Iter())

		// Decide se adiciona os tickers do ano anterior no ano corrente.
		for ticker, p := range prev {
			if c, found := curr[ticker]; found {
				// Temos o ticker tanto no ano anterior como no ano corrente,
				// então apenas atualiza a situação dele no ano corrente.
				c.SituacaoAnterior = p.SituacaoCorrente
				curr[ticker] = c
			} else if p.SituacaoCorrente.IsZero() || p.Vencimento.Year() == bem.AnoCorrente {
				// Situação zerada ou vencida no ano anterior não adiciona o ticker no ano corrente.
				delete(prev, ticker)
			} else {
				// Se não encontrou o ticker do ano anterior no ano corrente, e a situação
				// anterior dele não está nem zerada nem vencida, adiciona no ano corrente.
				p.SituacaoAnterior = p.SituacaoCorrente
				curr[ticker] = p
			}
		}

		prev = lo.Assign(prev, curr)
		for _, v := range sortKeysIter(prev) {
			bem.Tickers = append(bem.Tickers, v)
		}

		if len(bem.Tickers) > 0 {
			bens = append(bens, bem)
		}
	}

	return bens
}

type BemDireito struct {
	AnoAnterior int
	AnoCorrente int
	Tickers     []BensDireitoTicker
}

type BensDireitoTicker struct {
	Grupo, Codigo    string
	Ticker           string
	SituacaoAnterior decimal.Decimal
	SituacaoCorrente decimal.Decimal
	Discriminacao    string
	Vencimento       time.Time
}

type BensDireitosStrategy interface {
	SituacaoAnoCorrente(iter.Seq[OperacaoConsolidada]) map[string]BensDireitoTicker
}

type BemDireitoStrategy struct {
	Curr  map[string]BensDireitoTicker
	Param *Param
	g     GrupoCodigo
}

func (v *BemDireitoStrategy) SituacaoAnoCorrente(ano iter.Seq[OperacaoConsolidada]) map[string]BensDireitoTicker {
	v.Curr = map[string]BensDireitoTicker{}
	Filter := func(o OperacaoConsolidada) bool { return o.Tfg.IsBemDireito() }
	Partition := func(o OperacaoConsolidada) string { return o.Ticker + o.Tfg.BensDireitos.ID() }
	for _, ano := range it.PartitionBy(it.Filter(ano, Filter), Partition) {
		v.VisitBemDireito(ano)
	}
	return v.Curr
}

func (v *BemDireitoStrategy) VisitBemDireito(ano []OperacaoConsolidada) {
	v.g = lo.FirstOrEmpty(ano).Tfg.BensDireitos
	if v.g.Grupo == GrupoParticipacaoSocietaria && v.g.Codigo == CodigoAcoes {
		v.VisitBemDireito0301(ano)
	} else if v.g.Grupo == GrupoAplicacaoInvestimento && v.g.Codigo == CodigoOpcoes {
		v.VisitBemDireito0404(ano, "COMPRADAS")
	} else if v.g.Grupo == GrupoOutros && v.g.Codigo == CodigoJSCPNaoPagos {
		v.VisitBemDireito9907(ano)
	}
}

func (v *BemDireitoStrategy) VisitBemDireito0301(ano []OperacaoConsolidada) {
	op := lo.LastOrEmpty(ano)
	bn := iterQtd(iterFilterByTipo(slices.Values(ano), BONIFICACAO))
	id := op.Ticker + ".acao"
	v.Curr[id] = BensDireitoTicker{
		Grupo:            v.g.Grupo,
		Codigo:           v.g.Codigo,
		Ticker:           op.Ticker,
		SituacaoAnterior: decimal.Zero,
		SituacaoCorrente: op.Agg.ValorTotal,
		Discriminacao: fmt.Sprintf("%s AÇÕES %s COM PREÇO MÉDIO DE R$ %s", op.Agg.Qtd, op.Ticker, v.Param.FormatDecimal(op.Agg.PrecoMedio)) +
			lo.Ternary(bn.IsPositive(), fmt.Sprintf(" ONDE %s AÇÕES SÃO PROVENIENTES DE BONIFICAÇÃO", bn), ""),
	}
}

func (v *BemDireitoStrategy) VisitBemDireito9907(ano []OperacaoConsolidada) {
	op := lo.LastOrEmpty(ano)
	id := op.Ticker + ".jscpnp"
	v.Curr[id] = BensDireitoTicker{
		Vencimento:       op.Vencimento,
		Grupo:            v.g.Grupo,
		Codigo:           v.g.Codigo,
		Ticker:           op.Ticker,
		SituacaoAnterior: decimal.Zero,
		SituacaoCorrente: iterLucros(slices.Values(ano)),
		Discriminacao:    fmt.Sprintf("%s RENDIMENTOS CREDITADOS MAS NÃO PAGOS", op.Ticker),
	}
}

func (v *BemDireitoStrategy) VisitBemDireito0404(ano []OperacaoConsolidada, label string) {
	f := lo.Filter(ano, func(o OperacaoConsolidada, _ int) bool { return o.Vencimento.Year() > o.Data.Year() })
	for _, ano := range lo.PartitionBy(f, func(o OperacaoConsolidada) string { return o.Serie }) {
		qt := iterQtd(slices.Values(ano))
		tt := lo.Reduce(ano, func(agg decimal.Decimal, o OperacaoConsolidada, _ int) decimal.Decimal {
			return agg.Add(o.Premio.Mul(o.Qtd))
		}, decimal.Zero)
		op := lo.LastOrEmpty(ano)
		id := op.Serie + ".opcao"
		v.Curr[id] = BensDireitoTicker{
			Vencimento:       op.Vencimento,
			Grupo:            v.g.Grupo,
			Codigo:           v.g.Codigo,
			Ticker:           op.Serie,
			SituacaoAnterior: decimal.Zero,
			SituacaoCorrente: tt,
			Discriminacao:    fmt.Sprintf("%s OPÇÕES %s SÉRIE %s VENCIMENTO %s", qt, label, op.Serie, op.Vencimento.Format("02/01/2006")),
		}
	}
}

type DividaOnusReaisStrategy struct {
	Curr  map[string]BensDireitoTicker
	Param *Param
	g     GrupoCodigo
}

func (v *DividaOnusReaisStrategy) SituacaoAnoCorrente(ano iter.Seq[OperacaoConsolidada]) map[string]BensDireitoTicker {
	v.Curr = map[string]BensDireitoTicker{}
	Filter := func(o OperacaoConsolidada) bool { return o.Tfg.IsDividaOnusReais() }
	Partition := func(o OperacaoConsolidada) string { return o.Ticker + o.Tfg.DividaOnusReais.ID() }
	for _, ano := range it.PartitionBy(it.Filter(ano, Filter), Partition) {
		v.VisitDividaOnusReais(ano)
	}
	return v.Curr
}

func (v *DividaOnusReaisStrategy) VisitDividaOnusReais(ano []OperacaoConsolidada) {
	v.g = lo.FirstOrEmpty(ano).Tfg.DividaOnusReais
	if v.g.Codigo == CodigoOutrasDividasOnusReais {
		v.VisitDividaOnusReais0016(ano)
	}
}

func (v *DividaOnusReaisStrategy) VisitDividaOnusReais0016(ano []OperacaoConsolidada) {
	vv := BemDireitoStrategy{Curr: v.Curr, Param: v.Param, g: v.g}
	vv.VisitBemDireito0404(ano, "VENDIDAS")
}
