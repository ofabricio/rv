package data

import (
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var DefaultParam = Param{
	SeparadorDecimal:  ",",
	FormatoData:       time.DateOnly,
	MostrarValorExato: false,
	FiltrarAno:        0,
	FiltrarTicker:     "",
}

type Param struct {
	// Mostra os valores exatos.
	// Exemplo: "0,123456" em vez de "0,12".
	MostrarValorExato bool

	// Formata valores usando esse separador de casas decimais.
	// Exemplo: "1234,56" (padrão).
	SeparadorDecimal string

	// Formata datas usando esse formato.
	// Exemplo: "02/01/2006" (padrão).
	FormatoData string

	// Filtra operações por ano (0 para mostrar todos os anos).
	FiltrarAno int

	// Filtra operações por ticker (vazio para mostrar todos os tickers).
	FiltrarTicker string
}

func (p *Param) FormatDecimal(v decimal.Decimal) string {
	if p.MostrarValorExato {
		return strings.Replace(v.String(), ".", p.SeparadorDecimal, 1)
	}
	return strings.Replace(v.StringFixed(2), ".", p.SeparadorDecimal, 1)
}

func (p *Param) FormatDate(t time.Time) string {
	return t.Format(p.FormatoData)
}

func (p *Param) FilterOperacao(o OperacaoConsolidada) bool {
	return p.FilterYear(o.Data.Year()) && p.FilterTicker(o.Ticker)
}

func (p *Param) FilterYear(year int) bool {
	return p.FiltrarAno == 0 || year == p.FiltrarAno
}

func (p *Param) FilterTicker(ticker string) bool {
	return p.FiltrarTicker == "" || ticker == p.FiltrarTicker
}
