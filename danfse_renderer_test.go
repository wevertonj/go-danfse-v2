package danfse_test

import (
	"os"
	"os/exec"
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
