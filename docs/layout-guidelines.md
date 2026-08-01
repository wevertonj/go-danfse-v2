# Diretrizes de Leiaute do DANFSe

Este documento descreve as medidas, posicionamento, cores e tipografia usadas na geração do PDF do DANFSe. A implementação de referência está em `danfse_renderer.go`.

## Página

| Atributo | Valor |
|---|---|
| Formato | A4 (210 × 297 mm) |
| Margem de conteúdo | 0,30 cm em todos os lados |
| Borda externa | Retângulo de 20,70 × 29,40 cm, espessura 1 pt, origem (0,15; 0,15) |
| Linhas divisórias internas | Espessura 0,5 pt |

Conversão de medidas usada internamente:

```
1 cm = 28,34645669 pt
```

---

## Seções e Posicionamento Vertical

O grid é **fixo até a descrição do serviço**: acima de Y = 14,43 toda seção tem Y constante em `danfse_renderer.go`. A descrição do serviço é o único campo elástico de coluna (1 → 8 linhas); cada linha além da primeira desloca em `Δ = (linhas − 1) × passoLinha` tudo entre 14,43 e o canhoto — os Y da tabela abaixo marcados com `+ Δ` são o **grid base** (`dy()` no renderer). O canhoto não desloca.

| Y da divisória (cm) | Y do título (cm) | Seção |
|---|---|---|
| — | 0,30 | Cabeçalho (fundo cinza, faixa de 1,17 cm) |
| 1,47 | 1,55 | Chave de acesso + QR code |
| — | 3,72 | EMITENTE DA NFS-e / Situação / Finalidade |
| 4,34 | 4,41 | PRESTADOR / FORNECEDOR |
| 6,92 | 6,99 | TOMADOR / ADQUIRENTE |
| 8,86 | 8,93 | DESTINATÁRIO DA OPERAÇÃO |
| 10,80 | 10,87 | INTERMEDIÁRIO DA OPERAÇÃO |
| 12,74 | 12,81 | SERVIÇO PRESTADO (descrição do serviço em 14,16, elástica) |
| 14,43 + Δ | 14,50 + Δ | TRIBUTAÇÃO MUNICIPAL (ISSQN) |
| 17,02 + Δ | 17,09 + Δ | TRIBUTAÇÃO FEDERAL (EXCETO CBS) |
| 18,32 + Δ | 18,39 + Δ | TRIBUTAÇÃO IBS / CBS |
| 20,90 + Δ | 20,97 + Δ | VALOR TOTAL DA NFS-E |
| 22,27 + Δ | 22,34 + Δ | INFORMAÇÕES COMPLEMENTARES (conteúdo em 22,68 + Δ) |
| 28,10 / 28,77 | 28,17 | Canhoto: data de cientificação e assinatura — **fixo**, nunca desloca |

> **Consequência do grid fixo:** a folga entre o valor de um campo e o rótulo seguinte é de ~0,34 cm — não há altura para uma segunda linha em campo de coluna. Por isso texto que não cabe é **truncado**, nunca quebrado (ver `## Ajuste de Texto às Colunas`); as únicas exceções são os dois campos elásticos (`## Campos Elásticos`).

---

## Cabeçalho

- **Fundo:** cinza claro — `RGB(242, 242, 242)`.
- **Dimensões:** 20,40 × 1,17 cm, origem (0,30; 0,30).
- **Logo:** alinhado à esquerda, largura fixa de 4,00 cm, altura proporcional (max 1,16 cm), centralizado verticalmente na faixa de 1,16 cm; origem X = 0,49 cm.
- **Título central:** fonte `LiberationSans-Bold` 9 pt — "DANFSe v2.0" e "Documento Auxiliar da NFS-e", centralizados na faixa X de 5,41–15,60 cm.
- **Aviso de homologação:** exibido apenas quando `ambGer = "2"`, texto "NFS-e SEM VALIDADE JURÍDICA" em vermelho `RGB(255, 0, 0)`, fonte `LiberationSans-Bold` 9 pt.
- **Metadados (canto direito):** fonte `LiberationSans` 8 pt (município/UF) e 6 pt (ambiente gerador, tipo de ambiente), iniciando em X = 15,62 cm.

---

## Colunas de Dados

O leiaute usa três colunas horizontais recorrentes:

| Coluna | X inicial (cm) |
|---|---|
| 1 | 0,40 |
| 2 | 5,41 |
| 3 | 10,51 |
| 4 (quando presente) | 15,62 |

---

## Ajuste de Texto às Colunas

Todo campo de conteúdo variável passa por `fitCell` (`fit.go`) antes de ser desenhado: `pdf.Cell` do gopdf não quebra nem limita o texto, então um valor longo invadiria a coluna vizinha ou a borda da página.

