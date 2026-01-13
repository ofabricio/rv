package data

import (
	"encoding/json/v2"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/ofabricio/rv/pdf"
	"github.com/shopspring/decimal"
)

func ImportarNotas(w io.Writer) error {
	var inp ImportadorNotas
	ops, err := inp.Importar(".")
	if err != nil {
		return err
	}
	for _, op := range ops {
		_ = json.MarshalWrite(w, &op)
		w.Write([]byte("\n"))
	}
	return nil
}

type ImportadorNotas struct{}

func (t *ImportadorNotas) Importar(dir string) ([]Operacao, error) {

	var ops []Operacao
	files, err := t.getPDFFiles(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		n, err := pdf.ParseNota(file)
		if err != nil {
			return nil, err
		}
		ops = append(ops, t.processaNota(n)...)
	}

	return ops, nil
}

func (t *ImportadorNotas) getPDFFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".pdf") || strings.HasSuffix(entry.Name(), ".PDF") {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func (t *ImportadorNotas) processaNota(n pdf.Nota) []Operacao {

	totalTaxas := n.TotalLiquido.Sub(n.ComprasAVista)

	sumTaxas := decimal.Zero
	var ops []Operacao
	for _, neg := range n.Negociacoes {
		if neg.CV == "C" {
			taxa := totalTaxas.Div(n.ComprasAVista).Mul(neg.ValorTotal).Round(2)
			sumTaxas = sumTaxas.Add(taxa)
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

	if !totalTaxas.Equal(sumTaxas) {
		slog.Error("Erro ao distribuir taxas proporcionalmente entre ativos", slog.String("taxas", totalTaxas.String()), slog.String("soma", sumTaxas.String()))
	}

	return ops
}
