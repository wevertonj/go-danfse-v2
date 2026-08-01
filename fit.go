package danfse

import (
	"strings"

	"github.com/signintech/gopdf"
)

// Ajuste de texto às colunas do DANFSe conforme a NT 008/2026 v1.02:
//
//   - item 2.1: "a quantidade de caracteres para cada campo sugerida no item
//     2.4.5 não tem caráter obrigatório, podendo-se utilizar quantidade diversa,
//     acrescido de reticências (...), quando o campo não suportar a totalidade
//     de caracteres do texto a ser inserido";
//   - item 2.1 + 2.4: os tamanhos de fonte do item 2.4 são MÍNIMOS — encolher a
//     fonte para caber o texto não é alternativa válida.
//
// O corte é feito por largura medida na fonte corrente, e não por contagem de
// caracteres: as quantidades do item 2.4.5 pressupõem Arial/Microsoft Sans
// Serif e aqui a fonte embutida é a Liberation Sans.
const (
	// ellipsis é a reticência da NT, grafada com três pontos nos exemplos do
	// item 2.4.5 ("...Administ...", "...papel..."), não com o caractere "…".
	ellipsis = "..."

	// traco preenche o campo sem informação, conforme a nota 12 do item 2.4.5
	// da NT: "os campos sem informações no XML devem ser preenchidos com um
	// traço (-)". Não confundir com a linha inteira sem dados, que as notas 1 e
	// 5 permitem suprimir.
	traco = "-"

	// colGutter é a folga (cm) entre o fim do texto e a borda da coluna vizinha.
	colGutter = 0.15

	// Bordas (cm) usadas como limite direito das colunas do grid.
	colX2   = 5.41
	colX3   = 10.51
	colX4   = 15.62
	colXEnd = 20.60 // margem interna direita da página
	colXQR  = 17.48 // início do quadro do QR code, na faixa da chave de acesso

	// passoLinha é o avanço vertical (cm) de uma linha de MultiCell a 7 pt,
	// medido no bloco de informações complementares (17 linhas de 22,640 a
	// 26,340 cm de yMin).
	passoLinha = 0.2313

	// maxLinhasDescServ acomoda o orçamento da NT para o xDescServ (1.297
	// chars, tabela 2.4.5) na largura útil da página (~170 chars por linha de
	// 7 pt). A descrição do serviço é elástica (item 2.1 permite mais de uma
	// linha) e desloca os blocos seguintes para baixo até este teto.
	maxLinhasDescServ = 8

	// minLinhasInfoCompl preserva o orçamento da NT para o xOutInf (1.997
	// chars): mesmo cedendo espaço à descrição, as informações complementares
	// nunca ficam com menos linhas que isto.
	minLinhasInfoCompl = 12

	// yInfoCompl e yCanhoto delimitam o bloco de informações complementares no
	// grid sem deslocamento: conteúdo em 22,68 cm e canhoto fixo em 28,10 cm
	// (o canhoto nunca desloca — o DANFSe é de página única).
	yInfoCompl = 22.68
	yCanhoto   = 28.10

	// folgaInfoCompl é o respiro (cm) entre a última linha possível das
	// informações complementares e a linha superior do canhoto.
	folgaInfoCompl = 0.10
)

// linhasInfoCompl devolve quantas linhas de 7 pt as informações complementares
// comportam depois de descerem desloc cm (o crescimento da descrição do
// serviço): o que sobra até o canhoto, nunca menos que o orçamento da NT. O
// MultiCell do gopdf descarta em silêncio as linhas que não couberem no Rect;
// cortar por esta conta é o que permite marcar a supressão com reticências
// (NT, item 2.1).
func linhasInfoCompl(desloc float64) int {
	linhas := int((yCanhoto - (yInfoCompl + desloc) - folgaInfoCompl)/passoLinha + 1e-9)
	if linhas < minLinhasInfoCompl {
		return minLinhasInfoCompl
	}
	return linhas
}

// ouTraco devolve o traço da nota 12 quando o campo não veio no XML. Campo em
// branco no DANFSe impresso não se distingue de campo que o emissor deixou de
// renderizar; o traço declara que o dado não existe. Conteúdo só de espaços
// imprime igual a campo vazio, então conta como ausente.
//
// fitCell aplica isto em todo campo de conteúdo variável — os campos do
// cabeçalho, que juntam rótulo e valor na mesma célula ("Município: CCCC / CC"),
// chamam direto para não perder o rótulo.
func ouTraco(v string) string {
	if strings.TrimSpace(v) == "" {
		return traco
	}
	return v
}

// truncText devolve text encurtado com reticências até caber em maxW pontos,
// medido na fonte corrente de pdf. Texto que já cabe volta intacto; quando nem
// as reticências cabem, devolve "" — não há largura para representar o dado.
func truncText(pdf *gopdf.GoPdf, text string, maxW float64) string {
	if text == "" {
		return text
	}
	if w, err := pdf.MeasureTextWidth(text); err != nil || w <= maxW {
		return text
	}

	runes := []rune(text)
	// Busca binária pelo maior prefixo que ainda caiba com as reticências.
	lo, hi, melhor := 0, len(runes)-1, ""
	for lo <= hi {
		meio := (lo + hi) / 2
		cand := strings.TrimRight(string(runes[:meio]), " ") + ellipsis
		w, err := pdf.MeasureTextWidth(cand)
		if err != nil {
			return text
		}
		if w <= maxW {
			melhor = cand
			lo = meio + 1
		} else {
			hi = meio - 1
		}
	}
	if melhor == ellipsis {
		// Só as reticências couberam: não sobra nada do dado para exibir.
		return ""
	}
	return melhor
}

// limitLines corta text nas maxLines primeiras linhas de largura width (pt),
// marcando o corte com as reticências da NT na última linha mantida. Usado nos
// campos de bloco (MultiCell), que descartam o excedente sem sinal algum: sem
// isto o texto termina no meio de uma frase como se fosse esse o conteúdo.
func limitLines(pdf *gopdf.GoPdf, text string, width float64, maxLines int) string {
	linhas, err := pdf.SplitText(text, width)
	if err != nil || len(linhas) <= maxLines {
		return text
	}
	// Conta caracteres em vez de rejuntar as linhas: o SplitText do gopdf quebra
	// por largura, inclusive no meio de palavras, e recompor com espaços
	// introduziria cortes que não existem no texto original.
	cabem := 0
	for _, l := range linhas[:maxLines] {
		cabem += len([]rune(l))
	}
	runes := []rune(text)
	if cabem > len(runes) {
		cabem = len(runes)
	}
	cabem -= len(ellipsis) // espaço para as reticências na última linha
	if cabem < 0 {
		cabem = 0
	}
	return strings.TrimRight(string(runes[:cabem]), " ") + ellipsis
}

// fitCell escreve text na posição corrente do cursor, truncado para não passar
// de limitX (cm) — a borda da coluna vizinha ou da página. Substitui pdf.Cell
// em todo campo de conteúdo variável: gopdf não quebra nem limita o texto de
// Cell, então sem isto um valor longo invade a coluna ao lado (issue #2).
//
// Campo ausente sai com o traço da nota 12; o traço vem antes do truncamento
// porque a regra é sobre o dado que faltou no XML, não sobre o que não coube
// (nesse caso a NT manda reticências).
func fitCell(pdf *gopdf.GoPdf, limitX float64, text string) {
	maxW := cmToPt(limitX-colGutter) - pdf.GetX()
	pdf.Cell(nil, truncText(pdf, ouTraco(text), maxW))
}
