# ADR-002 — Tabela IBGE embutida no binário

**Data:** 2026-05-11  
**Status:** Aceito

---

## Contexto

O renderer precisa resolver o nome do município a partir do código IBGE de 7 dígitos (`cLocIncid`, `cMun`) quando o campo `xMun` não está presente no XML. As alternativas consideradas foram:

1. Consultar a API pública do IBGE em tempo de execução (`https://servicodados.ibge.gov.br/api/v1/localidades/municipios`)
2. Usar `embed.FS` com um arquivo JSON/CSV gerado offline
3. Embutir o mapa como `var ibgeCities = map[string]string{...}` diretamente no código-fonte Go

## Decisão

O mapa é embutido em `ibge.go` como `var ibgeCities = map[string]string{...}`.

Razões:

- **Sem dependência de rede em tempo de execução** — a biblioteca funciona em ambientes isolados (CI, contêineres sem saída de rede, edge)
- **Zero latência** — a resolução é O(1) via mapa em memória
- **Sem `embed.FS`** — elimina o arquivo de dados separado; o código é a fonte da verdade
- **Atualização controlada** — `cmd/update_ibge/main.go` regenera o arquivo quando necessário, mantendo o histórico de mudanças visível no Git

## Consequências

- O binário da biblioteca cresce com a tabela completa de municípios brasileiros (~5570 entradas). O impacto é negligenciável (~300 KB de código Go compilado).
- Quando o IBGE altera códigos ou nomes de municípios, é necessário executar `cmd/update_ibge` e commitar o novo `ibge.go`. Veja o procedimento em `docs/runbook.md`.
- O fallback de resolução em `formatMunUF` segue a ordem: `xMun` → `ibgeCities[cMun]` → `cMun` raw → `uf` → `"-"`.
