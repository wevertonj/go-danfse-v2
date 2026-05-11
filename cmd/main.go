package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Municipio struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

func main() {
	fmt.Println("Fetching IBGE cities data...")
	resp, err := http.Get("https://servicodados.ibge.gov.br/api/v1/localidades/municipios")
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		os.Exit(1)
	}

	var muns []Municipio
	if err := json.Unmarshal(body, &muns); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create("ibge.go")
	if err != nil {
		fmt.Printf("Error creating ibge.go file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintln(f, "package danfse")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "// ibgeCities contains the mapping of 7-digit IBGE codes to City names.")
	fmt.Fprintln(f, "var ibgeCities = map[string]string{")
	for _, m := range muns {
		fmt.Fprintf(f, "\t\"%d\": \"%s\",\n", m.ID, m.Nome)
	}
	fmt.Fprintln(f, "}")
	fmt.Printf("Successfully updated ibge.go with %d cities.\n", len(muns))
}
