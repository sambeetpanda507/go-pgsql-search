package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sambeetpanda507/advance-search/models"
	"gorm.io/gorm"
)

type NewController struct {
	DB *gorm.DB
}

func (c NewController) GetNewsFromFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	file, err := os.Open("assets/train.csv")
	if err != nil {
		respondError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	news := []models.News{}
	for i, record := range records {
		if i == 0 {
			continue
		}
		if i == 200 {
			break
		}
		row := models.News{Title: record[1], Description: record[2]}
		news = append(news, row)
	}

	// Write news to database
	result := c.DB.Create(news)
	if result.Error != nil {
		respondError(w, http.StatusInternalServerError, map[string]string{"message": result.Error.Error()})
		return
	}

	fmt.Fprintf(w, "Rows affected: %d", result.RowsAffected)
}

func respondError(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
