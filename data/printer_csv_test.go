package data

import (
	"strings"
	"testing"
)

func TestCSVPrinter(t *testing.T) {

	give := strings.NewReader(`
		{ "Data": "2025-01-01", "Ticker": "PETR4", "Tipo": "Compra", "Qtd": "100", "ValorUnitario": "10" }
		{ "Data": "2025-02-01", "Ticker": "PETR4", "Tipo": "Compra", "Qtd": "100", "ValorUnitario": "12" }
	`)

	then := strings.TrimSpace(`
ID,Ticker,Data,Operação,Qtd,V. Unit.,V. Total,Taxas,Qtd Ac.,V. Total Ac.,PM,V. Compra,Lucro
1,PETR4,2025-01-01,Compra,100,"10,00","1000,00","0,00",100,"1000,00","10,00","0,00","0,00"
2,PETR4,2025-02-01,Compra,100,"12,00","1200,00","0,00",200,"2200,00","11,00","0,00","0,00"
	`)

	c := NewCarteira()
	c.Read(give)

	var got strings.Builder

	p := CSVPrinter{c}
	p.PrintOperacoesComAcoes(&got)

	if got := strings.TrimSpace(got.String()); got != then {
		t.Errorf("\nGot:\n%s\nExp:\n%s", got, then)
	}
}
