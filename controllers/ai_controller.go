package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"gorm.io/gorm"
)

type AIController struct {
	DB *gorm.DB
}

func (c *AIController) HandleEmbedding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicaiton/json")
	var reqBody struct {
		Prompt string `json:"prompt"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var prompt string
	if reqBody.Prompt == "" {
		prompt = "I love my India"
	} else {
		prompt = reqBody.Prompt
	}

	body := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}{
		Model:  "all-minilm",
		Prompt: prompt,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	modelResp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	defer modelResp.Body.Close()
	respBody, err := io.ReadAll(modelResp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	type embeddingRespose struct {
		Embeddings []float64 `json:"embedding"`
	}

	var embeddingResp embeddingRespose
	err = json.Unmarshal(respBody, &embeddingResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Embedding []float64 `json:"embedding"`
		Size      int       `json:"size"`
	}{
		Embedding: embeddingResp.Embeddings,
		Size:      len(embeddingResp.Embeddings),
	})
}
