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

	AcaoEmolumentosB3 decimal.Decimal
	AcaoLiquidacaoB3  decimal.Decimal
	AcaoSwingTradeIR  decimal.Decimal
	AcaoDayTradeIR    decimal.Decimal
	AcaoLimiteIsento  decimal.Decimal
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

	// Coleta todas as operações que entram na guia BENS E DIREITOS.
	oprs := lo.Filter(s.Operacoes, func(o Operacao, _ int) bool { return tipoOprConfig[o.Tipo].IsBemDireito })

	// Encontra o intervalo de anos das operações.
	nowYear := time.Now().Local().Year()
	endYear := lo.LastOrEmpty(oprs).Data.Year()
	endYear = lo.Ternary(len(oprs) == 0, nowYear, endYear)
	iniYear := lo.FirstOrEmpty(oprs).Data.Year()
	iniYear = iniYear - (endYear-iniYear+1)%2
	iniYear = lo.Ternary(len(oprs) == 0, endYear-1, iniYear)

	// Cria uma lista de pares de anos cujos anos intermediários são repetidos.
	// Por exemplo, o range [2022, 2023, 2024, 2025] se torna [2022, 2023, 2023, 2024, 2024, 2025].
	years := lo.RangeFrom(iniYear, endYear-iniYear+1)
	for i := len(years) - 2; i >= 1; i -= 2 {
		years = slices.Insert(years, i, years[i-1], years[i])
	}

	// Cria um histórico de operações por ticker por ano. Apenas a última
	// operação é relevante, pois ela possui o valor total acumulado.
	history := map[int]map[string]Operacao{}
	for _, o := range oprs {
		if history[o.Data.Year()] == nil {
			history[o.Data.Year()] = map[string]Operacao{}
		}
		history[o.Data.Year()][o.Ticker] = o
	}

	var bens []BensDireitos
	visited := map[string]BensDireitoTicker{}
	for _, pair := range lo.Chunk(years, 2) {
		iniYear, endYear := pair[0], pair[1]
		bem := BensDireitos{
			AnoAnterior: iniYear,
			AnoCorrente: endYear,
		}
		for ticker, opr := range history[iniYear] {
			visited[ticker] = BensDireitoTicker{Ticker: ticker, SituacaoAnterior: opr.ValorTotalAc, SituacaoCorrente: opr.ValorTotalAc, Qtd: opr.QtdAc, PrecoMedio: opr.PrecoMedio}
		}
		for ticker, opr := range history[endYear] {
			visited[ticker] = BensDireitoTicker{Ticker: ticker, SituacaoAnterior: visited[ticker].SituacaoAnterior, SituacaoCorrente: opr.ValorTotalAc, Qtd: opr.QtdAc, PrecoMedio: opr.PrecoMedio}
		}
		for key, b := range sortKeysIter(visited) {
			if _, ok := history[bem.AnoCorrente][key]; b.SituacaoAnterior.IsZero() && b.SituacaoCorrente.IsZero() && !ok {
				// Situações zeradas no ano anterior e corrente cujo ano corrente não existe não são mostradas.
				continue
			}
			bem.Tickers = append(bem.Tickers, b)
		}
		bens = append(bens, bem)
	}

	return bens
}

type BensDireitos struct {
	AnoAnterior int
	AnoCorrente int
	Tickers     []BensDireitoTicker
}

type BensDireitoTicker struct {
	Ticker           string
	Qtd              decimal.Decimal
	PrecoMedio       decimal.Decimal
	SituacaoAnterior decimal.Decimal
	SituacaoCorrente decimal.Decimal
}

