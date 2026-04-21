package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// ── Test helpers ──────────────────────────────────────────────────────────────

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	{
		api.POST("/profiles", createProfile)
		api.GET("/profiles", listProfiles)
		api.GET("/profiles/:id", getProfile)
		api.DELETE("/profiles/:id", deleteProfile)
	}
	return r
}

// resetStore clears the in-memory store between tests.
func resetStore() {
	storeMu.Lock()
	store = map[string]*Profile{}
	storeMu.Unlock()
}

// seedProfile inserts a profile directly into the store, bypassing HTTP.
func seedProfile(p *Profile) {
	storeMu.Lock()
	store[p.Name] = p
	storeMu.Unlock()
}

// ── classifyAge ───────────────────────────────────────────────────────────────

func TestClassifyAge(t *testing.T) {
	cases := []struct {
		age  int
		want string
	}{
		{0, "child"},
		{12, "child"},
		{13, "teenager"},
		{19, "teenager"},
		{20, "adult"},
		{59, "adult"},
		{60, "senior"},
		{100, "senior"},
	}
	for _, tc := range cases {
		got := classifyAge(tc.age)
		if got != tc.want {
			t.Errorf("classifyAge(%d) = %q; want %q", tc.age, got, tc.want)
		}
	}
}

// ── POST /api/profiles ────────────────────────────────────────────────────────

func TestCreateProfile_MissingName(t *testing.T) {
	defer resetStore()
	r := setupRouter()

	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/api/profiles", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateProfile_Idempotent(t *testing.T) {
	defer resetStore()

	// Pre-seed a profile so createProfile short-circuits to the idempotency path.
	existing := &Profile{
		ID:   "fixed-id",
		Name: "Alice",
	}
	seedProfile(existing)

	r := setupRouter()
	body := bytes.NewBufferString(`{"name":"Alice"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/profiles", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for existing profile, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	data := resp["data"].(map[string]any)
	if data["id"] != "fixed-id" {
		t.Errorf("expected existing id %q, got %q", "fixed-id", data["id"])
	}
}

// ── GET /api/profiles ─────────────────────────────────────────────────────────

func TestListProfiles_Empty(t *testing.T) {
	defer resetStore()
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/profiles", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	data := resp["data"].([]any)
	if len(data) != 0 {
		t.Errorf("expected empty list, got %d items", len(data))
	}
}

func TestListProfiles_NonEmpty(t *testing.T) {
	defer resetStore()
	seedProfile(&Profile{ID: "id-1", Name: "Bob"})
	seedProfile(&Profile{ID: "id-2", Name: "Carol"})

	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/profiles", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	data := resp["data"].([]any)
	if len(data) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(data))
	}
}

// ── GET /api/profiles/:id ─────────────────────────────────────────────────────

func TestGetProfile_Found(t *testing.T) {
	defer resetStore()
	seedProfile(&Profile{ID: "abc123", Name: "Dave"})

	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/profiles/abc123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	data := resp["data"].(map[string]any)
	if data["id"] != "abc123" {
		t.Errorf("expected id %q, got %v", "abc123", data["id"])
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	defer resetStore()
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// ── DELETE /api/profiles/:id ──────────────────────────────────────────────────

func TestDeleteProfile_Found(t *testing.T) {
	defer resetStore()
	seedProfile(&Profile{ID: "del-1", Name: "Eve"})

	r := setupRouter()
	req := httptest.NewRequest(http.MethodDelete, "/api/profiles/del-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Confirm it's gone
	storeMu.RLock()
	_, stillThere := store["Eve"]
	storeMu.RUnlock()
	if stillThere {
		t.Error("profile should have been removed from the store")
	}
}

func TestDeleteProfile_NotFound(t *testing.T) {
	defer resetStore()
	r := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/api/profiles/ghost", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
