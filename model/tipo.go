package model

const (
	COMPRA            TipoOpr = "Compra"
	VENDA             TipoOpr = "Venda"
	BONIFICACAO       TipoOpr = "Bonificação"
	DESDOBRAMENTO     TipoOpr = "Desdobramento"
	GRUPAMENTO        TipoOpr = "Grupamento"
	LEILAO_FRACAO     TipoOpr = "Leilão Fração"
	DIVIDENDOS        TipoOpr = "Dividendos"
	JSCP              TipoOpr = "JSCP"
	REND_TRIB         TipoOpr = "Rend. Trib."
	REDUCAO_CAPITAL   TipoOpr = "Red. Cap."
	SUBSCRICAO_COMPRA TipoOpr = "Subscrição Compra"
	SUBSCRICAO_VENDA  TipoOpr = "Subscrição Venda"
	SUBSCRICAO_EX     TipoOpr = "Subscrição (EX)"
	PUT_VENDA         TipoOpr = "Venda PUT"
	PUT_VENDA_EX      TipoOpr = "Venda PUT (EX)"
	PUT_VENDA_NE      TipoOpr = "Venda PUT (NE)"
	PUT_COMPRA        TipoOpr = "Compra PUT"
	PUT_COMPRA_EX     TipoOpr = "Compra PUT (EX)"
	PUT_COMPRA_NE     TipoOpr = "Compra PUT (NE)"
	CALL_VENDA        TipoOpr = "Venda CALL"
	CALL_VENDA_EX     TipoOpr = "Venda CALL (EX)"
	CALL_VENDA_NE     TipoOpr = "Venda CALL (NE)"
	CALL_COMPRA       TipoOpr = "Compra CALL"
	CALL_COMPRA_EX    TipoOpr = "Compra CALL (EX)"
	CALL_COMPRA_NE    TipoOpr = "Compra CALL (NE)"
)

var tipoOprConfig = map[TipoOpr]TipoOprConfig{
	COMPRA:            {},
	VENDA:             {Codigo: "03 ── Isenção até 20k", IsRendimentoIsentoAte20k: true, IsRendimentoTributavelApos20k: true, IsPrejuizoAbativel: true},
	BONIFICACAO:       {Codigo: "18 ── Bonificação", IsRendimentoIsentoNaoTributavel: true},
	DESDOBRAMENTO:     {Codigo: "01"},
	GRUPAMENTO:        {Codigo: "01"},
	LEILAO_FRACAO:     {Codigo: "99 ── Leilão de Fração", IsRendimentoIsentoNaoTributavel: true, IsPrejuizoAbativel: true},
	DIVIDENDOS:        {Codigo: "09 ── Dividendos", IsRendimentoIsentoNaoTributavel: true},
	JSCP:              {Codigo: "10 ── JSCP", IsRendimentoSujeitoTributacaoExclusiva: true},
	REND_TRIB:         {Codigo: "99 ── Rendimento Tributável", IsRendimentoSujeitoTributacaoExclusiva: true},
	REDUCAO_CAPITAL:   {Codigo: "99 ── Redução do Capital - Restituição em Espécie", IsPrejuizoAbativel: true},
	SUBSCRICAO_COMPRA: {Codigo: "03"},
	SUBSCRICAO_VENDA:  {Codigo: "03", IsLucroTributavel: true, IsPrejuizoAbativel: true},
	SUBSCRICAO_EX:     {Codigo: "03"},
	PUT_VENDA:         {IsLucroTributavel: false},
	PUT_VENDA_EX:      {IsPrejuizoAbativel: true},
	PUT_VENDA_NE:      {IsLucroTributavel: true},
	PUT_COMPRA:        {IsPrejuizoAbativel: false},
	PUT_COMPRA_EX:     {IsLucroTributavel: true, IsPrejuizoAbativel: true},
	PUT_COMPRA_NE:     {IsPrejuizoAbativel: true},
	CALL_VENDA:        {IsLucroTributavel: false},
	CALL_VENDA_EX:     {IsLucroTributavel: true, IsPrejuizoAbativel: true},
	CALL_VENDA_NE:     {IsLucroTributavel: true},
	CALL_COMPRA:       {IsPrejuizoAbativel: false},
	CALL_COMPRA_EX:    {IsLucroTributavel: false, IsPrejuizoAbativel: false},
	CALL_COMPRA_NE:    {IsPrejuizoAbativel: true},
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
