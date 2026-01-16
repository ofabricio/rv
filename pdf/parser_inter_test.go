package pdf

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestParseInter(t *testing.T) {

	d, err := os.ReadFile("testdata/inter_nota_compra.give")
	if err != nil {
		panic(err)
	}

	got, err := ParseNotaContent(string(d))
	if err != nil {
		t.Fatalf("Erro ao parsear nota: %v", err)
	}

	exp := Nota{
		DataPregao: time.Date(2025, 4, 2, 0, 0, 0, 0, time.UTC),
		Negociacoes: []Negociacao{
			{
				CV:            "C",
				Titulo:        "SAPR4",
				Qtd:           decimal.RequireFromString("200"),
				ValorUnitario: decimal.RequireFromString("5.57"),
				ValorTotal:    decimal.RequireFromString("1114"),
				CD:            "D",
			},
		},
		Debentures:               decimal.RequireFromString("0"),
		VendasAVista:             decimal.RequireFromString("0"),
		ComprasAVista:            decimal.RequireFromString("1114"),
		OpcoesCompras:            decimal.RequireFromString("0"),
		OpcoesVendas:             decimal.RequireFromString("0"),
		OperacoesATermo:          decimal.RequireFromString("0"),
		ValorOperTitulosPubl:     decimal.RequireFromString("0"),
		ValorDasOperacoes:        decimal.RequireFromString("1114"),
		ValorLiquidoDasOperacoes: decimal.RequireFromString("1114"),
		TaxaLiquidacao:           decimal.RequireFromString("0.27"),
		TaxaRegistro:             decimal.RequireFromString("0"),
		TotalCBLC:                decimal.RequireFromString("1114.27"),
		TotalBovespaSoma:         decimal.RequireFromString("0.05"),
		TaxaTermoOpcoes:          decimal.RequireFromString("0"),
		TaxaANA:                  decimal.RequireFromString("0"),
		Emolumentos:              decimal.RequireFromString("0.05"),
		TaxaTransfAtivos:         decimal.RequireFromString("0"),
		TotalCustosDespesas:      decimal.RequireFromString("0"),
		TaxaOperacional:          decimal.RequireFromString("0"),
		Execucao:                 decimal.RequireFromString("0"),
		TaxaCustodia:             decimal.RequireFromString("0"),
		Impostos:                 decimal.RequireFromString("0"),
		IRRF:                     decimal.RequireFromString("0"),
		Outros:                   decimal.RequireFromString("0"),
		IRRFBase:                 decimal.RequireFromString("0"),
		TotalLiquido:             decimal.RequireFromString("1114.32"),
		LiquidoPara:              time.Date(2025, 4, 4, 0, 0, 0, 0, time.UTC),
	}

	if fmt.Sprintf("%+v", got) != fmt.Sprintf("%+v", exp) {
		t.Fatalf("\nGot:\n%+v\nExp:\n%+v", got, exp)
	}
}
