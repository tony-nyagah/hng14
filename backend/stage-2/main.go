package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// --- Models ---

type Profile struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Gender             string    `json:"gender"`
	GenderProbability  float64   `json:"gender_probability"`
	Age                int       `json:"age"`
	AgeGroup           string    `json:"age_group"`
	CountryID          string    `json:"country_id"`
	CountryName        string    `json:"country_name"`
	CountryProbability float64   `json:"country_probability"`
	CreatedAt          time.Time `json:"created_at"`
}

// SeedWrapper matches the JSON structure: {"profiles": [...]}
type SeedWrapper struct {
	Profiles []Profile `json:"profiles"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Page    int         `json:"page,omitempty"`
	Limit   int         `json:"limit,omitempty"`
	Total   int         `json:"total,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var db *sql.DB

// --- Database & Seeding ---

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./insighta.db")
	if err != nil {
		log.Fatal(err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS profiles (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE,
		gender TEXT,
		gender_probability REAL,
		age INTEGER,
		age_group TEXT,
		country_id TEXT,
		country_name TEXT,
		country_probability REAL,
		created_at TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_profiles_filters ON profiles(gender, age, country_id, age_group);
	`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
}

func seedDatabase() {
	jsonFile, err := os.Open("seed.json")
	if err != nil {
		log.Println("⚠️ seed.json not found. Place the file in the root directory.")
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var wrapper SeedWrapper
	if err := json.Unmarshal(byteValue, &wrapper); err != nil {
		log.Printf("❌ Error parsing JSON: %v", err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Error starting transaction: %v", err)
		return
	}

	stmt, _ := tx.Prepare(`INSERT OR IGNORE INTO profiles (id, name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	count := 0
	for _, p := range wrapper.Profiles {
		// FIX: Handle the dual return values from NewV7()
		if p.ID == "" {
			u, err := uuid.NewV7()
			if err != nil {
				p.ID = uuid.NewString() // Fallback
			} else {
				p.ID = u.String()
			}
		}

		if p.CreatedAt.IsZero() {
			p.CreatedAt = time.Now().UTC()
		}

		_, err := stmt.Exec(p.ID, p.Name, p.Gender, p.GenderProbability, p.Age, p.AgeGroup, p.CountryID, p.CountryName, p.CountryProbability, p.CreatedAt)
		if err == nil {
			count++
		}
	}
	tx.Commit()
	fmt.Printf("✅ Seeding complete. Processed %d records.\n", count)
}

// --- Rule-Based NLP Parser ---

func parseNLP(q string) (map[string]string, bool) {
	q = strings.ToLower(q)
	f := make(map[string]string)
	found := false

	if strings.Contains(q, "female") {
		f["gender"] = "female"
		found = true
	} else if strings.Contains(q, "male") {
		f["gender"] = "male"
		found = true
	}

	if strings.Contains(q, "young") {
		f["min_age"], f["max_age"] = "16", "24"
		found = true
	}
	if strings.Contains(q, "teenager") {
		f["age_group"] = "teenager"
		found = true
	}
	if strings.Contains(q, "adult") {
		f["age_group"] = "adult"
		found = true
	}

	words := strings.Fields(q)
	for i, w := range words {
		if w == "above" && i+1 < len(words) {
			if val, err := strconv.Atoi(words[i+1]); err == nil {
				f["min_age"] = strconv.Itoa(val + 1)
				found = true
			}
		}
	}

	countries := map[string]string{"nigeria": "NG", "kenya": "KE", "angola": "AO", "benin": "BJ"}
	for name, code := range countries {
		if strings.Contains(q, name) {
			f["country_id"] = code
			found = true
		}
	}

	return f, found
}

// --- Handlers ---

func profilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	params := r.URL.Query()
	filters := make(map[string]string)

	if strings.Contains(r.URL.Path, "/search") {
		q := params.Get("q")
		if q == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Missing or empty parameter"})
			return
		}
		var ok bool
		filters, ok = parseNLP(q)
		if !ok {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Unable to interpret query"})
			return
		}
	} else {
		keys := []string{"gender", "age_group", "country_id", "min_age", "max_age", "min_gender_probability", "min_country_probability"}
		for _, k := range keys {
			if v := params.Get(k); v != "" {
				filters[k] = v
			}
		}
	}

	page, _ := strconv.Atoi(params.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(params.Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	sortBy := params.Get("sort_by")
	if sortBy != "age" && sortBy != "gender_probability" && sortBy != "created_at" {
		sortBy = "created_at"
	}
	order := strings.ToLower(params.Get("order"))
	if order != "asc" {
		order = "desc"
	}

	sqlBase := "FROM profiles WHERE 1=1"
	var args []interface{}

	filterSpecs := map[string]string{
		"gender": "AND gender = ?", "age_group": "AND age_group = ?", "country_id": "AND country_id = ?",
		"min_age": "AND age >= ?", "max_age": "AND age <= ?",
		"min_gender_probability": "AND gender_probability >= ?", "min_country_probability": "AND country_probability >= ?",
	}

	for k, clause := range filterSpecs {
		if v, ok := filters[k]; ok {
			sqlBase += " " + clause
			args = append(args, v)
		}
	}

	var total int
	db.QueryRow("SELECT COUNT(*) "+sqlBase, args...).Scan(&total)

	finalQuery := fmt.Sprintf("SELECT * %s ORDER BY %s %s LIMIT %d OFFSET %d", sqlBase, sortBy, order, limit, (page-1)*limit)
	rows, err := db.Query(finalQuery, args...)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Server failure"})
		return
	}
	defer rows.Close()

	results := []Profile{}
	for rows.Next() {
		var p Profile
		rows.Scan(&p.ID, &p.Name, &p.Gender, &p.GenderProbability, &p.Age, &p.AgeGroup, &p.CountryID, &p.CountryName, &p.CountryProbability, &p.CreatedAt)
		results = append(results, p)
	}

	json.NewEncoder(w).Encode(APIResponse{
		Status: "success",
		Page:   page,
		Limit:  limit,
		Total:  total,
		Data:   results,
	})
}

func main() {
	initDB()
	seedDatabase()

	http.HandleFunc("/api/profiles", profilesHandler)
	http.HandleFunc("/api/profiles/search", profilesHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8070"
	}
	fmt.Printf("🚀 Intelligence Engine running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
