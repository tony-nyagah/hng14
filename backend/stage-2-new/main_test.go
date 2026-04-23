package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAdvancedFiltering ensures combined filters (gender + country) work correctly
func TestAdvancedFiltering(t *testing.T) {
	initDB()
	seedDatabase() // Ensure data is present

	req, _ := http.NewRequest("GET", "/api/profiles?gender=male&country_id=NG&min_age=25", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profilesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response.Status != "success" {
		t.Errorf("Expected success status, got %s", response.Status)
	}

	// Check if data matches filters
	data := response.Data.([]interface{})
	for _, item := range data {
		p := item.(map[string]interface{})
		if p["gender"] != "male" {
			t.Errorf("Filter failed: expected male, got %v", p["gender"])
		}
		if p["country_id"] != "NG" {
			t.Errorf("Filter failed: expected NG, got %v", p["country_id"])
		}
		if p["age"].(float64) < 25 {
			t.Errorf("Filter failed: expected age >= 25, got %v", p["age"])
		}
	}
}

// TestPagination ensures the total count and limit logic are correct
func TestPagination(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/profiles?page=1&limit=5", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profilesHandler)

	handler.ServeHTTP(rr, req)

	var response APIResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response.Limit != 5 {
		t.Errorf("Expected limit 5, got %v", response.Limit)
	}

	data := response.Data.([]interface{})
	if len(data) > 5 {
		t.Errorf("Pagination failed: returned %d items, limit was 5", len(data))
	}
}

// TestNLPParsing checks if the rule-based engine interprets queries correctly
func TestNLPParsing(t *testing.T) {
	testCases := []struct {
		query    string
		expected map[string]string
	}{
		{"young males from nigeria", map[string]string{"gender": "male", "min_age": "16", "max_age": "24", "country_id": "NG"}},
		{"females above 30", map[string]string{"gender": "female", "min_age": "31"}},
		{"adult males from kenya", map[string]string{"gender": "male", "age_group": "adult", "country_id": "KE"}},
	}

	for _, tc := range testCases {
		filters, ok := parseNLP(tc.query)
		if !ok {
			t.Errorf("NLP failed to interpret query: %s", tc.query)
		}

		for key, val := range tc.expected {
			if filters[key] != val {
				t.Errorf("Query [%s]: Expected %s=%s, got %s", tc.query, key, val, filters[key])
			}
		}
	}
}

// TestInvalidQueryResponse ensures 400 error on uninterpretable search
func TestInvalidQueryResponse(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/profiles/search?q=somethingrandom", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(profilesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid query, got %v", status)
	}

	var response APIResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response.Status != "error" {
		t.Errorf("Expected error status, got %s", response.Status)
	}
}
