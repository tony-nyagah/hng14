package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ── Test helpers ──────────────────────────────────────────────────────────────

func setupTestDB(t *testing.T) {
	t.Helper()
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&Profile{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
}

func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	setupTestDB(t)
	return httptest.NewServer(setupRouter())
}

func post(t *testing.T, srv *httptest.Server, path string, body any) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	resp, err := http.Post(srv.URL+path, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

func get(t *testing.T, srv *httptest.Server, path string) *http.Response {
	t.Helper()
	resp, err := http.Get(srv.URL + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	return resp
}

func decode(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	defer resp.Body.Close()
	var m map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return m
}

func seedProfile(t *testing.T, srv *httptest.Server, name, gender string, age int, countryID string) map[string]any {
	t.Helper()
	resp := post(t, srv, "/api/profiles", map[string]any{
		"name":       name,
		"gender":     gender,
		"age":        age,
		"country_id": countryID,
	})
	body := decode(t, resp)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Fatalf("seedProfile %q: unexpected status %d", name, resp.StatusCode)
	}
	return body
}

func assertStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		t.Errorf("status: got %d, want %d", resp.StatusCode, want)
	}
}

func assertField(t *testing.T, body map[string]any, key, want string) {
	t.Helper()
	got := fmt.Sprint(body[key])
	if got != want {
		t.Errorf("body[%q]: got %q, want %q", key, got, want)
	}
}

func assertDataField(t *testing.T, body map[string]any, key, want string) {
	t.Helper()
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("body[\"data\"] is not an object: %v", body["data"])
	}
	got := fmt.Sprint(data[key])
	if got != want {
		t.Errorf("data[%q]: got %q, want %q", key, got, want)
	}
}

func getDataList(t *testing.T, body map[string]any) []map[string]any {
	t.Helper()
	raw, ok := body["data"].([]any)
	if !ok {
		t.Fatalf("body[\"data\"] is not an array: %v", body["data"])
	}
	out := make([]map[string]any, len(raw))
	for i, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("data[%d] is not an object", i)
		}
		out[i] = m
	}
	return out
}

// ── POST /api/profiles ────────────────────────────────────────────────────────

func TestCreateProfile_FullData(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp := post(t, srv, "/api/profiles", map[string]any{
		"name":                "Alice Kamau",
		"gender":              "female",
		"gender_probability":  0.95,
		"age":                 28,
		"age_group":           "adult",
		"country_id":          "KE",
		"country_name":        "Kenya",
		"country_probability": 0.88,
	})
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusCreated)
	assertField(t, body, "status", "success")
	assertDataField(t, body, "gender", "female")
	assertDataField(t, body, "age", "28")
	assertDataField(t, body, "country_id", "KE")
	assertDataField(t, body, "age_group", "adult")
	assertDataField(t, body, "country_name", "Kenya")
}

func TestCreateProfile_DerivesAgeGroup(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	cases := []struct {
		age       int
		wantGroup string
	}{
		{5, "child"},
		{15, "teenager"},
		{30, "adult"},
		{65, "senior"},
	}
	for _, tc := range cases {
		resp := post(t, srv, "/api/profiles", map[string]any{
			"name":   fmt.Sprintf("AgeTest%d", tc.age),
			"gender": "male",
			"age":    tc.age,
		})
		body := decode(t, resp)
		assertStatus(t, resp, http.StatusCreated)
		assertDataField(t, body, "age_group", tc.wantGroup)
	}
}

func TestCreateProfile_DerivesCountryName(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp := post(t, srv, "/api/profiles", map[string]any{
		"name":       "Bob Okafor",
		"gender":     "male",
		"age":        35,
		"country_id": "NG",
	})
	body := decode(t, resp)
	assertStatus(t, resp, http.StatusCreated)
	assertDataField(t, body, "country_name", "Nigeria")
}

func TestCreateProfile_Idempotent(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	payload := map[string]any{"name": "Idempotent User", "gender": "male", "age": 25, "country_id": "GH"}
	post(t, srv, "/api/profiles", payload) // first → 201

	resp := post(t, srv, "/api/profiles", payload) // second → 200
	assertStatus(t, resp, http.StatusOK)
}

