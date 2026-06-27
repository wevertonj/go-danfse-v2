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
