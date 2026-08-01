package danfse

import (
	"bytes"
	"strings"
	"testing"

	"github.com/signintech/gopdf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMeasurePdf devolve um PDF com a fonte de conteúdo carregada no tamanho
// pedido, para medir larguras de texto como o renderer faz.
func newMeasurePdf(t *testing.T, size int) *gopdf.GoPdf {
	t.Helper()
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	require.NoError(t, pdf.AddTTFFontByReader("LiberationSans", bytes.NewReader(embeddedFontReg)))
	require.NoError(t, pdf.SetFont("LiberationSans", "", size))
	return pdf
}

func measure(t *testing.T, pdf *gopdf.GoPdf, s string) float64 {
	t.Helper()
	w, err := pdf.MeasureTextWidth(s)
	require.NoError(t, err)
	return w
}

// TestTruncText cobre o truncamento prescrito pela NT 008/2026 v1.02, item 2.1:
// quando o campo não suporta o texto, usa-se quantidade menor de caracteres
// acrescida de reticências — mantendo os tamanhos de fonte do item 2.4, que são
// mínimos (nunca encolher a fonte para caber).
func TestTruncText(t *testing.T) {
	pdf := newMeasurePdf(t, 7)
	maxW := cmToPt(9.96) // largura útil do campo NOME / NOME EMPRESARIAL

	t.Run("should return short text unchanged", func(t *testing.T) {
		assert.Equal(t, "PRESTADOR DE TESTE LTDA", truncText(pdf, "PRESTADOR DE TESTE LTDA", maxW))
	})

	t.Run("should return an empty string unchanged", func(t *testing.T) {
		assert.Equal(t, "", truncText(pdf, "", maxW))
	})

	t.Run("should truncate long text with an ellipsis that fits the width", func(t *testing.T) {
		longo := strings.Repeat("EMPRESA BRASILEIRA DE SERVICOS ", 20)

		got := truncText(pdf, longo, maxW)

		assert.True(t, strings.HasSuffix(got, "..."), "should end with the NT ellipsis, got %q", got)
		assert.LessOrEqual(t, measure(t, pdf, got), maxW)
		assert.True(t, strings.HasPrefix(longo, strings.TrimSuffix(got, "...")),
			"should keep the original prefix, got %q", got)
	})

	t.Run("should not truncate the 77-character name the NT sizes the field for", func(t *testing.T) {
		// NT 2.4.5: campo xNome de 80 chars, "reticências caso supere 77 caracteres".
		nome := "Companhia Brasileira de Desenvolvimento de Sistemas e Servicos Digitais Ltda."
		require.Len(t, []rune(nome), 77)

		assert.Equal(t, nome, truncText(pdf, nome, maxW))
	})

	t.Run("should truncate a 77-character name in upper case, which is wider than the column", func(t *testing.T) {
		// A referência de 77 caracteres da NT pressupõe Microsoft Sans Serif; em
		// Liberation Sans e caixa alta (como os nomes empresariais chegam no XML)
		// 77 caracteres medem ~11,1 cm e estouram os 10,19 cm do campo. Por isso o
		// corte é por largura medida, nunca por contagem de caracteres.
		nome := "COMPANHIA BRASILEIRA DE DESENVOLVIMENTO DE SISTEMAS E SERVICOS DIGITAIS LTDA"

		got := truncText(pdf, nome, maxW)

		assert.NotEqual(t, nome, got)
		assert.LessOrEqual(t, measure(t, pdf, got), maxW)
	})

	t.Run("should keep every field width of the layout within its column", func(t *testing.T) {
		// Larguras úteis reais do grid (cm) — coluna simples, dupla e página inteira.
		for _, larguraCm := range []float64{4.83, 4.86, 4.96, 9.94, 10.05} {
			w := cmToPt(larguraCm)
			got := truncText(pdf, strings.Repeat("ABCDEFGHIJ", 200), w)
			assert.LessOrEqual(t, measure(t, pdf, got), w, "largura %.2f cm", larguraCm)
		}
	})

	t.Run("should drop the text when not even the ellipsis fits", func(t *testing.T) {
		assert.Equal(t, "", truncText(pdf, "QUALQUER COISA", 1.0))
	})

	t.Run("should keep the original text when it already fits the line budget", func(t *testing.T) {
		curto := "Inf. Cont.: contrato 4472/2026 | Totais Aproximados dos Tributos: R$ 1.234,56"

		assert.Equal(t, curto, limitLines(pdf, curto, cmToPt(20.00), linhasInfoCompl(0)))
	})

	t.Run("should cut the block text to the lines that fit, ending with an ellipsis", func(t *testing.T) {
		longo := strings.Repeat("informacao complementar do emitente para conferencia fiscal ", 200)

		got := limitLines(pdf, longo, cmToPt(20.00), linhasInfoCompl(0))

		assert.True(t, strings.HasSuffix(got, "..."))
		assert.True(t, strings.HasPrefix(longo, strings.TrimSuffix(got, "...")),
			"should keep the original text, without re-joining the split lines")
		linhas, err := pdf.SplitText(got, cmToPt(20.00))
		require.NoError(t, err)
		assert.LessOrEqual(t, len(linhas), linhasInfoCompl(0))
	})

	t.Run("should measure with the current font size", func(t *testing.T) {
		// Mesmo texto e largura, fonte maior: trunca mais cedo (nunca encolher a fonte).
		texto := strings.Repeat("MUNICIPIO ", 10)
		pdf8 := newMeasurePdf(t, 8)

		got7 := truncText(pdf, texto, cmToPt(4.98))
		got8 := truncText(pdf8, texto, cmToPt(4.98))

		assert.Less(t, len(got8), len(got7))
		assert.LessOrEqual(t, measure(t, pdf8, got8), cmToPt(4.98))
	})
}

// TestOuTraco cobre a nota 12 do item 2.4.5 da NT 008/2026 v1.02: "os campos
// sem informações no XML devem ser preenchidos com um traço (-)". Campo em
// branco e campo preenchido com traço são indistinguíveis para quem lê o
// DANFSe impresso — mas o branco também é indistinguível de um campo que o
// emissor deixou de renderizar.
func TestOuTraco(t *testing.T) {
	t.Run("should fill an absent value with the NT dash", func(t *testing.T) {
		assert.Equal(t, "-", ouTraco(""))
	})

	t.Run("should keep an informed value untouched", func(t *testing.T) {
		assert.Equal(t, "PRESTADOR DE TESTE LTDA", ouTraco("PRESTADOR DE TESTE LTDA"))
	})

	t.Run("should treat blank content as absent", func(t *testing.T) {
		// Espaços vindos do XML imprimem exatamente como o campo vazio.
		assert.Equal(t, "-", ouTraco("   "))
	})
}

// TestFitCell cobre a nota 12 no ponto único por onde passam os 89 campos de
// conteúdo variável do DANFSe: o campo sem dado sai com o traço, não em branco.
func TestFitCell(t *testing.T) {
	t.Run("should write the NT dash when the field has no data", func(t *testing.T) {
		pdf := newMeasurePdf(t, 7)
		pdf.SetX(cmToPt(0.40))
		x := pdf.GetX()

		fitCell(pdf, colX3, "")

		assert.InDelta(t, measure(t, pdf, "-"), pdf.GetX()-x, 0.01,
			"o cursor deve avançar a largura do traço")
	})

	t.Run("should write the informed text as it came from the XML", func(t *testing.T) {
		pdf := newMeasurePdf(t, 7)
		pdf.SetX(cmToPt(0.40))
		x := pdf.GetX()

		fitCell(pdf, colX3, "PRESTADOR DE TESTE LTDA")

		assert.InDelta(t, measure(t, pdf, "PRESTADOR DE TESTE LTDA"), pdf.GetX()-x, 0.01)
	})
}

// TestLinhasInfoCompl cobre a repartição do pool de linhas entre os dois campos
// elásticos (issue #2, parte 2): a descrição do serviço desloca os blocos para
// baixo e as informações complementares ficam com o espaço que resta até o
// canhoto (fixo em 28,10 cm) — nunca menos que as 12 linhas do orçamento da NT
// para o xOutInf (1.997 chars).
func TestLinhasInfoCompl(t *testing.T) {
	t.Run("should use every idle line when the description keeps one line", func(t *testing.T) {
		// Sem deslocamento o bloco absorve as 6 linhas ociosas do rodapé: as 17
		// que cabiam nos 4,00 cm do Rect original mais o vão até o canhoto.
		assert.Equal(t, 23, linhasInfoCompl(0))
	})

	t.Run("should give back one line per description line", func(t *testing.T) {
		for extra := 1; extra <= maxLinhasDescServ-1; extra++ {
			assert.Equal(t, 23-extra, linhasInfoCompl(float64(extra)*passoLinha),
				"descrição com %d linha(s) extra", extra)
		}
	})

	t.Run("should never drop below the NT budget of twelve lines", func(t *testing.T) {
		assert.Equal(t, minLinhasInfoCompl, linhasInfoCompl(10.0))
	})
}
