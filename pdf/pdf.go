package pdf

import (
	"bytes"
	"io"

	"github.com/ledongthuc/pdf"
)

func ImportarNota(pdfData []byte) (Nota, error) {

	v, err := ReadContent(bytes.NewReader(pdfData), len(pdfData))
	if err != nil {
		return Nota{}, err
	}

	nota, err := ParseNota(v)
	if err != nil {
		return Nota{}, err
	}

	return nota, nil
}

func ReadContent(at io.ReaderAt, size int) (string, error) {

	r, err := pdf.NewReader(at, int64(size))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
