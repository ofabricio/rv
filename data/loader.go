package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"slices"
	"time"
)

func ReadOperacoes(r io.Reader) ([]OperacaoDesconsolidada, error) {
	lines, err := ReadLinesSorted(r)
	if err != nil {
		return nil, err
	}
	prevYear := 1900
	ooc := make([]OperacaoDesconsolidada, 0, 128)
	cfg := DefaultConfig2025
	tfg := DefaultTonfig2025
	for line := range lines {
		var tipo struct {
			Tipo string
			Data time.Time `json:",format:DateOnly"`
		}
		if err := json.Unmarshal(line, &tipo); err != nil {
			return nil, fmt.Errorf("error: %s; content: %s", err, line)
		}
		// Na mudança de ano atualiza para a configuração correspondente.
		// Isso sobrescreve todas as alterações de configurações definidas
		// no arquivo fonte, de modo que é necessário repeti-las todo novo
		// ano para que continuem tendo efeito.
		if currYear := tipo.Data.Year(); currYear != prevYear {
			cfg = GetDefaultConfig(currYear, cfg)
			tfg = GetDefaultTonfig(currYear, tfg)
			prevYear = currYear
		}
		switch tipo.Tipo {
		case "Config":
			if err := json.Unmarshal(line, &cfg); err != nil {
				return nil, fmt.Errorf("error: %s; content: %s", err, line)
			}
		default:
			var opr Operacao
			if err := json.Unmarshal(line, &opr); err != nil {
				return nil, fmt.Errorf("error: %s; content: %s", err, line)
			}
			ooc = append(ooc, OperacaoDesconsolidada{Opr: opr, Cfg: cfg, Tfg: tfg[opr.Tipo]})
		}
	}
	return ooc, nil
}

func ReadLinesSorted(r io.Reader) (iter.Seq[[]byte], error) {

	lines, err := ReadLines(r)
	if err != nil {
		return nil, err
	}

	type Line struct {
		Data time.Time
		Line []byte
	}

	var data []Line
	for line := range lines {
		var v struct {
			Data time.Time `json:",format:DateOnly"`
		}
		if err := json.Unmarshal(line, &v); err != nil {
			return nil, fmt.Errorf("error: %s; content: %s", err, line)
		}
		data = append(data, Line{v.Data, line})
	}

	slices.SortStableFunc(data, func(a, b Line) int { return a.Data.Compare(b.Data) })

	it := func(yield func([]byte) bool) {
		for _, line := range data {
			if !yield(line.Line) {
				return
			}
		}
	}

	return it, nil
}

func ReadLines(r io.Reader) (iter.Seq[[]byte], error) {

	d, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	it := func(yield func([]byte) bool) {
		for line := range bytes.SplitSeq(d, []byte("\n")) {
			line = bytes.TrimSpace(line)
			if len(line) == 0 || bytes.HasPrefix(line, []byte("//")) {
				continue
			}
			if !yield(line) {
				return
			}
		}
	}

	return it, nil
}
