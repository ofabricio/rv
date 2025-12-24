package data

import (
	"bytes"
	"os"
	"testing"
)

func TestStateLoad(t *testing.T) {

	tt := []struct {
		Give string
		Then string
	}{
		{
			Give: "testdata/_carteira.give",
			Then: "testdata/_carteira.then",
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
			Give: "testdata/dividendos.give",
			Then: "testdata/dividendos.then",
		},
		{
			Give: "testdata/bonificacao.give",
			Then: "testdata/bonificacao.then",
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
	}

	for _, tc := range tt {

		t.Run(tc.Then, func(t *testing.T) {

			var got bytes.Buffer

			var c Carteira
			c.Load(tc.Give)
			c.Print(&got)

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
