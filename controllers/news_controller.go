package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/pgvector/pgvector-go"
	"github.com/sambeetpanda507/advance-search/models"
	"github.com/sambeetpanda507/advance-search/utils"
	"gorm.io/gorm"
)

type NewsController struct {
	DB *gorm.DB
}

func (c NewsController) GetNewsFromFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if news data already exists
	var count int64
	c.DB.Model(&models.News{}).Count(&count)
	if count > 500 {
		json.NewEncoder(w).Encode(map[string]any{"message": "Ok", "count": count})
		return
	}

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
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errorCount int
	batch := make(chan struct{}, 10)
	for i, record := range records {
		if i == 0 {
			continue
		}

		if i == 1000 {
			break
		}

		combinedText := record[1] + " " + record[2]
		if combinedText == " " {
			continue
		}

		wg.Add(1)
		batch <- struct{}{} // Block if more than 10 goroutines are running
		go func(i int, record []string, text string) {
			defer wg.Done()
			defer func() { <-batch }() // Release the batch
			fmt.Printf("Running for i = %d\n", i)
			embedding, err := utils.GetEmbedding(text)
			if err != nil {
				mu.Lock()
				errorCount++
				fmt.Printf("Error at i = %d: %v\n", i, err)
				mu.Unlock()
				return
			}

			row := models.News{Title: record[1], Description: record[2], Embedding: pgvector.NewVector(embedding)}
			mu.Lock()
			news = append(news, row)
			mu.Unlock()
		}(i, record, combinedText)
	}

	wg.Wait()
	fmt.Printf("Total embedding error = %d\n", errorCount)

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

func (c NewsController) GetAllNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	pageStr := query.Get("page")
	limitStr := query.Get("limit")
	search := query.Get("search")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	type NewsData struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	news := []NewsData{}
	var result *gorm.DB
	if len(search) > 0 {
		if err := c.DB.Exec("SET pg_trgrm.similarity_threshold = 0.2; -- update the threshold from 0.3 to 0.2").Error; err != nil {
			log.Fatalf("Error while setting similarity_threshold, %s", err.Error())
		}

		embedding, err := utils.GetEmbedding(search)
		if err != nil {
			log.Fatalf("Error while creating embedding, %s", err.Error())
		}

		sql := `
			SELECT * FROM (
				SELECT
					ID,
					TITLE,
					DESCRIPTION,
					WORD_SIMILARITY  (?, COALESCE(TITLE, '') || ' ' || COALESCE(DESCRIPTION,'')) AS SIMILARITY_RANK,
					TS_RANK(
						SEARCH_VECTOR,
						WEBSEARCH_TO_TSQUERY('english', ?)
					) AS RANK,
					EMBEDDING <=> ? AS DISTANCE
				FROM
					NEWS
				WHERE
					EMBEDDING <=> ? < 0.8
					OR SEARCH_VECTOR @@ WEBSEARCH_TO_TSQUERY('english', ?)
					OR WORD_SIMILARITY(?, COALESCE(TITLE,'') || ' ' || COALESCE(DESCRIPTION, '')) > 0.25
			) AS T
			ORDER BY
				T.DISTANCE ASC,
				(T.RANK * 2 + T.SIMILARITY_RANK) DESC
			OFFSET
				?	
			LIMIT
				?;
		`
		result = c.DB.Raw(
			sql,
			search,
			search,
			pgvector.NewVector(embedding),
			pgvector.NewVector(embedding),
			search,
			search,
			page*limit,
			limit,
		).Scan(&news)
	} else {
		result = c.DB.Model(&models.News{}).Limit(limit).Offset(page * limit).Find(&news)
	}

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": result.Error.Error()})
	}

	json.NewEncoder(w).Encode(news)
}

func (c NewsController) HandleFillEmbedding(w http.ResponseWriter, r *http.Request) {
	var newsItems []models.News
	c.DB.Where("embedding IS NULL").FindInBatches(&newsItems, 100, func(tx *gorm.DB, batch int) error {
		for _, item := range newsItems {
			fmt.Printf("Processing ID %d: %s\n", item.ID, item.Title)
			combinedText := item.Title + " " + item.Description
			if combinedText == " " {
				continue
			}

			embedding, err := utils.GetEmbedding(combinedText)
			if err != nil {
				fmt.Printf("Error generating embedding for %d\n", item.ID)
				continue
			}

			err = tx.Model(&item).Select("Embedding").Updates(models.News{Embedding: pgvector.NewVector(embedding)}).Error
			if err != nil {
				fmt.Printf("Error updating embedding for %d\n", item.ID)
				continue
			}

			fmt.Printf("Successfully updated embedding for id = %d\n", item.ID)
		}

		return nil
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ok")
}
