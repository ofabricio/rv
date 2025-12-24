package data

import (
	"cmp"
	"io"
	"iter"
	"maps"
	"slices"
	"time"

	"github.com/samber/lo"
	"github.com/samber/lo/it"
	"github.com/shopspring/decimal"
)

type Carteira struct {
	Acoes Acoes
	Param Param
}

type Acoes []OperacaoAnual

type OperacaoAnual struct {
	Ano int
	Ops []OperacaoMensal
	Cfg Config
}

type OperacaoMensal struct {
	Mes int
	Ops []OperacaoConsolidada
	Cfg Config
}

func (c *Carteira) Load(file string) {
	ops, err := LoadOperacoes(file)
	if err != nil {
		panic(err)
	}
	c.Consolidar(ops)
}

func (c *Carteira) Print(w io.Writer) {
	p := TablePrinter{c}
	p.PrintOperacoesComAcoes(w)
	p.PrintBensDireitos(w)
	p.PrintDividaOnusReais(w)
	p.PrintRendimentosIsentosNaoTributaveis(w)
	p.PrintRendimentosSujeitosTributacaoExclusiva(w)
	p.PrintOperacoesComunsDayTrade(w)
}

func (c *Carteira) Consolidar(all []OperacaoOuConfig) {
	c.Param = DefaultParam
	var con Consolidador
	var cfg Config = DefaultConfig
	for _, year := range lo.PartitionBy(all, func(v OperacaoOuConfig) int { return v.Data.Year() }) {
		oprs := lo.FilterMap(year, func(o OperacaoOuConfig, _ int) (Operacao, bool) { return o.Opr, !o.IsCfg() })
		cfgs := lo.FilterMap(year, func(o OperacaoOuConfig, _ int) (ConfigOpcional, bool) { return o.Cfg, o.IsCfg() })
		cfg = lo.Reduce(cfgs, func(agg Config, c ConfigOpcional, _ int) Config { return c.Merge(agg) }, cfg)
		if len(oprs) == 0 {
			continue
		}
		con.Cfg = cfg
		var om []OperacaoMensal
		for _, month := range lo.PartitionBy(oprs, func(o Operacao) time.Month { return o.Data.Month() }) {
			var opc []OperacaoConsolidada
			for _, o := range month {
				opc = append(opc, con.Consolidar(&o))
			}
			om = append(om, OperacaoMensal{
				Mes: int(month[0].Data.Month()),
				Ops: opc,
				Cfg: cfg,
			})
		}
		c.Acoes = append(c.Acoes, OperacaoAnual{
			Ano: oprs[0].Data.Year(),
			Ops: om,
			Cfg: cfg,
		})
	}
}

type OperacaoOuConfig struct {
	Data time.Time
	Opr  Operacao
	Cfg  ConfigOpcional
}

func (o *OperacaoOuConfig) IsCfg() bool {
	return !o.Cfg.Data.IsZero()
}

func (a Acoes) Iter() iter.Seq[OperacaoConsolidada] {
	return func(yield func(OperacaoConsolidada) bool) {
		for _, ano := range a {
			for o := range ano.Iter() {
				if !yield(o) {
					return
				}
			}
		}
	}
}

func (a Acoes) IterI() iter.Seq2[int, OperacaoConsolidada] {
	return func(yield func(int, OperacaoConsolidada) bool) {
		i := 0
		for o := range a.Iter() {
			if !yield(i, o) {
				return
			}
			i++
		}
	}
}

func (a *OperacaoAnual) Iter() iter.Seq[OperacaoConsolidada] {
	return func(yield func(OperacaoConsolidada) bool) {
		for _, mes := range a.Ops {
			for op := range mes.Iter() {
				if !yield(op) {
					return
				}
			}
		}
	}
}

func (a *OperacaoMensal) Iter() iter.Seq[OperacaoConsolidada] {
	return func(yield func(OperacaoConsolidada) bool) {
		for _, op := range a.Ops {
			if !yield(op) {
				return
			}
		}
	}
}

func (a *OperacaoAnual) LucrosIsentosVenda() decimal.Decimal {
	total := decimal.Zero
	for _, mes := range a.Ops {
		vendaNoMes := mes.TotalVendas()
		lucroNoMes := mes.LucroVendas()
		if vendaNoMes.LessThanOrEqual(mes.Cfg.LimiteVendaIsenta) && lucroNoMes.IsPositive() {
			total = total.Add(lucroNoMes)
		}
	}
	return total
}

func (m *OperacaoMensal) LucroVendas() decimal.Decimal {
	return iterLucroVendas(m.Iter())
}

func (m *OperacaoMensal) LucroTributavelOuAbativelAcoes() decimal.Decimal {
	return iterLucroTributavelOuAbativelAcao(m.Iter())
}

func (m *OperacaoMensal) LucroTributavelOuAbativelOpcao() decimal.Decimal {
	return iterLucroTributavelOuAbativelOpcao(m.Iter())
}

func (m *OperacaoMensal) TotalVendas() decimal.Decimal {
	return iterTotalVendas(m.Iter())
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

func iterValorTotal(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal { return agg.Add(o.ValorTotal) }, decimal.Zero)
}

func iterTotalVendas(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return iterValorTotal(iterFilterByTipo(ops, VENDA))
}

func iterLucroVendas(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return iterLucros(iterFilterByTipo(ops, VENDA))
}

func iterLucros(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal { return agg.Add(o.Lucro) }, decimal.Zero)
}

func iterLucroTributavelOuAbativelAcao(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal {
		if o.IsAcao() {
			if o.Tfg.LucroTributavel && o.Lucro.IsPositive() {
				return agg.Add(o.Lucro)
			}
			if o.Tfg.PrejuizoAbativel && o.Lucro.IsNegative() {
				return agg.Add(o.Lucro)
			}
		}
		return agg
	}, decimal.Zero)
}

func iterLucroTributavelOuAbativelOpcao(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal {
		if o.IsOpcao() {
			if o.Tfg.LucroTributavel && o.Lucro.IsPositive() {
				return agg.Add(o.Lucro)
			}
			if o.Tfg.PrejuizoAbativel && o.Lucro.IsNegative() {
				return agg.Add(o.Lucro)
			}
		}
		return agg
	}, decimal.Zero)
}

func iterQtd(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal { return agg.Add(o.Qtd) }, decimal.Zero)
}

func iterFilterByTipo(ops iter.Seq[OperacaoConsolidada], t Tipo) iter.Seq[OperacaoConsolidada] {
	return it.Filter(ops, func(o OperacaoConsolidada) bool { return o.Tipo == t })
}
