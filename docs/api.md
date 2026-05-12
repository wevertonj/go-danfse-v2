# API Pública — go-danfse-v2

> Esta biblioteca não expõe endpoints HTTP. Este documento descreve a API Go pública do pacote `danfse`.

---

## Interface `DanfseRenderer`

```go
type DanfseRenderer interface {
    Render(xmlData []byte) ([]byte, error)
}
```

### `Render(xmlData []byte) ([]byte, error)`

Recebe o XML bruto de uma NFS-e e retorna o PDF do DANFSe.

**Parâmetros:**

| Nome | Tipo | Descrição |
|---|---|---|
| `xmlData` | `[]byte` | Conteúdo completo do arquivo XML da NFS-e no padrão SPED NFS-e v1.01 |

**Retorno:**

| Tipo | Descrição |
|---|---|
| `[]byte` | Bytes do PDF gerado. Os primeiros 5 bytes são sempre `%PDF-` |
| `error` | Não-nil se o XML for inválido, as fontes não forem encontradas ou a renderização falhar |

**Exemplo:**

```go
pdfBytes, err := renderer.Render(xmlData)
if err != nil {
    log.Fatal(err)
}
os.WriteFile("danfse.pdf", pdfBytes, 0644)
```

---

## Construtor `NewDanfseRenderer`

```go
func NewDanfseRenderer(fontReg, fontBold, logo string) DanfseRenderer
```

Instancia um `DanfseRenderer` configurado com os caminhos dos recursos externos.

**Parâmetros:**

| Nome | Tipo | Descrição |
|---|---|---|
| `fontReg` | `string` | Caminho para o arquivo TrueType da fonte regular (ex.: `fonts/LiberationSans-Regular.ttf`) |
| `fontBold` | `string` | Caminho para o arquivo TrueType da fonte negrito (ex.: `fonts/LiberationSans-Bold.ttf`) |
| `logo` | `string` | Caminho para a imagem PNG do logo NFS-e (ex.: `assets/logo-nfse.png`) |

Os caminhos são relativos ao diretório de trabalho do processo chamador (ou absolutos).

**Exemplo mínimo:**

```go
renderer := danfse.NewDanfseRenderer(
    "fonts/LiberationSans-Regular.ttf",
    "fonts/LiberationSans-Bold.ttf",
    "assets/logo-nfse.png",
)
```

---

## Função `ParseXML`

```go
func ParseXML(data []byte) (*NFSe, error)
```

Desserializa o XML da NFS-e na struct tipada `NFSe`. Útil quando o consumidor precisa acessar os campos individualmente sem gerar o PDF.

**Parâmetros:**

| Nome | Tipo | Descrição |
|---|---|---|
| `data` | `[]byte` | Conteúdo do XML da NFS-e |

**Retorno:**

| Tipo | Descrição |
|---|---|
| `*NFSe` | Struct populada com todos os campos do documento |
| `error` | Não-nil se o XML for malformado ou não corresponder ao namespace esperado |

---

## Struct `NFSe` (principais campos)

```go
type NFSe struct {
    InfNFSe InfNFSeNode
}

type InfNFSeNode struct {
    ID        string      // Identificador único da NFS-e
    NNFSe     string      // Número da NFS-e
    CLocIncid string      // Código IBGE do município de incidência
    AmbGer    string      // Ambiente gerador: "1" = Produção, "2" = Homologação
    DhProc    string      // Data/hora de processamento (ISO 8601)
    Emit      EmitNode    // Dados do emitente
    Valores   Valores     // Valores consolidados (total líquido, retenções, ISSQN)
    DPS       DPSNode     // Documento de Prestação de Serviço
}
```

Para a hierarquia completa das structs, consulte [docs/architecture.md](architecture.md).

---

## Compatibilidade

| Versão Go | Suporte |
|---|---|
| 1.26+ | Sim (versão mínima exigida pelo `go.mod`) |

O módulo segue o padrão SPED NFS-e **NT 008/2026** (publicada em 05/05/2026). Leiautes de versões anteriores podem ser parcialmente compatíveis, mas não são cobertos pelos testes.