func TestCreateProfile_UpsertStaleProfile(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	// Insert a stale profile with empty data (simulates a previous failed attempt)
	db.Exec(
		`INSERT INTO profiles (id,name,gender,gender_probability,age,age_group,country_id,country_name,country_probability,created_at)
		 VALUES (?,?,?,?,?,?,?,?,?,datetime('now'))`,
		newUUID(), "Stale Profile", "", 0, 0, "child", "", "", 0,
	)

	// POST with full data should UPDATE the stale record
	resp := post(t, srv, "/api/profiles", map[string]any{
		"name":       "Stale Profile",
		"gender":     "female",
		"age":        30,
		"country_id": "TZ",
	})
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	assertDataField(t, body, "gender", "female")
	assertDataField(t, body, "age", "30")
	assertDataField(t, body, "country_id", "TZ")
}

func TestCreateProfile_MissingName(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp := post(t, srv, "/api/profiles", map[string]any{"gender": "male"})
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusBadRequest)
	assertField(t, body, "status", "error")
}

func TestCreateProfile_ResponseKeyOrder(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp := post(t, srv, "/api/profiles", map[string]any{
		"name": "KeyOrder Test", "gender": "male", "age": 20, "country_id": "UG",
	})
	body := decode(t, resp)

	assertField(t, body, "status", "success")
	if _, ok := body["data"]; !ok {
		t.Error("response missing 'data' field")
	}
}

// ── GET /api/profiles ─────────────────────────────────────────────────────────

func TestListProfiles_PaginationEnvelope(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	for i := 0; i < 10; i++ {
		seedProfile(t, srv, fmt.Sprintf("PagUser%d", i), "male", 20+i, "NG")
	}

	resp := get(t, srv, "/api/profiles?page=1&limit=5")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	assertField(t, body, "status", "success")

	if body["page"] == nil {
		t.Error("missing 'page' field")
	}
	if body["limit"] == nil {
		t.Error("missing 'limit' field")
	}
	if body["total"] == nil {
		t.Error("missing 'total' field")
	}
	if body["data"] == nil {
		t.Error("missing 'data' field")
	}

	items := getDataList(t, body)
	if len(items) != 5 {
		t.Errorf("data length: got %d, want 5", len(items))
	}

	if fmt.Sprint(body["page"]) != "1" {
		t.Errorf("page: got %v, want 1", body["page"])
	}
	if fmt.Sprint(body["limit"]) != "5" {
		t.Errorf("limit: got %v, want 5", body["limit"])
	}
}

func TestListProfiles_Page2DifferentFromPage1(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	for i := 0; i < 12; i++ {
		seedProfile(t, srv, fmt.Sprintf("StableUser%d", i), "female", 25+i, "KE")
	}

	ids := func(page int) []string {
		body := decode(t, get(t, srv, fmt.Sprintf("/api/profiles?page=%d&limit=5", page)))
		items := getDataList(t, body)
		out := make([]string, len(items))
		for i, item := range items {
			out[i] = fmt.Sprint(item["id"])
		}
		return out
	}

	p1 := ids(1)
	p2 := ids(2)

	if len(p1) == 0 || len(p2) == 0 {
		t.Fatal("pages should not be empty")
	}
	for _, id1 := range p1 {
		for _, id2 := range p2 {
			if id1 == id2 {
				t.Errorf("duplicate profile id %q found in both page 1 and page 2", id1)
			}
		}
	}
}

func TestListProfiles_LimitCap(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	for i := 0; i < 60; i++ {
		seedProfile(t, srv, fmt.Sprintf("CapUser%d", i), "male", 25, "GH")
	}

	resp := get(t, srv, "/api/profiles?limit=200")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)

	items := getDataList(t, body)
	if len(items) > 50 {
		t.Errorf("data length %d exceeds cap of 50", len(items))
	}

	limit := body["limit"].(float64)
	if limit > 50 {
		t.Errorf("limit in response is %v, should be <= 50", limit)
	}
}

func TestListProfiles_DefaultPagination(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	for i := 0; i < 15; i++ {
		seedProfile(t, srv, fmt.Sprintf("DefUser%d", i), "female", 30, "ET")
	}

	resp := get(t, srv, "/api/profiles")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	items := getDataList(t, body)
	if len(items) > 10 {
		t.Errorf("default limit should be 10, got %d items", len(items))
	}
}

// ── Filtering ─────────────────────────────────────────────────────────────────

