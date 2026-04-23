package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ── Model ─────────────────────────────────────────────────────────────────────

type Profile struct {
	ID                 string    `gorm:"primarykey;type:varchar(36)" json:"id"`
	Name               string    `gorm:"uniqueIndex;not null;type:varchar(255)" json:"name"`
	Gender             string    `gorm:"type:varchar(10);index" json:"gender"`
	GenderProbability  float64   `gorm:"index" json:"gender_probability"`
	Age                int       `gorm:"index" json:"age"`
	AgeGroup           string    `gorm:"type:varchar(20);index" json:"age_group"`
	CountryID          string    `gorm:"type:varchar(2);index" json:"country_id"`
	CountryName        string    `gorm:"type:varchar(100)" json:"country_name"`
	CountryProbability float64   `gorm:"index" json:"country_probability"`
	CreatedAt          time.Time `gorm:"index" json:"created_at"`
}

// ── Database ──────────────────────────────────────────────────────────────────

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		pass := getEnv("DB_PASSWORD", "postgres")
		name := getEnv("DB_NAME", "hng14_stage2")
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			host, port, user, pass, name,
		)
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Profile{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	log.Println("Database connected and migrated")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ── Seeding ───────────────────────────────────────────────────────────────────

type seedEntry struct {
	Name               string  `json:"name"`
	Gender             string  `json:"gender"`
	GenderProbability  float64 `json:"gender_probability"`
	Age                int     `json:"age"`
	AgeGroup           string  `json:"age_group"`
	CountryID          string  `json:"country_id"`
	CountryName        string  `json:"country_name"`
	CountryProbability float64 `json:"country_probability"`
}

type seedFile struct {
	Profiles []seedEntry `json:"profiles"`
}

func seedDB() {
	seedPath := getEnv("SEED_FILE", "seed_profiles.json")
	f, err := os.Open(seedPath)
	if err != nil {
		log.Printf("seed file not found: %v (skipping)", err)
		return
	}
	defer f.Close()

	var data seedFile
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		log.Fatalf("failed to parse seed file: %v", err)
	}

	inserted := 0
	for _, e := range data.Profiles {
		res := db.Exec(
			`INSERT INTO profiles (id,name,gender,gender_probability,age,age_group,country_id,country_name,country_probability,created_at)
			 VALUES (?,?,?,?,?,?,?,?,?,?) ON CONFLICT (name) DO NOTHING`,
			newUUID(), e.Name, e.Gender, e.GenderProbability, e.Age, e.AgeGroup,
			e.CountryID, e.CountryName, e.CountryProbability, time.Now().UTC(),
		)
		if res.Error == nil && res.RowsAffected > 0 {
			inserted++
		}
	}

	var total int64
	db.Model(&Profile{}).Count(&total)
	log.Printf("Seed complete: inserted %d new, total %d profiles in DB", inserted, total)
}

func newUUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}

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

// ── External API types ────────────────────────────────────────────────────────

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

// ── Country maps ──────────────────────────────────────────────────────────────

