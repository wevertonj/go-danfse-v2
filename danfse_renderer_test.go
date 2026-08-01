package danfse_test

import (
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wevertonj/go-danfse-v2"
)

// pdfText renderiza o XML e extrai o texto via pdftotext (pula se ausente).
func pdfText(t *testing.T, xmlData []byte) string {
	t.Helper()
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext indisponível")
	}
	pdf, err := danfse.NewDanfseRenderer().Render(xmlData)
	require.NoError(t, err)
	f, err := os.CreateTemp(t.TempDir(), "danfse-*.pdf")
	require.NoError(t, err)
	_, err = f.Write(pdf)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	out, err := exec.Command("pdftotext", "-layout", f.Name(), "-").Output()
	require.NoError(t, err)
	return string(out)
}

// TestRenderer_SemValidadeJuridica_DependeDoTpAmb garante que o aviso
// "SEM VALIDADE JURÍDICA" aparece só em homologação (tpAmb=2), não em produção
// (tpAmb=1) — independentemente do AmbGer (que é 2 = SEFIN Nacional em ambos).
func TestRenderer_SemValidadeJuridica_DependeDoTpAmb(t *testing.T) {
	prod, err := os.ReadFile("testdata/mock_danfse.xml") // mock é produção: tpAmb=1, ambGer=2
	require.NoError(t, err)

	assert.NotContains(t, pdfText(t, prod), "SEM VALIDADE JURÍDICA",
		"produção (tpAmb=1) não deve exibir o aviso")

	homolog := []byte(strings.Replace(string(prod), "<tpAmb>1</tpAmb>", "<tpAmb>2</tpAmb>", 1))
	assert.Contains(t, pdfText(t, homolog), "SEM VALIDADE JURÍDICA",
		"homologação (tpAmb=2) deve exibir o aviso")
}

// TestRenderer_HeaderUF cobre a resolução da UF do cabeçalho (issue #1): vem do
// <UF> do emitente, é derivada do código IBGE quando ausente e vira erro quando
// não é derivável — nunca um valor fixo inventado. O separador município–UF
// segue a NT 008/2026 v1.02, tabela 2.4.5 ("Município: CCCC / CC").
func TestRenderer_HeaderUF(t *testing.T) {
	mock, err := os.ReadFile("testdata/mock_danfse.xml")
	require.NoError(t, err)

	t.Run("should print the emitter UF from the XML separated by a slash", func(t *testing.T) {
		assert.Contains(t, pdfText(t, mock), "Município: São Paulo / SP")
	})

	t.Run("should derive the UF from the IBGE code when the XML omits it", func(t *testing.T) {
		semUF := []byte(strings.Replace(string(mock), "<UF>SP</UF>", "", 1))
		assert.Contains(t, pdfText(t, semUF), "Município: São Paulo / SP",
			"should derive SP from cMun 3550308 instead of falling back to a hardcoded UF")
	})

	t.Run("should fail when the UF is absent and not derivable from cMun", func(t *testing.T) {
		semUF := strings.Replace(string(mock), "<UF>SP</UF>", "", 1)
		cMunInvalido := strings.Replace(semUF, "<cMun>3550308</cMun>", "<cMun>9999999</cMun>", 1)

		pdfBytes, err := danfse.NewDanfseRenderer().Render([]byte(cMunInvalido))

		require.Error(t, err, "should not emit a DANFSe with a fabricated UF")
		assert.Nil(t, pdfBytes)
		assert.Contains(t, err.Error(), "UF")
		assert.Contains(t, err.Error(), "9999999")
	})
}

