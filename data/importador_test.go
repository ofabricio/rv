package data

import (
	"fmt"
	"testing"
	"time"

	"github.com/ofabricio/rv/pdf"
	"github.com/shopspring/decimal"
)

func TestImportadorNotasProcessaNota(t *testing.T) {

	give := pdf.Nota{
		DataPregao: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		Negociacoes: []pdf.Negociacao{
			{
				CV:            "C",
				Titulo:        "MELK3",
				ONPN:          "ON NM",
				Qtd:           decimal.RequireFromString("100"),
				Preco:         decimal.RequireFromString("3.83"),
				ValorOperacao: decimal.RequireFromString("383"),
				DC:            "D",
			},
			{
				CV:            "C",
				Titulo:        "PETR4",
				ONPN:          "PN EDJ N2",
				Qtd:           decimal.RequireFromString("100"),
				Preco:         decimal.RequireFromString("30.15"),
				ValorOperacao: decimal.RequireFromString("3015"),
				DC:            "D",
			},
			{
				CV:            "C",
				Titulo:        "SYNE3",
				ONPN:          "ON NM",
				Qtd:           decimal.RequireFromString("200"),
				Preco:         decimal.RequireFromString("4.71"),
				ValorOperacao: decimal.RequireFromString("942"),
				DC:            "D",
			},
		},
		Debentures:               decimal.RequireFromString("0"),
		VendasAVista:             decimal.RequireFromString("0"),
		ComprasAVista:            decimal.RequireFromString("4340"),
		OpcoesCompras:            decimal.RequireFromString("0"),
		OpcoesVendas:             decimal.RequireFromString("0"),
		OperacoesATermo:          decimal.RequireFromString("0"),
		ValorOperTitulosPubl:     decimal.RequireFromString("0"),
		ValorDasOperacoes:        decimal.RequireFromString("4340"),
		ValorLiquidoDasOperacoes: decimal.RequireFromString("4340"),
		TaxaLiquidacao:           decimal.RequireFromString("0.97"),
		TaxaRegistro:             decimal.RequireFromString("0"),
		TotalCBLC:                decimal.RequireFromString("4340.97"),
		TotalBovespaSoma:         decimal.RequireFromString("0.32"),
		TaxaTermoOpcoes:          decimal.RequireFromString("0"),
		TaxaANA:                  decimal.RequireFromString("0"),
		Emolumentos:              decimal.RequireFromString("0.21"),
		TaxaTransfAtivos:         decimal.RequireFromString("0.11"),
		TotalCustosDespesas:      decimal.RequireFromString("0"),
		TaxaOperacional:          decimal.RequireFromString("0"),
		Execucao:                 decimal.RequireFromString("0"),
		TaxaCustodia:             decimal.RequireFromString("0"),
		Impostos:                 decimal.RequireFromString("0"),
		IRRF:                     decimal.RequireFromString("0"),
		Outros:                   decimal.RequireFromString("0"),
		IRRFBase:                 decimal.RequireFromString("0"),
		TotalLiquido:             decimal.RequireFromString("4341.29"),
		LiquidoPara:              time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
	}

	var inp ImportadorNotas
	got := inp.processaNota(give)

	then := []Operacao{
		{
			Data:          time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
			Tipo:          COMPRA,
			Ticker:        "MELK3",
			Qtd:           decimal.RequireFromString("100"),
			ValorUnitario: decimal.RequireFromString("3.83"),
			Taxas:         decimal.RequireFromString("0.11"),
		},
		{
			Data:          time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
			Tipo:          COMPRA,
			Ticker:        "PETR4",
			Qtd:           decimal.RequireFromString("100"),
			ValorUnitario: decimal.RequireFromString("30.15"),
			Taxas:         decimal.RequireFromString("0.90"),
		},
		{
			Data:          time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
			Tipo:          COMPRA,
			Ticker:        "SYNE3",
			Qtd:           decimal.RequireFromString("200"),
			ValorUnitario: decimal.RequireFromString("4.71"),
			Taxas:         decimal.RequireFromString("0.28"),
		},
	}

	if fmt.Sprintf("%+v", got) != fmt.Sprintf("%+v", then) {
		t.Fatalf("\nGot:\n%+v\nExp:\n%+v", got, then)
	}
}
