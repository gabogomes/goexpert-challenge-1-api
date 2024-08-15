package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbTimeout  = 10 * time.Millisecond
	apiTimeout = 200 * time.Millisecond
	httpPort   = ":8080"
)

type Usdbrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Cotacao struct {
	Usdbrl Usdbrl `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", Handler)
	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(httpPort, nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), apiTimeout)
	defer cancel()

	cotacao, err := FetchCotacao(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch cotacao", http.StatusInternalServerError)
		log.Println("Error fetching cotacao:", err)
		return
	}

	err = SaveCotacao(ctx, cotacao.Usdbrl.Bid)
	if err != nil {
		http.Error(w, "Failed to save cotacao", http.StatusInternalServerError)
		log.Println("Error saving cotacao:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": cotacao.Usdbrl.Bid})
}

func FetchCotacao(ctx context.Context) (*Cotacao, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cotacao Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func SaveCotacao(ctx context.Context, bid string) error {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)", bid)
	if err != nil {
		return err
	}

	return nil
}
