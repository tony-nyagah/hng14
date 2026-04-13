package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type GenderizeResponse struct {
	Name        string  `json:"name"`
	Gender      *string `json:"gender"`
	Probability float64 `json:"probability"`
	Count       *int    `json:"count"`
}

type ClassifyData struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	SampleSize  int     `json:"sample_size"`
	IsConfident bool    `json:"is_confident"`
	ProcessedAt string  `json:"processed_at"`
}

type SuccessResponse struct {
	Status string       `json:"status"`
	Data   ClassifyData `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func errResponse(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Status: "error", Message: msg})
}

func classifyHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		errResponse(w, http.StatusBadRequest, "Missing or empty name parameter")
		return
	}

	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		errResponse(w, http.StatusBadGateway, "Failed to reach upstream API")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "Failed to read upstream response")
		return
	}

	var gr GenderizeResponse
	if err := json.Unmarshal(body, &gr); err != nil {
		errResponse(w, http.StatusInternalServerError, "Failed to parse upstream response")
		return
	}

	if gr.Gender == nil || gr.Count == nil || *gr.Count == 0 {
		errResponse(w, http.StatusOK, "No prediction available for the provided name")
		return
	}

	isConfident := gr.Probability >= 0.7 && *gr.Count >= 100

	writeJSON(w, http.StatusOK, SuccessResponse{
		Status: "success",
		Data: ClassifyData{
			Name:        gr.Name,
			Gender:      *gr.Gender,
			Probability: gr.Probability,
			SampleSize:  *gr.Count,
			IsConfident: isConfident,
			ProcessedAt: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8059"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/classify", classifyHandler)

	fmt.Printf("Server running on :%s\n", port)
	http.ListenAndServe(":"+port, mux)
}