func TestFilterByGender(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "FilterMale1", "male", 30, "NG")
	seedProfile(t, srv, "FilterMale2", "male", 25, "NG")
	seedProfile(t, srv, "FilterFemale1", "female", 28, "NG")

	resp := get(t, srv, "/api/profiles?gender=male")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		if p["gender"] != "male" {
			t.Errorf("expected gender=male, got %v", p["gender"])
		}
	}
}

func TestFilterByCountryID(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "NGUser1", "male", 30, "NG")
	seedProfile(t, srv, "NGUser2", "female", 25, "NG")
	seedProfile(t, srv, "KEUser1", "male", 28, "KE")

	resp := get(t, srv, "/api/profiles?country_id=NG")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		if p["country_id"] != "NG" {
			t.Errorf("expected country_id=NG, got %v", p["country_id"])
		}
	}
}

func TestFilterByAgeGroup(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "AgAdult1", "male", 30, "TZ")
	seedProfile(t, srv, "AgChild1", "female", 8, "TZ")
	seedProfile(t, srv, "AgTeen1", "male", 16, "TZ")

	resp := get(t, srv, "/api/profiles?age_group=adult")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		if p["age_group"] != "adult" {
			t.Errorf("expected age_group=adult, got %v", p["age_group"])
		}
	}
}

func TestFilterByMinAge(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "MinAge20", "male", 20, "GH")
	seedProfile(t, srv, "MinAge30", "male", 30, "GH")
	seedProfile(t, srv, "MinAge15", "female", 15, "GH")

	resp := get(t, srv, "/api/profiles?min_age=20")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		age := p["age"].(float64)
		if age < 20 {
			t.Errorf("expected age >= 20, got %v", age)
		}
	}
}

func TestFilterByMaxAge(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "MaxAge20", "female", 20, "SD")
	seedProfile(t, srv, "MaxAge30", "female", 30, "SD")
	seedProfile(t, srv, "MaxAge40", "male", 40, "SD")

	resp := get(t, srv, "/api/profiles?max_age=30")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		age := p["age"].(float64)
		if age > 30 {
			t.Errorf("expected age <= 30, got %v", age)
		}
	}
}

func TestFilterByMinGenderProbability(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	db.Exec(`INSERT INTO profiles (id,name,gender,gender_probability,age,age_group,country_id,country_name,country_probability,created_at)
		VALUES (?,?,?,?,?,?,?,?,?,datetime('now'))`, newUUID(), "HighGP", "male", 0.95, 30, "adult", "NG", "Nigeria", 0.8)
	db.Exec(`INSERT INTO profiles (id,name,gender,gender_probability,age,age_group,country_id,country_name,country_probability,created_at)
		VALUES (?,?,?,?,?,?,?,?,?,datetime('now'))`, newUUID(), "LowGP", "female", 0.4, 25, "adult", "KE", "Kenya", 0.7)

	resp := get(t, srv, "/api/profiles?min_gender_probability=0.9")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		gp := p["gender_probability"].(float64)
		if gp < 0.9 {
			t.Errorf("expected gender_probability >= 0.9, got %v", gp)
		}
	}
}

func TestFilterCombined(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	// Profiles that match all filters
	seedProfile(t, srv, "ComboA", "male", 30, "NG")
	seedProfile(t, srv, "ComboB", "male", 35, "NG")
	// Profiles that should NOT match
	seedProfile(t, srv, "ComboC", "female", 30, "NG") // wrong gender
	seedProfile(t, srv, "ComboD", "male", 15, "NG")   // age too low
	seedProfile(t, srv, "ComboE", "male", 30, "KE")   // wrong country

	resp := get(t, srv, "/api/profiles?gender=male&country_id=NG&min_age=25")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		if p["gender"] != "male" {
			t.Errorf("combined filter: expected gender=male, got %v", p["gender"])
		}
		if p["country_id"] != "NG" {
			t.Errorf("combined filter: expected country_id=NG, got %v", p["country_id"])
		}
		age := p["age"].(float64)
		if age < 25 {
			t.Errorf("combined filter: expected age >= 25, got %v", age)
		}
	}
}

// ── Sorting ───────────────────────────────────────────────────────────────────

