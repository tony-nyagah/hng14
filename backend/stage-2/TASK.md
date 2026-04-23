# Stage 2 Backend Task: Intelligence Query Engine

| Field    | Details                                     |
|----------|---------------------------------------------|
| Deadline | 4/23/2026 7:00pm                            |
| Task     | Stage 2 (BACKEND): Intelligence Query Engine |
| Points   | 100 (Pass mark: 75)                         |

---

## Overview

Building on the Stage 1 Profile Intelligence Service, this stage upgrades the system into a **Queryable Intelligence Engine** for Insighta Labs — a demographic intelligence company whose clients (marketing teams, product teams, growth analysts) need to segment users, identify patterns, and query large datasets quickly.

The existing system stores profile data enriched from Genderize, Agify, and Nationalize. It currently lacks:

- Effective filtering
- Combined multi-condition queries
- Sorting and pagination
- Natural language query support

---

## Database Schema

The `profiles` table must follow this structure **exactly**:

| Field                | Type             | Notes                              |
|----------------------|------------------|------------------------------------|
| `id`                 | UUID v7          | Primary key                        |
| `name`               | VARCHAR + UNIQUE | Person's full name                 |
| `gender`             | VARCHAR          | `"male"` or `"female"`             |
| `gender_probability` | FLOAT            | Confidence score                   |
| `age`                | INT              | Exact age                          |
| `age_group`          | VARCHAR          | `child`, `teenager`, `adult`, `senior` |
| `country_id`         | VARCHAR(2)       | ISO code (e.g. `NG`, `KE`)         |
| `country_name`       | VARCHAR          | Full country name                  |
| `country_probability`| FLOAT            | Confidence score                   |
| `created_at`         | TIMESTAMP        | Auto-generated, UTC                |

---

## Data Seeding

- Seed the database with the **2026 profiles** from the provided seed file.
- Re-running the seed must **not create duplicate records** (idempotent by `name`).

---

## Functional Requirements

### 1. Advanced Filtering — `GET /api/profiles`

Supported query parameters:

| Parameter               | Description                        |
|-------------------------|------------------------------------|
| `gender`                | `male` or `female`                 |
| `age_group`             | `child`, `teenager`, `adult`, `senior` |
| `country_id`            | ISO 2-letter code                  |
| `min_age`               | Minimum age (inclusive)            |
| `max_age`               | Maximum age (inclusive)            |
| `min_gender_probability`| Minimum gender confidence score    |
| `min_country_probability`| Minimum country confidence score  |

All filters are **combinable** — results must match **all** provided conditions.

**Example:** `/api/profiles?gender=male&country_id=NG&min_age=25`

---

### 2. Sorting

| Parameter | Values                                   |
|-----------|------------------------------------------|
| `sort_by` | `age` \| `created_at` \| `gender_probability` |
| `order`   | `asc` \| `desc`                          |

**Example:** `/api/profiles?sort_by=age&order=desc`

---

### 3. Pagination

| Parameter | Default | Max |
|-----------|---------|-----|
| `page`    | `1`     | —   |
| `limit`   | `10`    | `50`|

**Response format:**
```json
{
  "status": "success",
  "page": 1,
  "limit": 10,
  "total": 2026,
  "data": [ ... ]
}
```

---

### 4. Natural Language Query — `GET /api/profiles/search`

Converts plain English queries into structured filters. Pagination (`page`, `limit`) applies.

**Example:** `/api/profiles/search?q=young males from nigeria`

**Mapping rules (rule-based only — no AI/LLMs):**

| Query example                          | Interpreted as                                          |
|----------------------------------------|---------------------------------------------------------|
| `"young males"`                        | `gender=male` + `min_age=16` + `max_age=24`             |
| `"females above 30"`                   | `gender=female` + `min_age=30`                          |
| `"people from angola"`                 | `country_id=AO`                                         |
| `"adult males from kenya"`             | `gender=male` + `age_group=adult` + `country_id=KE`     |
| `"male and female teenagers above 17"` | `age_group=teenager` + `min_age=17`                     |

**Notes:**
- `"young"` maps to ages **16–24** for parsing only — it is **not** a stored `age_group`
- Uninterpretable queries return:
  ```json
  { "status": "error", "message": "Unable to interpret query" }
  ```

---

### 5. Query Validation

Invalid query parameters return:
```json
{ "status": "error", "message": "Invalid query parameters" }
```

---

### 6. Performance

- Must handle 2026 records efficiently
- Pagination must be properly implemented
- Avoid unnecessary full-table scans (use indexes where appropriate)

---

## Error Responses

All errors follow this structure:

```json
{ "status": "error", "message": "<error message>" }
```

| Status Code | Meaning                           |
|-------------|-----------------------------------|
| `400`       | Missing or empty parameter        |
| `422`       | Invalid parameter type            |
| `404`       | Profile not found                 |
| `500`/`502` | Server failure                    |

---

## Additional Requirements

- CORS header: `Access-Control-Allow-Origin: *`
- All timestamps in **UTC ISO 8601**
- All IDs in **UUID v7**
- Response structure must match exactly — grading is partially automated

---

## Evaluation Criteria

| Criterion                | Points |
|--------------------------|--------|
| Filtering Logic          | 20     |
| Combined Filters         | 15     |
| Pagination               | 15     |
| Sorting                  | 10     |
| Natural Language Parsing | 20     |
| README Explanation       | 10     |
| Query Validation         | 5      |
| Performance              | 5      |
| **Total**                | **100**|

---

## Submission Format

Submit via `/submit` in `#stage-2-backend`:

1. **GitHub repository link** (must be public)
2. **Public API base URL** (e.g. `https://yourapp.domain.app`)

Ensure:
- Database is seeded with all 2026 profiles
- All endpoints are working
- README is complete
