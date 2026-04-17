package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ── External API response types ──────────────────────────────────────────────

type GenderizeResponse struct {
	Gender      *string `json:"gender"`
	Probability float64 `json:"probability"`
	Count       *int    `json:"count"`
}

type AgifyResponse struct {
	Age *int `json:"age"`
}

type NationalizeCountry struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type NationalizeResponse struct {
	Country []NationalizeCountry `json:"country"`
}

// ── Domain model ─────────────────────────────────────────────────────────────

type Profile struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Gender            string    `json:"gender"`
	GenderProbability float64   `json:"gender_probability"`
	SampleSize        int       `json:"sample_size"`
	Age               int       `json:"age"`
	AgeGroup          string    `json:"age_group"`
	CountryID         string    `json:"country_id"`
	CreatedAt         time.Time `json:"created_at"`
}

// ── In-memory store (swap out for a real DB) ─────────────────────────────────

var (
	store   = map[string]*Profile{} // keyed by lowercase name
	storeMu sync.RWMutex
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func classifyAge(age int) string {
	switch {
	case age <= 12:
		return "child"
	case age <= 19:
		return "teenager"
	case age <= 59:
		return "adult"
	default:
		return "senior"
	}
}

func fetchJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func createProfile(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "name is required"})
		return
	}

	// Idempotency check
	storeMu.RLock()
	existing, found := store[body.Name]
	storeMu.RUnlock()
	if found {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": existing})
		return
	}

	// Call all three APIs concurrently
	var (
		gr   GenderizeResponse
		ar   AgifyResponse
		nr   NationalizeResponse
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	fetch := func(url string, target any) {
		defer wg.Done()
		if err := fetchJSON(url, target); err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}

	wg.Add(3)
	go fetch(fmt.Sprintf("https://api.genderize.io?name=%s", body.Name), &gr)
	go fetch(fmt.Sprintf("https://api.agify.io?name=%s", body.Name), &ar)
	go fetch(fmt.Sprintf("https://api.nationalize.io?name=%s", body.Name), &nr)
	wg.Wait()

	if len(errs) > 0 {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "failed to reach upstream APIs"})
		return
	}

	// Pick top country
	topCountry := ""
	var topProb float64
	for _, ct := range nr.Country {
		if ct.Probability > topProb {
			topProb = ct.Probability
			topCountry = ct.CountryID
		}
	}

	gender := ""
	if gr.Gender != nil {
		gender = *gr.Gender
	}
	sampleSize := 0
	if gr.Count != nil {
		sampleSize = *gr.Count
	}
	age := 0
	if ar.Age != nil {
		age = *ar.Age
	}

	p := &Profile{
		ID:                uuid.New().String(),
		Name:              body.Name,
		Gender:            gender,
		GenderProbability: gr.Probability,
		SampleSize:        sampleSize,
		Age:               age,
		AgeGroup:          classifyAge(age),
		CountryID:         topCountry,
		CreatedAt:         time.Now().UTC(),
	}

	storeMu.Lock()
	store[body.Name] = p
	storeMu.Unlock()

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": p})
}

func listProfiles(c *gin.Context) {
	storeMu.RLock()
	profiles := make([]*Profile, 0, len(store))
	for _, p := range store {
		profiles = append(profiles, p)
	}
	storeMu.RUnlock()
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": profiles})
}

func getProfile(c *gin.Context) {
	id := c.Param("id")
	storeMu.RLock()
	defer storeMu.RUnlock()
	for _, p := range store {
		if p.ID == id {
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": p})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Profile not found"})
}

func deleteProfile(c *gin.Context) {
	id := c.Param("id")
	storeMu.Lock()
	defer storeMu.Unlock()
	for name, p := range store {
		if p.ID == id {
			delete(store, name)
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Profile deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Profile not found"})
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8060"
	}

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/profiles", createProfile)
		api.GET("/profiles", listProfiles)
		api.GET("/profiles/:id", getProfile)
		api.DELETE("/profiles/:id", deleteProfile)
	}

	r.Run(":" + port)
}
