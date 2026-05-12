# Arquitetura — go-danfse-v2

## Visão geral

`go-danfse-v2` é uma biblioteca Go que gera o DANFSe (Documento Auxiliar da Nota Fiscal de Serviço Eletrônica) em formato PDF a partir de um XML da NFS-e no padrão SPED NFS-e v1.01 (NT 008/2026).

O fluxo de execução é linear e sem estado compartilhado entre chamadas:

```
XML ([]byte)
    └─► ParseXML()           parser.go
            └─► NFSe (struct tipada)
                    └─► renderer.Render()    danfse_renderer.go
                                └─► PDF ([]byte)
```

---

## Componentes

### 1. Parser (`parser.go`)

Recebe um `[]byte` com o XML da NFS-e e retorna a struct `NFSe` populada.

- Namespace XML reconhecido: `http://www.sped.fazenda.gov.br/nfse`
- Função pública: `ParseXML(data []byte) (*NFSe, error)`
- Todas as structs de domínio (emitente, tomador, serviço, valores, IBS/CBS, etc.) são definidas neste arquivo
- Não há estado global nem cache — cada chamada é independente

Hierarquia principal das structs:

```
NFSe
└── InfNFSe
    ├── Emit        (dados do emitente — CNPJ, endereço, contato)
    ├── Valores     (valores consolidados da NFS-e)
    └── DPS
        └── InfDPS
            ├── Prest   (prestador)
            ├── Toma    (tomador)
            ├── Interm  (intermediário)
            ├── Serv    (serviço prestado + localização)
            ├── IBSCBS  (tributação IBS/CBS — NT 008/2026)
            └── Valores (valores da DPS)
```

### 2. Renderer (`danfse_renderer.go`)

Recebe o XML bruto (`[]byte`), invoca o parser internamente e produz o PDF do DANFSe.

- Interface pública: `DanfseRenderer` com método `Render(xmlData []byte) ([]byte, error)`
- Construtor: `NewDanfseRenderer(fontReg, fontBold, logo string) DanfseRenderer`
- Dependência de arquivos externos (resolvidos em tempo de execução):
  - `fontReg` — caminho para `LiberationSans-Regular.ttf` (ou equivalente)
  - `fontBold` — caminho para `LiberationSans-Bold.ttf` (ou equivalente)
  - `logo` — caminho para a imagem PNG do logo NFS-e
- Usa `signintech/gopdf` para composição do PDF e `skip2/go-qrcode` para o QR Code
- O leiaute visual segue `docs/layout-guidelines.md`

Funções auxiliares internas (não exportadas):

| Função | Propósito |
|---|---|
| `cmToPt(cm)` | Converte centímetros para pontos (1 cm = 28,34645669 pt) |
| `formatDoc(doc)` | Formata CNPJ (14 dígitos) ou CPF (11 dígitos) com pontuação |
| `formatCEP(cep)` | Formata CEP no padrão `XX.XXX-XXX` |
| `formatCurrency(val)` | Converte string numérica para `R$ X.XXX,XX` |
| `formatMunUF(xMun, cMun, uf)` | Resolve nome do município via `xMun` ou fallback para `ibgeCities[cMun]` |
| `formatTribISSQN(code)` | Decodifica código de tributação ISSQN |
| `formatRegEspTrib(code)` | Decodifica regime especial de tributação |
| `formatTpRetISSQN(code)` | Decodifica tipo de retenção ISSQN |
| `formatTpRetPisCofins(code)` | Decodifica tipo de retenção PIS/COFINS/CSLL |
| `formatOpSimpNac(code)` | Decodifica opção Simples Nacional |
| `formatRegApTribSN(code)` | Decodifica regime de apuração Simples Nacional |
| `formatTpSusp(code)` | Decodifica tipo de suspensão de exigibilidade |

### 3. Tabela IBGE (`ibge.go`)

Mapa `map[string]string` embutido no binário que associa o código IBGE de 7 dígitos ao nome do município. Usado como fallback quando o campo `xMun` não está presente no XML.

- Variável interna: `ibgeCities`
- Atualização via `cmd/update_ibge/main.go` (consulta a API pública do IBGE e sobrescreve `ibge.go`)

### 4. Utilitário de atualização IBGE (`cmd/update_ibge/main.go`)

Programa autônomo (não faz parte da biblioteca). Consulta `https://servicodados.ibge.gov.br/api/v1/localidades/municipios` e regenera `ibge.go` na raiz do módulo.

---

## Dependências externas

| Pacote | Versão mínima | Uso |
|---|---|---|
| `github.com/signintech/gopdf` | v0.36.0 | Renderização do PDF (A4, fontes TrueType, imagens, geometria) |
| `github.com/skip2/go-qrcode` | v0.0.0-20200617... | Geração do QR Code de autenticação da NFS-e |
| `github.com/stretchr/testify` | v1.11.1 | Asserções nos testes (apenas `_test`) |

---

## Política de testes

Os testes vivem em `parser_test.go` e `danfse_renderer_test.go`, ambos no pacote `danfse_test` (caixa-preta). Os arquivos de fixture ficam em `testdata/`.

> **Regra crítica:** se a implementação nova quebrar testes existentes de outros pacotes, **corrija a implementação, nunca os testes**. Não é permitido reescrever, deletar, comentar, marcar como `t.Skip` ou alterar testes/features que já funcionam para acomodar código novo. Se um teste antigo está realmente errado, **pare e pergunte ao usuário** antes de tocar nele.

---

## Decisões arquiteturais

As decisões de design relevantes estão registradas em `docs/adr/`:

- [ADR-001](adr/001-gopdf-renderizacao.md) — Escolha do `gopdf` para geração do PDF
- [ADR-002](adr/002-ibge-mapa-embutido.md) — Tabela IBGE embutida no binário
