# Stage 2 Backend — Intelligence Query Engine

A queryable demographic intelligence API built with **Go**, **Gin**, **GORM**, and **PostgreSQL**. Seeded with 2026 profiles, it supports advanced filtering, sorting, pagination, and a rule-based natural language query interface.

---

## Architecture

| Component | Technology |
|-----------|-----------|
| Language  | Go 1.24   |
| Framework | Gin       |
| ORM       | GORM      |
| Database  | PostgreSQL 16 |
| Container | Docker (multi-stage, non-root) |

---

## Quick Start

### With Docker Compose

```bash
cp .env.example .env
docker compose up --build
```

The API will be available at `http://localhost:8060`.

### Locally (requires PostgreSQL running)

```bash
cp .env.example .env
# edit .env with your DB credentials
go run .
```

---

## Environment Variables

| Variable      | Default        | Description              |
|---------------|----------------|--------------------------|
| `DB_HOST`     | `localhost`    | PostgreSQL host          |
| `DB_PORT`     | `5432`         | PostgreSQL port          |
| `DB_USER`     | `postgres`     | Database user            |
| `DB_PASSWORD` | `postgres`     | Database password        |
| `DB_NAME`     | `hng14_stage2` | Database name            |
| `PORT`        | `8060`         | API listen port          |
| `DATABASE_URL`| —              | Full DSN (overrides above)|
| `SEED_FILE`   | `seed_profiles.json` | Path to seed data  |

---

## Endpoints

### `POST /api/profiles`

Create a profile (idempotent by name — returns existing if found).

```json
{ "name": "ella" }
```

Returns `201` on creation, `200` if already exists.

---

### `GET /api/profiles`

List profiles with optional filtering, sorting, and pagination.

**Filter params:**

| Param                   | Example           |
|-------------------------|-------------------|
| `gender`                | `male` / `female` |
| `age_group`             | `adult`           |
| `country_id`            | `NG`              |
| `min_age`               | `25`              |
| `max_age`               | `40`              |
| `min_gender_probability`| `0.8`             |
| `min_country_probability`| `0.5`            |
| `sort_by`               | `age` / `created_at` / `gender_probability` |
| `order`                 | `asc` / `desc`    |
| `page`                  | `1`               |
| `limit`                 | `10` (max `50`)   |

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

Natural language query. Rule-based — no AI or LLMs.

**Examples:**

| Query | Interpreted as |
|-------|---------------|
| `young males from nigeria` | `gender=male`, `min_age=16`, `max_age=24`, `country_id=NG` |
| `females above 30` | `gender=female`, `min_age=30` |
| `people from angola` | `country_id=AO` |
| `adult males from kenya` | `gender=male`, `age_group=adult`, `country_id=KE` |
| `male and female teenagers above 17` | `age_group=teenager`, `min_age=17` |

Uninterpretable queries return:
```json
{ "status": "error", "message": "Unable to interpret query" }
```

Supports `page` and `limit` query params for pagination.

---

### `GET /api/profiles/:id`

Get a single profile by UUID. Returns `404` if not found.

---

### `DELETE /api/profiles/:id`

Delete a profile by UUID. Returns `404` if not found.

---

## NLQ Parsing Rules

The natural language parser uses regex-based rules:

- **Gender:** keywords like `male`, `males`, `men`, `female`, `females`, `women`
- **Age groups:** `child/children`, `teenager/teen`, `adult`, `senior/elderly`
- **"young":** maps to `min_age=16`, `max_age=24` (not a stored age group)
- **Age ranges:** `above N` → `min_age`, `below N` → `max_age`, `between N and M`
- **Country:** `from <country name>` → ISO-2 code lookup (65+ countries supported)
- **Both genders:** `male and female ...` → no gender filter applied

If no recognizable tokens are found, returns `Unable to interpret query`.

---

## Database Schema

```sql
CREATE TABLE profiles (
  id                  VARCHAR(36) PRIMARY KEY,
  name                VARCHAR(255) UNIQUE NOT NULL,
  gender              VARCHAR(10),
  gender_probability  FLOAT,
  age                 INT,
  age_group           VARCHAR(20),
  country_id          VARCHAR(2),
  country_name        VARCHAR(100),
  country_probability FLOAT,
  created_at          TIMESTAMPTZ
);
```

Indexes on: `gender`, `age`, `age_group`, `country_id`, `gender_probability`, `country_probability`, `created_at` — ensures efficient filtering without full-table scans.

---

## Data Seeding

The app auto-seeds on startup from `seed_profiles.json` (2026 profiles). Re-running is idempotent — uses `ON CONFLICT (name) DO NOTHING`.

---

## Error Responses

All errors follow:
```json
{ "status": "error", "message": "<message>" }
```

| Code | Meaning |
|------|---------|
| `400` | Missing/empty required param |
| `404` | Profile not found |
| `422` | Invalid param type/value |
| `502` | Upstream API failure |
