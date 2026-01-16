package pdf

func ParseClear(content string) (Nota, error) {

	BNF := `
		root =
			FIND(GROUP(negociacao+)):Negociacoes
			FIND(data):DataPregao
			FIND(preco):Debentures
			FIND(preco):VendasAVista
			FIND(preco):ComprasAVista
			FIND(preco):OpcoesCompras
			FIND(preco):OpcoesVendas
			FIND(preco):OperacoesATermo
			FIND(preco):ValorOperTitulosPubl
			FIND(preco):ValorDasOperacoes
			FIND(preco):TotalCBLC
			FIND(preco):ValorLiquidoDasOperacoes
			FIND(preco):TaxaLiquidacao
			FIND(preco):TaxaRegistro
			FIND(preco):TotalBovespaSoma
			FIND(preco):TaxaTermoOpcoes
			FIND(preco):TaxaANA
			FIND(preco):Emolumentos
			FIND(preco):TaxaTransfAtivos
			FIND(preco):TotalCustosDespesas
			FIND(preco):TaxaOperacional
			FIND(preco):Execucao
			FIND(preco):TaxaCustodia
			FIND(preco):Impostos
			FIND(preco):IRRF
			FIND(preco):Outros
			FIND(preco):IRRFBase
			FIND(preco):TotalLiquido
			FIND(data):LiquidoPara

		negociacao = GROUP( q S neg S cv:CV S merc S prazo S titulo:Titulo S obs S qtd:Qtd S preco:ValorUnitario S preco:ValorTotal S cd:CD )
				 q = TEXT()
			   neg = '1-BOVESPA'
			  merc = 'VISTA'
			 prazo = TEXT()
			titulo = JOIN(ONPN | '\w+'r S TEXT(' ') titulo)
			  ONPN = 'ON' S TEXT(' ') 'NM'
				   | 'PN' S TEXT(' ') 'EDJ' S TEXT(' ') 'N2'
			   obs = '@' '#'?
			   qtd = '\d+'r
				cv = 'C' | 'V'
				cd = 'C' | 'D'

		data  = '\d{2}\/\d{2}\/\d{4}'r
		preco = JOIN(( '\d+'r '.'i? )+ ','i TEXT('.') '\d+'r)
		S = WS*
	`
	return parse(content, BNF)
}
