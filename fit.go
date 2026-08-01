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

	// colGutter é a folga (cm) entre o fim do texto e a borda da coluna vizinha.
	colGutter = 0.15

	// Bordas (cm) usadas como limite direito das colunas do grid.
	colX2   = 5.41
	colX3   = 10.51
	colX4   = 15.62
	colXEnd = 20.60 // margem interna direita da página
	colXQR  = 17.48 // início do quadro do QR code, na faixa da chave de acesso

	// maxLinhasInfoCompl é quanto cabe nos 4,00 cm de altura do bloco de
	// informações complementares, a ~0,232 cm por linha de 7 pt. O MultiCell do
	// gopdf descarta em silêncio as linhas que não couberem no Rect; cortar aqui
	// é o que permite marcar a supressão com reticências (NT, item 2.1).
	maxLinhasInfoCompl = 17
)

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
func fitCell(pdf *gopdf.GoPdf, limitX float64, text string) {
	maxW := cmToPt(limitX-colGutter) - pdf.GetX()
	pdf.Cell(nil, truncText(pdf, text, maxW))
}
