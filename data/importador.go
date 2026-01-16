package data

import (
	"encoding/json/v2"
	"io"
	"iter"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/ofabricio/rv/pdf"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

func ImportarNotasDir(dir string, w io.Writer) error {

	var inp ImportadorNotas
	ops, err := inp.ImportarDir(dir)
	if err != nil {
		return err
	}

	for _, op := range ops {
		if err := json.MarshalWrite(w, &op); err != nil {
			return err
		}
		w.Write([]byte("\n"))
	}

	return nil
}

func ImportarNota(pdfData []byte, w io.Writer) error {

	var inp ImportadorNotas

	ops, err := inp.ImportarNota(pdfData)
	if err != nil {
		return err
	}

	for _, op := range ops {
		if err := json.MarshalWrite(w, &op); err != nil {
			return err
		}
		w.Write([]byte("\n"))
	}

	return nil
}

type ImportadorNotas struct{}

func (t *ImportadorNotas) ImportarDir(dir string) ([]Operacao, error) {

	files, err := t.getPDFFiles(dir)
	if err != nil {
		return nil, err
	}

	var oprs []Operacao
	for file := range files {

		d, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		ops, err := t.ImportarNota(d)
		if err != nil {
			return nil, err
		}

		oprs = append(oprs, ops...)
	}

	slices.SortStableFunc(oprs, func(a, b Operacao) int { return a.Data.Compare(b.Data) })

	return oprs, nil
}

func (t *ImportadorNotas) ImportarNota(pdfData []byte) ([]Operacao, error) {

	nota, err := pdf.ImportarNota(pdfData)
	if err != nil {
		return nil, err
	}

	return t.processar(nota), nil
}

func (t *ImportadorNotas) processar(n pdf.Nota) []Operacao {

	totalTaxas := n.TotalLiquido.Sub(n.ValorLiquidoDasOperacoes)
	somasTaxas := decimal.Zero

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

	var ops []Operacao
	for _, neg := range negs {
		tp := COMPRA
		if neg.CV == "V" {
			tp = VENDA
		}
		taxa := totalTaxas.Div(n.ValorDasOperacoes).Mul(neg.ValorTotal).Round(2)
		somasTaxas = somasTaxas.Add(taxa)
		ops = append(ops, Operacao{
			Data:          n.DataPregao,
			Tipo:          tp,
			Ticker:        neg.Titulo,
			Qtd:           neg.Qtd,
			ValorUnitario: neg.ValorUnitario,
			Taxas:         taxa,
		})
	}

	if !totalTaxas.Equal(somasTaxas) {
		slog.Error("Erro ao distribuir taxas proporcionalmente entre ativos", slog.String("taxas", totalTaxas.String()), slog.String("soma", somasTaxas.String()))
	}

	return ops
}

func (t *ImportadorNotas) getPDFFiles(dir string) (iter.Seq[string], error) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	it := func(yield func(string) bool) {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".pdf") || strings.HasSuffix(entry.Name(), ".PDF") {
				if !yield(entry.Name()) {
					return
				}
			}
		}
	}

	return it, nil
}
