# go-danfse-v2

Gerador de DANFSe (Documento Auxiliar da Nota Fiscal de Serviço Eletrônica) em Go, compatível com o leiaute da **NT 008/2026** (publicada em 05/05/2026) do SPED NFS-e Nacional.

> **Aviso:** Esta biblioteca foi implementada com auxílio de IA generativa e ainda não passou por validação formal junto às prefeituras ou à Receita Federal. Use-a com cautela em ambiente de produção e valide o PDF gerado contra os requisitos do seu município.

## Instalação

```bash
go get github.com/wevertonj/go-danfse-v2
```

**Requisitos:**
- Go 1.26+
- Fontes TrueType para renderização (inclua `LiberationSans-Regular.ttf` e `LiberationSans-Bold.ttf` ou equivalentes)

## Exemplo mínimo

```go
package main

import (
    "log"
    "os"

    danfse "github.com/wevertonj/go-danfse-v2"
)

func main() {
    xmlData, err := os.ReadFile("nfse.xml")
    if err != nil {
        log.Fatal(err)
    }

    renderer := danfse.NewDanfseRenderer(
        "fonts/LiberationSans-Regular.ttf",
        "fonts/LiberationSans-Bold.ttf",
        "assets/logo-nfse.png",
    )

    pdfBytes, err := renderer.Render(xmlData)
    if err != nil {
        log.Fatal(err)
    }

    os.WriteFile("danfse.pdf", pdfBytes, 0644)
}
```

O arquivo `examples/basic/main.go` contém um exemplo completo com XML de mock.

## Arquitetura

A biblioteca é composta por dois componentes principais:

### Parser (`parser.go`)

Responsável por desserializar o XML da NFS-e (namespace `http://www.sped.fazenda.gov.br/nfse`) em structs Go tipadas. Suporta o leiaute completo do padrão SPED NFS-e v1.01, incluindo:

- Dados do emitente (`emit`)
- Dados do tomador (`toma`)
- Informações do serviço (`serv`) e valores (`valores`)
- Código IBGE do município (`cLocIncid`) para resolução do nome da cidade

### Renderer (`danfse_renderer.go`)

Responsável por gerar o PDF do DANFSe a partir dos dados parseados. Usa [`signintech/gopdf`](https://github.com/signintech/gopdf) para renderização e [`skip2/go-qrcode`](https://github.com/skip2/go-qrcode) para o QR Code de autenticação. O leiaute segue as diretrizes visuais da NT 008/2026 — veja `docs/layout-guidelines.md` para detalhes.

A interface pública é:

```go
type DanfseRenderer interface {
    Render(xmlData []byte) ([]byte, error)
}

func NewDanfseRenderer(fontReg, fontBold, logo string) DanfseRenderer
```

### Utilitário IBGE (`ibge.go` + `cmd/update_ibge`)

O mapa de códigos IBGE para nomes de municípios é embutido em `ibge.go` como um `map[string]string`. Para atualizar a tabela a partir da API oficial do IBGE:

```bash
go run ./cmd/update_ibge/main.go
```

O arquivo `ibge.go` na raiz da lib será sobrescrito com os dados atualizados.

## Estrutura do repositório

```
go-danfse-v2/
├── parser.go             # Structs e desserialização do XML NFS-e
├── danfse_renderer.go    # Geração do PDF (DANFSe)
├── ibge.go               # Tabela de municípios IBGE
├── assets/               # Logo oficial NFS-e
├── fonts/                # Fontes TrueType (LiberationSans)
├── testdata/             # XML de exemplo para testes
├── examples/
│   └── basic/main.go     # Exemplo de uso completo
├── cmd/
│   └── update_ibge/      # Utilitário para atualizar tabela IBGE
└── docs/
    └── layout-guidelines.md  # Diretrizes visuais do DANFSe
```

## Conformidade

- Namespace XML: `http://www.sped.fazenda.gov.br/nfse`
- Versão do leiaute: NFS-e v1.01 (XSD oficial da Sefin Nacional)
- Referência normativa: NT 008/2026 (publicada em 05/05/2026)

## Licença

BSD 3-Clause — veja [LICENSE](LICENSE).
