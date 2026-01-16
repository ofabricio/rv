package pdf

func ParseInter(content string) (Nota, error) {

	BNF := `
		root =
			FIND(data):DataPregao
			FIND(GROUP(negociacao+)):Negociacoes
			FIND(data):LiquidoPara
			FIND(preco):TotalLiquido
			FIND(preco):TotalCBLC
			FIND(preco):ValorLiquidoDasOperacoes
			FIND(preco):TaxaLiquidacao
			FIND(preco):TaxaRegistro
			FIND(preco):Emolumentos
			FIND(preco):TaxaANA
			FIND(preco):TaxaTermoOpcoes
			FIND(preco):TotalBovespaSoma
			FIND(preco):TotalCustosDespesas
			FIND(preco):Outros
			FIND(preco):IRRF
			FIND(preco):IRRFBase
			FIND(preco):Impostos
			FIND(preco):Execucao
			FIND(preco):TaxaCustodia
			FIND(preco):TaxaOperacional
			FIND(preco):ValorOperTitulosPubl
			FIND(preco):ValorDasOperacoes
			FIND(preco):OperacoesATermo
			FIND(preco):OpcoesCompras
			FIND(preco):OpcoesVendas
			FIND(preco):Debentures
			FIND(preco):ComprasAVista
			FIND(preco):VendasAVista
			(
			     FIND(preco):PIS FIND(preco):COFINS FIND(preco):TaxaTransfAtivos
			  |  FIND(preco):PIS FIND(preco):COFINS FIND(preco):TaxaTransfAtivos?
			  | (FIND(preco):PIS FIND(preco):COFINS FIND(preco):TaxaTransfAtivos)?
			)

		negociacao  = negociacao1 | negociacao2
		negociacao1 = GROUP(q S neg S merc S prazo S cv:CV S qtd:Qtd S preco:ValorUnitario S preco:ValorTotal S cd:CD S titulo:Titulo)
		negociacao2 = GROUP(q S neg S merc S cv:CV S qtd:Qtd? S preco:ValorUnitario S preco:ValorTotal S cd:CD S titulo:Titulo)
				  q = TEXT()
			    neg = 'Bovespa'
			   merc = 'VIS'
			  prazo = TEXT()
			 titulo = JOIN(REVERSE(ONPN S TEXT(' ') nome))
			   nome = MATCH(ANYNOT(NL)+)
			   ONPN = JOIN( ( 'PN' | 'ON' | 'UNT' ) S TEXT(' ') ('ED' S TEXT(' '))? ( 'N2' | 'N1' | 'NM' )
						  | 'CI' )
			    obs = '@'
			    qtd = '\d+'r
				 cv = 'C' | 'V'
				 cd = 'C' | 'D'

		data  = '\d{2}\/\d{2}\/\d{4}'r
		preco = JOIN(( '\d+'r '.'i? )+ ','i TEXT('.') '\d+'r)
		S = WS*
	`
	return parse(content, BNF)
}
