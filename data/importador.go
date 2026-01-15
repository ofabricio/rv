package data

import (
	"encoding/json/v2"
	"io"
	"iter"
	"log/slog"
	"os"
	"strings"

	"github.com/ofabricio/rv/pdf"
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

	totalTaxas := n.TotalLiquido.Sub(n.ComprasAVista)
	somasTaxas := decimal.Zero

	var ops []Operacao
	for _, neg := range n.Negociacoes {
		if neg.CV == "C" {
			taxa := totalTaxas.Div(n.ComprasAVista).Mul(neg.ValorTotal).Round(2)
			somasTaxas = somasTaxas.Add(taxa)
			ops = append(ops, Operacao{
				Data:          n.DataPregao,
				Tipo:          COMPRA,
				Ticker:        neg.Titulo,
				Qtd:           neg.Qtd,
				ValorUnitario: neg.ValorUnitario,
				Taxas:         taxa,
			})
		}
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
