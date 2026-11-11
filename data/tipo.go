package data

import (
	"encoding/json/v2"
	"time"
)

const (
	COMPRA               Tipo = "Compra"
	VENDA                Tipo = "Venda"
	BONIFICACAO          Tipo = "Bonificação"
	DESDOBRAMENTO        Tipo = "Desdobramento"
	GRUPAMENTO           Tipo = "Grupamento"
	LEILAO_FRACAO        Tipo = "Leilão Fração"
	DIVIDENDOS           Tipo = "Dividendos"
	REEMBOLSO            Tipo = "Reembolso"
	JSCP                 Tipo = "JSCP"
	JSCP_NAO_PAGO        Tipo = "JSCP Não Pago"
	REND_TRIB            Tipo = "Rend. Trib."
	ALUGUEL              Tipo = "Aluguel"
	REDUCAO_CAPITAL      Tipo = "Red. Cap."
	SUBSCRICAO_COMPRA    Tipo = "Compra Subscrição"
	SUBSCRICAO_VENDA     Tipo = "Venda Subscrição"
	SUBSCRICAO_COMPRA_EX Tipo = "Compra Subscrição (EX)"
	PUT_VENDA            Tipo = "Venda PUT"
	PUT_VENDA_EX         Tipo = "Venda PUT (EX)"
	PUT_VENDA_NE         Tipo = "Venda PUT (NE)"
	PUT_COMPRA           Tipo = "Compra PUT"
	PUT_COMPRA_EX        Tipo = "Compra PUT (EX)"
	PUT_COMPRA_NE        Tipo = "Compra PUT (NE)"
	CALL_VENDA           Tipo = "Venda CALL"
	CALL_VENDA_EX        Tipo = "Venda CALL (EX)"
	CALL_VENDA_NE        Tipo = "Venda CALL (NE)"
	CALL_COMPRA          Tipo = "Compra CALL"
	CALL_COMPRA_EX       Tipo = "Compra CALL (EX)"
	CALL_COMPRA_NE       Tipo = "Compra CALL (NE)"
)

type Tipo string

type DateOnly time.Time

func (d DateOnly) V() time.Time {
	return time.Time(d)
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d *DateOnly) MarshalJSON() ([]byte, error) {
	data := d.V().Format(time.DateOnly)
	return json.Marshal(data)
}
