# rv

Criei esta ferramenta para meu uso pessoal, ela controla minhas operações com ações do mercado brasileiro,
calcula lucros, prejuízos, impostos, DARF e consolida os dados para meu informe anual de Imposto de Renda.

⚠️ Se deseja usá-la, **é por sua conta e risco**.

⚠️ Esta ferramenta **não oferece qualquer garantia de exatidão fiscal ou conformidade legal**.
Seu uso **é de inteira responsabilidade do usuário.**

[![webui](https://img.shields.io/badge/web-ui-blueviolet.svg?logo=github&labelColor=2b3137)](https://ofabricio.github.io/rv/)

## Funcionalidades

- Registro e consolidação de operações
  - Compra / Venda
  - Swing Trade / Day Trade
  - Cálculo automático de preço médio, lucros e prejuízos, IR e DARF

- Suporte a eventos corporativos e proventos
  - Dividendos
  - JSCP / JSCP Não Pago
  - Reembolso
  - Bonificação
  - Grupamento / Desdobramento
  - Rendimentos Tributáveis
  - Redução de Capital
  - Leilão de Fração
  - Subscrição
  - Aluguel

- Geração de informe anual de Imposto de Renda
  - Bens e Direitos
  - Dívida e Ônus Reais
  - Rendimentos Isentos e Não Tributáveis
  - Rendimentos Sujeitos à Tributação Exclusiva/Definitiva
  - Renda Variável: Operações Comuns/Day Trade

- Suporte a Opções
  - Compra/Venda de PUT
  - Compra/Venda de CALL

- Importação de Notas de Corretagem
  - [Clear](https://corretora.clear.com.br)

- Exportação em CSV

## Compilação

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

## WebUI

Para acessar a versão web clique no badge abaixo. Essa versão possui algumas limitações.

[![webui](https://img.shields.io/badge/web-ui-blueviolet.svg?logo=github&labelColor=2b3137)](https://ofabricio.github.io/rv/)

## Como usar

Basta adicionar as operações no arquivo [db.ndjson](/db.ndjson), que está na raiz deste repositório,
executar a ferramenta `./rv` e ver o resultado exibido no terminal.
Veja um exemplo de resultado abaixo.

**Nota:** o arquivo `db.ndjson` precisa estar na mesma pasta do executável; copie o arquivo modelo deste repositório ou crie um arquivo novo.

No momento só é possível adicionar as operações editando manualmente o arquivo `db.ndjson`.
Importação de notas de corretagem ou extratos da B3 também estão nos planos para facilitar esse processo.

Veja exemplos de uso de cada tipo de operação nos arquivos `.give` do diretório [/data/testdata](/data/testdata).

## Exemplo de resultado

```
$ ./rv
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                         OPERAÇÕES COM AÇÕES                                                         │
├────────┬────────────┬────────────┬─────┬──────────┬──────────┬───────┬──────┬─────────┬──────────────┬───────┬───────────┬──────────┤
│ Ticker │    Data    │  Operação  │ Qtd │ V. Unit. │ V. Total │ Taxas │ IRRF │ Qtd Ac. │ V. Total Ac. │  PM   │ V. Compra │  Lucro   │
├────────┼────────────┼────────────┼─────┼──────────┼──────────┼───────┼──────┼─────────┼──────────────┼───────┼───────────┼──────────┤
│ PETR4  │ 2024-01-01 │ Compra     │ 300 │    30,00 │  9000,00 │     - │    - │     300 │      9000,00 │ 30,00 │         - │        - │
│ VALE3  │ 2025-01-01 │ Compra     │ 200 │    40,00 │  8000,00 │     - │    - │     200 │      8000,00 │ 40,00 │         - │        - │
│ PETR4  │ 2025-02-01 │ Dividendos │   - │        - │        - │     - │    - │     300 │      9000,00 │ 30,00 │         - │  1000,00 │
│ VALE3  │ 2025-03-01 │ JSCP       │   - │        - │        - │     - │    - │     200 │      8000,00 │ 40,00 │         - │  1000,00 │
│ PETR4  │ 2025-04-01 │ Dividendos │   - │        - │        - │     - │    - │     300 │      9000,00 │ 30,00 │         - │   800,00 │
│ VALE3  │ 2025-05-01 │ Venda      │  50 │    65,00 │  3250,00 │     - │ 0,16 │     150 │      6000,00 │ 40,00 │   2000,00 │  1250,00 │
│ VALE3  │ 2025-06-01 │ Venda      │ 100 │   220,00 │ 22000,00 │     - │ 1,10 │      50 │      2000,00 │ 40,00 │   4000,00 │ 18000,00 │
└────────┴────────────┴────────────┴─────┴──────────┴──────────┴───────┴──────┴─────────┴──────────────┴───────┴───────────┴──────────┘
┌──────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                           BENS E DIREITOS                                            │
├────────┬──────────────────┬──────────────────┬─────────┬─────────────────────────────────────────────┤
│ Ticker │ Situação em 2023 │ Situação em 2024 │ Grp Cód │                Discriminação                │
├────────┼──────────────────┼──────────────────┼─────────┼─────────────────────────────────────────────┤
│ PETR4  │       0,00       │     9000,00      │  03 01  │ 300 AÇÕES PETR4 COM PREÇO MÉDIO DE R$ 30,00 │
└────────┴──────────────────┴──────────────────┴─────────┴─────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                           BENS E DIREITOS                                            │
├────────┬──────────────────┬──────────────────┬─────────┬─────────────────────────────────────────────┤
│ Ticker │ Situação em 2024 │ Situação em 2025 │ Grp Cód │                Discriminação                │
├────────┼──────────────────┼──────────────────┼─────────┼─────────────────────────────────────────────┤
│ PETR4  │     9000,00      │     9000,00      │  03 01  │ 300 AÇÕES PETR4 COM PREÇO MÉDIO DE R$ 30,00 │
├────────┼──────────────────┼──────────────────┼─────────┼─────────────────────────────────────────────┤
│ VALE3  │       0,00       │     2000,00      │  03 01  │ 50 AÇÕES VALE3 COM PREÇO MÉDIO DE R$ 40,00  │
└────────┴──────────────────┴──────────────────┴─────────┴─────────────────────────────────────────────┘
┌────────┬─────────────────────────────────────────┐
│  2025  │  RENDIMENTOS ISENTOS E NÃO TRIBUTÁVEIS  │
├────────┼─────────┬────────┬──────────────────────┤
│ Ticker │  Valor  │ Código │      Descrição       │
├────────┼─────────┼────────┼──────────────────────┤
│ PETR4  │ 1800,00 │   09   │ Dividendos           │
│ (SOMA) │ 1250,00 │   20   │ Isenção até R$ 20000 │
└────────┴─────────┴────────┴──────────────────────┘
┌────────┬─────────────────────────────────────────────┐
│  2025  │ RENDIMENTOS SUJEITOS À TRIBUTAÇÃO EXCLUSIVA │
├────────┼──────────────┬─────────────┬────────────────┤
│ Ticker │    Valor     │   Código    │   Descrição    │
├────────┼──────────────┼─────────────┼────────────────┤
│ VALE3  │      1000,00 │     10      │ JSCP           │
└────────┴──────────────┴─────────────┴────────────────┘
┌───────┬────────────────────────────────────────────────────────────────────┐
│ 2025  │              OPERAÇÕES COMUNS/DAY-TRADE [SWING TRADE]              │
├───────┼──────────┬────────┬───────────┬──────────┬───────────────┬─────────┤
│  Mês  │  Ações   │ Opções │ Acumulado │ IR (15%) │ IRRF (0.005%) │  DARF   │
├───────┼──────────┼────────┼───────────┼──────────┼───────────────┼─────────┤
│  JUN  │ 18000,00 │   0,00 │  18000,00 │  2700,00 │          1,10 │ 2698,90 │
├───────┼──────────┼────────┼───────────┼──────────┼───────────────┼─────────┤
│ Total │ 18000,00 │   0,00 │      0,00 │  2700,00 │          1,10 │ 2698,90 │
└───────┴──────────┴────────┴───────────┴──────────┴───────────────┴─────────┘
```

## Documentação

Para ver tudo o que pode ser feito com esta ferramenta, use o comando de ajuda `./rv -h`:

```sh
$ ./rv -h

Usage: rv [flags]

-date-format string
      formato de data (ex. 02/01/2006) (default "2006-01-02")
-exact-value
      mostra valores exatos (ex. 0,14952765 em vez de 0,15) (default false)
-file string
      arquivo de operações (default "db.ndjson")
-filter-ticker string
      filtra operações por ticker 
-filter-year int
      filtra operações por ano (YYYY para mostrar o ano em questão; 0 para mostrar todos os anos; -1 para mostrar apenas o ano atual) (default 0)
-format string
      mostra resultado no formato especificado (table, csv) (default "table")
-notas
      importa notas de corretagem de arquivos pdf que estiverem no diretório corrente
```

### Importar notas de corretagem

Para importar notas de corretagem, coloque os arquivos PDF no mesmo diretório do arquivo executável e execute:

```sh
./rv -notas
```

O resultado é exibido no terminal. Para salvar esse resultado no seu arquivo `db.ndjson` redirecione a saída:

```sh
./rv -notas >> db.ndjson
```

Em seguida visualize o resultado com `./rv`

### Exportar para CSV

Para exportar o histórico de operações para CSV:

```sh
./rv -format csv
```
