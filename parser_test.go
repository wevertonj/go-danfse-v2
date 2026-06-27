package danfse_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wevertonj/go-danfse-v2"
	"os"
	"testing"
)

func TestParseXML(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/mock_danfse.xml")
	require.NoError(t, err)

	nfse, err := danfse.ParseXML(xmlData)
	require.NoError(t, err)

	assert.Equal(t, "NFS355030826048BJSEBRKWX5Z541700000000000459123847291", nfse.InfNFSe.ID)
	assert.Equal(t, "459123", nfse.InfNFSe.NNFSe)
	assert.Equal(t, "8BJSEBRKWX5Z54", nfse.InfNFSe.Emit.CNPJ)
	assert.Equal(t, "PRESTADOR DE TESTE LTDA", nfse.InfNFSe.Emit.XNome)
	assert.Equal(t, "58APBW6N000109", nfse.InfNFSe.DPS.InfDPS.Toma.CNPJ)
	assert.Equal(t, "14500.50", nfse.InfNFSe.Valores.VLiq)
	// Campos de nível infNFSe usados na seção "Serviço Prestado" do DANFSe.
	assert.Equal(t, "São Paulo", nfse.InfNFSe.XLocPrestacao)
	assert.Equal(t, "Análise e desenvolvimento de sistemas.", nfse.InfNFSe.XTribNac)
}
