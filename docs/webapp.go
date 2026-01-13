//go:build js && wasm

package main

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"syscall/js"

	"github.com/ofabricio/rv/data"
)

func main() {
	js.Global().Set("rvRun", rvRun())
	js.Global().Set("rvNota", rvNota())
	<-make(chan struct{})
}

func rvRun() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {

		INP := args[0].String()
		ARG := args[1].String()

		// Apaga as flags previamente definidas.
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = append([]string{"rv"}, strings.Split(ARG, " ")...)

		var buf bytes.Buffer
		c := data.NewCarteira()
		if err := c.CommandLine(strings.NewReader(INP), &buf); err != nil {
			return err.Error()
		}

		return buf.String()
	})
}

func rvNota() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {

		pdfLen := args[0].Get("length").Int()
		pdfData := make([]byte, pdfLen)
		js.CopyBytesToGo(pdfData, args[0])

		var buf bytes.Buffer
		if err := data.ProcessarNotaPDF(pdfData, &buf); err != nil {
			return err.Error()
		}

		return buf.String()
	})
}
