package data

import (
	"strings"
	"testing"
)

func TestCSVPrinter(t *testing.T) {

	give := strings.NewReader(`
		{ "Data": "2025-01-01", "Ticker": "PETR4", "Tipo": "Compra", "Qtd": "100", "ValorUnitario": "10" }
		{ "Data": "2025-02-01", "Ticker": "PETR4", "Tipo": "Compra", "Qtd": "100", "ValorUnitario": "12" }
		{ "Data": "2025-03-01", "Ticker": "PETR4", "Serie": "PETRM100", "Tipo": "Venda PUT", "Vencimento": "2025-04-01", "Qtd": "100", "ValorUnitario": "12", "Premio": "1" }
		{ "Data": "2025-04-01", "Ticker": "PETR4", "Serie": "PETRM100", "Tipo": "Venda PUT (EX)", "Vencimento": "2025-04-01", "Qtd": "100", "ValorUnitario": "12", "Premio": "1", "ValorExercicio": "10.52" }
	`)

	then := strings.TrimSpace(`
ID,Ticker,Série,Data,Vencimento,Operação,Qtd,Fração,V. Unit.,V. Total,Taxas,IRRF,Qtd Ac.,V. Total Ac.,PM,V. Compra,Lucro,Prêmio,V. Exercício
1,PETR4,,2025-01-01,,Compra,100,"0,00","10,00","1000,00","0,00","0,00",100,"1000,00","10,00","0,00","0,00","0,00","0,00"
2,PETR4,,2025-02-01,,Compra,100,"0,00","12,00","1200,00","0,00","0,00",200,"2200,00","11,00","0,00","0,00","0,00","0,00"
3,PETR4,PETRM100,2025-03-01,2025-04-01,Venda PUT,100,"0,00","12,00","1200,00","0,00","0,00",200,"2200,00","11,00","0,00","100,00","1,00","0,00"
4,PETR4,PETRM100,2025-04-01,2025-04-01,Venda PUT (EX),100,"0,00","12,00","1200,00","0,00","0,00",300,"3300,00","11,00","0,00","100,00","1,00","10,52"
`)

	c := NewCarteira()
	c.Read(give)

	var got strings.Builder

	p := PrinterCSV{c}
	p.PrintOperacoesComAcoes(&got)

	if got := strings.TrimSpace(got.String()); got != then {
		t.Errorf("\nGot:\n%s\nExp:\n%s", got, then)
	}
}
