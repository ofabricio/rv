package data

import (
	"fmt"
	"slices"
	"time"

	"github.com/samber/lo"
	"github.com/samber/lo/it"
	"github.com/shopspring/decimal"
)

func (c *Carteira) bensDireitos(bemOuDivida bool) []BemDireito {

	var bens []BemDireito

	prev := map[string]BensDireitoTicker{}
	for _, ano := range c.Acoes {

		bem := BemDireito{
			AnoAnterior: ano.Ano - 1,
			AnoCorrente: ano.Ano,
		}

		f := it.Filter(ano.Iter(), func(o OperacaoConsolidada) bool {
			return o.Tfg.IsBemDireito() == bemOuDivida || o.Tfg.IsDividaOnusReais() == !bemOuDivida
		})

		curr := map[string]BensDireitoTicker{}

		v := BemDireitoVisitor{Param: &c.Param, Curr: curr}
		for _, ano := range it.PartitionBy(f, func(o OperacaoConsolidada) string {
			return o.Ticker + lo.Ternary(bemOuDivida, o.Tfg.BensDireitos, o.Tfg.DividaOnusReais).ID()
		}) {
			v.Visit(ano)
		}

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

	if len(bens)%2 == 0 && len(bens) > 0 {
		bens = bens[1:] // Remove a primeira posição se o número de anos for par.
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

type BemDireitoVisitor struct {
	Curr  map[string]BensDireitoTicker
	Param *Param
	g     GrupoCodigo
}

func (v *BemDireitoVisitor) Visit(ano []OperacaoConsolidada) {
	o := lo.FirstOrEmpty(ano)
	switch {
	case o.Tfg.IsBemDireito():
		v.g = o.Tfg.BensDireitos
		v.VisitBemDireito(ano)
	case o.Tfg.IsDividaOnusReais():
		v.g = o.Tfg.DividaOnusReais
		v.VisitDividaOnusReais(ano)
	}
}

func (v *BemDireitoVisitor) VisitBemDireito(ano []OperacaoConsolidada) {
	if v.g.Grupo == GrupoParticipacaoSocietaria && v.g.Codigo == CodigoAcoes {
		v.VisitBemDireito0301(ano)
	} else if v.g.Grupo == GrupoAplicacaoInvestimento && v.g.Codigo == CodigoOpcoes {
		v.VisitBemDireito0404(ano, "COMPRADAS")
	} else if v.g.Grupo == GrupoOutros && v.g.Codigo == CodigoJSCPNaoPagos {
		v.VisitBemDireito9907(ano)
	}
}

func (v *BemDireitoVisitor) VisitDividaOnusReais(ano []OperacaoConsolidada) {
	if v.g.Codigo == CodigoOutrasDividasOnusReais {
		v.VisitDividaOnusReais0016(ano)
	}
}

func (v *BemDireitoVisitor) VisitBemDireito0301(ano []OperacaoConsolidada) {
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

func (v *BemDireitoVisitor) VisitBemDireito9907(ano []OperacaoConsolidada) {
	op := lo.LastOrEmpty(ano)
	id := op.Ticker + ".jscpnp"
	v.Curr[id] = BensDireitoTicker{
		Grupo:            v.g.Grupo,
		Codigo:           v.g.Codigo,
		Ticker:           op.Ticker,
		SituacaoAnterior: decimal.Zero,
		SituacaoCorrente: iterLucros(slices.Values(ano)),
		Discriminacao:    fmt.Sprintf("%s RENDIMENTOS CREDITADOS MAS NÃO PAGOS", op.Ticker),
	}
}

func (v *BemDireitoVisitor) VisitDividaOnusReais0016(ano []OperacaoConsolidada) {
	v.VisitBemDireito0404(ano, "VENDIDAS")
}

func (v *BemDireitoVisitor) VisitBemDireito0404(ano []OperacaoConsolidada, label string) {
	f := it.Filter(slices.Values(ano), func(o OperacaoConsolidada) bool { return o.Vencimento.Year() > o.Data.Year() })
	for _, ano := range it.PartitionBy(f, func(o OperacaoConsolidada) string { return o.Serie }) {
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
