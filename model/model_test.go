package model

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestStateLoad(t *testing.T) {

	tt := []struct {
		Give string
		Then string
	}{
		{
			Give: "testdata/_me_give.ndjson",
			Then: "testdata/_me_then.txt",
		},
		{
			Give: "testdata/opcoes_give.ndjson",
			Then: "testdata/opcoes_then.txt",
		},
		{
			Give: "testdata/bens_e_direitos_give.ndjson",
			Then: "testdata/bens_e_direitos_then.txt",
		},
		{
			Give: "testdata/subscricao_give.ndjson",
			Then: "testdata/subscricao_then.txt",
		},
		{
			Give: "testdata/reducao_capital_give.ndjson",
			Then: "testdata/reducao_capital_then.txt",
		},
		{
			Give: "testdata/rend_trib_give.ndjson",
			Then: "testdata/rend_trib_then.txt",
		},
		{
			Give: "testdata/bonificacao_give.ndjson",
			Then: "testdata/bonificacao_then.txt",
		},
		{
			Give: "testdata/dividendos_give.ndjson",
			Then: "testdata/dividendos_then.txt",
		},
		{
			Give: "testdata/jscp_give.ndjson",
			Then: "testdata/jscp_then.txt",
		},
		{
			Give: "testdata/grupamento_give.ndjson",
			Then: "testdata/grupamento_then.txt",
		},
		{
			Give: "testdata/desdobramento_give.ndjson",
			Then: "testdata/desdobramento_then.txt",
		},
		{
			Give: "testdata/rendimentos_give.ndjson",
			Then: "testdata/rendimentos_then.txt",
		},
		{
			Give: "testdata/compra_venda_give.ndjson",
			Then: "testdata/compra_venda_then.txt",
		},
		{
			Give: "testdata/all_give.ndjson",
			Then: "testdata/all_then.txt",
		},
		{
			Give: "testdata/readme_give.ndjson",
			Then: "testdata/readme_then.txt",
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

		if strings.TrimSpace(got.String()) != strings.TrimSpace(string(exp)) {
			t.Errorf("\nFile: %s\nGot:\n%s\nExp:\n%s", tc.Then, got.String(), exp)
		}
	}
}
