package data

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCarteira(t *testing.T) {

	tt := []struct {
		Give string
		Then string
	}{
		{
			Give: "testdata/_carteira.give",
			Then: "testdata/_carteira.then",
		},
		{
			Give: "testdata/dedo_duro.give",
			Then: "testdata/dedo_duro.then",
		},
		{
			Give: "testdata/dedo_duro2.give",
			Then: "testdata/dedo_duro2.then",
		},
		{
			Give: "testdata/alteracao_config.give",
			Then: "testdata/alteracao_config.then",
		},
		{
			Give: "testdata/compra_venda.give",
			Then: "testdata/compra_venda.then",
		},
		{
			Give: "testdata/compra_venda2.give",
			Then: "testdata/compra_venda2.then",
		},
		{
			Give: "testdata/compra_venda3.give",
			Then: "testdata/compra_venda3.then",
		},
		{
			Give: "testdata/dividendos.give",
			Then: "testdata/dividendos.then",
		},
		{
			Give: "testdata/reembolso.give",
			Then: "testdata/reembolso.then",
		},
		{
			Give: "testdata/bonificacao.give",
			Then: "testdata/bonificacao.then",
		},
		{
			Give: "testdata/bonificacao2.give",
			Then: "testdata/bonificacao2.then",
		},
		{
			Give: "testdata/grupamento.give",
			Then: "testdata/grupamento.then",
		},
		{
			Give: "testdata/desdobramento.give",
			Then: "testdata/desdobramento.then",
		},
		{
			Give: "testdata/rendimentos.give",
			Then: "testdata/rendimentos.then",
		},
		{
			Give: "testdata/rend_trib.give",
			Then: "testdata/rend_trib.then",
		},
		{
			Give: "testdata/reducao_capital.give",
			Then: "testdata/reducao_capital.then",
		},
		{
			Give: "testdata/readme_file.give",
			Then: "testdata/readme_file.then",
		},
		{
			Give: "testdata/jscp.give",
			Then: "testdata/jscp.then",
		},
		{
			Give: "testdata/aluguel.give",
			Then: "testdata/aluguel.then",
		},
		{
			Give: "testdata/bens_e_direitos.give",
			Then: "testdata/bens_e_direitos.then",
		},
		{
			Give: "testdata/subscricao.give",
			Then: "testdata/subscricao.then",
		},
		{
			Give: "testdata/opcoes_venda_put.give",
			Then: "testdata/opcoes_venda_put.then",
		},
		{
			Give: "testdata/opcoes_venda_call.give",
			Then: "testdata/opcoes_venda_call.then",
		},
		{
			Give: "testdata/opcoes_compra_put.give",
			Then: "testdata/opcoes_compra_put.then",
		},
		{
			Give: "testdata/opcoes_compra_call.give",
			Then: "testdata/opcoes_compra_call.then",
		},
		{
			Give: "testdata/opcoes_bens_e_direitos.give",
			Then: "testdata/opcoes_bens_e_direitos.then",
		},
		{
			Give: "testdata/opcoes_divida_onus_reais.give",
			Then: "testdata/opcoes_divida_onus_reais.then",
		},
		{
			Give: "testdata/opcoes_serie_duplicada.give",
			Then: "testdata/opcoes_serie_duplicada.then",
		},
		{
			Give: "testdata/day_trade.give",
			Then: "testdata/day_trade.then",
		},
	}

	for _, tc := range tt {

		t.Run(tc.Then, func(t *testing.T) {

			var got bytes.Buffer

			c := NewCarteira()
			c.Load(tc.Give)
			c.Print("table", &got)

			// Then.

			exp, err := os.ReadFile(tc.Then)
			if err != nil {
				t.Errorf("failed to read then file: %v", err)
			}

			if !bytes.Equal(bytes.TrimSpace(got.Bytes()), bytes.TrimSpace(exp)) {
				t.Errorf("\nFile: %s\nGot:\n%s\nExp:\n%s", tc.Then, got.Bytes(), exp)
			}
		})
	}
}

func TestCarteiraValorizacao(t *testing.T) {

	r := strings.NewReader(`
		{ "Data": "2026-01-05", "Ticker": "MELK3", "Tipo": "Compra", "Qtd": "100", "ValorUnitario": "3.83", "Taxas": "0.11" }
		{ "Data": "2026-01-05", "Ticker": "PETR4", "Tipo": "Compra", "Qtd": "100", "ValorUnitario": "30.15", "Taxas": "0.90" }
		{ "Data": "2026-01-05", "Ticker": "SYNE3", "Tipo": "Compra", "Qtd": "200", "ValorUnitario": "4.71", "Taxas": "0.28" }
		{ "Data": "2026-01-15", "Ticker": "MELK3", "Tipo": "Compra", "Qtd": "700", "ValorUnitario": "3.7", "Taxas": "0.76" }
		{ "Data": "2026-01-15", "Ticker": "SYNE3", "Tipo": "Compra", "Qtd": "200", "ValorUnitario": "4.83", "Taxas": "0.29" }
	`)

	c := NewCarteira()
	c.IFunc = func(s string) (string, error) {
		return map[string]string{
			"MELK3": "3.61",
			"PETR4": "32.04",
			"SYNE3": "4.80",
		}[s], nil
	}

	if err := c.Read(r); err != nil {
		t.Fatalf("failed to read operations: %v", err)
	}

	var got strings.Builder
	p := PrinterTable{c}
	p.PrintValorizacao(&got)

	exp := strings.TrimSpace(`
┌───────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                VALORIZAÇÃO                                                │
├────────┬───────────┬───────┬─────────┬──────────┬────────────┬────────┬───────┬───────────────────────────┤
│ Ticker │ Investido │  PM   │ Cotação │ Variação │ Valorizado │ Ganho  │ Vende │          Compra           │
├────────┼───────────┼───────┼─────────┼──────────┼────────────┼────────┼───────┼───────────────────────────┤
│ MELK3  │   2973,87 │  3,72 │    3,61 │   -2,89% │    2888,00 │ -85,87 │   -23 │ -2 PETR4 -17 SYNE3        │
│ PETR4  │   3015,90 │ 30,16 │   32,04 │    6,24% │    3204,00 │ 188,10 │     5 │ 44 MELK3 33 SYNE3         │
│ SYNE3  │   1908,57 │  4,77 │    4,80 │    0,60% │    1920,00 │  11,43 │     2 │ 2 MELK3                   │
├────────┼───────────┼───────┼─────────┼──────────┼────────────┼────────┼───────┼───────────────────────────┤
│ Total  │   7898,34 │       │         │    1,44% │    8012,00 │ 113,66 │     7 │ 31 MELK3 3 PETR4 23 SYNE3 │
└────────┴───────────┴───────┴─────────┴──────────┴────────────┴────────┴───────┴───────────────────────────┘
	`)

	if strings.TrimSpace(got.String()) != exp {
		t.Errorf("\nGot:\n%s\nExp:\n%s", got.String(), exp)
	}
}