// countryNameToCode maps lowercase country names (and common variants) to ISO-2 codes.
var countryNameToCode = map[string]string{
	// Africa (seed data + variants)
	"nigeria": "NG", "kenya": "KE", "tanzania": "TZ", "uganda": "UG",
	"ghana": "GH", "ethiopia": "ET", "angola": "AO", "sudan": "SD",
	"south africa": "ZA", "egypt": "EG", "morocco": "MA", "senegal": "SN",
	"cameroon": "CM", "zimbabwe": "ZW", "zambia": "ZM", "mozambique": "MZ",
	"madagascar": "MG", "mali": "ML", "niger": "NE", "burkina faso": "BF",
	"somalia": "SO", "south sudan": "SS", "rwanda": "RW", "burundi": "BI",
	"benin": "BJ", "togo": "TG", "sierra leone": "SL", "liberia": "LR",
	"malawi": "MW", "namibia": "NA", "botswana": "BW", "gabon": "GA",
	"republic of the congo": "CG", "congo": "CG",
	"dr congo": "CD", "democratic republic of congo": "CD", "democratic republic of the congo": "CD",
	"guinea": "GN", "guinea-bissau": "GW", "equatorial guinea": "GQ",
	"central african republic": "CF", "chad": "TD", "mauritania": "MR",
	"gambia": "GM", "djibouti": "DJ", "eritrea": "ER", "comoros": "KM",
	"cape verde": "CV", "sao tome and principe": "ST", "seychelles": "SC",
	"mauritius": "MU", "lesotho": "LS", "swaziland": "SZ", "eswatini": "SZ",
	"libya": "LY", "algeria": "DZ", "tunisia": "TN", "western sahara": "EH",
	"ivory coast": "CI", "cote d ivoire": "CI",
	// Rest of world
	"united states": "US", "usa": "US", "united kingdom": "GB", "uk": "GB",
	"canada": "CA", "australia": "AU", "germany": "DE", "france": "FR",
	"spain": "ES", "italy": "IT", "brazil": "BR", "india": "IN",
	"china": "CN", "japan": "JP", "russia": "RU", "mexico": "MX",
	"indonesia": "ID", "pakistan": "PK", "bangladesh": "BD", "philippines": "PH",
	"vietnam": "VN", "thailand": "TH", "iran": "IR", "turkey": "TR",
	"saudi arabia": "SA", "iraq": "IQ", "afghanistan": "AF", "nepal": "NP",
	"sri lanka": "LK", "malaysia": "MY", "singapore": "SG", "argentina": "AR",
	"colombia": "CO", "chile": "CL", "peru": "PE", "venezuela": "VE",
	"ukraine": "UA", "poland": "PL", "netherlands": "NL", "sweden": "SE",
	"norway": "NO", "denmark": "DK", "finland": "FI", "switzerland": "CH",
	"austria": "AT", "belgium": "BE", "portugal": "PT", "greece": "GR",
	"new zealand": "NZ", "ireland": "IE", "israel": "IL", "jordan": "JO",
	"uae": "AE", "united arab emirates": "AE", "qatar": "QA",
}

// countryCodeToName is the reverse map, built at init time.
var countryCodeToName = map[string]string{}

func init() {
	for name, code := range countryNameToCode {
		if existing, ok := countryCodeToName[code]; !ok || len(name) > len(existing) {
			countryCodeToName[code] = toTitle(name)
		}
	}
}

func toTitle(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// ── Filters ───────────────────────────────────────────────────────────────────

type ProfileFilters struct {
	Gender                *string
	AgeGroup              *string
	CountryID             *string
	MinAge                *int
	MaxAge                *int
	MinGenderProbability  *float64
	MinCountryProbability *float64
	SortBy                string
	Order                 string
	Page                  int
	Limit                 int
}

func applyFilters(q *gorm.DB, f ProfileFilters) *gorm.DB {
	if f.Gender != nil {
		q = q.Where("gender = ?", *f.Gender)
	}
	if f.AgeGroup != nil {
		q = q.Where("age_group = ?", *f.AgeGroup)
	}
	if f.CountryID != nil {
		q = q.Where("country_id = ?", *f.CountryID)
	}
	if f.MinAge != nil {
		q = q.Where("age >= ?", *f.MinAge)
	}
	if f.MaxAge != nil {
		q = q.Where("age <= ?", *f.MaxAge)
	}
	if f.MinGenderProbability != nil {
		q = q.Where("gender_probability >= ?", *f.MinGenderProbability)
	}
	if f.MinCountryProbability != nil {
		q = q.Where("country_probability >= ?", *f.MinCountryProbability)
	}
	return q
}

func applySortPagination(q *gorm.DB, f ProfileFilters) *gorm.DB {
	validSort := map[string]bool{"age": true, "created_at": true, "gender_probability": true}
	sortBy := "created_at"
	if validSort[f.SortBy] {
		sortBy = f.SortBy
	}
	order := "asc"
	if strings.ToLower(f.Order) == "desc" {
		order = "desc"
	}
	q = q.Order(fmt.Sprintf("%s %s", sortBy, order))

	page := f.Page
	if page < 1 {
		page = 1
	}
	limit := f.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	q = q.Offset((page - 1) * limit).Limit(limit)
	return q
}

func parseQueryFilters(c *gin.Context) (ProfileFilters, bool) {
	f := ProfileFilters{
		SortBy: c.Query("sort_by"),
		Order:  c.Query("order"),
		Page:   1,
		Limit:  10,
	}

	if v := c.Query("gender"); v != "" {
		if v != "male" && v != "female" {
			return f, false
		}
		f.Gender = &v
	}
	if v := c.Query("age_group"); v != "" {
		valid := map[string]bool{"child": true, "teenager": true, "adult": true, "senior": true}
		if !valid[v] {
			return f, false
		}
		f.AgeGroup = &v
	}
	if v := c.Query("country_id"); v != "" {
		upper := strings.ToUpper(v)
		f.CountryID = &upper
	}
	if v := c.Query("min_age"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return f, false
		}
		f.MinAge = &n
	}
	if v := c.Query("max_age"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return f, false
		}
		f.MaxAge = &n
	}
	if v := c.Query("min_gender_probability"); v != "" {
		fv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return f, false
		}
		f.MinGenderProbability = &fv
	}
	if v := c.Query("min_country_probability"); v != "" {
		fv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return f, false
		}
		f.MinCountryProbability = &fv
	}
	if v := c.Query("sort_by"); v != "" {
		valid := map[string]bool{"age": true, "created_at": true, "gender_probability": true}
		if !valid[v] {
			return f, false
		}
	}
	if v := c.Query("order"); v != "" {
		if v != "asc" && v != "desc" {
			return f, false
		}
	}
	if v := c.Query("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return f, false
		}
		f.Page = n
	}
	if v := c.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return f, false
		}
		f.Limit = n
	}
	return f, true
}

