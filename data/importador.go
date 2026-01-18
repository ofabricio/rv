package data

import (
	"encoding/json/v2"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ofabricio/rv/pdf"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

func ProcessarNotasPDF(dir string, w io.Writer) error {

	ops, err := ImportarNotasPDF(dir)
	if err != nil {
		return err
	}

	return outputOperacoes(ops, w)
}

func ProcessarNotaPDF(pdfData []byte, w io.Writer) error {

	ops, err := ImportarNotaPDF(pdfData)
	if err != nil {
		return err
	}

	return outputOperacoes(ops, w)
}

func outputOperacoes(ops []Operacao, w io.Writer) error {
	for _, op := range ops {
		if err := outputOperacao(op, w); err != nil {
			return err
		}
	}
	return nil
}

func outputOperacao(op Operacao, w io.Writer) error {
	if err := json.MarshalWrite(w, &op); err != nil {
		return err
	}
	_, err := w.Write([]byte("\n"))
	return err
}

func ImportarNotasPDF(dir string) ([]Operacao, error) {

	var oprs []Operacao

	walk := func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if de.IsDir() {
			return nil
		}

		if !strings.EqualFold(filepath.Ext(de.Name()), ".pdf") {
			return nil
		}

		d, err1 := os.ReadFile(path)
		if err1 != nil {
			return err1
		}

		ops, err := ImportarNotaPDF(d)
		if err != nil {
			return err
		}

		oprs = append(oprs, ops...)
		return nil
	}

	if err := filepath.WalkDir(dir, walk); err != nil {
		return nil, err
	}

	slices.SortStableFunc(oprs, func(a, b Operacao) int { return a.Data.Compare(b.Data) })

	return oprs, nil
}

func ImportarNotaPDF(pdfData []byte) ([]Operacao, error) {

	n, err := pdf.ParseNotaPDF(pdfData)
	if err != nil {
		return nil, err
	}

	return ImportarNota(n)
}

func ImportarNota(n pdf.Nota) ([]Operacao, error) {

	// Agrupa negociações iguais.
	var negs []pdf.Negociacao
	for _, same := range lo.PartitionBy(n.Negociacoes, func(n pdf.Negociacao) string { return n.Titulo + n.CV }) {
		if len(same) == 1 {
			negs = append(negs, same...)
		} else {
			neg := lo.FirstOrEmpty(same)
			neg.Qtd = lo.Reduce(same, func(acc decimal.Decimal, n pdf.Negociacao, _ int) decimal.Decimal { return acc.Add(n.Qtd) }, decimal.Zero)
			neg.ValorUnitario = lo.Reduce(same, func(acc decimal.Decimal, n pdf.Negociacao, _ int) decimal.Decimal { return acc.Add(n.ValorUnitario) }, decimal.Zero)
			neg.ValorUnitario = neg.ValorUnitario.Div(decimal.NewFromInt(int64(len(same))))
			neg.ValorTotal = neg.Qtd.Mul(neg.ValorUnitario)
			negs = append(negs, neg)
		}
	}

	totalTaxas := n.TotalLiquido.Sub(n.ValorLiquidoDasOperacoes)
	somasTaxas := decimal.Zero

	var ops []Operacao
	for _, neg := range negs {
		taxa := totalTaxas.Div(n.ValorDasOperacoes).Mul(neg.ValorTotal).Round(2)
		somasTaxas = somasTaxas.Add(taxa)
		ops = append(ops, Operacao{
			Data:          n.DataPregao,
			Tipo:          lo.Ternary(neg.CV == "C", COMPRA, VENDA),
			Ticker:        neg.Titulo,
			Qtd:           neg.Qtd,
			ValorUnitario: neg.ValorUnitario,
			Taxas:         taxa,
		})
	}

	if !totalTaxas.Equal(somasTaxas) {
		slog.Error("Erro ao distribuir taxas proporcionalmente entre ativos",
			slog.String("taxas", totalTaxas.String()),
			slog.String("soma", somasTaxas.String()),
			slog.String("formula", "TotalTaxas / ValorDasOperacoes * ValorTotalAtivo"),
		)
	}

	return ops, nil
}
