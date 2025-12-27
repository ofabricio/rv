package data

import (
	"github.com/shopspring/decimal"
)

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

type OperacaoConsolidada struct {
	Operacao
	Agg Agregado
	Tfg Tonfig
}

type Consolidador struct {
	Map map[string]Agregado
	Agg Agregado
	Cfg Config
}

func (c *Consolidador) Consolidar(o *Operacao) OperacaoConsolidada {
	if c.Map == nil {
		c.Map = make(map[string]Agregado)
	}
	c.Agg = c.Map[o.Ticker]
	o.Accept(c)
	c.Map[o.Ticker] = c.Agg
	return OperacaoConsolidada{*o, c.Agg, c.Cfg.Tonfig[o.Tipo]}
}

func (c *Consolidador) VisitCompra(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Lucro = o.Premio.Mul(o.Qtd)
	c.Agg.Qtd = c.Agg.Qtd.Add(o.Qtd)
	c.Agg.ValorTotal = c.Agg.ValorTotal.Add(o.ValorTotal).Sub(o.Premio.Mul(o.Qtd))
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitVenda(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Sub(o.Taxas)
	o.ValorCompra = c.Agg.PrecoMedio.Mul(o.Qtd)
	o.Lucro = o.ValorTotal.Add(o.Premio.Mul(o.Qtd)).Sub(o.ValorCompra)
	c.Agg.Qtd = c.Agg.Qtd.Sub(o.Qtd)
	c.Agg.ValorTotal = c.Agg.ValorTotal.Sub(o.ValorCompra)
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitBonificacao(o *Operacao) {
	if !o.Fator.IsZero() {
		o.Fracao = c.Agg.Qtd.Mul(o.Fator)
		o.Qtd = o.Fracao.Truncate(0)
		o.Fracao = o.Fracao.Sub(o.Qtd).Abs()
	} else {
		o.Fracao = o.Qtd.Sub(o.Qtd.Truncate(0)).Abs()
		o.Qtd = o.Qtd.Truncate(0)
	}
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.ValorTotal
	c.Agg.Qtd = c.Agg.Qtd.Add(o.Qtd)
	c.Agg.ValorTotal = o.ValorTotal.Add(c.Agg.ValorTotal)
	if c.Cfg.AlterarPrecoMedioNaBonificacao {
		c.Agg.CalcPrecoMedio()
	}
}

func (c *Consolidador) VisitDesdobramento(o *Operacao) {
	c.Agg.Qtd = c.Agg.Qtd.Mul(o.Fator)
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitGrupamento(o *Operacao) {
	qtdAcFra := c.Agg.Qtd.Div(o.Fator)
	qtdAcInt := qtdAcFra.Truncate(0)
	o.Fracao = qtdAcFra.Sub(qtdAcInt)
	c.Agg.Qtd = qtdAcInt
	c.Agg.ValorTotal = qtdAcInt.Mul(c.Agg.ValorTotal.Div(qtdAcFra))
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitLeilaoFracao(o *Operacao) {
	o.Lucro = o.ValorUnitario.Mul(o.Fracao).Sub(c.Agg.PrecoMedio.Mul(o.Fracao))
}

func (*Consolidador) VisitDividendos(*Operacao) {}

func (*Consolidador) VisitReembolso(*Operacao) {}

func (*Consolidador) VisitJSCP(*Operacao) {}

func (*Consolidador) VisitJSCPNaoPago(*Operacao) {}

func (*Consolidador) VisitAluguel(*Operacao) {}

func (*Consolidador) VisitRendTrib(*Operacao) {}

func (c *Consolidador) VisitSubscricaoCompra(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	c.Agg.ValorTotal = c.Agg.ValorTotal.Add(o.ValorTotal)
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitSubscricaoVenda(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Sub(o.Taxas)
	o.Lucro = o.ValorTotal
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitSubscricaoExercicio(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	c.Agg.Qtd = c.Agg.Qtd.Add(o.Qtd)
	c.Agg.ValorTotal = c.Agg.ValorTotal.Add(o.ValorTotal)
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitReducaoCapital(o *Operacao) {
	if o.Fator.IsPositive() {
		redcap := c.Agg.ValorTotal.Mul(o.Fator)
		restituicao := c.Agg.Qtd.Mul(o.ValorUnitario)
		o.ValorTotal = restituicao
		o.ValorCompra = redcap
		o.Lucro = restituicao.Sub(redcap)
		c.Agg.ValorTotal = c.Agg.ValorTotal.Sub(redcap)
	} else {
		c.Agg.ValorTotal = c.Agg.ValorTotal.Sub(o.ValorTotal)
	}
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitVendaPut(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Lucro = o.Premio.Mul(o.Qtd)
}

func (c *Consolidador) VisitVendaPutExercida(o *Operacao) {
	c.VisitCompra(o)
}

func (c *Consolidador) VisitVendaPutNaoExercida(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Lucro = o.Premio.Mul(o.Qtd)
}

func (c *Consolidador) VisitCompraPut(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
}

func (c *Consolidador) VisitCompraPutExercida(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.ValorUnitario.Sub(o.ValorExercicio).Sub(o.Premio).Mul(o.Qtd)
	c.Agg.Qtd = c.Agg.Qtd.Sub(o.Qtd)
	c.Agg.ValorTotal = c.Agg.ValorTotal.Sub(o.ValorTotal)
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitCompraPutNaoExercida(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
}

func (c *Consolidador) VisitCompraCall(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
}

func (c *Consolidador) VisitCompraCallExercida(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
	c.Agg.Qtd = c.Agg.Qtd.Add(o.Qtd)
	c.Agg.ValorTotal = c.Agg.ValorTotal.Add(o.ValorTotal).Add(o.Premio.Mul(o.Qtd))
	c.Agg.CalcPrecoMedio()
}

func (c *Consolidador) VisitCompraCallNaoExercida(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Neg().Mul(o.Qtd)
}

func (c *Consolidador) VisitVendaCall(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.Lucro = o.Premio.Mul(o.Qtd)
}

func (c *Consolidador) VisitVendaCallExercida(o *Operacao) {
	c.VisitVenda(o)
}

func (c *Consolidador) VisitVendaCallNaoExercida(o *Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.Lucro = o.Premio.Mul(o.Qtd)
}
