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

O grid é **fixo**: toda seção tem Y constante em `danfse_renderer.go`, independente do conteúdo. Nenhum bloco empurra o seguinte.

| Y da divisória (cm) | Y do título (cm) | Seção |
|---|---|---|
| — | 0,30 | Cabeçalho (fundo cinza, faixa de 1,17 cm) |
| 1,47 | 1,55 | Chave de acesso + QR code |
| — | 3,72 | EMITENTE DA NFS-e / Situação / Finalidade |
| 4,34 | 4,41 | PRESTADOR / FORNECEDOR |
| 6,92 | 6,99 | TOMADOR / ADQUIRENTE |
| 8,86 | 8,93 | DESTINATÁRIO DA OPERAÇÃO |
| 10,80 | 10,87 | INTERMEDIÁRIO DA OPERAÇÃO |
| 12,74 | 12,81 | SERVIÇO PRESTADO |
| 14,43 | 14,50 | TRIBUTAÇÃO MUNICIPAL (ISSQN) |
| 17,02 | 17,09 | TRIBUTAÇÃO FEDERAL (EXCETO CBS) |
| 18,32 | 18,39 | TRIBUTAÇÃO IBS / CBS |
| 20,90 | 20,97 | VALOR TOTAL DA NFS-E |
| 22,27 | 22,34 | INFORMAÇÕES COMPLEMENTARES (conteúdo em 22,68) |
| 28,10 / 28,77 | 28,17 | Canhoto: data de cientificação e assinatura |

> **Consequência do grid fixo:** a folga entre o valor de um campo e o rótulo seguinte é de ~0,34 cm — não há altura para uma segunda linha. Por isso texto que não cabe é **truncado**, nunca quebrado (ver `## Ajuste de Texto às Colunas`).

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

| Bloco | `Rect` | Constante | Verificação |
|---|---|---|---|
| Informações complementares | 20,00 × 4,00 cm em (0,40; 22,68) | `maxLinhasInfoCompl` = 17 | 17ª linha medida em `yMax` = 26,62 cm, contra o fim do `Rect` em 26,68 cm — a 18ª (~26,86) seria descartada |

Passo medido entre linhas de 7 pt: **0,23–0,24 cm**.

> A referência de 77 caracteres da tabela 2.4.5 pressupõe Microsoft Sans Serif. Em Liberation Sans e caixa alta — como os nomes empresariais chegam no XML — 77 caracteres medem ~11,1 cm e estouram os 10,19 cm do campo; por isso o corte é sempre por largura medida.

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
| TRIBUTAÇÃO MUNICIPAL (ISSQN) | 14,43 | 0,64 | 4,90 |
| TRIBUTAÇÃO FEDERAL (EXCETO CBS) | 17,02 | 0,64 | 4,90 |
| TRIBUTAÇÃO IBS / CBS | 18,32 | 0,64 | 4,90 |
| VALOR TOTAL DA NFS-E | 20,90 | 0,67 | 4,90 |

> `INFORMAÇÕES COMPLEMENTARES` é a única seção **sem** fundo cinza — só a divisória em 22,27 cm.

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
- `fit.go` — truncamento por largura (`truncText`, `fitCell`, `limitLines`).

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
| Grid 100% fixo (Y constante por seção) | um bloco maior **não** empurra o seguinte: ele sobrepõe | nunca acrescentar linha a um campo sem recalcular todo o grid abaixo |
| Descrição do serviço só comporta **uma** linha (13,79 → 14,43) | a 2ª linha caía sobre o bloco de tributação municipal | uma linha com reticências; não reintroduzir `SplitText` aqui |
| Tamanhos de fonte do item 2.4 da NT são mínimos | encolher a fonte para caber texto é não conformidade | truncar, nunca reduzir |
| `pdftotext -layout` não revela sobreposição | teste de transbordo passa falsamente | usar `-bbox` e comparar `xMax`/`yMax` |
| Repo usa Liberation Sans; a NT pede Arial e Microsoft Sans Serif | métricas diferentes das que calibram os limites de caracteres da NT | truncar por `MeasureTextWidth` |
