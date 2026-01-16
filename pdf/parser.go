package pdf

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ofabricio/bnf"
	"github.com/shopspring/decimal"
)

func ParseNotaContent(content string) (Nota, error) {
	if strings.Contains(content, "CLEAR CTVM S/A") {
		return ParseClear(content)
	}
	if strings.Contains(content, "INTER DTVM LTDA") {
		return ParseInter(content)
	}
	return Nota{}, fmt.Errorf("parser para este tipo de nota não implementado")
}

func parse(content, BNF string) (Nota, error) {

	b := bnf.Compile(BNF)
	v := bnf.Parse(b, content)
	if v.Type == "Error" {
		return Nota{}, fmt.Errorf("não foi possível parsear o arquivo pdf")
	}

	var n Nota
	for _, node := range v.Next {
		switch node.Type {
		case "Negociacoes":
			for _, neg := range node.Next {
				var v Negociacao
				for _, item := range neg.Next {
					switch item.Type {
					case "Titulo":
						v.Titulo = mapTituloTicker[item.Text]
					case "CV", "CD":
						setField(&v, item.Type, item.Text)
					case "Qtd", "ValorUnitario", "ValorTotal":
						setField(&v, item.Type, decimal.RequireFromString(item.Text))
					}
				}
				n.Negociacoes = append(n.Negociacoes, v)
			}
		case "DataPregao", "LiquidoPara":
			t, _ := time.Parse("02/01/2006", node.Text)
			setField(&n, node.Type, t)
		case "PIS", "COFINS":
		default:
			setField(&n, node.Type, decimal.RequireFromString(node.Text))
		}
	}

	return n, nil
}

// setField usa reflection para encontrar o campo na struct e setar o valor.
// O nome do campo na struct tem que ser igual ao nome definido no BNF.
func setField(strct any, field string, v any) {
	reflect.ValueOf(strct).Elem().FieldByName(field).Set(reflect.ValueOf(v))
}

type Nota struct {
	Negociacoes []Negociacao
	DataPregao  time.Time
	// Resumo dos Negócios.
	Debentures           decimal.Decimal
	VendasAVista         decimal.Decimal
	ComprasAVista        decimal.Decimal
	OpcoesCompras        decimal.Decimal
	OpcoesVendas         decimal.Decimal
	OperacoesATermo      decimal.Decimal
	ValorOperTitulosPubl decimal.Decimal
	ValorDasOperacoes    decimal.Decimal
	// Resumo Financeiro.
	// Clearing.
	ValorLiquidoDasOperacoes decimal.Decimal
	TaxaLiquidacao           decimal.Decimal
	TaxaRegistro             decimal.Decimal
	TotalCBLC                decimal.Decimal
	// Bolsa.
	TotalBovespaSoma decimal.Decimal
	TaxaTermoOpcoes  decimal.Decimal
	TaxaANA          decimal.Decimal
	Emolumentos      decimal.Decimal
	TaxaTransfAtivos decimal.Decimal
	// Custos Operacionais.
	TotalCustosDespesas decimal.Decimal
	TaxaOperacional     decimal.Decimal
	Execucao            decimal.Decimal
	TaxaCustodia        decimal.Decimal
	Impostos            decimal.Decimal
	IRRF                decimal.Decimal
	Outros              decimal.Decimal
	IRRFBase            decimal.Decimal
	//
	TotalLiquido decimal.Decimal
	LiquidoPara  time.Time
}

type Negociacao struct {
	CV            string
	Titulo        string
	Qtd           decimal.Decimal
	ValorUnitario decimal.Decimal
	ValorTotal    decimal.Decimal
	CD            string
}

var mapTituloTicker = map[string]string{
	"BRASIL ON NM":        "BBAS3",
	"GERDAU MET PN N1":    "GOAU4",
	"GERDAU MET PN N2":    "GOAU4",
	"ISHARES BOVA CI":     "BOVA11",
	"IT NOW IDIV CI":      "DIVO11",
	"KLABIN S/A PN ED N2": "KLBN4",
	"KLABIN S/A PN N2":    "KLBN4",
	"MARFRIG ON NM":       "MBRF3",
	"MELNICK ON NM":       "MELK3",
	"PETROBRAS PN EDJ N2": "PETR4",
	"SANEPAR PN N2":       "SAPR4",
	"SYN PROP TEC ON NM":  "SYNE3",
	"TAESA UNT N2":        "TAEE11",
}
