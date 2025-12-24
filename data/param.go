package data

import (
	"strings"

	"github.com/shopspring/decimal"
)

var DefaultParam = Param{
	SeparadorDecimal: ",",
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
}

func (p *Param) FormatDecimal(v decimal.Decimal) string {
	return strings.Replace(v.StringFixed(2), ".", p.SeparadorDecimal, 1)
}
