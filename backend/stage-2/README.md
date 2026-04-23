# Stage 2 Backend — Intelligence Query Engine

A queryable demographic intelligence API built with **Go** (standard library `net/http`) and **SQLite**. Auto-seeded with 2026 profiles on startup, it supports advanced filtering, sorting, pagination, and a rule-based natural language query interface.

---

## Architecture

| Component | Technology |
|-----------|-----------|
| Language  | Go 1.26   |
| HTTP      | `net/http` (standard library) |
| Database  | SQLite (file: `insighta.db`) |
| Container | Docker (multi-stage) |

---

## Quick Start

### With Docker

```bash
docker build -t insighta-api .
docker run -p 8070:8070 insighta-api
```

The API will be available at `http://localhost:8070`.

### Locally

```bash
go run .
```

Requires `seed.json` in the working directory — auto-seeded on first run.

---

## Environment Variables

| Variable | Default | Description         |
|----------|---------|---------------------|
| `PORT`   | `8070`  | API listen port     |

---

## Endpoints

### `GET /api/profiles`

List profiles with optional filtering, sorting, and pagination. All filter params are combinable — results match **all** provided conditions.

**Query parameters:**

| Param                    | Example           | Description                          |
|--------------------------|-------------------|--------------------------------------|
| `gender`                 | `male` / `female` | Filter by gender                     |
| `age_group`              | `adult`           | `child`, `teenager`, `adult`, `senior` |
| `country_id`             | `NG`              | ISO 2-letter country code            |
| `min_age`                | `25`              | Minimum age (inclusive)              |
| `max_age`                | `40`              | Maximum age (inclusive)              |
| `min_gender_probability` | `0.8`             | Minimum gender confidence score      |
| `min_country_probability`| `0.5`             | Minimum country confidence score     |
| `sort_by`                | `age`             | `age` / `created_at` / `gender_probability` |
| `order`                  | `desc`            | `asc` / `desc` (default: `desc`)     |
| `page`                   | `1`               | Page number (default: `1`)           |
| `limit`                  | `10`              | Results per page (default: `10`, max: `50`) |

**Example:** `GET /api/profiles?gender=male&country_id=NG&min_age=25&sort_by=age&order=desc`

**Response:**
```json
{
  "status": "success",
  "page": 1,
  "limit": 10,
  "total": 312,
  "data": [ ... ]
}
```

---

### `GET /api/profiles/search?q=<query>`

Natural language query converted to structured filters. Rule-based — no AI or LLMs. Supports `page` and `limit` for pagination.

**Examples:**

| Query | Interpreted as |
|-------|---------------|
| `young males from nigeria` | `gender=male`, `min_age=16`, `max_age=24`, `country_id=NG` |
| `females above 30` | `gender=female`, `min_age=31` |
| `people from angola` | `country_id=AO` |
| `adult males from kenya` | `gender=male`, `age_group=adult`, `country_id=KE` |
| `teenagers above 17` | `age_group=teenager`, `min_age=18` |

Uninterpretable queries return `400`:
```json
{ "status": "error", "message": "Unable to interpret query" }
```

---

## NLQ Parsing Rules

The parser uses keyword matching on the lowercased query string:

- **Gender:** `female` → `gender=female`; `male` → `gender=male` (female checked first)
- **Age groups:** `young` → `min_age=16, max_age=24`; `teenager` → `age_group=teenager`; `adult` → `age_group=adult`
- **Age ranges:** `above N` → `min_age=N+1`
- **Country:** keyword lookup against a built-in map (Nigeria, Kenya, Angola, Benin supported)

If no recognizable tokens are found, returns `Unable to interpret query`.

---

## Database Schema

```sql
CREATE TABLE IF NOT EXISTS profiles (
  id                  TEXT PRIMARY KEY,
  name                TEXT UNIQUE,
  gender              TEXT,
  gender_probability  REAL,
  age                 INTEGER,
  age_group           TEXT,
  country_id          TEXT,
  country_name        TEXT,
  country_probability REAL,
  created_at          TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_profiles_filters
  ON profiles(gender, age, country_id, age_group);
```

---

## Data Seeding

The app auto-seeds on startup from `seed.json` (2026 profiles). Uses `INSERT OR IGNORE` — re-running is idempotent (no duplicates by `name`).

---

## Error Responses

All errors follow:
```json
{ "status": "error", "message": "<message>" }
```

| Code | Meaning                          |
|------|----------------------------------|
| `400` | Missing/empty param or uninterpretable NLQ query |
| `500` | Server/database failure          |
