package data

import (
	"slices"

	"github.com/samber/lo"
	"github.com/samber/lo/it"
	"github.com/shopspring/decimal"
)

func (c *Carteira) BensDireitos() []BemDireito {
	return c.bensDireitos(&BemDireitoStrategy{Param: &c.Param})
}

func (c *Carteira) DividaOnusReais() []BemDireito {
	return c.bensDireitos(&DividaOnusReaisStrategy{Param: &c.Param})
}

func (c *Carteira) OperacoesComunsDayTrade() []RendimentosTributaveis {

	var rts []RendimentosTributaveis

	lucroAcumuladoPelosAnos := decimal.Zero
	for _, ano := range c.Acoes {

		var rt RendimentosTributaveis
		rt.Ano = ano.Ano
		rt.SwingTradeIR = ano.Cfg.SwingTradeIR

		for _, mes := range ano.Ops {

			lucroNoMesAcoes := mes.LucroTributavelOuAbativelAcoes()
			lucroNoMesOpcao := mes.LucroTributavelOuAbativelOpcao()

			if mes.TotalVendas().GreaterThan(ano.Cfg.LimiteVendaIsenta) {
				lucroNoMesAcoes = lucroNoMesAcoes.Add(mes.LucroVendas())
			}

			if lucroNoMesAcoes.IsZero() && lucroNoMesOpcao.IsZero() {
				continue
			}

			lucroAcumuladoPelosAnos = lucroAcumuladoPelosAnos.Add(lucroNoMesAcoes).Add(lucroNoMesOpcao)
			rtm := RendimentoTributavelMensal{Mes: mes.Mes, Lucro: lucroNoMesAcoes, LucroOp: lucroNoMesOpcao, LucroAc: lucroAcumuladoPelosAnos}
			if lucroAcumuladoPelosAnos.IsPositive() {
				rtm.IR = lucroAcumuladoPelosAnos.Mul(ano.Cfg.SwingTradeIR)
				lucroAcumuladoPelosAnos = decimal.Zero
			}
			rt.Meses = append(rt.Meses, rtm)
			rt.TotalAcoes = rt.TotalAcoes.Add(lucroNoMesAcoes)
			rt.TotalOpcao = rt.TotalOpcao.Add(lucroNoMesOpcao)
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
	Ticker       string
	Ano          int
	Meses        []RendimentoTributavelMensal
	SwingTradeIR decimal.Decimal
	TotalAcoes   decimal.Decimal
	TotalOpcao   decimal.Decimal
	TotalAc      decimal.Decimal
	TotalIR      decimal.Decimal
}

type RendimentoTributavelMensal struct {
	Mes     int
	Lucro   decimal.Decimal
	LucroOp decimal.Decimal
	LucroAc decimal.Decimal
	IR      decimal.Decimal
}

type RendimentosIsentosAteLimite struct {
	Ticker string
	Ano    int
	Meses  []RendimentoIsentosAteLimiteMensal
	Total  decimal.Decimal
}

type RendimentoIsentosAteLimiteMensal struct {
	Mes   int
	Valor decimal.Decimal
}

func (c *Carteira) RendimentosIsentosNaoTributaveis() []RendimentosIsentosNaoTributaveis {

	var rs []RendimentosIsentosNaoTributaveis

	for _, ano := range c.Acoes {
		var r RendimentosIsentosNaoTributaveis
		r.Ano = ano.Ano
		for _, ops := range it.PartitionBy(ano.Iter(), func(o OperacaoConsolidada) string { return o.Ticker }) {
			for _, ops := range lo.PartitionBy(ops, func(o OperacaoConsolidada) string { return o.Tfg.RendimentoIsentoNaoTributavel.Codigo }) {
				ref := lo.FirstOrEmpty(ops)
				cfg := ref.Tfg.RendimentoIsentoNaoTributavel
				if cfg.Codigo == "" {
					continue
				}
				lucro := decimal.Zero
				if cfg.Codigo == CodigoIsencaoAte20k {
					lucro = ano.LucrosIsentosVenda()
				} else {
					lucro = iterLucros(slices.Values(ops))
				}
				if lucro.IsPositive() {
					r.Rendimentos = append(r.Rendimentos, RendimentoIsentoNaoTributavel{
						Ticker: ref.Ticker,
						Valor:  lucro,
						Codigo: cfg.Codigo + " ── " + cfg.Descr,
					})
				}
			}
		}
		if len(r.Rendimentos) > 0 {
			rs = append(rs, r)
		}
	}

	return rs
}

func (c *Carteira) RendimentosSujeitosTributacaoExclusiva() []RendimentosIsentosNaoTributaveis {

	var rs []RendimentosIsentosNaoTributaveis

	for _, ano := range c.Acoes {
		var r RendimentosIsentosNaoTributaveis
		r.Ano = ano.Ano
		for _, ops := range it.PartitionBy(ano.Iter(), func(o OperacaoConsolidada) string { return o.Ticker }) {
			for _, ops := range lo.PartitionBy(ops, func(o OperacaoConsolidada) string { return o.Tfg.RendimentoSujeitoTributacaoExclusiva.Codigo }) {
				ref := lo.FirstOrEmpty(ops)
				cfg := ref.Tfg.RendimentoSujeitoTributacaoExclusiva
				if cfg.Codigo == "" {
					continue
				}
				lucro := iterLucros(slices.Values(ops))
				if lucro.IsPositive() {
					r.Rendimentos = append(r.Rendimentos, RendimentoIsentoNaoTributavel{
						Ticker: ref.Ticker,
						Valor:  lucro,
						Codigo: cfg.Codigo + " ── " + cfg.Descr,
					})
				}
			}
		}
		if len(r.Rendimentos) > 0 {
			rs = append(rs, r)
		}
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
