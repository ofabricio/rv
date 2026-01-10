package pdf

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ofabricio/bnf"
	"github.com/shopspring/decimal"
)

func ParseNota(pdfFile string) (Nota, error) {
	v, err := Content(pdfFile)
	if err != nil {
		return Nota{}, fmt.Errorf("erro ao abrir arquivo %s: %w", pdfFile, err)
	}
	return ParseContent(v)
}

func ParseContent(content string) (Nota, error) {
	if strings.Contains(content, "CLEAR CTVM S/A") {
		return ParseClear(content)
	}
	return Nota{}, fmt.Errorf("parser para este tipo de nota não implementado")
}

func ParseClear(content string) (Nota, error) {
	BNF := `
		root =
			FIND(GROUP(negociacao+)):Negociacoes
			FIND(data):DataPregao
			FIND(preco):Debentures
			FIND(preco):VendasAVista
			FIND(preco):ComprasAVista
			FIND(preco):OpcoesCompras
			FIND(preco):OpcoesVendas
			FIND(preco):OperacoesATermo
			FIND(preco):ValorOperTitulosPubl
			FIND(preco):ValorDasOperacoes
			FIND(preco):TotalCBLC
			FIND(preco):ValorLiquidoDasOperacoes
			FIND(preco):TaxaLiquidacao
			FIND(preco):TaxaRegistro
			FIND(preco):TotalBovespaSoma
			FIND(preco):TaxaTermoOpcoes
			FIND(preco):TaxaANA
			FIND(preco):Emolumentos
			FIND(preco):TaxaTransfAtivos
			FIND(preco):TotalCustosDespesas
			FIND(preco):TaxaOperacional
			FIND(preco):Execucao
			FIND(preco):TaxaCustodia
			FIND(preco):Impostos
			FIND(preco):IRRF
			FIND(preco):Outros
			FIND(preco):IRRFBase
			FIND(preco):TotalLiquido
			FIND(data):LiquidoPara

		negociacao = GROUP( q S neg S cv:CV S merc S prazo S titulo:Titulo S obs S qtd:Qtd S preco:ValorUnitario S preco:ValorTotal S cd:CD )
				 q = TEXT()
			   neg = '1-BOVESPA'
			  merc = 'VISTA'
			 prazo = TEXT()
			titulo = JOIN(ONPN | '\w+'r S TEXT(' ') titulo)
			  ONPN = 'ON' S TEXT(' ') 'NM'
				   | 'PN' S TEXT(' ') 'EDJ' S TEXT(' ') 'N2'
			   obs = '@'
			   qtd = '\d+'r
				cv = 'C' | 'V'
				cd = 'C' | 'D'

		data  = '\d{2}\/\d{2}\/\d{4}'r
		preco = JOIN(( '\d+'r '.'i? )+ ','i TEXT('.') '\d+'r)
		S = WS*
	`
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

var mapTituloTicker = map[string]string{
	"MELNICK ON NM":       "MELK3",
	"PETROBRAS PN EDJ N2": "PETR4",
	"SYN PROP TEC ON NM":  "SYNE3",
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
