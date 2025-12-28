package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

func LoadOperacoes(file string) ([]OperacaoOuConfig, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadOperacoes(f)
}

func ReadOperacoes(r io.Reader) ([]OperacaoOuConfig, error) {
	d, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	ooc := make([]OperacaoOuConfig, 0, 128)
	cfg := DefaultConfig
	for line := range bytes.SplitSeq(d, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.HasPrefix(line, []byte("//")) {
			continue
		}
		var tipo struct {
			Tipo string
			Data time.Time `json:",format:DateOnly"`
		}
		if err := json.Unmarshal(line, &tipo); err != nil {
			return nil, fmt.Errorf("error: %s; content: %s", err, line)
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
			ooc = append(ooc, OperacaoOuConfig{Data: tipo.Data, Opr: opr, Cfg: cfg})
		}
	}
	return ooc, nil
}
