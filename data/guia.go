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

func (c *Carteira) OperacoesComunsDayTradeSwingTrade() []RendimentosTributaveis {
	ir := func(ano OperacaoAnual) decimal.Decimal {
		return ano.Cfg.SwingTradeIR
	}
	lucros := func(mes OperacaoMensal) (decimal.Decimal, decimal.Decimal) {
		swng := iterSwingTrades(mes.Iter())
		acoes := filterAcoes(swng)
		opcao := filterOpcao(swng)
		lucroNoMesAcoes := iterLucroTributavelOuAbativel(acoes).Sub(LucroIsentoVendaMes(acoes, mes.Cfg.LimiteVendaIsenta))
		lucroNoMesOpcao := iterLucroTributavelOuAbativel(opcao)
		return lucroNoMesAcoes, lucroNoMesOpcao
	}
	return c.operacoesComunsDayTrade(ir, lucros)
}

func (c *Carteira) OperacoesComunsDayTradeDayTrade() []RendimentosTributaveis {
	ir := func(ano OperacaoAnual) decimal.Decimal {
		return ano.Cfg.DayTradeIR
	}
	lucros := func(mes OperacaoMensal) (decimal.Decimal, decimal.Decimal) {
		dayt := iterDayTrades(mes.Iter())
		lucroNoMesAcoes := iterLucrosPejuizos(filterAcoes(dayt))
		lucroNoMesOpcao := iterLucrosPejuizos(filterOpcao(dayt))
		return lucroNoMesAcoes, lucroNoMesOpcao
	}
	return c.operacoesComunsDayTrade(ir, lucros)
}

func (c *Carteira) operacoesComunsDayTrade(ir func(OperacaoAnual) decimal.Decimal, lucros func(OperacaoMensal) (decimal.Decimal, decimal.Decimal)) []RendimentosTributaveis {

	var rts []RendimentosTributaveis

	lucroAcumuladoPelosAnos := decimal.Zero
	for _, ano := range c.Acoes {

		var rt RendimentosTributaveis
		rt.Ano = ano.Ano
		rt.IR = ir(ano)

		for _, mes := range ano.Ops {

			lucroNoMesAcoes, lucroNoMesOpcao := lucros(mes)

			if lucroNoMesAcoes.IsZero() && lucroNoMesOpcao.IsZero() {
				continue
			}

			lucroAcumuladoPelosAnos = lucroAcumuladoPelosAnos.Add(lucroNoMesAcoes).Add(lucroNoMesOpcao)
			rtm := RendimentoTributavelMensal{Mes: mes.Mes, Lucro: lucroNoMesAcoes, LucroOp: lucroNoMesOpcao, LucroAc: lucroAcumuladoPelosAnos}
			if lucroAcumuladoPelosAnos.IsPositive() {
				rtm.IR = lucroAcumuladoPelosAnos.Mul(ir(ano))
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
	Ticker     string
	Ano        int
	Meses      []RendimentoTributavelMensal
	IR         decimal.Decimal
	TotalAcoes decimal.Decimal
	TotalOpcao decimal.Decimal
	TotalAc    decimal.Decimal
	TotalIR    decimal.Decimal
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

		process := func(ops []OperacaoConsolidada) {
			ref := lo.FirstOrEmpty(ops)
			lucro := decimal.Zero
			if ref.Tfg.IsLimiteIsentoAplicavel() {
				lucro = LucroIsentoVendaAno(slices.Values(ops), ano.Cfg.LimiteVendaIsenta)
			} else {
				lucro = iterLucros(slices.Values(ops))
			}
			if lucro.IsPositive() {
				cfg := ref.Tfg.RendimentoIsentoNaoTributavel
				r.Rendimentos = append(r.Rendimentos, RendimentoIsentoNaoTributavel{
					Ticker: lo.Ternary(cfg.Agregar, "Total", ref.Ticker),
					Valor:  lucro,
					Codigo: cfg.Codigo,
					Descr:  cfg.Descr,
				})
			}
		}

		rnIseTri := it.Filter(iterSwingTrades(ano.Iter()), func(o OperacaoConsolidada) bool { return o.Tfg.IsRendimentoIsentoNaoTributavel() })
		agregado, desagreg := lo.FilterReject(slices.Collect(rnIseTri), func(o OperacaoConsolidada, _ int) bool { return o.Tfg.RendimentoIsentoNaoTributavel.Agregar })
		for _, ops := range lo.PartitionBy(desagreg, func(o OperacaoConsolidada) string { return o.Ticker + o.Tfg.RendimentoIsentoNaoTributavel.Codigo }) {
			process(ops)
		}
		for _, ops := range lo.PartitionBy(agregado, func(o OperacaoConsolidada) string { return o.Tfg.RendimentoIsentoNaoTributavel.Codigo }) {
			process(ops)
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
		f := it.Filter(ano.Iter(), func(o OperacaoConsolidada) bool { return o.Tfg.IsRendimentoSujeitoTributacaoExclusiva() })
		for _, ops := range it.PartitionBy(f, func(o OperacaoConsolidada) string {
			return o.Ticker + o.Tfg.RendimentoSujeitoTributacaoExclusiva.Codigo
		}) {
			ref := lo.FirstOrEmpty(ops)
			cfg := ref.Tfg.RendimentoSujeitoTributacaoExclusiva
			lucro := iterLucrosPejuizos(slices.Values(ops))
			if lucro.IsPositive() {
				r.Rendimentos = append(r.Rendimentos, RendimentoIsentoNaoTributavel{
					Ticker: ref.Ticker,
					Valor:  lucro,
					Codigo: cfg.Codigo + " ── " + cfg.Descr,
				})
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
	Totais      []RendimentoIsentoNaoTributavel
}

type RendimentoIsentoNaoTributavel struct {
	Ticker string
	Valor  decimal.Decimal
	Codigo string
	Descr  string
}
