package data

import (
	"github.com/shopspring/decimal"
)

const (
	GrupoParticipacaoSocietaria = "03"
	GrupoAplicacaoInvestimento  = "04"
	GrupoOutros                 = "99"

	// Bens e Direitos.
	CodigoAcoes        = "01"
	CodigoOpcoes       = "04"
	CodigoJSCPNaoPagos = "07"

	// Rendimento Isento e Não Tributável.
	CodigoDividendos    = "09"
	CodigoBonificacao   = "18"
	CodigoIsencaoAte20k = "20"
	CodigoRINTOutros    = "99"

	// Rendimento Sujeito à Tributação Exclusiva/Definitiva.
	CodigoRendimentos = "06"
	CodigoJSCP        = "10"
	CodigoRSTEOutros  = "99"

	// Dívida e Ônus Reais.
	CodigoOutrasDividasOnusReais = "16"
)

func GetDefaultConfig(year int, def Config) Config {
	if v, ok := DefaultConfigByYear[year]; ok {
		return v
	}
	return def
}

func GetDefaultTonfig(year int, def TonfigMap) TonfigMap {
	if v, ok := DefaultTonfigByYear[year]; ok {
		return v
	}
	return def
}

var DefaultConfigByYear = map[int]Config{}
var DefaultTonfigByYear = map[int]TonfigMap{}

type Config struct {
	// Determina a estratégia de Preço Médio usada na Bonificação.
	//
	// Se true, as Bonificações vão alterar seu preço médio.
	// É como se você tivesse pagado o valor que recebeu de bonificação.
	// Significa que na hora de vender você pagará menos impostos.
	//
	// Se false (padrão), as Bonificações NÃO vão alterar seu preço médio.
	// É como se você tivesse recebido as ações de bonificação de graça.
	// Significa que na hora de vender você pagará mais impostos.
	//
	// Não há consenso sobre qual estratégia é a correta.
	// Boatos afirmam que a maioria prefere a estratégia que reduz o imposto a pagar (true).
	//
	// Para nunca ter problemas com a Receita Federal use false, pois você pagará mais impostos.
	//
	// Nota do dev:
	// Não me faz muito sentido não alterar o preço médio, pois se fizermos uma compra
	// de 1 única cota após uma bonificação o preço médio já será reajustado novamente.
	// Então parece que não alterar o preço médio teria diferença apenas ao fazer uma
	// venda imediatamente após a bonificação. E se é esse o caso, poderíamos resolver
	// comprando uma única cota antes de vender. Não está claro essa lógica.
	AlterarPrecoMedioNaBonificacao bool

	LimiteVendaIsenta decimal.Decimal
	SwingTradeIR      decimal.Decimal
	SwingTradeIRRF    decimal.Decimal // Dedo Duro.
	DayTradeIR        decimal.Decimal
	DayTradeIRRF      decimal.Decimal // Dedo Duro.
}

type TonfigMap map[Tipo]Tonfig

type Tonfig struct {
	LucroTributavel                      bool
	PrejuizoAbativel                     bool
	BensDireitos                         GrupoCodigo
	DividaOnusReais                      GrupoCodigo
	RendimentoIsentoNaoTributavel        GrupoCodigo
	RendimentoSujeitoTributacaoExclusiva GrupoCodigo
}

func (c Tonfig) IsBemDireito() bool {
	return c.BensDireitos.Grupo != ""
}

func (c Tonfig) IsDividaOnusReais() bool {
	return c.DividaOnusReais.Codigo != ""
}

func (c Tonfig) IsRendimentoIsentoNaoTributavel() bool {
	return c.RendimentoIsentoNaoTributavel.Codigo != ""
}

func (c Tonfig) IsRendimentoSujeitoTributacaoExclusiva() bool {
	return c.RendimentoSujeitoTributacaoExclusiva.Codigo != ""
}

func (c Tonfig) IsLimiteIsentoAplicavel() bool {
	return c.RendimentoIsentoNaoTributavel.Codigo == CodigoIsencaoAte20k
}

type GrupoCodigo struct {
	Grupo   string
	Codigo  string
	Descr   string
	Agregar bool
}

func (g GrupoCodigo) ID() string {
	return g.Grupo + g.Codigo
}

var DefaultConfig2025 = Config{
	AlterarPrecoMedioNaBonificacao: true,
	LimiteVendaIsenta:              decimal.RequireFromString("20000"),
	SwingTradeIR:                   decimal.RequireFromString("0.15"),
	SwingTradeIRRF:                 decimal.RequireFromString("0.00005"),
	DayTradeIR:                     decimal.RequireFromString("0.20"),
	DayTradeIRRF:                   decimal.RequireFromString("0.01"),
}

