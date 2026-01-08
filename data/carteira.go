package data

import (
	"cmp"
	"fmt"
	"io"
	"iter"
	"maps"
	"os"
	"slices"
	"time"

	"github.com/samber/lo"
	"github.com/samber/lo/it"
	"github.com/shopspring/decimal"
)

func NewCarteira() *Carteira {
	return &Carteira{
		Param: DefaultParam,
	}
}

type Carteira struct {
	Acoes Acoes
	Param Param
}

type Acoes []OperacaoAnual

type OperacaoAnual struct {
	Ano int
	Cfg Config
	Ops []OperacaoMensal
}

type OperacaoMensal struct {
	Mes int
	Cfg Config
	Ops []OperacaoConsolidada
}

func (c *Carteira) Load(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	ops, err := ReadOperacoes(f)
	if err != nil {
		f.Close()
		panic(err)
	}
	c.Consolidar(ops)
}

func (c *Carteira) Read(r io.Reader) {
	ops, err := ReadOperacoes(r)
	if err != nil {
		panic(err)
	}
	c.Consolidar(ops)
}

func (c *Carteira) Consolidar(all []OperacaoDesconsolidada) {
	var con Consolidador
	for _, year := range lo.PartitionBy(all, func(v OperacaoDesconsolidada) int { return v.Opr.Data.Year() }) {
		opa := make([]OperacaoMensal, 0, 12)
		for _, month := range lo.PartitionBy(year, func(o OperacaoDesconsolidada) time.Month { return o.Opr.Data.Month() }) {
			opm := make([]OperacaoConsolidada, 0, len(month))
			for _, o := range month {
				opm = append(opm, con.Consolidar(o))
			}
			opa = append(opa, OperacaoMensal{
				Mes: int(month[0].Opr.Data.Month()),
				Cfg: year[0].Cfg,
				Ops: opm,
			})
		}
		c.Acoes = append(c.Acoes, OperacaoAnual{
			Ano: year[0].Opr.Data.Year(),
			Cfg: year[0].Cfg,
			Ops: opa,
		})
	}
}

func (c *Carteira) Print(frmt string, w io.Writer) {
	var p interface {
		PrintOperacoesComAcoes(io.Writer)
		PrintBensDireitos(io.Writer)
		PrintDividaOnusReais(io.Writer)
		PrintRendimentosIsentosNaoTributaveis(io.Writer)
		PrintRendimentosSujeitosTributacaoExclusiva(io.Writer)
		PrintOperacoesComunsDayTrade(io.Writer)
	}
	switch frmt {
	case "table":
		p = &PrinterTable{c}
	case "csv":
		p = &PrinterCSV{c}
	}
	p.PrintOperacoesComAcoes(w)
	p.PrintBensDireitos(w)
	p.PrintDividaOnusReais(w)
	p.PrintRendimentosIsentosNaoTributaveis(w)
	p.PrintRendimentosSujeitosTributacaoExclusiva(w)
	p.PrintOperacoesComunsDayTrade(w)
}

type OperacaoDesconsolidada struct {
	Opr Operacao
	Cfg Config
	Tfg Tonfig
}

type OperacaoConsolidada struct {
	Operacao
	Agg Agregado
	Tfg Tonfig
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

func LucroIsentoVendaAno(ops iter.Seq[OperacaoConsolidada], limit decimal.Decimal) decimal.Decimal {
	lucros := decimal.Zero
	for _, mes := range it.PartitionBy(ops, func(o OperacaoConsolidada) time.Month { return o.Data.Month() }) {
		lucros = lucros.Add(LucroIsentoVendaMes(slices.Values(mes), limit))
	}
	return lucros
}

func LucroIsentoVendaMes(ops iter.Seq[OperacaoConsolidada], limit decimal.Decimal) decimal.Decimal {
	if iterTotalVendas(ops).LessThanOrEqual(limit) {
		return iterLucroVendas(ops)
	}
	return decimal.Zero
}

func IRRFIsentoVendaMes(ops iter.Seq[OperacaoConsolidada], limit decimal.Decimal) decimal.Decimal {
	if iterTotalVendas(ops).LessThanOrEqual(limit) {
		return iterIRRF(ops)
	}
	return decimal.Zero
}

func iterSwingTrades(ops iter.Seq[OperacaoConsolidada]) iter.Seq[OperacaoConsolidada] {
	return it.Filter(ops, func(o OperacaoConsolidada) bool { return !o.DayTrade })
}

func iterDayTrades(ops iter.Seq[OperacaoConsolidada]) iter.Seq[OperacaoConsolidada] {
	return it.Filter(ops, func(o OperacaoConsolidada) bool { return o.DayTrade })
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

func iterTotalVendas(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return iterValorTotalBruto(iterFilterByTipo(ops, VENDA))
}

func iterLucroVendas(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return iterLucros(iterFilterByTipo(ops, VENDA))
}

func iterLucrosPejuizos(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal { return agg.Add(o.Lucro) }, decimal.Zero)
}

func iterLucros(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal {
		if o.Lucro.IsPositive() {
			return agg.Add(o.Lucro)
		}
		return agg
	}, decimal.Zero)
}

func iterLucroTributavelOuAbativel(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal {
		if o.Tfg.LucroTributavel && o.Lucro.IsPositive() {
			return agg.Add(o.Lucro)
		}
		if o.Tfg.PrejuizoAbativel && o.Lucro.IsNegative() {
			return agg.Add(o.Lucro)
		}
		return agg
	}, decimal.Zero)
}

func iterIRRF(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal { return agg.Add(o.IRRF) }, decimal.Zero)
}

func iterValorTotalBruto(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal { return agg.Add(o.ValorTotalBruto()) }, decimal.Zero)
}

func iterQtd(ops iter.Seq[OperacaoConsolidada]) decimal.Decimal {
	return it.Reduce(ops, func(agg decimal.Decimal, o OperacaoConsolidada) decimal.Decimal { return agg.Add(o.QtdInt()) }, decimal.Zero)
}

func iterFilterByTipo(ops iter.Seq[OperacaoConsolidada], t Tipo) iter.Seq[OperacaoConsolidada] {
	return it.Filter(ops, func(o OperacaoConsolidada) bool { return o.Tipo == t })
}

func filterAcoes(ops iter.Seq[OperacaoConsolidada]) iter.Seq[OperacaoConsolidada] {
	return it.Filter(ops, func(o OperacaoConsolidada) bool { return o.IsAcao() })
}

func filterOpcao(ops iter.Seq[OperacaoConsolidada]) iter.Seq[OperacaoConsolidada] {
	return it.Filter(ops, func(o OperacaoConsolidada) bool { return o.IsOpcao() })
}
