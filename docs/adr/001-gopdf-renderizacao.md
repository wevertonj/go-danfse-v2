# ADR-001 — Uso do `gopdf` para geração do PDF

**Data:** 2026-05-11  
**Status:** Aceito

---

## Contexto

O DANFSe exige um PDF A4 com posicionamento preciso de texto, suporte a fontes TrueType, inserção de imagem (logo PNG) e QR Code. A biblioteca precisa operar sem dependências externas de sistema (como LaTeX, Ghostscript ou wkhtmltopdf) e gerar o arquivo inteiramente em memória.

## Decisão

A biblioteca usa [`github.com/signintech/gopdf`](https://github.com/signintech/gopdf) para renderização do PDF.

Características determinantes para a escolha:

- API de posicionamento absoluto em pontos (pt), compatível com as diretrizes de leiaute da NT 008/2026
- Suporte nativo a fontes TrueType (sem conversão prévia)
- Suporte a imagens PNG embutidas
- Saída para `io.Writer` ou `[]byte`, sem escrita obrigatória em disco
- Dependência pura Go — sem CGO, sem binários externos

## Consequências

- O leiaute é construído com coordenadas absolutas (`cmToPt`); qualquer alteração de seção deve manter a coerência das posições verticais documentadas em `docs/layout-guidelines.md`.
- Fontes TrueType devem estar disponíveis nos caminhos configurados em `NewDanfseRenderer`. A ausência dos arquivos causa erro em tempo de execução em `Render`.
- Testes que validam o PDF verificam apenas os primeiros bytes (`%PDF-`) e o tamanho mínimo do arquivo; o conteúdo visual é auditado manualmente.
