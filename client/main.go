package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverURL     = "http://localhost:8080/cotacao"
	clientTimeout = 300 * time.Millisecond
)

type BidResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	cotacao, err := FetchCotacaoClient(ctx)
	if err != nil {
		log.Println("Error fetching cotacao:", err)
		return
	}

	// Log para verificar se o campo Bid está presente
	log.Printf("Cotacao Bid: %s\n", cotacao.Bid)

	err = SaveCotacaoToFile(cotacao)
	if err != nil {
		log.Println("Error saving cotacao to file:", err)
	} else {
		log.Println("Cotacao salva com sucesso no arquivo 'cotacao.txt'")
	}
}

func FetchCotacaoClient(ctx context.Context) (*BidResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cotacao BidResponse
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func SaveCotacaoToFile(cotacao *BidResponse) error {
	fileContent := "Dólar: " + cotacao.Bid

	// Log para verificar o conteúdo que será salvo
	log.Printf("Conteúdo a ser salvo no arquivo: %s\n", fileContent)

	err := os.WriteFile("cotacao.txt", []byte(fileContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
