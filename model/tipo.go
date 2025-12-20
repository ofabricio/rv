package model

const (
	COMPRA               TipoOpr = "Compra"
	VENDA                TipoOpr = "Venda"
	BONIFICACAO          TipoOpr = "Bonificação"
	DESDOBRAMENTO        TipoOpr = "Desdobramento"
	GRUPAMENTO           TipoOpr = "Grupamento"
	LEILAO_FRACAO        TipoOpr = "Leilão Fração"
	DIVIDENDOS           TipoOpr = "Dividendos"
	JSCP                 TipoOpr = "JSCP"
	REND_TRIB            TipoOpr = "Rend. Trib."
	REDUCAO_CAPITAL      TipoOpr = "Red. Cap."
	SUBSCRICAO_COMPRA    TipoOpr = "Subscrição Compra"
	SUBSCRICAO_VENDA     TipoOpr = "Subscrição Venda"
	SUBSCRICAO_EXERCICIO TipoOpr = "Subscrição Exercício"
	PUT_VENDA            TipoOpr = "Venda de PUT"
	PUT_VENDA_EX         TipoOpr = "Venda de PUT EX"
	PUT_VENDA_NE         TipoOpr = "Venda de PUT NE"
	PUT_COMPRA           TipoOpr = "Compra de PUT"
	PUT_COMPRA_EX        TipoOpr = "Compra de PUT EX"
	PUT_COMPRA_NE        TipoOpr = "Compra de PUT NE"
	CALL_VENDA           TipoOpr = "Venda de CALL"
	CALL_VENDA_EX        TipoOpr = "Venda de CALL EX"
	CALL_VENDA_NE        TipoOpr = "Venda de CALL NE"
	CALL_COMPRA          TipoOpr = "Compra de CALL"
	CALL_COMPRA_EX       TipoOpr = "Compra de CALL EX"
	CALL_COMPRA_NE       TipoOpr = "Compra de CALL NE"
)

var tipoOprConfig = map[TipoOpr]TipoOprConfig{
	COMPRA:               {},
	VENDA:                {Codigo: "03 ── Isenção até 20k", IsRendimentoIsentoAte20k: true, IsRendimentoTributavelApos20k: true, IsPrejuizoAbativel: true},
	BONIFICACAO:          {Codigo: "18 ── Bonificação", IsRendimentoIsentoNaoTributavel: true},
	DESDOBRAMENTO:        {Codigo: "01"},
	GRUPAMENTO:           {Codigo: "01"},
	LEILAO_FRACAO:        {Codigo: "99 ── Leilão de Fração", IsRendimentoIsentoNaoTributavel: true, IsPrejuizoAbativel: true},
	DIVIDENDOS:           {Codigo: "09 ── Dividendos", IsRendimentoIsentoNaoTributavel: true},
	JSCP:                 {Codigo: "10 ── JSCP", IsRendimentoSujeitoTributacaoExclusiva: true},
	REND_TRIB:            {Codigo: "99 ── Rendimento Tributável", IsRendimentoSujeitoTributacaoExclusiva: true},
	REDUCAO_CAPITAL:      {Codigo: "99 ── Redução do Capital - Restituição em Espécie", IsPrejuizoAbativel: true},
	SUBSCRICAO_COMPRA:    {Codigo: "03"},
	SUBSCRICAO_VENDA:     {Codigo: "03", IsLucroTributavel: true, IsPrejuizoAbativel: true},
	SUBSCRICAO_EXERCICIO: {Codigo: "03"},
	PUT_VENDA:            {},
	PUT_VENDA_EX:         {IsPrejuizoAbativel: true},
	PUT_VENDA_NE:         {IsLucroTributavel: true},
	PUT_COMPRA:           {},
	PUT_COMPRA_EX:        {IsLucroTributavel: true, IsPrejuizoAbativel: true},
	PUT_COMPRA_NE:        {IsPrejuizoAbativel: true},
	CALL_VENDA:           {},
	CALL_VENDA_EX:        {IsLucroTributavel: true, IsPrejuizoAbativel: true},
	CALL_VENDA_NE:        {IsLucroTributavel: true},
	CALL_COMPRA:          {},
	CALL_COMPRA_EX:       {IsLucroTributavel: true, IsPrejuizoAbativel: true},
	CALL_COMPRA_NE:       {IsPrejuizoAbativel: true},
}

type TipoOprConfig struct {
	IsRendimentoIsentoNaoTributavel        bool   // Entra na guia RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS.
	IsRendimentoIsentoAte20k               bool   // Entra na guia RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS se até 20k.
	IsRendimentoSujeitoTributacaoExclusiva bool   // Entra na guia RENDIMENTOS SUJEITOS À TRIBUTAÇÃO EXCLUSIVA/DEFINITIVA.
	IsRendimentoTributavelApos20k          bool   // Entra na guia OPERAÇÕES COMUNS/DAY-TRADE quando for lucro de vendas maior que 20k.
	IsPrejuizoAbativel                     bool   // Entra na guia OPERAÇÕES COMUNS/DAY-TRADE quando for prejuízo.
	IsLucroTributavel                      bool   // Entra na guia OPERAÇÕES COMUNS/DAY-TRADE quando for lucro.
	Codigo                                 string // Código na guia.
}

type TipoOpr string