| Regra | Valor |
|---|---|
| Base normativa | NT 008/2026 v1.02, item 2.1 — quantidade menor de caracteres "acrescido de reticências (...)" quando o campo não suporta o texto |
| Reticências | `...` (três pontos), como nos exemplos do item 2.4.5 da NT — não o caractere `…` |
| Fonte | **mantida**; os tamanhos do item 2.4 são mínimos e não podem ser reduzidos para caber texto |
| Critério de corte | largura medida com `MeasureTextWidth`, nunca contagem de caracteres |
| Folga até a coluna vizinha | `colGutter` = 0,15 cm |

Limites direitos usados por `fitCell` (constantes em `fit.go`):

| Constante | X (cm) | Campos |
|---|---|---|
| `colX2` | 5,41 | coluna 1 simples (número da NFS-e, número da DPS, BC ISSQN, …) |
| `colX3` | 10,51 | coluna 2 e campos de duas colunas iniciados em 0,40 (nome, endereço) |
| `colX4` | 15,62 | coluna 3 simples (município/UF, IM, retenções) |
| `colXEnd` | 20,60 | coluna 4 e campos de largura total (descrições, e-mail) |
| `colXQR` | 17,48 | chave de acesso — barrada pelo quadro do QR code |

Campos de bloco (`MultiCell`) usam `limitLines`: o gopdf descarta em silêncio as linhas que não cabem na altura do `Rect`, e o corte prévio marca a supressão com reticências.

> A referência de 77 caracteres da tabela 2.4.5 pressupõe Microsoft Sans Serif. Em Liberation Sans e caixa alta — como os nomes empresariais chegam no XML — 77 caracteres medem ~11,1 cm e estouram os 10,19 cm do campo; por isso o corte é sempre por largura medida.

---

## Campos sem Informação

Base normativa: NT 008/2026 v1.02, **nota 12 do item 2.4.5** — "os campos sem informações no XML devem ser preenchidos com um traço (-)". Campo em branco é indistinguível de campo que o emissor deixou de renderizar; o traço declara que o dado não existe.

| Regra | Valor |
|---|---|
| Preenchimento | `traco` = `-` (`fit.go`) |
| Ponto único | `ouTraco` dentro do `fitCell` — cobre os 89 campos de conteúdo variável de uma vez |
| Conteúdo só de espaços | conta como ausente (imprime igual a campo vazio) |
| Ordem | traço **antes** do truncamento: a nota 12 trata do dado que faltou no XML; o que não coube leva reticências (item 2.1) |
| Células que juntam rótulo e valor | `Município: <valor> / <UF>`, `Ambiente Gerador: <valor>` e `Tipo de Ambiente: <valor>` chamam `ouTraco` no valor antes de compor — passar a célula inteira ao `fitCell` nunca dispararia a regra, porque o rótulo já a torna não vazia |
| Bloco de informações complementares | vazio → traço em `Cell`; com texto → `MultiCell` + `limitLines` |

Fora do alcance da nota 12:

| Campo | Motivo |
|---|---|
| `DATA CIENTIFICAÇÃO` / `IDENTIFICAÇÃO E ASSINATURA` (canhoto) | em branco de propósito, para preenchimento à mão |
| `EMITENTE DA NFS-e` | a célula não tem valor ligado ao XML em `danfse_renderer.go` — campo não implementado, não campo vazio |
| Linha inteira sem dados | notas 1 e 5 do item 2.4.5 permitem **suprimir a linha**; é regra distinta e não implementada |

Fixture da regra: `testdata/mock_danfse_minimo.xml` (só chave de acesso, município/UF do emitente e `tpAmb`).

---

## Campos Elásticos

Base normativa: NT 008/2026 v1.02, item 2.1 — mais de uma linha é permitido ("desde que a leitura possa ser feita de forma clara") e os tamanhos do 2.4.5 não são obrigatórios. Só dois campos têm orçamento da NT acima de 1 linha; todos os demais truncam em 1 linha via `fitCell`.

| Campo | Orçamento da NT | Linhas | Regra (`fit.go` + `danfse_renderer.go`) |
|---|---|---|---|
| Descrição do serviço (`xDescServ`) | 1.297 chars | 1 → `maxLinhasDescServ` = 8 | `SplitText` conta as linhas; acima de 8, `limitLines` corta com `...`; `Δ = (linhas − 1) × passoLinha` desloca tudo entre 14,43 e o canhoto |
| Informações complementares (`xOutInf`) | 1.997 chars | `linhasInfoCompl(Δ)`: 23 − linhas extras da descrição, nunca < `minLinhasInfoCompl` = 12 | `Rect` de 20,00 cm × (28,10 − 0,05 − (22,68 + Δ)); recebe o saldo do pool quando a descrição encolhe |

