package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Item struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func ListItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.QueryContext(r.Context(), `SELECT id, name, created_at FROM items ORDER BY id DESC`)
		if err != nil {
			http.Error(w, "query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		items := []Item{}
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Name, &it.CreatedAt); err != nil {
				http.Error(w, "scan failed", http.StatusInternalServerError)
				return
			}
			items = append(items, it)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func CreateItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		res, err := db.ExecContext(r.Context(), `INSERT INTO items (name) VALUES (?)`, body.Name)
		if err != nil {
			http.Error(w, "insert failed", http.StatusInternalServerError)
			return
		}

		id, _ := res.LastInsertId()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Item{ID: id, Name: body.Name})
	}
}