func (b BensDireitoTicker) Discriminacao() string {
	return fmt.Sprintf("%s AÇÕES %s COM PREÇO MÉDIO DE R$ %s", b.Qtd.String(), b.Ticker, b.PrecoMedio.StringFixed(2))
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

	lucroAcumuladoAnos := decimal.Zero
	for year, oprs := range s.Operacoes.GroupByYear() {

		var rt RendimentosTributaveis
		rt.Ano = year

		for month, oprs := range oprs.GroupByMonth() {

			vendaNoMes := oprs.ValorTotalAcumulado()
			lucroNoMes := decimal.Zero

			for _, o := range oprs {
				if tipoOprConfig[o.Tipo].IsRendimentoTributavelApos20k && vendaNoMes.GreaterThan(s.Settings.AcaoLimiteIsento) && o.Lucro.IsPositive() {
					lucroNoMes = lucroNoMes.Add(o.Lucro)
				}
				if tipoOprConfig[o.Tipo].IsLucroTributavel && o.Lucro.IsPositive() {
					lucroNoMes = lucroNoMes.Add(o.Lucro)
				}
				if tipoOprConfig[o.Tipo].IsPrejuizoAbativel && o.Lucro.IsNegative() {
					lucroNoMes = lucroNoMes.Add(o.Lucro)
				}
			}

			if !lucroNoMes.IsZero() {
				lucroAcumuladoAnos = lucroAcumuladoAnos.Add(lucroNoMes)
				rtm := RendimentoTributavelMensal{Mes: month, Lucro: lucroNoMes, LucroAc: lucroAcumuladoAnos}
				if lucroAcumuladoAnos.IsPositive() {
					rtm.IR = lucroAcumuladoAnos.Mul(s.Settings.AcaoSwingTradeIR)
					lucroAcumuladoAnos = decimal.Zero
				}
				rt.Meses = append(rt.Meses, rtm)
				rt.Total = rt.Total.Add(lucroNoMes)
				rt.TotalAc = lucroAcumuladoAnos
				rt.TotalIR = rt.TotalIR.Add(rtm.IR)
			}
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
	TotalAc decimal.Decimal
	TotalIR decimal.Decimal
}

type RendimentoTributavelMensal struct {
	Mes     string
	Lucro   decimal.Decimal
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

func (o Operacoes) GetID(id int64) Operacao {
	if id == 0 {
		return Operacao{}
	}
	return o[id-1]
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
	QtdAc         decimal.Decimal
	ValorTotalAc  decimal.Decimal
	PrecoMedio    decimal.Decimal
	ValorCompra   decimal.Decimal
	Lucro         decimal.Decimal // Lucro ou prejuízo da operação de Venda, Bonificação, Grupamento, Subscrição Compra, Redução de Capital.
	Fracao        decimal.Decimal // Parte fracionária resultante de Bonificação, Grupamento ou Desdobramento.
	Fator         decimal.Decimal // Fator de Bonificação, Grupamento ou Desdobramento e Redução de Capital.

	// Opções.
	// Strike     string
	// Vencimento time.Time
	// Premio     decimal.Decimal
	// Exercido   bool
}

func (o *Operacao) CalcCompra(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.QtdAc = p.QtdAc.Add(o.Qtd)
	o.ValorTotalAc = p.ValorTotalAc.Add(o.ValorTotal)
	o.CalcPrecoMedio()
}

func (o *Operacao) CalcVenda(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Sub(o.Taxas)
	o.ValorCompra = p.PrecoMedio.Mul(o.Qtd)
	o.QtdAc = p.QtdAc.Sub(o.Qtd)
	o.ValorTotalAc = p.ValorTotalAc.Sub(o.ValorCompra)
	o.Lucro = o.ValorTotal.Sub(o.ValorCompra)
	o.CalcPrecoMedio()
}

func (o *Operacao) CalcBonificacao(p Operacao, alterarPrecoMedio bool) {
	if !o.Fator.IsZero() {
		o.Fracao = p.QtdAc.Mul(o.Fator)
		o.Qtd = o.Fracao.Truncate(0)
		o.Fracao = o.Fracao.Sub(o.Qtd).Abs()
	}
	o.QtdAc = p.QtdAc.Add(o.Qtd)
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.ValorTotalAc = o.ValorTotal.Add(p.ValorTotalAc)
	o.Lucro = o.ValorTotal
	if alterarPrecoMedio {
		o.CalcPrecoMedio()
	}
}

func (o *Operacao) CalcDesdobramento(p Operacao) {
	o.QtdAc = p.QtdAc.Mul(o.Fator)
	o.CalcPrecoMedio()
}

func (o *Operacao) CalcGrupamento(p Operacao) {
	qtdAcFra := p.QtdAc.Div(o.Fator)
	qtdAcInt := qtdAcFra.Truncate(0)
	o.Fracao = qtdAcFra.Sub(qtdAcInt)
	o.QtdAc = qtdAcInt
	o.ValorTotalAc = qtdAcInt.Mul(p.ValorTotalAc.Div(qtdAcFra))
	o.CalcPrecoMedio()
}

func (o *Operacao) CalcLeilaoFracao(p Operacao) {
	o.Lucro = o.ValorUnitario.Mul(o.Qtd).Sub(o.PrecoMedio.Mul(o.Qtd))
}

func (o *Operacao) CalcDividendos(p Operacao) {
}

func (o *Operacao) CalcJSCP(p Operacao) {
}

func (o *Operacao) CalcRendTrib(p Operacao) {
}

func (o *Operacao) CalcSubscricaoCompra(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Add(o.Taxas)
	o.ValorTotalAc = p.ValorTotalAc.Add(o.ValorTotal)
	o.CalcPrecoMedio()
}

func (o *Operacao) CalcSubscricaoVenda(p Operacao) {
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd).Sub(o.Taxas)
	o.Lucro = o.ValorTotal
	o.CalcPrecoMedio()
}

func (o *Operacao) CalcSubscricaoExercicio(p Operacao) {
	o.QtdAc = p.QtdAc.Add(o.Qtd)
	o.ValorTotal = o.ValorUnitario.Mul(o.Qtd)
	o.ValorTotalAc = p.ValorTotalAc.Add(o.ValorTotal)
	o.CalcPrecoMedio()
}

func (o *Operacao) CalcReducaoCapital(p Operacao) {
	if o.Fator.IsPositive() {
		redcap := p.ValorTotalAc.Mul(o.Fator)
		restituicao := p.QtdAc.Mul(o.ValorUnitario)
		o.ValorTotal = restituicao
		o.ValorCompra = redcap
		o.Lucro = restituicao.Sub(redcap)
		o.ValorTotalAc = p.ValorTotalAc.Sub(redcap)
	} else {
		o.ValorTotalAc = p.ValorTotalAc.Sub(o.ValorTotal)
	}
	o.CalcPrecoMedio()
}

func (o *Operacao) Inherit(p Operacao) {
	o.ValorTotalAc = p.ValorTotalAc
	o.QtdAc = p.QtdAc
	o.PrecoMedio = p.PrecoMedio
}

func (o *Operacao) CalcPrecoMedio() {
	if o.QtdAc.IsZero() {
		o.PrecoMedio = decimal.Zero
	} else {
		o.PrecoMedio = o.ValorTotalAc.Div(o.QtdAc)
	}
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