| Constante | Valor | Origem |
|---|---|---|
| `passoLinha` | 0,2313 cm | avanço de linha do `MultiCell` a 7 pt, medido (17 linhas de `yMin` 22,640 a 26,340) |
| `yInfoCompl` | 22,68 | Y base do conteúdo das informações complementares |
| `yCanhoto` | 28,10 | linha superior do canhoto — limite rígido do pool |
| `folgaInfoCompl` | 0,10 cm | respiro antes do canhoto no cálculo de `linhasInfoCompl` |

Pool de 24 linhas (1 da descrição + 17 históricas + 6 ociosas do rodapé) contra 20 de orçamento da NT (8 + 12): os dois orçamentos cabem simultaneamente, por isso a alocação é gulosa — descrição primeiro, resto para as informações complementares.

---

## Chave de Acesso e QR Code

- **Rótulo:** `LiberationSans-Bold` 7 pt, Y = 1,55 cm, X = 0,40 cm.
- **Chave (texto):** `LiberationSans` 7 pt, Y = 1,90 cm.
- **QR Code:** 1,90 × 1,90 cm, origem (17,48; 1,60), sem borda, resolução 256 px. URL de consulta: `https://www.nfse.gov.br/ConsultaPublica/?tpc=1&chave=<chave>`.
- **Texto de autenticidade:** `LiberationSans` 6 pt, linhas em Y = 3,65 / 3,85 / 4,05 cm, X = 15,62 cm.

---

## Tipografia

| Uso | Fonte | Tamanho |
|---|---|---|
| Rótulos de seção (títulos cinza) | `LiberationSans-Bold` | 7 pt |
| Sub-rótulos de campo | `LiberationSans-Bold` | 6 pt |
| Valores de dados | `LiberationSans` | 7 pt |
| Metadados de cabeçalho | `LiberationSans` | 6–8 pt |
| Título principal | `LiberationSans-Bold` | 9 pt |

Arquivos de fonte esperados (configuráveis via `NewDanfseRenderer`):
- `fonts/LiberationSans-Regular.ttf`
- `fonts/LiberationSans-Bold.ttf`

---

## Cores

| Elemento | Cor |
|---|---|
| Texto padrão | `RGB(0, 0, 0)` — preto |
| Fundo de títulos de seção | `RGB(242, 242, 242)` — cinza 5% |
| Aviso de homologação | `RGB(255, 0, 0)` — vermelho |

---

## Campos com Destaque de Seção (fundo cinza)

As células de título das seções principais recebem fundo cinza aplicado **antes** de desenhar a linha divisória, para que a linha sobreponha o fundo:

| Seção | Y (cm) | Altura (cm) | Largura (cm) |
|---|---|---|---|
| EMITENTE DA NFS-e | 3,67 | 0,67 | 4,90 |
| PRESTADOR / FORNECEDOR | 4,34 | 0,64 | 4,90 |
| TOMADOR / ADQUIRENTE | 6,92 | 0,64 | 4,90 |
| DESTINATÁRIO DA OPERAÇÃO | 8,86 | 0,64 | 4,90 |
| INTERMEDIÁRIO DA OPERAÇÃO | 10,80 | 0,64 | 4,90 |
| SERVIÇO PRESTADO | 12,74 | 0,64 | 4,90 |
| TRIBUTAÇÃO MUNICIPAL (ISSQN) | 14,43 + Δ | 0,64 | 4,90 |
| TRIBUTAÇÃO FEDERAL (EXCETO CBS) | 17,02 + Δ | 0,64 | 4,90 |
| TRIBUTAÇÃO IBS / CBS | 18,32 + Δ | 0,64 | 4,90 |
| VALOR TOTAL DA NFS-E | 20,90 + Δ | 0,67 | 4,90 |

> `INFORMAÇÕES COMPLEMENTARES` é a única seção **sem** fundo cinza — só a divisória em 22,27 + Δ cm. Os `Rect` de fundo deslocam junto com as divisórias (`dy()` cobre ambos): fundo e linha nunca desalinham.

---

## Formatação de Dados

| Dado | Formato |
|---|---|
| CNPJ | `XX.XXX.XXX/XXXX-XX` |
| CPF | `XXX.XXX.XXX-XX` |
| CEP | `XX.XXX-XXX` |
| Telefone (10 dígitos) | `(XX) XXXX-XXXX` |
| Telefone (11 dígitos) | `(XX) XXXXX-XXXX` |
| Valores monetários | `R$ X.XXX,XX` (separador milhar `.`, decimal `,`) |
| Código de Tributação Nacional | `XX.XX.XX` (6 dígitos formatados com pontos) |
| Cód. Tributação Nacional + Municipal | `XX.XX.XX / NNN` |

