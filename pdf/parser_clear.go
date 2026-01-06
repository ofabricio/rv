package pdf

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ofabricio/bnf"
	"github.com/shopspring/decimal"
)

func ParseNota(pdfFile string) (Nota, error) {
	ctt, err := Content(pdfFile)
	if err != nil {
		return Nota{}, err
	}
	return ParseContent(ctt)
}

func ParseContent(content string) (Nota, error) {
	return ParseClear(content)
}

func ParseClear(content string) (Nota, error) {
	BNF := `
		root = findNegociacoes:Negociacoes
			findData:DataPregao
			findPreco:Debentures
			findPreco:VendasAVista
			findPreco:ComprasAVista
			findPreco:OpcoesCompras
			findPreco:OpcoesVendas
			findPreco:OperacoesATermo
			findPreco:ValorOperTitulosPubl
			findPreco:ValorDasOperacoes
			findPreco:TotalCBLC
			findPreco:ValorLiquidoDasOperacoes
			findPreco:TaxaLiquidacao
			findPreco:TaxaRegistro
			findPreco:TotalBovespaSoma
			findPreco:TaxaTermoOpcoes
			findPreco:TaxaANA
			findPreco:Emolumentos
			findPreco:TaxaTransfAtivos
			findPreco:TotalCustosDespesas
			findPreco:TaxaOperacional
			findPreco:Execucao
			findPreco:TaxaCustodia
			findPreco:Impostos
			findPreco:IRRF
			findPreco:Outros
			findPreco:IRRFBase
			findPreco:TotalLiquido
			findData:LiquidoPara

		findData  = MATCH(ANYNOT(data)*)i  data
		findPreco = MATCH(ANYNOT(preco)*)i preco
		findNegociacoes =
			MATCH(ANYNOT('D/C')+)i 'D/C'i
			GROUP(negociacao+)

		negociacao = GROUP(Q neg CV merc prazo titulo ONPN obs qtd preco preco CD)
				Q = TEXT()
			neg = WS* '1-BOVESPA'
			merc = WS* 'VISTA'
			prazo = WS* TEXT()
			titulo = WS* MATCH(ANYNOT(ONPN)+)
			ONPN = JOIN('ON' WS+ TEXT(' ') 'NM' | 'PN' WS+ TEXT(' ') 'EDJ' WS+ TEXT(' ') 'N2')
			obs = WS* '@'
			qtd = WS* '\d+'r
				CV = WS* ('C' | 'V')
				CD = WS* ('C' | 'D')

		data  = '\d{2}\/\d{2}\/\d{4}'r
		preco = WS*'((\d+\.?)+,\d+)'r
	`
	b := bnf.Compile(BNF)
	v := bnf.Parse(b, content)
	if v.Type == "Error" {
		return Nota{}, fmt.Errorf("não foi possível parsear o arquivo pdf")
	}
	var n Nota
	for _, negociacao := range v.Next[0].Next {
		var neg Negociacao
		_ = negociacao.Next[0].Text                          // Q
		_ = negociacao.Next[1].Text                          // Negociação (ex. 1-IBOVESPA)
		neg.CV = negociacao.Next[2].Text                     // C/V
		_ = negociacao.Next[3].Text                          // Tipo mercado (ex. VISTA)
		_ = negociacao.Next[4].Text                          // Prazo
		titulo := strings.TrimSpace(negociacao.Next[5].Text) // Titulo
		neg.ONPN = negociacao.Next[6].Text                   // ON PN
		neg.Titulo = mapTituloTicker[strings.ToUpper(titulo)+" "+strings.ToUpper(neg.ONPN)]
		_ = negociacao.Next[7].Text       // Obs
		qtd := negociacao.Next[8].Text    // Qtd
		preco := negociacao.Next[9].Text  // Preço
		valOp := negociacao.Next[10].Text // Valor Operacao
		neg.DC = negociacao.Next[11].Text // DC
		if err := errors.Join(
			parseDecimal(preco, &neg.Preco, "Preço"),
			parseDecimal(valOp, &neg.ValorOperacao, "Valor Operação"),
			parseDecimal(qtd, &neg.Qtd, "Quantidade"),
		); err != nil {
			return n, err
		}
		n.Negociacoes = append(n.Negociacoes, neg)
	}
	return n, errors.Join(
		parseData(v.Next[1].Text, &n.DataPregao, "Data Pregão"),
		parseDecimal(v.Next[2].Text, &n.Debentures, "Debentures"),
		parseDecimal(v.Next[3].Text, &n.VendasAVista, "Vendas à Vista"),
		parseDecimal(v.Next[4].Text, &n.ComprasAVista, "Compras à Vista"),
		parseDecimal(v.Next[5].Text, &n.OpcoesCompras, "Opções - compras"),
		parseDecimal(v.Next[6].Text, &n.OpcoesVendas, "Opções - vendas"),
		parseDecimal(v.Next[7].Text, &n.OperacoesATermo, "Operações à termo"),
		parseDecimal(v.Next[8].Text, &n.ValorOperTitulosPubl, "Valor das oper. c/ títulos públ. (v. nom.)"),
		parseDecimal(v.Next[9].Text, &n.ValorDasOperacoes, "Valor das operações"),
		parseDecimal(v.Next[10].Text, &n.TotalCBLC, "Total CBLC"),
		parseDecimal(v.Next[11].Text, &n.ValorLiquidoDasOperacoes, "Valor Líquido das Operações"),
		parseDecimal(v.Next[12].Text, &n.TaxaLiquidacao, "Taxa Liquidação"),
		parseDecimal(v.Next[13].Text, &n.TaxaRegistro, "Taxa Registro"),
		parseDecimal(v.Next[14].Text, &n.TotalBovespaSoma, "Total Bovespa Soma"),
		parseDecimal(v.Next[15].Text, &n.TaxaTermoOpcoes, "Taxa de termo/opções"),
		parseDecimal(v.Next[16].Text, &n.TaxaANA, "Taxa A.N.A."),
		parseDecimal(v.Next[17].Text, &n.Emolumentos, "Emolumentos"),
		parseDecimal(v.Next[18].Text, &n.TaxaTransfAtivos, "Taxa de Transf. de Ativos"),
		parseDecimal(v.Next[19].Text, &n.TotalCustosDespesas, "Total Custos/Despesas"),
		parseDecimal(v.Next[20].Text, &n.TaxaOperacional, "Taxa Operacional"),
		parseDecimal(v.Next[21].Text, &n.Execucao, "Execução"),
		parseDecimal(v.Next[22].Text, &n.TaxaCustodia, "Taxa de Custódia"),
		parseDecimal(v.Next[23].Text, &n.Impostos, "Impostos"),
		parseDecimal(v.Next[24].Text, &n.IRRF, "I.R.R.F. s/ operações, base"),
		parseDecimal(v.Next[25].Text, &n.Outros, "Outros"),
		parseDecimal(v.Next[26].Text, &n.IRRFBase, "IRRFBase"),
		parseDecimal(v.Next[27].Text, &n.TotalLiquido, "Total Líquido"),
		parseData(v.Next[28].Text, &n.LiquidoPara, "Líquido para"),
	)
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
	ONPN          string
	Qtd           decimal.Decimal
	Preco         decimal.Decimal
	ValorOperacao decimal.Decimal
	DC            string
}

func parseData(v string, out *time.Time, label string) error {
	t, err := time.Parse("02/01/2006", v)
	if err != nil {
		return fmt.Errorf("erro ao parsear %s: %w", label, err)
	}
	*out = t
	return nil
}

func parseDecimal(v string, out *decimal.Decimal, label string) error {
	v = strings.ReplaceAll(v, ".", "")
	v = strings.ReplaceAll(v, ",", ".")
	d, err := decimal.NewFromString(v)
	if err != nil {
		return fmt.Errorf("erro ao parsear %s: %w", label, err)
	}
	*out = d
	return nil
}

var mapTituloTicker = map[string]string{
	"MELNICK ON NM":       "MELK3",
	"PETROBRAS PN EDJ N2": "PETR4",
	"SYN PROP TEC ON NM":  "SYNE3",
}
