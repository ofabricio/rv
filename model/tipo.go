package model

const (
	COMPRA               TipoOpr = "COMPRA"
	VENDA                TipoOpr = "VENDA"
	BONIFICACAO          TipoOpr = "BONIFICAÇÃO"
	DESDOBRAMENTO        TipoOpr = "DESDOBRAMENTO"
	GRUPAMENTO           TipoOpr = "GRUPAMENTO"
	LEILAO_FRACAO        TipoOpr = "LEILÃO FRAÇÃO"
	DIVIDENDOS           TipoOpr = "DIVIDENDOS"
	JSCP                 TipoOpr = "JSCP"
	REND_TRIB            TipoOpr = "REND. TRIB."
	REDUCAO_CAPITAL      TipoOpr = "RED. CAP."
	SUBSCRICAO_COMPRA    TipoOpr = "SUBSCRIÇÃO COMPRA"
	SUBSCRICAO_VENDA     TipoOpr = "SUBSCRIÇÃO VENDA"
	SUBSCRICAO_EXERCICIO TipoOpr = "SUBSCRIÇÃO EXERCÍCIO"
	PUT_VENDA            TipoOpr = "VENDA DE PUT"
	PUT_VENDA_EX         TipoOpr = "VENDA DE PUT EX"
	PUT_VENDA_NE         TipoOpr = "VENDA DE PUT NE"
	PUT_COMPRA           TipoOpr = "COMPRA DE PUT"
	PUT_COMPRA_EX        TipoOpr = "COMPRA DE PUT EX"
	PUT_COMPRA_NE        TipoOpr = "COMPRA DE PUT NE"
	CALL_VENDA           TipoOpr = "VENDA DE CALL"
	CALL_VENDA_EX        TipoOpr = "VENDA DE CALL EX"
	CALL_VENDA_NE        TipoOpr = "VENDA DE CALL NE"
	CALL_COMPRA          TipoOpr = "COMPRA DE CALL"
	CALL_COMPRA_EX       TipoOpr = "COMPRA DE CALL EX"
	CALL_COMPRA_NE       TipoOpr = "COMPRA DE CALL NE"
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
