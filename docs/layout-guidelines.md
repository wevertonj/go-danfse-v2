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

| Y (cm) | Seção |
|---|---|
| 0,30 | Cabeçalho (fundo cinza) |
| 1,47 | Linha separadora do cabeçalho |
| 1,55 | Chave de acesso (início) |
| 3,67 | Emitente da NFS-e / Situação / Finalidade |
| 4,34 | Linha divisória / Prestador–Fornecedor |
| 6,92 | Linha divisória / Tomador–Adquirente |
| ~9,00 | Linha divisória / Serviço |
| ~12,00 | Linha divisória / Valores |
| ~15,00 | Linha divisória / Tributação |
| ~18,00 | Linha divisória / Retenção PIS/COFINS/CSLL |
| ~21,00 | Linha divisória / Observações / Discriminação |

> Os valores marcados com `~` são aproximados; o posicionamento exato depende do conteúdo dinâmico (número de linhas de discriminação, etc.).

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
| Demais seções | variável | ~0,64 | 4,90 |

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

- NT 008/2026 — publicada em 05/05/2026 pela Sefin Nacional.
- Leiaute XML NFS-e v1.01 — XSD oficial disponível em `doc-externa/`.
- `danfse_renderer.go` — implementação de referência.
