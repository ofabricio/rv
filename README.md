# rv

Esta ferramenta faz o controle de operações de ações do mercado brasileiro, gerando as guias de declaração de IR (Bens e Direitos, Operações Comuns/Day Trade, etc.)

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
  - JSCP
  - Leilão de Fração
  - Redução de Capital
  - Rend. Trib.
  - Subscrição

- Opções (em breve...)
  - Compra/Venda de PUT
  - Compra/Venda de CALL

## Build & Run

Se tiver Go [instalado](https://go.dev/dl):

```sh
git clone --depth 1 https://github.com/ofabricio/rv.git ; cd rv
GOEXPERIMENT=jsonv2 go run .
```

**Nota:** para gerar um executável use `go build` em vez de `go run .`

---

Se tiver docker [instalado](https://www.docker.com):

```sh
git clone --depth 1 https://github.com/ofabricio/rv.git ; cd rv
docker run --rm -v $PWD:/src -w /src -e GOEXPERIMENT=jsonv2 golang:alpine go run .
```

**Nota:** para gerar um executável use `go build` em vez de `go run .` e adicione `-e GOOS=darwin -e GOARCH=amd64` e altere os [valores](https://go.dev/doc/install/source#environment) conforme seu sistema operacional.

## Como usar

Basta adicionar as operações e acompanhar a evolução.

No momento só é possível adicionar as operações editando manualmente o arquivo [db.ndjson](/db.ndjson), que fica na raiz do projeto. No futuro será criado uma interface no terminal ou no navegador para facilitar isso. Importação de notas de corretagem ou extratos da B3 também estão nos planos.

Veja exemplos de uso de cada tipo de operação nos arquivos `*.give` do diretório [/model/testdata](/model/testdata).

Para visualizar todas as operações:

```sh
./rv
```

## Exemplo de output

```
$ ./rv
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                         OPERAÇÕES COM AÇÕES                                                          │
├────┬────────┬────────────┬───────────────┬─────┬──────────┬──────────┬───────┬─────────┬──────────────┬───────┬───────────┬──────────┤
│ ID │ Ticker │    Data    │      Opr      │ Qtd │ V. Unit. │ V. Total │ Taxas │ Qtd Ac. │ V. Total Ac. │  PM   │ V. Compra │  Lucro   │
├────┼────────┼────────────┼───────────────┼─────┼──────────┼──────────┼───────┼─────────┼──────────────┼───────┼───────────┼──────────┤
│ 1  │ PETR4  │ 01/01/2024 │ Compra        │ 100 │    30,00 │  3000,00 │  0,00 │     100 │      3000,00 │ 30,00 │      0,00 │     0,00 │
│ 2  │ PETR4  │ 01/02/2024 │ Compra        │ 200 │    30,50 │  6100,00 │  0,00 │     300 │      9100,00 │ 30,33 │      0,00 │     0,00 │
│ 3  │ VALE3  │ 01/03/2025 │ Compra        │ 200 │    60,00 │ 12000,00 │  0,00 │     200 │     12000,00 │ 60,00 │      0,00 │     0,00 │
│ 4  │ VALE3  │ 01/04/2025 │ Desdobramento │   0 │     0,00 │     0,00 │  0,00 │     400 │     12000,00 │ 30,00 │      0,00 │     0,00 │
│ 5  │ VALE3  │ 01/05/2025 │ Compra        │ 100 │    61,00 │  6100,00 │  0,00 │     500 │     18100,00 │ 36,20 │      0,00 │     0,00 │
│ 6  │ PETR4  │ 01/06/2025 │ Dividendos    │   0 │     0,00 │     0,00 │  0,00 │     300 │      9100,00 │ 30,33 │      0,00 │  1000,00 │
│ 7  │ VALE3  │ 01/07/2025 │ JSCP          │   0 │     0,00 │     0,00 │  0,00 │     500 │     18100,00 │ 36,20 │      0,00 │  1000,00 │
│ 8  │ VALE3  │ 01/08/2025 │ Venda         │ 400 │    65,00 │ 26000,00 │  0,00 │     100 │      3620,00 │ 36,20 │  14480,00 │ 11520,00 │
└────┴────────┴────────────┴───────────────┴─────┴──────────┴──────────┴───────┴─────────┴──────────────┴───────┴───────────┴──────────┘
┌────────────────────────────────────────────────────────────────────────────────────────────┐
│                          BENS E DIREITOS ── Grupo 03 ── Código 01                          │
├────────┬──────────────────┬──────────────────┬─────────────────────────────────────────────┤
│ Ticker │ Situação em 2024 │ Situação em 2025 │                Discriminação                │
├────────┼──────────────────┼──────────────────┼─────────────────────────────────────────────┤
│ PETR4  │     9100,00      │     9100,00      │ 300 AÇÕES PETR4 COM PREÇO MÉDIO DE R$ 30,33 │
├────────┼──────────────────┼──────────────────┼─────────────────────────────────────────────┤
│ VALE3  │       0,00       │     3620,00      │ 100 AÇÕES VALE3 COM PREÇO MÉDIO DE R$ 36,20 │
└────────┴──────────────────┴──────────────────┴─────────────────────────────────────────────┘
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