// TestRenderer_QRCode cobre a geração do código QR de verificação (issue #1):
// nenhuma falha da cadeia (codificação, PNG, holder, desenho) pode ser engolida —
// emitir um DANFSe sem QR tira do documento a verificação de autenticidade.
func TestRenderer_QRCode(t *testing.T) {
	mock, err := os.ReadFile("testdata/mock_danfse.xml")
	require.NoError(t, err)

	t.Run("should fail when the access key overflows the QR code capacity", func(t *testing.T) {
		// A URL de consulta leva a chave inteira; acima da capacidade da versão 40
		// (5.596 dígitos em modo numérico, nível M) o encoder falha.
		chaveGigante := "NFS" + strings.Repeat("1", 6000)
		xmlGigante := []byte(strings.Replace(string(mock),
			`Id="NFS355030826048BJSEBRKWX5Z541700000000000459123847291"`,
			`Id="`+chaveGigante+`"`, 1))

		pdfBytes, err := danfse.NewDanfseRenderer().Render(xmlGigante)

		require.Error(t, err, "should not emit a DANFSe without the verification QR code")
		assert.Nil(t, pdfBytes)
		assert.Contains(t, err.Error(), "qrcode")
	})

	t.Run("should render a valid PDF when the QR code succeeds", func(t *testing.T) {
		pdfBytes, err := danfse.NewDanfseRenderer().Render(mock)

		require.NoError(t, err)
		require.NotEmpty(t, pdfBytes)
		assert.Equal(t, "%PDF-", string(pdfBytes[:5]))
	})
}

// palavra é uma palavra do PDF com sua caixa delimitadora em cm.
type palavra struct {
	texto      string
	xMin, xMax float64
	yMin, yMax float64
}

// pdfWords extrai as palavras do PDF com coordenadas via `pdftotext -bbox`.
// Diferente de `-layout`, o bbox revela sobreposição e transbordo: o texto
// invasor sai igual no modo layout, mas com xMax fora da coluna.
func pdfWords(t *testing.T, xmlData []byte) []palavra {
	t.Helper()
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext indisponível")
	}
	pdf, err := danfse.NewDanfseRenderer().Render(xmlData)
	require.NoError(t, err)
	f, err := os.CreateTemp(t.TempDir(), "danfse-*.pdf")
	require.NoError(t, err)
	_, err = f.Write(pdf)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	out, err := exec.Command("pdftotext", "-bbox", f.Name(), "-").Output()
	require.NoError(t, err)

	re := regexp.MustCompile(`<word xMin="([\d.]+)" yMin="([\d.]+)" xMax="([\d.]+)" yMax="([\d.]+)">([^<]*)</word>`)
	const ptPorCm = 28.34645669
	var palavras []palavra
	for _, m := range re.FindAllStringSubmatch(string(out), -1) {
		num := func(s string) float64 {
			v, err := strconv.ParseFloat(s, 64)
			require.NoError(t, err)
			return v / ptPorCm
		}
		palavras = append(palavras, palavra{
			texto: m[5], xMin: num(m[1]), xMax: num(m[3]), yMin: num(m[2]), yMax: num(m[4]),
		})
	}
	require.NotEmpty(t, palavras, "pdftotext -bbox não retornou palavras")
	return palavras
}

