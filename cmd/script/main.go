package main

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/ofabricio/rv/data"
)

func main() {

	// Script para atualizar todos os arquivos de testes.
	//
	// Este script lê os arquivos .give na pasta data/testdata,
	// gera a saída e grava no arquivo .then correspondente.
	//
	// Para rodar o script execute na pasta raiz do projeto:
	//   go run ./cmd/script/main.go

	target := "data/testdata"

	dir, err := os.ReadDir(target)
	if err != nil {
		panic(err)
	}

	for _, f := range dir {
		name, ext, _ := strings.Cut(f.Name(), ".")
		if ext == "give" {
			src := path.Join(target, name) + ".give"
			dst := path.Join(target, name) + ".then"

			var buf bytes.Buffer
			c := data.NewCarteira()
			c.Load(src)
			c.Print("table", &buf)

			if err := os.WriteFile(dst, bytes.TrimSpace(buf.Bytes()), 0644); err != nil {
				panic(err)
			}
		}
	}
}
