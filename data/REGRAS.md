# Regras

Este arquivo contém as regras de cálculo e declaração de IR de cada operação.

## Ações

- Bonificação
    - Deve ser declarado em **RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS**; a quantidade deve ser declarada na discriminação de **BENS E DIREITOS**.

- Dedo-Duro (IRRF)
    - Nas vendas abaixo de 20 mil o dedo-duro não é recolhido, você pode verificar isso somando os valores nas notas de corretagem.
    Ele aparece na nota, mas é desconsiderado quando o valor de negociação é inferior a 20 mil.
    Caso o valor das vendas brutas no mês ultrapasse esse montante a corretora desconta ele da conta.
    Logo não devemos informar o dedo-duro em meses onde a soma das vendas seja inferior a 20 mil no Swing Trade.
    ([ref](https://www.youtube.com/watch?v=31zeJEf4bo4&lc=Ugyv0R7MepCxtWgipdZ4AaABAg.9hNmE_x5Hx19hO87rLk4g5))
    - Resumo
        1. Vendas no mês <= 20 mil: não descontar IRRF.
        1. Vendas no mês > 20 mil com prejuízo: restituir o IRRF.
        1. Vendas no mês > 20 mil com lucro: descontar o IRRF no lucro.

## Opções

1. Não há isenção de IR, sempre paga IR de 15% (swing trade) ou 20% (day trade) sobre o lucro independente do valor.
1. Só se deve realizar a verificação de lucro/prejuízo após o encerramento/vencimento da operação. E tem até o último dia útil do mês seguinte para efetuar o pagamento do IR.
1. Como calcular lucro (lembrar de descontar as taxas nos lucros):
    - VENDA DE PUT
        1. Exercido `Lucro = (valor da ação no dia do vencimento - strike) + prêmio` **Lucro isento de IR**.
        2. Não Exercido `Lucro = prêmio`. Pagar IR.
        3. Saindo da operação (COMPRANDO SUA PUT): Prêmio $100, Custo $15 = Lucro $85. Pagar IR.
    - COMPRA DE PUT
        1. Exercido `Lucro = strike - (valor da ação no dia do vencimento + custo)`. Pagar IR.
        2. Não Exercido `Custo`. Abatível no IR.
        3. Saindo da operação (VENDENDO SUA PUT): Gastei $100, Ganhei $200, Lucro: $100. Pagar IR.
    - VENDA DE CALL
        1. Exercido `Lucro = (strike + prêmio) - valor da ação no dia do vencimento`. Pagar IR.
        2. Não Exercido `Lucro = prêmio`. Pagar IR.
        3. Saindo da operação (COMPRANDO SUA CALL): Prêmio $100, Custo $15 = Lucro $85. Pagar IR.
    - COMPRA DE CALL
        1. Exercido `Lucro = valor da ação no dia do vencimento - (strike + prêmio)`. Pagar IR.
        2. Não Exercido `Custo`. Abatível no IR.
        3. Saindo da operação (VENDENDO SUA CALL): Gastei $100, Ganhei $200, Lucro: $100. Pagar IR.
1. Opções vendidas (VENDA DE PUT ou CALL) devem ser declaradas em **Dívida e Ônus Reais** sob o **Código 16**, **APENAS** se o **encerramento/vencimento** for após a virada do ano. Exemplo: vendi em 20/12, mas encerrei em 20/01.
1. Opções compradas (COMPRA DE PUT ou CALL) devem ser declaradas em **Bens e Direitos** sob o **Grupo 04 e Código 04**, **APENAS** se o **encerramento/vencimento** for após a virada do ano. Exemplo: comprei em 20/12, mas encerrei em 20/01.
1. Somar lucros ou prejuízos (`valor_venda - valor_compra - taxa`) para cada operação no **mês** e declarar em **Operações Comuns/Day Trade**. Pagar IR sobre esse montante de lucro no mês.
1. Código DARF: 6015
