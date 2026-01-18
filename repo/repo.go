package repo

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ofabricio/bnf"
)

func GetTickerInfo(ticker string) (string, error) {

	res, err := http.Get(fmt.Sprintf("https://www.google.com/finance/quote/%s:BVMF", ticker))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	d, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	b := bnf.Compile(`
		r = FIND('R$'i v)
		v = '\d+\.\d+'r
	`)
	v := bnf.Parse(b, string(d))
	return v.Text, nil
}

func GetTickerInfoInv10(ticker string) (string, error) {

	res, err := http.Get("https://investidor10.com.br/acoes/" + ticker)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	d, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	b := bnf.Compile(`
		r = FIND('_card cotacao')i FIND(v)
		v = JOIN('\d+'r ','i TEXT('.') '\d+'r)
	`)
	v := bnf.Parse(b, string(d))
	return v.Text, nil
}
