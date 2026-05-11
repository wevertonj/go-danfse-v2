package danfse_test

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wevertonj/go-danfse-v2"
)

func TestRenderer_Render(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/mock_danfse.xml")
	require.NoError(t, err)

	renderer := danfse.NewDanfseRenderer("fonts/LiberationSans-Regular.ttf", "fonts/LiberationSans-Bold.ttf", "assets/logo-nfse.png")
	
	pdfBytes, err := renderer.Render(xmlData)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)

	assert.True(t, len(pdfBytes) > 100)
	assert.Equal(t, "%PDF-", string(pdfBytes[:5]))
}