// ── NLQ Parser ────────────────────────────────────────────────────────────────

var (
	reMale    = regexp.MustCompile(`\b(male|males|men|man|boy|boys)\b`)
	reFemale  = regexp.MustCompile(`\b(female|females|women|woman|girl|girls)\b`)
	reChild   = regexp.MustCompile(`\b(child|children|kid|kids)\b`)
	reTeen    = regexp.MustCompile(`\b(teenager|teenagers|teen|teens|adolescent|adolescents)\b`)
	reAdult   = regexp.MustCompile(`\b(adult|adults)\b`)
	reSenior  = regexp.MustCompile(`\b(senior|seniors|elderly)\b`)
	reYoung   = regexp.MustCompile(`\byoung\b`)
	reAbove   = regexp.MustCompile(`(?:above|over|older than|at least)\s+(\d+)`)
	reBelow   = regexp.MustCompile(`(?:below|under|younger than|at most)\s+(\d+)`)
	reBetween = regexp.MustCompile(`between\s+(\d+)\s+and\s+(\d+)`)
	reFrom    = regexp.MustCompile(`\bfrom\s+([a-z][a-z\s\-]*)`)
	reStop    = regexp.MustCompile(`\s+(above|below|over|under|between|who|with|aged?|and|\d).*$`)
)