func TestSortByAgeDesc(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	for _, age := range []int{20, 40, 30, 50, 25} {
		seedProfile(t, srv, fmt.Sprintf("Sort%d", age), "male", age, "ET")
	}

	resp := get(t, srv, "/api/profiles?sort_by=age&order=desc&limit=10")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	items := getDataList(t, body)
	for i := 1; i < len(items); i++ {
		prev := items[i-1]["age"].(float64)
		curr := items[i]["age"].(float64)
		if prev < curr {
			t.Errorf("not sorted desc at index %d: %v < %v", i, prev, curr)
		}
	}
}

func TestSortByAgeAsc(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	for _, age := range []int{50, 10, 35, 22} {
		seedProfile(t, srv, fmt.Sprintf("SortAsc%d", age), "female", age, "ZM")
	}

	resp := get(t, srv, "/api/profiles?sort_by=age&order=asc&limit=10")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	items := getDataList(t, body)
	for i := 1; i < len(items); i++ {
		prev := items[i-1]["age"].(float64)
		curr := items[i]["age"].(float64)
		if prev > curr {
			t.Errorf("not sorted asc at index %d: %v > %v", i, prev, curr)
		}
	}
}

func TestSortOrderCaseInsensitive(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "CaseOrder1", "male", 30, "NG")

	// Both ASC and asc should work without error
	for _, order := range []string{"ASC", "DESC", "Asc", "Desc"} {
		resp := get(t, srv, "/api/profiles?sort_by=age&order="+order)
		assertStatus(t, resp, http.StatusOK)
	}
}

// ── Query Validation ──────────────────────────────────────────────────────────

func TestQueryValidation_InvalidSortBy(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp := get(t, srv, "/api/profiles?sort_by=invalid_field")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusBadRequest)
	assertField(t, body, "status", "error")
	assertField(t, body, "message", "Invalid query parameters")
}

func TestQueryValidation_InvalidOrder(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp := get(t, srv, "/api/profiles?order=sideways")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusBadRequest)
	assertField(t, body, "status", "error")
}

func TestQueryValidation_InvalidGender(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp := get(t, srv, "/api/profiles?gender=attack_helicopter")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusBadRequest)
	assertField(t, body, "status", "error")
}

func TestQueryValidation_InvalidAgeGroup(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp := get(t, srv, "/api/profiles?age_group=boomer")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusBadRequest)
	assertField(t, body, "status", "error")
}

func TestQueryValidation_NonNumericAge(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	for _, param := range []string{"min_age=abc", "max_age=xyz", "page=one", "limit=ten"} {
		resp := get(t, srv, "/api/profiles?"+param)
		body := decode(t, resp)
		assertStatus(t, resp, http.StatusBadRequest)
		assertField(t, body, "status", "error")
	}
}

func TestQueryValidation_NonFloatProbability(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	for _, param := range []string{"min_gender_probability=high", "min_country_probability=low"} {
		resp := get(t, srv, "/api/profiles?"+param)
		body := decode(t, resp)
		assertStatus(t, resp, http.StatusBadRequest)
		assertField(t, body, "status", "error")
	}
}

func TestQueryValidation_ErrorEnvelopeKeyOrder(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	// Raw bytes check — status must come before message in JSON output
	raw, _ := http.Get(srv.URL + "/api/profiles?sort_by=bad")
	defer raw.Body.Close()
	var buf bytes.Buffer
	buf.ReadFrom(raw.Body)
	s := buf.String()

	statusIdx := -1
	messageIdx := -1
	for i := 0; i < len(s)-8; i++ {
		if s[i:i+8] == `"status"` && statusIdx == -1 {
			statusIdx = i
		}
		if s[i:i+9] == `"message"` && messageIdx == -1 {
			messageIdx = i
		}
	}
	if statusIdx == -1 || messageIdx == -1 {
		t.Fatal("response missing status or message field")
	}
	if statusIdx > messageIdx {
		t.Errorf("'status' key appears after 'message' in JSON — spec requires status first")
	}
}

// ── GET /api/profiles/search (NLQ) ────────────────────────────────────────────

func TestNLQ_EmptyQuery(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	// Empty q param
	resp := get(t, srv, "/api/profiles/search?q=")
	body := decode(t, resp)
	assertStatus(t, resp, http.StatusBadRequest)
	assertField(t, body, "status", "error")

	// Whitespace-only q (URL-encoded spaces)
	resp2 := get(t, srv, "/api/profiles/search?q=%20%20%20")
	body2 := decode(t, resp2)
	assertStatus(t, resp2, http.StatusBadRequest)
	assertField(t, body2, "status", "error")
}

