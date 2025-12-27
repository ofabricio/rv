package data

import (
	"time"

	"github.com/shopspring/decimal"
)

type Operacao struct {
	Data          time.Time `json:",format:DateOnly"` // Opções: Data de Compra/Venda da opção ou Data de Vencimento ou Data de Encerramento.
	Ticker        string
	Tipo          Tipo
	Qtd           decimal.Decimal
	ValorUnitario decimal.Decimal // Opções: Strike.
	Taxas         decimal.Decimal
	ValorTotal    decimal.Decimal
	ValorCompra   decimal.Decimal
	Lucro         decimal.Decimal // Lucro ou prejuízo da operação de Venda, Bonificação, Grupamento, Subscrição Compra, Redução de Capital, Opções.
	Fracao        decimal.Decimal // Parte fracionária resultante de Bonificação, Grupamento ou Desdobramento.
	Fator         decimal.Decimal // Fator de Bonificação, Grupamento ou Desdobramento e Redução de Capital.

	// Opções.
	Serie          string
	ValorExercicio decimal.Decimal // Valor da ação no dia do exercício da opção.
	Premio         decimal.Decimal
	Vencimento     time.Time `json:",format:DateOnly"`
}

func (o *Operacao) IsAcao() bool {
	return o.Serie == ""
}

func (o *Operacao) IsOpcao() bool {
	return o.Serie != ""
}

type Operavel interface {
	Accept(Visitor)
}

type Visitor interface {
	VisitCompra(*Operacao)
	VisitVenda(*Operacao)
	VisitBonificacao(*Operacao)
	VisitDesdobramento(*Operacao)
	VisitGrupamento(*Operacao)
	VisitLeilaoFracao(*Operacao)
	VisitDividendos(*Operacao)
	VisitJSCP(*Operacao)
	VisitJSCPNaoPago(*Operacao)
	VisitAluguel(*Operacao)
	VisitRendTrib(*Operacao)
	VisitReducaoCapital(*Operacao)
	VisitSubscricaoCompra(*Operacao)
	VisitSubscricaoVenda(*Operacao)
	VisitSubscricaoExercicio(*Operacao)
	VisitVendaPut(*Operacao)
	VisitVendaPutExercida(*Operacao)
	VisitVendaPutNaoExercida(*Operacao)
	VisitCompraPut(*Operacao)
	VisitCompraPutExercida(*Operacao)
	VisitCompraPutNaoExercida(*Operacao)
	VisitCompraCall(*Operacao)
	VisitCompraCallExercida(*Operacao)
	VisitCompraCallNaoExercida(*Operacao)
	VisitVendaCall(*Operacao)
	VisitVendaCallExercida(*Operacao)
	VisitVendaCallNaoExercida(*Operacao)
}

func (o *Operacao) Accept(v Visitor) {
	switch o.Tipo {
	case COMPRA:
		v.VisitCompra(o)
	case VENDA:
		v.VisitVenda(o)
	case BONIFICACAO:
		v.VisitBonificacao(o)
	case DESDOBRAMENTO:
		v.VisitDesdobramento(o)
	case GRUPAMENTO:
		v.VisitGrupamento(o)
	case LEILAO_FRACAO:
		v.VisitLeilaoFracao(o)
	case DIVIDENDOS:
		v.VisitDividendos(o)
	case JSCP:
		v.VisitJSCP(o)
	case JSCP_NAO_PAGO:
		v.VisitJSCPNaoPago(o)
	case ALUGUEL:
		v.VisitAluguel(o)
	case REND_TRIB:
		v.VisitRendTrib(o)
	case REDUCAO_CAPITAL:
		v.VisitReducaoCapital(o)
	case SUBSCRICAO_COMPRA:
		v.VisitSubscricaoCompra(o)
	case SUBSCRICAO_VENDA:
		v.VisitSubscricaoVenda(o)
	case SUBSCRICAO_COMPRA_EX:
		v.VisitSubscricaoExercicio(o)
	case PUT_VENDA:
		v.VisitVendaPut(o)
	case PUT_VENDA_EX:
		v.VisitVendaPutExercida(o)
	case PUT_VENDA_NE:
		v.VisitVendaPutNaoExercida(o)
	case PUT_COMPRA:
		v.VisitCompraPut(o)
	case PUT_COMPRA_EX:
		v.VisitCompraPutExercida(o)
	case PUT_COMPRA_NE:
		v.VisitCompraPutNaoExercida(o)
	case CALL_COMPRA:
		v.VisitCompraCall(o)
	case CALL_COMPRA_EX:
		v.VisitCompraCallExercida(o)
	case CALL_COMPRA_NE:
		v.VisitCompraCallNaoExercida(o)
	case CALL_VENDA:
		v.VisitVendaCall(o)
	case CALL_VENDA_EX:
		v.VisitVendaCallExercida(o)
	case CALL_VENDA_NE:
		v.VisitVendaCallNaoExercida(o)
	}
}
