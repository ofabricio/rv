package data

import (
	"fmt"
	"time"

	"github.com/samber/lo"
)

var translateMonth = map[int]string{
	1:  "JAN",
	2:  "FEV",
	3:  "MAR",
	4:  "ABR",
	5:  "MAI",
	6:  "JUN",
	7:  "JUL",
	8:  "AGO",
	9:  "SET",
	10: "OUT",
	11: "NOV",
	12: "DEZ",
}

type OperacoesRowBuilder struct {
	Param Param
	ID    int
}

func (rb *OperacoesRowBuilder) Headers() []string {
	return []string{
		"ID",
		"Ticker",
		"Data",
		"Operação",
		"Qtd",
		"V. Unit.",
		"V. Total",
		"Taxas",
		"Qtd Ac.",
		"V. Total Ac.",
		"PM",
		"V. Compra",
		"Lucro",
	}
}

func (rb *OperacoesRowBuilder) Build(o *OperacaoConsolidada) []string {
	rb.ID++
	return []string{
		// ID
		fmt.Sprint(rb.ID),
		// Ticker
		lo.Ternary(o.Serie != "", fmt.Sprintf("%s %s", o.Ticker, o.Serie), o.Ticker),
		// Data
		lo.Ternary(!o.Vencimento.IsZero(), fmt.Sprintf("%s V %s", o.Data.Format(time.DateOnly), o.Vencimento.Format(time.DateOnly)), o.Data.Format(time.DateOnly)),
		// Operação
		string(o.Tipo),
		// Qtd
		lo.Ternary(o.Fracao.IsPositive(), fmt.Sprintf("%s (%s)", o.Qtd, rb.Param.FormatDecimal(o.Fracao)), o.Qtd.String()),
		// V. Unit.
		lo.Ternary(o.ValorExercicio.IsPositive(), fmt.Sprintf("(E %s) %s", rb.Param.FormatDecimal(o.ValorExercicio), rb.Param.FormatDecimal(o.ValorUnitario)), rb.Param.FormatDecimal(o.ValorUnitario)),
		// V. Total
		lo.Ternary(o.ValorExercicio.IsPositive(), fmt.Sprintf("(E %s) %s", rb.Param.FormatDecimal(o.ValorExercicio.Mul(o.Qtd)), rb.Param.FormatDecimal(o.ValorTotal)), rb.Param.FormatDecimal(o.ValorTotal)),
		// Taxas
		rb.Param.FormatDecimal(o.Taxas),
		// Qtd Ac.
		o.Agg.Qtd.String(),
		// V. Total Ac.
		rb.Param.FormatDecimal(o.Agg.ValorTotal),
		// PM
		rb.Param.FormatDecimal(o.Agg.PrecoMedio),
		// V. Compra
		rb.Param.FormatDecimal(o.ValorCompra),
		// Lucro
		lo.Ternary(o.Premio.IsPositive(), fmt.Sprintf("(P %s) %s", rb.Param.FormatDecimal(o.Premio), rb.Param.FormatDecimal(o.Lucro)), rb.Param.FormatDecimal(o.Lucro)),
	}
}
