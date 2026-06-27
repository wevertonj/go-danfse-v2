package danfse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUFFromIBGE(t *testing.T) {
	assert.Equal(t, "PR", ufFromIBGE("4106902")) // Curitiba
	assert.Equal(t, "SP", ufFromIBGE("3550308")) // São Paulo
	assert.Equal(t, "SP", ufFromIBGE("3509502")) // Campinas
	assert.Equal(t, "", ufFromIBGE(""))
	assert.Equal(t, "", ufFromIBGE("9"))
	assert.Equal(t, "", ufFromIBGE("9999999")) // prefixo inexistente
}

func TestFormatMunUF_DerivaUFQuandoAusente(t *testing.T) {
	// Sem UF no XML (caso do tomador): deriva do código IBGE.
	assert.Equal(t, "São Paulo / SP", formatMunUF("São Paulo", "3550308", ""))
	// UF explícita tem precedência sobre a derivada.
	assert.Equal(t, "Curitiba / PR", formatMunUF("Curitiba", "4106902", "PR"))
	// Sem município nem UF derivável: "-".
	assert.Equal(t, "-", formatMunUF("", "", ""))
}
