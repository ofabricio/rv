# Regras

Este arquivo contém as regras de cálculo e declaração de IR de cada operação.

## Ações

TODO

## Opções

1. Não há isenção de IR, sempre paga IR de 15% (swing trade) ou 20% (day trade) sobre o lucro independente do valor.
1. Só se deve realizar a verificação de lucro/prejuízo após o encerramento/vencimento da operação. E tem até o último dia útil do mês seguinte para efetuar o pagamento do IR.
1. Como calcular o lucro (lembrar de descontar as taxas nos lucros):
    1. Na VENDA DE PUT:
        1. Exercido:
            1. PM = `(PM - PrêmioUnitário)`
            1. Lucro = `(Strike - PM) * Qtd`. **Lucro isento de IR**.
        1. Se sair da operação (COMPRANDO SUA PUT):  Prêmio $100, Custo $15 = Lucro $85. Pagar IR.
        1. Virou pó: Pagar IR sobre o prêmio.
    1. Na VENDA DE CALL:
        1. Exercido:
            1. PM = `(PM - PrêmioUnitário)`
            1. Lucro = `(Strike - PM) * Qtd`. Pagar IR.
        1. Se sair da operação (COMPRANDO SUA CALL): Prêmio $100, Custo $15 = Lucro $85. Pagar IR.
        1. Virou pó: Pagar IR sobre o prêmio.
    1. Na COMPRA DE PUT ou CALL:
        1. Viou pó: Prejuízo. Sem lucro. Sem IR. Prejuízo a compensar.
        1. Se sair da operação (VENDENDO SUA CALL/PUT): Gastei $100, Ganhei $200, Lucro: $100. IR sobre o lucro.
        1. Exercido:
            1. PM = `(PM + PrêmioUnitário)`
            1. Lucro = `(Strike - PM) * Qtd`. Pagar IR.
1. Opções vendidas (VENDA DE PUT ou CALL) devem ser declaradas em **Dívida e Ônus Reais** sob o **Código 16**, **APENAS** se o **encerramento/vencimento** for após a virada do ano. Exemplo: vendi em 20/12, mas encerrei em 20/01.
1. Opções compradas (COMPRA DE PUT ou CALL) devem ser declaradas em **Bens e Direitos** sob o **Grupo 04 e Código 04**, **APENAS** se o **encerramento/vencimento** for após a virada do ano. Exemplo: comprei em 20/12, mas encerrei em 20/01.
1. Somar lucros ou prejuízos (`valor_venda - valor_compra - taxa`) para cada operação no **mês** e declarar em **Operações Comuns/Day Trade**. Pagar IR sobre esse montante de lucro no mês.
1. Código DARF: 6015