func TestNLQ_UninterpretableQuery(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	// All values are URL-safe or properly encoded
	for _, q := range []string{"xyzxyzxyz", "12345", "lorem+ipsum+dolor", "randomgibberish999"} {
		resp := get(t, srv, "/api/profiles/search?q="+q)
		body := decode(t, resp)
		assertStatus(t, resp, http.StatusBadRequest)
		assertField(t, body, "status", "error")
		assertField(t, body, "message", "Unable to interpret query")
	}
}

func TestNLQ_YoungMalesFromNigeria(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	// Profiles that match: male, NG, age 16-24
	seedProfile(t, srv, "YoungNGMale1", "male", 18, "NG")
	seedProfile(t, srv, "YoungNGMale2", "male", 22, "NG")
	// Should NOT match
	seedProfile(t, srv, "OldNGMale", "male", 40, "NG")
	seedProfile(t, srv, "YoungKEMale", "male", 20, "KE")
	seedProfile(t, srv, "YoungNGFem", "female", 19, "NG")

	resp := get(t, srv, "/api/profiles/search?q=young+males+from+nigeria")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	assertField(t, body, "status", "success")

	for _, p := range getDataList(t, body) {
		if p["gender"] != "male" {
			t.Errorf("NLQ young males NG: expected male, got %v", p["gender"])
		}
		if p["country_id"] != "NG" {
			t.Errorf("NLQ young males NG: expected NG, got %v", p["country_id"])
		}
		age := p["age"].(float64)
		if age < 16 || age > 24 {
			t.Errorf("NLQ young males NG: expected 16<=age<=24, got %v", age)
		}
	}
}

func TestNLQ_FemalesAbove30(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "FemAbove30A", "female", 35, "ET")
	seedProfile(t, srv, "FemAbove30B", "female", 50, "ET")
	seedProfile(t, srv, "FemBelow30", "female", 25, "ET")
	seedProfile(t, srv, "MaleAbove30", "male", 40, "ET")

	resp := get(t, srv, "/api/profiles/search?q=females+above+30")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		if p["gender"] != "female" {
			t.Errorf("NLQ females above 30: expected female, got %v", p["gender"])
		}
		age := p["age"].(float64)
		if age < 30 {
			t.Errorf("NLQ females above 30: expected age >= 30, got %v", age)
		}
	}
}

func TestNLQ_PeopleFromAngola(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "AOUser1", "male", 30, "AO")
	seedProfile(t, srv, "AOUser2", "female", 25, "AO")
	seedProfile(t, srv, "NGUser", "male", 30, "NG")

	resp := get(t, srv, "/api/profiles/search?q=people+from+angola")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	items := getDataList(t, body)
	if len(items) == 0 {
		t.Fatal("expected results for people from angola")
	}
	for _, p := range items {
		if p["country_id"] != "AO" {
			t.Errorf("NLQ from angola: expected AO, got %v", p["country_id"])
		}
	}
}

func TestNLQ_AdultMalesFromKenya(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "KEAdultMale1", "male", 30, "KE")
	seedProfile(t, srv, "KEAdultMale2", "male", 45, "KE")
	seedProfile(t, srv, "KEAdultFem", "female", 30, "KE") // wrong gender
	seedProfile(t, srv, "NGAdultMale", "male", 30, "NG")  // wrong country
	seedProfile(t, srv, "KETeenMale", "male", 16, "KE")   // wrong age group (teenager)

	resp := get(t, srv, "/api/profiles/search?q=adult+males+from+kenya")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		if p["gender"] != "male" {
			t.Errorf("NLQ adult males KE: expected male, got %v", p["gender"])
		}
		if p["country_id"] != "KE" {
			t.Errorf("NLQ adult males KE: expected KE, got %v", p["country_id"])
		}
		if p["age_group"] != "adult" {
			t.Errorf("NLQ adult males KE: expected adult, got %v", p["age_group"])
		}
	}
}

