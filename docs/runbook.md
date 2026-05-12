# Runbook — go-danfse-v2

Procedimentos operacionais para manutenção e uso da biblioteca.

---

## 1. Atualizar a tabela de municípios IBGE

Execute quando o IBGE publicar novos códigos ou renomear municípios.

A partir da **raiz do repositório**:

```bash
go run ./cmd/update_ibge/main.go
```

O script consulta `https://servicodados.ibge.gov.br/api/v1/localidades/municipios`, regenera `ibge.go` na raiz do módulo e encerra. Verifique o diff antes de commitar:

```bash
git diff ibge.go
```

Após validar, commite o arquivo atualizado:

```bash
git add ibge.go
git commit -m "chore: atualizar tabela IBGE"
```

---

## 2. Trocar ou atualizar fontes TrueType

As fontes são embutidas no binário via `//go:embed` a partir dos arquivos `fonts/LiberationSans-Regular.ttf` e `fonts/LiberationSans-Bold.ttf`. Para substituí-las:

1. Substitua os arquivos `.ttf` em `fonts/`
2. Recompile a biblioteca
3. Execute os testes para confirmar que o PDF continua sendo gerado corretamente:

```bash
go test ./...
```

Se a nova fonte não suportar caracteres acentuados presentes nos dados das NFS-e, o PDF exibirá substituições ou espaços em branco — verifique visualmente com `examples/basic/main.go`.

---

## 3. Rodar os testes

```bash
go test ./...
```

Os testes utilizam `testdata/mock_danfse.xml` como fixture. Nenhuma rede ou serviço externo é necessário.

Para rodar com saída detalhada:

```bash
go test -v ./...
```

> **Regra crítica:** se a implementação nova quebrar testes existentes de outros pacotes, **corrija a implementação, nunca os testes**. Não é permitido reescrever, deletar, comentar, marcar como `t.Skip` ou alterar testes/features que já funcionam para acomodar código novo. Se um teste antigo está realmente errado, **pare e pergunte ao usuário** antes de tocar nele.

---

## 4. Gerar o PDF de exemplo (smoke test visual)

```bash
cd /home/weverton/node/go-danfse-v2
go run examples/basic/main.go
```

O arquivo `preview.pdf` é gerado no diretório de trabalho. Abra-o em qualquer visualizador de PDF para conferir o leiaute visualmente.

---

## 5. Adicionar suporte a um novo campo da NFS-e

1. Identifique o campo no schema SPED NFS-e v1.01 (namespace `http://www.sped.fazenda.gov.br/nfse`)
2. Adicione o campo na struct correspondente em `parser.go` com a tag XML correta
3. Adicione a renderização do campo em `danfse_renderer.go` seguindo as posições definidas em `docs/layout-guidelines.md`
4. Adicione ou atualize o caso de teste em `parser_test.go` ou `danfse_renderer_test.go`
5. Se a mudança alterar o fluxo de execução ou a dependência entre componentes, atualize `docs/architecture.md`

---

## 6. Publicar nova versão do módulo

```bash
git tag v2.x.y
git push origin v2.x.y
```

Não há pipeline de CI configurado neste repositório. A validação é feita localmente com `go test ./...` antes do tag.
