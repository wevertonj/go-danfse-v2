# go-danfse-v2

[![Go Reference](https://pkg.go.dev/badge/github.com/wevertonj/go-danfse-v2.svg)](https://pkg.go.dev/github.com/wevertonj/go-danfse-v2)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

Gerador em Go para o Documento Auxiliar da Nota Fiscal de Serviço Eletrônica (DANFSe) no **Padrão Nacional (v2.0)**. 

Implementação baseada nas diretrizes de layout da **Nota Técnica Nº 008 – Versão 1.0 (Maio de 2026)**.

## ⚠️ Status do Projeto

* Como o layout é recente, o mapeamento de coordenadas foi validado contra os PDFs de exemplo oficiais da Receita Federal. **Ainda não foi comparado contra documentos reais gerados em produção.**
* A extração das coordenadas base foi realizada com auxílio de IA e revisada manualmente. A renderização utiliza desenho absoluto.

## 💡 Motivação

O serviço público deixará de fornecer o PDF da NFS-e via API a partir de **1º de julho de 2026**. O desenho e geração do DANFSe passam a ser de responsabilidade do emissor.

Este pacote efetua a montagem visual da nota sem depender de motores baseados em HTML, gerando o arquivo com coordenadas usando `gopdf`.

## ✨ Características

* **Sem CGO/Binários externos:** Geração via `gopdf`, dispensando instalação de `wkhtmltopdf` ou ferramentas baseadas no Chrome/Puppeteer.
* **Containers menores:** Adequado para imagens Docker baseadas em *scratch* ou ambientes serverless.
* **Layout aderente à norma:** Mantém fontes, margens, espessuras e posicionamento exigidos pela legislação.
* **Stateless:** A entrada (XML) e a saída (PDF) são manipuladas em `[]byte`.

## 📦 Instalação

```bash
go get github.com/wevertonj/go-danfse-v2
```

## 🚀 Como Usar

Exemplo de uso básico:

```go
package main

import (
	"log"
	"os"

	"github.com/wevertonj/go-danfse-v2"
)

func main() {
	renderer := danfse.NewDanfseRenderer()

	xmlData, err := os.ReadFile("nota_fiscal.xml")
	if err != nil {
		log.Fatalf("Falha ao ler XML: %v", err)
	}

	pdfBytes, err := renderer.Render(xmlData)
	if err != nil {
		log.Fatalf("Erro ao gerar DANFSe: %v", err)
	}

	os.WriteFile("danfse.pdf", pdfBytes, 0644)
}
```

## 🔌 Integração

A biblioteca não interage com o sistema de arquivos após sua distribuição, pois as fontes de texto (Liberation Sans) e a logo obrigatória oficial são embutidos no binário (`//go:embed`).

1. Instancie o `DanfseRenderer` na inicialização do serviço.
2. Forneça o `[]byte` do XML autorizado recebido em sua aplicação para o método `Render()`.
3. O payload `[]byte` do PDF gerado pode ser trafegado diretamente para APIs de email, storages de nuvem (S3) ou repassado ao cliente HTTP, dispensando alocação em disco.

## 🛠️ Testes

Usando os mocks embarcados no projeto:

```bash
go test ./...
go run examples/basic/main.go
```

## 📄 Licença

Este projeto é distribuído sob a licença BSD 3-Clause. Veja o arquivo `LICENSE` para mais detalhes.
