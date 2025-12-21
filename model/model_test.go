package model

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
			Give: "testdata/bens_e_direitos.give",
			Then: "testdata/bens_e_direitos.then",
		},
		{
			Give: "testdata/subscricao.give",
			Then: "testdata/subscricao.then",
		},
		{
			Give: "testdata/reducao_capital.give",
			Then: "testdata/reducao_capital.then",
		},
		{
			Give: "testdata/rend_trib.give",
			Then: "testdata/rend_trib.then",
		},
		{
			Give: "testdata/bonificacao.give",
			Then: "testdata/bonificacao.then",
		},
		{
			Give: "testdata/dividendos.give",
			Then: "testdata/dividendos.then",
		},
		{
			Give: "testdata/jscp.give",
			Then: "testdata/jscp.then",
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
			Give: "testdata/compra_venda.give",
			Then: "testdata/compra_venda.then",
		},
		{
			Give: "testdata/all.give",
			Then: "testdata/all.then",
		},
		{
			Give: "testdata/readme.give",
			Then: "testdata/readme.then",
		},
	}

	for _, tc := range tt {

		var got bytes.Buffer
		var s State

		s.Load(tc.Give)
		s.Print(&got)

		exp, err := os.ReadFile(tc.Then)
		if err != nil {
			t.Errorf("failed to read expected output: %v", err)
		}

		if !bytes.Equal(bytes.TrimSpace(got.Bytes()), bytes.TrimSpace(exp)) {
			t.Errorf("\nFile: %s\nGot:\n%s\nExp:\n%s", tc.Then, got.Bytes(), exp)
		}
	}
}
