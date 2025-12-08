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
)

type TipoOpr string

var tipoOprConfig = map[TipoOpr]struct {
	IsRendimentoIsentoNaoTributavel        bool   // Entra na guia RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS.
	IsRendimentoIsentoAte20k               bool   // Entra na guia RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS se até 20k.
	IsRendimentoSujeitoTributacaoExclusiva bool   // Entra na guia RENDIMENTOS SUJEITOS À TRIBUTAÇÃO EXCLUSIVA/DEFINITIVA.
	IsBemDireito                           bool   // Entra na guia BENS E DIREITOS.
	IsRendimentoTributavelApos20k          bool   // Entra na guia OPERAÇÕES COMUNS/DAY-TRADE quando for lucro de vendas maior que 20k.
	IsPrejuizoAbativel                     bool   // Entra na guia OPERAÇÕES COMUNS/DAY-TRADE quando for prejuízo.
	IsLucroTributavel                      bool   // Entra na guia OPERAÇÕES COMUNS/DAY-TRADE quando for lucro.
	Grupo                                  string // Grupo na guia.
	Codigo                                 string // Código na guia.
}{
	COMPRA:               {Grupo: "03", Codigo: "03", IsBemDireito: true},
	VENDA:                {Grupo: "04", Codigo: "03 ── Isenção até 20k", IsBemDireito: true, IsRendimentoIsentoAte20k: true, IsRendimentoTributavelApos20k: true, IsPrejuizoAbativel: true},
	BONIFICACAO:          {Codigo: "18 ── Bonificação", IsBemDireito: true, IsRendimentoIsentoNaoTributavel: true},
	DESDOBRAMENTO:        {Codigo: "01", IsBemDireito: true},
	GRUPAMENTO:           {Codigo: "01", IsBemDireito: true},
	LEILAO_FRACAO:        {Codigo: "99 ── Leilão de Fração", IsRendimentoIsentoNaoTributavel: true, IsPrejuizoAbativel: true},
	DIVIDENDOS:           {Codigo: "09 ── Dividendos", IsRendimentoIsentoNaoTributavel: true},
	JSCP:                 {Codigo: "10 ── JSCP", IsRendimentoSujeitoTributacaoExclusiva: true},
	REND_TRIB:            {Codigo: "99 ── Rendimento Tributável", IsRendimentoSujeitoTributacaoExclusiva: true},
	REDUCAO_CAPITAL:      {Codigo: "99 ── Redução do Capital - Restituição em Espécie", IsBemDireito: true, IsPrejuizoAbativel: true},
	SUBSCRICAO_COMPRA:    {Codigo: "03"},
	SUBSCRICAO_VENDA:     {Codigo: "03", IsLucroTributavel: true, IsPrejuizoAbativel: true},
	SUBSCRICAO_EXERCICIO: {Codigo: "03", IsBemDireito: true},
}