var DefaultTonfig2025 = TonfigMap{
	COMPRA: {
		BensDireitos: GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
	},
	VENDA: {
		BensDireitos:                  GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
		RendimentoIsentoNaoTributavel: GrupoCodigo{Codigo: CodigoIsencaoAte20k, Descr: "Isenção até R$ 20000", Agregar: true},
		PrejuizoAbativel:              true,
		LucroTributavel:               true,
	},
	BONIFICACAO: {
		BensDireitos:                  GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
		RendimentoIsentoNaoTributavel: GrupoCodigo{Codigo: CodigoBonificacao, Descr: "Bonificação"},
	},
	DESDOBRAMENTO: {
		BensDireitos: GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
	},
	GRUPAMENTO: {
		BensDireitos: GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
	},
	LEILAO_FRACAO: {
		RendimentoIsentoNaoTributavel: GrupoCodigo{Codigo: CodigoRINTOutros, Descr: "Leilão de Fração"},
		PrejuizoAbativel:              true,
	},
	DIVIDENDOS: {
		RendimentoIsentoNaoTributavel: GrupoCodigo{Codigo: CodigoDividendos, Descr: "Dividendos"},
	},
	REEMBOLSO: {
		RendimentoIsentoNaoTributavel: GrupoCodigo{Codigo: CodigoRINTOutros, Descr: "Reembolso de proventos", Agregar: true},
	},
	JSCP: {
		RendimentoSujeitoTributacaoExclusiva: GrupoCodigo{Codigo: CodigoJSCP, Descr: "JSCP"},
	},
	JSCP_NAO_PAGO: {
		BensDireitos:                         GrupoCodigo{Grupo: GrupoOutros, Codigo: CodigoJSCPNaoPagos},
		RendimentoSujeitoTributacaoExclusiva: GrupoCodigo{Codigo: CodigoJSCP, Descr: "JSCP Não Pago"},
	},
	REND_TRIB: {
		RendimentoSujeitoTributacaoExclusiva: GrupoCodigo{Codigo: CodigoRSTEOutros, Descr: "Rendimento Tributável"},
	},
	ALUGUEL: {
		RendimentoSujeitoTributacaoExclusiva: GrupoCodigo{Codigo: CodigoRendimentos, Descr: "Aluguel de Ações"},
	},
	REDUCAO_CAPITAL: {
		BensDireitos: GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
		// RendimentoSujeitoTributacaoExclusiva: GrupoCodigo{Codigo: CodigoRSTEOutros, Descr: "Redução do Capital - Restituição em Espécie"},
		PrejuizoAbativel: true,
	},
	SUBSCRICAO_COMPRA: {},
	SUBSCRICAO_VENDA: {
		LucroTributavel:  true,
		PrejuizoAbativel: true,
	},
	SUBSCRICAO_COMPRA_EX: {
		BensDireitos: GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
	},
	PUT_VENDA: {
		DividaOnusReais: GrupoCodigo{Codigo: CodigoOutrasDividasOnusReais},
		LucroTributavel: false,
	},
	PUT_VENDA_EX: {
		BensDireitos:     GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
		PrejuizoAbativel: true,
	},
	PUT_VENDA_NE: {
		LucroTributavel: true,
	},
	PUT_COMPRA: {
		BensDireitos:     GrupoCodigo{Grupo: GrupoAplicacaoInvestimento, Codigo: CodigoOpcoes},
		PrejuizoAbativel: false,
	},
	PUT_COMPRA_EX: {
		BensDireitos:     GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
		LucroTributavel:  true,
		PrejuizoAbativel: true,
	},
	PUT_COMPRA_NE: {
		PrejuizoAbativel: true,
	},
	CALL_VENDA: {
		DividaOnusReais: GrupoCodigo{Codigo: CodigoOutrasDividasOnusReais},
		LucroTributavel: false,
	},
	CALL_VENDA_EX: {
		BensDireitos:     GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
		LucroTributavel:  true,
		PrejuizoAbativel: true,
	},
	CALL_VENDA_NE: {
		LucroTributavel: true,
	},
	CALL_COMPRA: {
		BensDireitos:     GrupoCodigo{Grupo: GrupoAplicacaoInvestimento, Codigo: CodigoOpcoes},
		PrejuizoAbativel: false,
	},
	CALL_COMPRA_EX: {
		BensDireitos:     GrupoCodigo{Grupo: GrupoParticipacaoSocietaria, Codigo: CodigoAcoes},
		LucroTributavel:  false,
		PrejuizoAbativel: false,
	},
	CALL_COMPRA_NE: {
		PrejuizoAbativel: true,
	},
}
