package danfse_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wevertonj/go-danfse-v2"
)

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
