# rv

Esta ferramenta faz o controle de operações de ações do mercado brasileiro,
gerando as guias de declaração de IR (Bens e Direitos, Operações Comuns/Day Trade, etc).

Se deseja usá-la **é por sua conta e risco**. :warning:

## Features

Operações suportadas:

- Ações
  - Compra
  - Venda
  - Bonificação
  - Grupamento
  - Desdobramento
  - Dividendos
  - Reembolso
  - JSCP
  - JSCP Não Pago
  - Leilão de Fração
  - Redução de Capital
  - Rend. Trib.
  - Subscrição
  - Aluguel

- Opções
  - Compra/Venda de PUT
  - Compra/Venda de CALL

## Build

Escolha abaixo uma das opções de compilação da ferramenta; um arquivo
executável `rv` aparecerá na pasta onde o comando for executado.

### Com Go

Tendo [Go](https://go.dev/dl) instalado:

```
GOBIN=$(pwd) GOEXPERIMENT=jsonv2 go install github.com/ofabricio/rv@latest
```

### Com Docker

Tendo [Docker](https://www.docker.com) instalado:

```sh
docker run --rm -v $(pwd):/src -w /src golang:alpine sh -c "
    GOOS=darwin GOARCH=arm64 GOEXPERIMENT=jsonv2 go install github.com/ofabricio/rv@latest && cp /go/bin/*/rv ."
```

**Nota:** altere `GOOS=darwin GOARCH=arm64` para os [valores](https://go.dev/doc/install/source#environment) correspondentes ao seu sistema operacional.

## Como usar

Basta adicionar as operações no arquivo [db.ndjson](/db.ndjson), que está na raiz deste repositório,
rodar a ferramenta `./rv`, e ver o resultado no terminal.
A ferramenta fará os cálculos necessários e mostrará o resultado final consolidado.
Veja um exemplo abaixo.

**Nota:** o arquivo `db.ndjson` precisa estar na mesma pasta do executável; copie o arquivo modelo deste repositório ou crie um arquivo novo.

No momento só é possível adicionar as operações editando manualmente o arquivo `db.ndjson`.
No futuro será criado uma interface no terminal ou no navegador para facilitar isso.
Importação de notas de corretagem ou extratos da B3 também estão nos planos.

Veja exemplos de uso de cada tipo de operação nos arquivos `*.give` do diretório [/data/testdata](/data/testdata).

## Exemplo de resultado

```
$ ./rv
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                             OPERAÇÕES COM AÇÕES                                                             │
├────┬────────┬────────────┬──────────────────────┬─────┬──────────┬──────────┬───────┬─────────┬──────────────┬───────┬───────────┬──────────┤
│ ID │ Ticker │    Data    │       Operação       │ Qtd │ V. Unit. │ V. Total │ Taxas │ Qtd Ac. │ V. Total Ac. │  PM   │ V. Compra │  Lucro   │
├────┼────────┼────────────┼──────────────────────┼─────┼──────────┼──────────┼───────┼─────────┼──────────────┼───────┼───────────┼──────────┤
│ 1  │ PETR4  │ 2024-01-01 │ Compra               │ 100 │    30,00 │  3000,00 │  0,00 │     100 │      3000,00 │ 30,00 │      0,00 │     0,00 │
│ 2  │ PETR4  │ 2024-02-01 │ Compra               │ 200 │    30,50 │  6100,00 │  0,00 │     300 │      9100,00 │ 30,33 │      0,00 │     0,00 │
│ 3  │ VALE3  │ 2025-03-01 │ Compra               │ 200 │    60,00 │ 12000,00 │  0,00 │     200 │     12000,00 │ 60,00 │      0,00 │     0,00 │
│ 4  │ VALE3  │ 2025-04-01 │ Desdobramento (2,00) │   0 │     0,00 │     0,00 │  0,00 │     400 │     12000,00 │ 30,00 │      0,00 │     0,00 │
│ 5  │ VALE3  │ 2025-05-01 │ Compra               │ 100 │    61,00 │  6100,00 │  0,00 │     500 │     18100,00 │ 36,20 │      0,00 │     0,00 │
│ 6  │ PETR4  │ 2025-06-01 │ Dividendos           │   0 │     0,00 │     0,00 │  0,00 │     300 │      9100,00 │ 30,33 │      0,00 │  1000,00 │
│ 7  │ VALE3  │ 2025-07-01 │ JSCP                 │   0 │     0,00 │     0,00 │  0,00 │     500 │     18100,00 │ 36,20 │      0,00 │  1000,00 │
│ 8  │ VALE3  │ 2025-08-01 │ Venda                │ 400 │    65,00 │ 26000,00 │  0,00 │     100 │      3620,00 │ 36,20 │  14480,00 │ 11520,00 │
└────┴────────┴────────────┴──────────────────────┴─────┴──────────┴──────────┴───────┴─────────┴──────────────┴───────┴───────────┴──────────┘
┌──────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                           BENS E DIREITOS                                            │
├────────┬──────────────────┬──────────────────┬─────────┬─────────────────────────────────────────────┤
│ Ticker │ Situação em 2023 │ Situação em 2024 │ Grp Cód │                Discriminação                │
├────────┼──────────────────┼──────────────────┼─────────┼─────────────────────────────────────────────┤
│ PETR4  │       0,00       │     9100,00      │  03 01  │ 300 AÇÕES PETR4 COM PREÇO MÉDIO DE R$ 30,33 │
└────────┴──────────────────┴──────────────────┴─────────┴─────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                           BENS E DIREITOS                                            │
├────────┬──────────────────┬──────────────────┬─────────┬─────────────────────────────────────────────┤
│ Ticker │ Situação em 2024 │ Situação em 2025 │ Grp Cód │                Discriminação                │
├────────┼──────────────────┼──────────────────┼─────────┼─────────────────────────────────────────────┤
│ PETR4  │     9100,00      │     9100,00      │  03 01  │ 300 AÇÕES PETR4 COM PREÇO MÉDIO DE R$ 30,33 │
├────────┼──────────────────┼──────────────────┼─────────┼─────────────────────────────────────────────┤
│ VALE3  │       0,00       │     3620,00      │  03 01  │ 100 AÇÕES VALE3 COM PREÇO MÉDIO DE R$ 36,20 │
└────────┴──────────────────┴──────────────────┴─────────┴─────────────────────────────────────────────┘
┌────────┬───────────────────────────────────────┐
│  2025  │ RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS │
├────────┼──────────────┬────────────────────────┤
│ Ticker │    Valor     │         Código         │
├────────┼──────────────┼────────────────────────┤
│ PETR4  │      1000,00 │ 09 ── Dividendos       │
└────────┴──────────────┴────────────────────────┘
┌────────┬─────────────────────────────────────────────┐
│  2025  │ RENDIMENTOS SUJEITOS À TRIBUTAÇÃO EXCLUSIVA │
├────────┼────────────────────┬────────────────────────┤
│ Ticker │       Valor        │         Código         │
├────────┼────────────────────┼────────────────────────┤
│ VALE3  │            1000,00 │ 10 ── JSCP             │
└────────┴────────────────────┴────────────────────────┘
┌───────┬──────────────────────────────────────────┐
│ 2025  │        OPERAÇÕES COMUNS/DAY-TRADE        │
├───────┼──────────┬────────┬───────────┬──────────┤
│  Mês  │  Ações   │ Opções │ Acumulado │ IR (15%) │
├───────┼──────────┼────────┼───────────┼──────────┤
│  AGO  │ 11520,00 │   0,00 │  11520,00 │  1728,00 │
├───────┼──────────┼────────┼───────────┼──────────┤
│ Total │ 11520,00 │   0,00 │      0,00 │  1728,00 │
└───────┴──────────┴────────┴───────────┴──────────┘
```
