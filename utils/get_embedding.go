package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetEmbedding(text string) ([]float32, error) {
	body := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}{
		Model:  "all-minilm",
		Prompt: text,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return make([]float32, 0), err
	}

	modelResp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return make([]float32, 0), err
	}

	defer modelResp.Body.Close()
	respBytes, err := io.ReadAll(modelResp.Body)
	if err != nil {
		return make([]float32, 0), err
	}

	type parsedResponse struct {
		Embedding []float32 `json:"embedding"`
	}

	var parsedResp parsedResponse
	err = json.Unmarshal(respBytes, &parsedResp)
	if err != nil {
		return make([]float32, 0), err
	}

	if len(parsedResp.Embedding) == 0 {
		return make([]float32, 0), fmt.Errorf("Generated embedding of length 0")
	}

	return parsedResp.Embedding, nil
}