---

## Referências

- NT 008/2026 **v1.02** (14/07/2026) — Sefin Nacional: `https://www.gov.br/nfse/pt-br/biblioteca/documentacao-tecnica/rtc/nt-008-se-cgnfse-danfse-20260714-v1-02.pdf`
- Leiaute XML NFS-e v1.01 — XSD oficial no portal nacional da NFS-e.
- `danfse_renderer.go` — implementação de referência.
- `fit.go` — truncamento por largura (`truncText`, `fitCell`, `limitLines`) e alocação dos campos elásticos (`linhasInfoCompl`, `passoLinha`, `maxLinhasDescServ`).

---

## Segurança e Cibersegurança

| Vetor | Mitigação no contexto |
|---|---|
| Injeção via conteúdo do XML | todo texto do XML entra como **conteúdo** do PDF (`Cell`/`MultiCell`), nunca como estrutura ou ação; o gopdf escapa os literais de string do content stream |
| XXE / expansão de entidades no parse | `encoding/xml` da stdlib não resolve entidades externas nem DTD; não trocar por parser que resolva |
| Supressão indevida de dado fiscal | o truncamento **sempre** marca o corte com `...` (NT, item 2.1); nunca cortar em silêncio — o leitor precisa saber que o campo continua |
| Chave de acesso não validada na URL do QR | a chave vai concatenada em `https://www.nfse.gov.br/ConsultaPublica/?tpc=1&chave=<chave>` sem validação de formato — ver `T4.5` do plano |
| Dados inventados no documento | UF ausente e não derivável **aborta** o `Render` (issue #1); é proibido preencher campo fiscal com valor padrão fabricado |
| Menor privilégio | o renderer não lê disco nem rede: fontes e logo são embutidos via `go:embed`; a entrada é só o `[]byte` do XML |

---

## Desenvolvimento & Gotchas

| Gotcha | Impacto | Ação |
|---|---|---|
| `pdf.Cell(nil, ...)` não quebra nem limita texto | valor longo invade a coluna vizinha e a borda | todo dado variável passa por `fitCell` |
| `MultiCell` descarta linhas excedentes sem sinal | texto termina no meio de uma frase como se fosse o conteúdo | `limitLines` corta antes e acrescenta `...` |
| `SplitText`/`MultiCell` quebram por largura, **no meio da palavra** | informações complementares saem com `informaca` / `o emitente` em linhas seguidas | comportamento nativo do gopdf; `limitLines` conta caracteres em vez de rejuntar linhas, para não gravar no texto cortes que não existem no original |
| Grid fixo acima de 14,43 (Y constante por seção) | um bloco maior **não** empurra o seguinte: ele sobrepõe | nunca acrescentar linha a um campo sem passar pelo mecanismo `dy()`/`passoLinha` |
| Y do trecho ISSQN → informações complementares passam por `dy()` | um `SetY`/`Line`/`Rect` novo com `cmToPt` cru nesse trecho fica para trás quando a descrição cresce e sobrepõe o bloco vizinho | usar `dy()` em todo Y novo entre 14,43 e o canhoto; o detector par-a-par de bbox no teste da descrição elástica acusa o esquecimento |
| `desloc` calculado com `SplitText` na fonte corrente | contar linhas com outra fonte/tamanho descasa o Δ do texto real | `SetFont` de 7 pt **antes** do `SplitText` da descrição |
| Canhoto fixo em 28,10 | deslocá-lo quebraria a página única | canhoto fora do `dy()`; `linhasInfoCompl` já respeita o limite |
| Tamanhos de fonte do item 2.4 da NT são mínimos | encolher a fonte para caber texto é não conformidade | truncar, nunca reduzir |
| `pdftotext -layout` não revela sobreposição | teste de transbordo passa falsamente | usar `-bbox` e comparar `xMax`/`yMax` |
| Repo usa Liberation Sans; a NT pede Arial e Microsoft Sans Serif | métricas diferentes das que calibram os limites de caracteres da NT | truncar por `MeasureTextWidth` |
| Campo novo desenhado com `pdf.Cell` em vez de `fitCell` | escapa do truncamento **e** do traço da nota 12 | todo dado do XML entra por `fitCell`; célula que junta rótulo e valor aplica `ouTraco` só no valor |
| `SITUAÇÃO DA NFS-e` e `FINALIDADE` saem fixos em `1 - Normal` | valor não vem do XML — nenhum traço aparece ali porque nada é lido | tratar junto com a marca d'água de cancelamento/substituição (`T4.2` do plano) |
