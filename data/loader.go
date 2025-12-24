package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	var ooc []OperacaoOuConfig
	for line := range bytes.SplitSeq(d, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.HasPrefix(line, []byte("//")) {
			continue
		}
		var tipo struct {
			Tipo string
		}
		if err := json.Unmarshal(line, &tipo); err != nil {
			return nil, fmt.Errorf("error: %s; content: %s", err, line)
		}
		switch tipo.Tipo {
		case "Config":
			var v ConfigOpcional
			if err := json.Unmarshal(line, &v); err != nil {
				return nil, fmt.Errorf("error: %s; content: %s", err, line)
			}
			ooc = append(ooc, OperacaoOuConfig{Data: v.Data, Cfg: v})
		default:
			var v Operacao
			if err := json.Unmarshal(line, &v); err != nil {
				return nil, fmt.Errorf("error: %s; content: %s", err, line)
			}
			ooc = append(ooc, OperacaoOuConfig{Data: v.Data, Opr: v})
		}
	}
	return ooc, nil
}