// TestRenderer_TextosLongos cobre a issue #2: com todo campo variável preenchido
// acima da capacidade da coluna, nenhum texto pode invadir a coluna vizinha nem
// a borda da página. A NT 008/2026 v1.02 (item 2.1) manda truncar com
// reticências mantendo a fonte, cujos tamanhos do item 2.4 são mínimos.
func TestRenderer_TextosLongos(t *testing.T) {
	longo, err := os.ReadFile("testdata/mock_danfse_textos_longos.xml")
	require.NoError(t, err)
	palavras := pdfWords(t, longo)

	// Faixas Y (cm) dos campos que a NT dimensiona com a largura da página
	// inteira (20,40 cm) e que, portanto, atravessam legitimamente as colunas:
	// descrição do código de tributação, descrição do serviço e informações
	// complementares.
	larguraTotal := func(y float64) bool {
		return (y > 13.30 && y < 15.00) || y > 22.50
	}

	t.Run("should keep every word inside the page border", func(t *testing.T) {
		for _, p := range palavras {
			assert.LessOrEqual(t, p.xMax, 20.70, "palavra %q ultrapassa a borda direita", p.texto)
		}
	})

	t.Run("should not let a left-hand field invade the right-hand columns", func(t *testing.T) {
		// Nome, endereço e e-mail ocupam duas colunas (NT: 10,19 cm); nenhum
		// campo que começa à esquerda pode cruzar a divisória de 10,51 cm.
		// Fora da regra: o cabeçalho e a chave de acesso (y < 2,20), que têm
		// blocos próprios, e os campos de largura total.
		const tol = 0.05 // folga para o arredondamento do pdftotext
		for _, p := range palavras {
			if p.yMin < 2.20 || p.xMin >= 10.51-tol || larguraTotal(p.yMin) {
				continue
			}
			assert.LessOrEqual(t, p.xMax, 10.51+tol, "palavra %q (y=%.2f) invade a coluna da direita", p.texto, p.yMin)
		}
	})

	t.Run("should keep the access key clear of the QR code", func(t *testing.T) {
		for _, p := range palavras {
			if p.yMin < 1.70 || p.yMin > 2.20 {
				continue
			}
			assert.LessOrEqual(t, p.xMax, 17.48, "chave de acesso invade o quadro do QR code")
		}
	})

	t.Run("should mark the truncated fields with the NT ellipsis", func(t *testing.T) {
		var comReticencias int
		for _, p := range palavras {
			if strings.HasSuffix(p.texto, "...") {
				comReticencias++
			}
		}
		assert.GreaterOrEqual(t, comReticencias, 10,
			"os campos cortados devem terminar com as reticências da NT")
	})

	t.Run("should flag the suppressed complementary information with an ellipsis", func(t *testing.T) {
		// O MultiCell descarta as linhas que não cabem na altura do bloco sem
		// deixar sinal: o texto terminaria no meio de uma frase como se fosse
		// esse o conteúdo do XML. A NT (item 2.1) manda marcar com reticências.
		gigante := []byte(strings.Replace(string(longo),
			"<xOutInf>", "<xOutInf>"+strings.Repeat("informacao complementar do emitente ", 300), 1))

		// Faixa do bloco: do início do conteúdo (22,68) à linha do canhoto (28,10).
		var ultima palavra
		for _, p := range pdfWords(t, gigante) {
			if p.yMin > 22.50 && p.yMax <= 28.10 && p.yMin >= ultima.yMin {
				ultima = p
			}
		}

		assert.True(t, strings.HasSuffix(ultima.texto, "..."),
			"o bloco deve terminar com as reticências da NT, terminou em %q", ultima.texto)
		assert.LessOrEqual(t, ultima.yMax, 28.10, "as informações complementares invadem o canhoto")
	})

	t.Run("should keep the neighbouring columns readable", func(t *testing.T) {
		texto := pdfText(t, longo)

		assert.Contains(t, texto, "3550308 / 01.310-200", "código IBGE / CEP do prestador")
		assert.Contains(t, texto, "(11) 99999-9999", "telefone do prestador")
		assert.Contains(t, texto, "459123", "número da NFS-e")
	})
}

func TestRenderer_Render(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/mock_danfse.xml")
	require.NoError(t, err)

	renderer := danfse.NewDanfseRenderer()

	pdfBytes, err := renderer.Render(xmlData)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)

	assert.True(t, len(pdfBytes) > 100)
	assert.Equal(t, "%PDF-", string(pdfBytes[:5]))
}

// TestRenderer_ServicoPrestado_DescricaoELocal garante que a seção "Serviço
// Prestado" do DANFSe renderiza a descrição do código de tributação nacional
// (infNFSe/xTribNac) e o local da prestação (infNFSe/xLocPrestacao), que vêm do
// nível infNFSe — não do <serv> da DPS.
func TestRenderer_ServicoPrestado_DescricaoELocal(t *testing.T) {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext indisponível")
	}

	xmlData, err := os.ReadFile("testdata/mock_danfse.xml")
	require.NoError(t, err)
	pdf, err := danfse.NewDanfseRenderer().Render(xmlData)
	require.NoError(t, err)

	f, err := os.CreateTemp(t.TempDir(), "danfse-*.pdf")
	require.NoError(t, err)
	_, err = f.Write(pdf)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	out, err := exec.Command("pdftotext", "-layout", f.Name(), "-").Output()
	require.NoError(t, err)
	text := string(out)

	assert.Contains(t, text, "Análise e desenvolvimento de sistemas.",
		"a descrição do código de tributação nacional deve aparecer no DANFSe")
	assert.Contains(t, text, "São Paulo / SP",
		"o local da prestação deve aparecer com a UF derivada do código IBGE")
}
