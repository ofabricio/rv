package main

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/ofabricio/rv/data"
)

func main() {
	updateTests()
}

// Script para atualizar todos os arquivos de resultados de testes.
//
// Este script lê os arquivos .give na pasta data/testdata,
// gera a saída e grava no arquivo .then correspondente.
//
// Para rodar o script execute na pasta raiz do projeto:
//
//	go run ./script/main.go
func updateTests() {
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
