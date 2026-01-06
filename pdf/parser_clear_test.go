package pdf

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestParseContent(t *testing.T) {

	got, err := ParseClear(clearPdfContentExample)
	if err != nil {
		t.Fatalf("Erro ao parsear nota: %v", err)
	}

	exp := Nota{
		DataPregao: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		Negociacoes: []Negociacao{
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

	if fmt.Sprintf("%+v", got) != fmt.Sprintf("%+v", exp) {
		t.Fatalf("\nGot:\n%+v\nExp:\n%+v", got, exp)
	}
}

var clearPdfContentExample = strings.TrimSpace(`
Negociações
Negócios realizados
Q
Negociação
C/V
Tipo mercado
Prazo
Especificação do título
Obs. (*)
Quantidade
Preço / Ajuste
Valor Operação / Ajuste
D/C
1-BOVESPA
C
VISTA
MELNICK          ON 
NM
@
100
3,83
383,00
D
1-BOVESPA
C
VISTA
PETROBRAS          PN 
EDJ N2
@
100
30,15
3.015,00
D
1-BOVESPA
C
VISTA
SYN PROP TEC          ON 
NM
@
200
4,71
942,00
D
NOTA DE NEGOCIAÇÃO
Nr. nota
123123123
Folha
1
Data pregão
05/01/2026
CLEAR CTVM S/A
Avenida Ataulfo de Paiva, 153 - Sala 201 Leblon
12312-123
Rio de Janeiro - RJ
Tel. +55 11 4003-6245   Fax:
Internet: https://corretora.clear.com.br/ SAC: 0800-77-40404
e-mail:
C.N.P.J: 15.107.963/0001-66
Carta Patente:
Ouvidoria: Tel. 0800-722-3730
e-mail ouvidoria:
Cliente
12312312
FULANO DE TAL
RUA VINTE CINCO DE MARÇO, 123 - APT 123 - LUGAR
ONDEMORA
Tel. (11) 99999-8888
29000-100 CIDADE - UF
Participante destino do 
repasse
-
Cliente
0
Valor
Banco
Agência
Conta corrente
Acionista
Administrador
000.111.222-33
C.P.F./C.N.P.J/C.V.M./C.O.B.
Código cliente
 1-4   12312312 
Assessor
12312
Custodiante
C.I
N
Complemento nome
P. Vinc
N
0,00
0,00
4.340,00
0,00
0,00
0,00
0,00
4.340,00
Resumo dos Negócios
Debêntures
 
Vendas à vista
 
Compras à vista
 
Opções - compras
 
Opções - vendas
 
Operações à termo
 
Valor das oper. c/ títulos públ. (v. nom.)
 
Valor das operações
 
Especificações diversas
A coluna Q indica liquidação no Agente do Qualificado.
(*) Observações
A - Posição futuro
T - Liquidação pelo Bruto
2 - Corretora ou pessoa vinculada atuou na contra parte.
C - Clubes e fundos de Ações
I - POP
# - Negócio direto
P - Carteira Própria
8 - Liquidação Institucional
H - Home Broker
D - Day Trade
X - Box
F - Cobertura
Y - Desmanche de Box
B - Debêntures
L - Precatório
4.340,97
Total CBLC
D
4.340,00
Valor líquido das operações
D
0,97
Taxa de liquidação
D
0,00
Taxa de Registro
 
0,32
Total Bovespa / Soma
D
0,00
Taxa de termo/opções
 
0,00
Taxa A.N.A.
 
0,21
Emolumentos
D
0,11
Taxa de Transf. de Ativos
D
 
0,00
Total Custos / Despesas
 
 
0,00
Taxa Operacional
 
0,00
Execução
 
0,00
Taxa de Custódia
0,00
Impostos
 
0,00
I.R.R.F. s/ operações, base R$0,00
0,00
Outros
 
4.341,29
Líquido para 07/01/2026
D
Resumo Financeiro
Clearing
 
Bolsa
 
Custos Operacionais
Observação: (1) As operações a termo não são computadas no líquido da fatura.
Capitais e regiões metropolitanas: +55 11 4003-6245 Demais localidades: 0800-880-3710 SAC: 0800-77-40404 Ouvidoria: 0800-722-3730
`)
