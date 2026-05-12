package main

import (
	"fmt"
	"log"
	"os"

	danfse "github.com/wevertonj/go-danfse-v2"
)

func main() {
	fmt.Println("Reading mock XML data...")
	xmlData, err := os.ReadFile("testdata/mock_danfse.xml")
	if err != nil {
		log.Fatalf("Failed to read mock XML: %v", err)
	}

	fmt.Println("Initializing DANFSe Renderer...")
	renderer := danfse.NewDanfseRenderer()

	fmt.Println("Rendering PDF...")
	pdfBytes, err := renderer.Render(xmlData)
	if err != nil {
		log.Fatalf("Failed to render PDF: %v", err)
	}

	fmt.Println("Saving preview.pdf...")
	err = os.WriteFile("preview.pdf", pdfBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write preview.pdf: %v", err)
	}

	fmt.Println("Successfully generated preview.pdf!")
}
