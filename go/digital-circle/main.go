package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "031126"
	dbname   = "postgres"
)

type Entry struct {
	ID       int       `json:"id"`
	ColTexto string    `json:"col_texto"`
	ColDt    time.Time `json:"col_dt"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/tb01", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		//Novo endpoint inserido de acordo com solicitação do Haroldo para retornar os 10 últimos registros inseridos.
		if r.Method == http.MethodGet {
			rows, err := db.Query("SELECT id, col_texto, col_dt FROM TB01 ORDER BY col_dt DESC LIMIT 10")
			if err != nil {
				http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var entries []Entry
			for rows.Next() {
				var e Entry
				if err := rows.Scan(&e.ID, &e.ColTexto, &e.ColDt); err != nil {
					http.Error(w, "Erro ao ler dados", http.StatusInternalServerError)
					return
				}
				entries = append(entries, e)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(entries)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
			return
		}

		var e Entry
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
			return
		}

		var id int
		var colDt time.Time
		err = db.QueryRow("INSERT INTO TB01 (col_texto) VALUES ($1) RETURNING id, col_dt", e.ColTexto).Scan(&id, &colDt)
		if err != nil {
			http.Error(w, "Erro ao inserir no banco de dados", http.StatusInternalServerError)
			return
		}

		e.ID = id
		e.ColDt = colDt

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(e)
	})

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