func parseNLQ(q string) (ProfileFilters, bool) {
	f := ProfileFilters{Page: 1, Limit: 10}
	lower := strings.ToLower(strings.TrimSpace(q))
	if lower == "" {
		return f, false
	}

	interpreted := false

	// Gender
	hasMale := reMale.MatchString(lower)
	hasFemale := reFemale.MatchString(lower)
	if hasMale && !hasFemale {
		g := "male"
		f.Gender = &g
		interpreted = true
	} else if hasFemale && !hasMale {
		g := "female"
		f.Gender = &g
		interpreted = true
	} else if hasMale && hasFemale {
		interpreted = true // both genders mentioned — no gender filter, but valid query
	}

	// Age group (first match wins)
	switch {
	case reChild.MatchString(lower):
		ag := "child"
		f.AgeGroup = &ag
		interpreted = true
	case reTeen.MatchString(lower):
		ag := "teenager"
		f.AgeGroup = &ag
		interpreted = true
	case reAdult.MatchString(lower):
		ag := "adult"
		f.AgeGroup = &ag
		interpreted = true
	case reSenior.MatchString(lower):
		ag := "senior"
		f.AgeGroup = &ag
		interpreted = true
	}

	// "young" → 16–24 (only if no age_group detected)
	if f.AgeGroup == nil && reYoung.MatchString(lower) {
		min, max := 16, 24
		f.MinAge = &min
		f.MaxAge = &max
		interpreted = true
	}

	// Age ranges: between X and Y takes precedence over above/below
	if m := reBetween.FindStringSubmatch(lower); m != nil {
		n1, _ := strconv.Atoi(m[1])
		n2, _ := strconv.Atoi(m[2])
		f.MinAge = &n1
		f.MaxAge = &n2
		interpreted = true
	} else {
		if m := reAbove.FindStringSubmatch(lower); m != nil {
			n, _ := strconv.Atoi(m[1])
			f.MinAge = &n
			interpreted = true
		}
		if m := reBelow.FindStringSubmatch(lower); m != nil {
			n, _ := strconv.Atoi(m[1])
			f.MaxAge = &n
			interpreted = true
		}
	}

	// Country: extract text after "from"
	if m := reFrom.FindStringSubmatch(lower); m != nil {
		candidate := strings.TrimSpace(reStop.ReplaceAllString(m[1], ""))
		if candidate != "" {
			if code, ok := countryNameToCode[candidate]; ok {
				f.CountryID = &code
				interpreted = true
			} else {
				// Try progressively shorter word sequences (handles multi-word countries)
				words := strings.Fields(candidate)
				for i := len(words); i > 0; i-- {
					partial := strings.Join(words[:i], " ")
					if code, ok := countryNameToCode[partial]; ok {
						f.CountryID = &code
						interpreted = true
						break
					}
				}
			}
		}
	}

	return f, interpreted
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

	name := strings.TrimSpace(body.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "name is required"})
		return
	}

	// Idempotency check
	var existing Profile
	if err := db.Where("name = ?", name).First(&existing).Error; err == nil {
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
	go fetch(fmt.Sprintf("https://api.genderize.io?name=%s", name), &gr)
	go fetch(fmt.Sprintf("https://api.agify.io?name=%s", name), &ar)
	go fetch(fmt.Sprintf("https://api.nationalize.io?name=%s", name), &nr)
	wg.Wait()

	if len(errs) > 0 {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "failed to reach upstream APIs"})
		return
	}

	// Pick top country
	top := NationalizeCountry{}
	for _, ct := range nr.Country {
		if ct.Probability > top.Probability {
			top = ct
		}
	}

	gender := ""
	if gr.Gender != nil {
		gender = *gr.Gender
	}
	age := 0
	if ar.Age != nil {
		age = *ar.Age
	}

	countryName := countryCodeToName[top.CountryID]
	if countryName == "" {
		countryName = top.CountryID
	}

	p := Profile{
		ID:                 newUUID(),
		Name:               name,
		Gender:             gender,
		GenderProbability:  gr.Probability,
		Age:                age,
		AgeGroup:           classifyAge(age),
		CountryID:          top.CountryID,
		CountryName:        countryName,
		CountryProbability: top.Probability,
		CreatedAt:          time.Now().UTC(),
	}

	if err := db.Create(&p).Error; err != nil {
		// Race condition: another goroutine may have inserted the same name
		var existing2 Profile
		if db.Where("name = ?", name).First(&existing2).Error == nil {
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": existing2})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": p})
}

func listProfiles(c *gin.Context) {
	f, valid := parseQueryFilters(c)
	if !valid {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
		return
	}

	var total int64
	if err := applyFilters(db.Model(&Profile{}), f).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "database error"})
		return
	}

	var profiles []Profile
	if err := applySortPagination(applyFilters(db.Model(&Profile{}), f), f).Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "database error"})
		return
	}

	limit := f.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   f.Page,
		"limit":  limit,
		"total":  total,
		"data":   profiles,
	})
}

func searchProfiles(c *gin.Context) {
	q := c.Query("q")
	if strings.TrimSpace(q) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid query parameters"})
		return
	}

	f, ok := parseNLQ(q)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Unable to interpret query"})
		return
	}

	// Allow pagination overrides from explicit query params
	if v := c.Query("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
			return
		}
		f.Page = n
	}
	if v := c.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
			return
		}
		f.Limit = n
	}

	var total int64
	if err := applyFilters(db.Model(&Profile{}), f).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "database error"})
		return
	}

	var profiles []Profile
	if err := applySortPagination(applyFilters(db.Model(&Profile{}), f), f).Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "database error"})
		return
	}

	limit := f.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   f.Page,
		"limit":  limit,
		"total":  total,
		"data":   profiles,
	})
}

func getProfile(c *gin.Context) {
	id := c.Param("id")
	var p Profile
	if err := db.Where("id = ?", id).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": p})
}

func deleteProfile(c *gin.Context) {
	id := c.Param("id")
	result := db.Where("id = ?", id).Delete(&Profile{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to delete"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Profile deleted successfully"})
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	initDB()
	seedDB()

	port := getEnv("PORT", "8070")

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/profiles", createProfile)
		api.GET("/profiles/search", searchProfiles) // must be before /:id
		api.GET("/profiles", listProfiles)
		api.GET("/profiles/:id", getProfile)
		api.DELETE("/profiles/:id", deleteProfile)
	}

	log.Printf("Starting server on :%s", port)
	r.Run(":" + port)
}