func TestNLQ_MaleAndFemaleTeenagersAbove17(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "Teen18M", "male", 18, "UG")
	seedProfile(t, srv, "Teen19F", "female", 19, "UG")
	seedProfile(t, srv, "Teen15", "male", 15, "UG")  // age < 17
	seedProfile(t, srv, "Adult25", "male", 25, "UG") // not a teenager

	resp := get(t, srv, "/api/profiles/search?q=male+and+female+teenagers+above+17")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	for _, p := range getDataList(t, body) {
		if p["age_group"] != "teenager" {
			t.Errorf("NLQ teen above 17: expected teenager, got %v", p["age_group"])
		}
		age := p["age"].(float64)
		if age < 17 {
			t.Errorf("NLQ teen above 17: expected age >= 17, got %v", age)
		}
	}
}

func TestNLQ_PaginationApplies(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	for i := 0; i < 8; i++ {
		seedProfile(t, srv, fmt.Sprintf("NLQPag%d", i), "male", 20+i, "NG")
	}

	resp := get(t, srv, "/api/profiles/search?q=males+from+nigeria&page=1&limit=3")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	assertField(t, body, "status", "success")
	items := getDataList(t, body)
	if len(items) != 3 {
		t.Errorf("NLQ pagination: expected 3 items, got %d", len(items))
	}
	if body["total"] == nil {
		t.Error("NLQ response missing 'total' field")
	}
}

func TestNLQ_ResponseEnvelopeKeyOrder(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	seedProfile(t, srv, "EnvTest", "male", 30, "NG")

	raw, _ := http.Get(srv.URL + "/api/profiles/search?q=males+from+nigeria")
	defer raw.Body.Close()
	var buf bytes.Buffer
	buf.ReadFrom(raw.Body)
	s := buf.String()

	statusIdx := -1
	dataIdx := -1
	for i := 0; i < len(s)-6; i++ {
		if len(s) > i+8 && s[i:i+8] == `"status"` && statusIdx == -1 {
			statusIdx = i
		}
		if len(s) > i+6 && s[i:i+6] == `"data"` && dataIdx == -1 {
			dataIdx = i
		}
	}
	if statusIdx == -1 || dataIdx == -1 {
		t.Fatal("NLQ response missing status or data field")
	}
	if statusIdx > dataIdx {
		t.Errorf("'status' appears after 'data' — spec requires status first")
	}
}

// ── GET /api/profiles/:id ─────────────────────────────────────────────────────

func TestGetProfile_Found(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	created := decode(t, post(t, srv, "/api/profiles", map[string]any{
		"name": "GetByID Test", "gender": "male", "age": 28, "country_id": "GH",
	}))
	id := created["data"].(map[string]any)["id"].(string)

	resp := get(t, srv, "/api/profiles/"+id)
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	assertField(t, body, "status", "success")
	assertDataField(t, body, "id", id)
}

func TestGetProfile_NotFound(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp := get(t, srv, "/api/profiles/00000000-0000-0000-0000-000000000000")
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusNotFound)
	assertField(t, body, "status", "error")
}

// ── DELETE /api/profiles/:id ──────────────────────────────────────────────────

func TestDeleteProfile_Success(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	created := decode(t, post(t, srv, "/api/profiles", map[string]any{
		"name": "DeleteMe", "gender": "female", "age": 22, "country_id": "TZ",
	}))
	id := created["data"].(map[string]any)["id"].(string)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/api/profiles/"+id, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusOK)
	assertField(t, body, "status", "success")

	// Confirm it's gone
	assertStatus(t, get(t, srv, "/api/profiles/"+id), http.StatusNotFound)
}

func TestDeleteProfile_NotFound(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/api/profiles/00000000-0000-0000-0000-000000000000", nil)
	resp, _ := http.DefaultClient.Do(req)
	body := decode(t, resp)

	assertStatus(t, resp, http.StatusNotFound)
	assertField(t, body, "status", "error")
}

// ── CORS ──────────────────────────────────────────────────────────────────────

func TestCORSHeader(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/profiles?limit=1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	got := resp.Header.Get("Access-Control-Allow-Origin")
	if got != "*" {
		t.Errorf("CORS header: got %q, want %q", got, "*")
	}
}

func TestCORSPreflightOPTIONS(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodOptions, srv.URL+"/api/profiles", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("OPTIONS status: got %d, want 204", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("OPTIONS response missing CORS header")
	}
}
