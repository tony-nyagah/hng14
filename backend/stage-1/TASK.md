# Stage 1 Backend Task: Data Persistence & API Design Assessment

| Field    | Details                                  |
|----------|------------------------------------------|
| Deadline | 4/18/2026 1:59am                         |
| Task     | Data Persistence & API Design Assessment |

---

## Overview

At this stage, you are expected to:

- Work with multiple external APIs
- Process and structure data
- Persist data in a database
- Design clean and usable APIs

This stage is partially automated and partially reviewed. Only candidates who meet the required quality threshold will move to Stage 2.

---

## Objective

By completing this stage, you must demonstrate that you can:

- Integrate multiple third-party APIs
- Design and implement a database schema
- Store and retrieve structured data
- Build multiple RESTful endpoints
- Handle duplicate data intelligently (idempotency)
- Return clean, consistent JSON responses

---

## Core Concept

You are building a **Profile Intelligence Service**.

Your system will:

- Accept a name
- Enrich it using external APIs
- Store the result
- Allow retrieval and management of stored data

---

## External APIs

You must integrate with (free, no key required):

- **Genderize:** `https://api.genderize.io?name={name}`
- **Agify:** `https://api.agify.io?name={name}`
- **Nationalize:** `https://api.nationalize.io?name={name}`

---

## Functional Requirements

### Processing Rules

- Call all three APIs using the provided name and aggregate the responses
- Extract `gender`, `gender_probability`, and `count` from Genderize. Rename `count` to `sample_size`
- Extract age from Agify. Classify `age_group`:
  - `0–12` → child
  - `13–19` → teenager
  - `20–59` → adult
  - `60+` → senior
- Extract country list from Nationalize. Pick the country with the highest probability as `country_id`
- Store the processed result with a UUID v7 `id` and UTC `created_at` timestamp

---

## API Endpoints

### 1. `POST /api/profiles`

**Request body:**
```json
{ "name": "ella" }
```

**Success response (201):**
```json
{
  "status": "success",
  "data": {
    "id": "b3f9c1e2-7d4a-4c91-9c2a-1f0a8e5b6d12",
    "name": "ella",
    "gender": "female",
    "gender_probability": 0.99,
    "sample_size": 1234,
    "age": 46,
    "age_group": "adult",
    "country_id": "PH",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

> If a profile with the same name already exists, return the existing record (idempotent) with status `200`.

---

### 2. `GET /api/profiles`

Returns all stored profiles.

**Success response (200):**
```json
{
  "status": "success",
  "data": [
    {
      "id": "b3f9c1e2-7d4a-4c91-9c2a-1f0a8e5b6d12",
      "name": "ella",
      "gender": "female",
      "gender_probability": 0.99,
      "sample_size": 1234,
      "age": 46,
      "age_group": "adult",
      "country_id": "PH",
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

---

### 3. `GET /api/profiles/{id}`

Returns a single profile by ID.

**Success response (200):**
```json
{
  "status": "success",
  "data": {
    "id": "b3f9c1e2-7d4a-4c91-9c2a-1f0a8e5b6d12",
    "name": "ella",
    "gender": "female",
    "gender_probability": 0.99,
    "sample_size": 1234,
    "age": 46,
    "age_group": "adult",
    "country_id": "PH",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

**Not found response (404):**
```json
{
  "status": "error",
  "message": "Profile not found"
}
```

---

### 4. `DELETE /api/profiles/{id}`

Deletes a profile by ID.

**Success response (200):**
```json
{
  "status": "success",
  "message": "Profile deleted successfully"
}
```

---

## Evaluation Criteria

- Correct aggregation of data from all three external APIs
- Proper `age_group` classification
- Correct `country_id` selection (highest probability)
- Idempotent `POST` — same name returns existing record
- UUID v7 used for `id`
- All responses use consistent JSON structure with `status` field
- Data is persisted across requests (use a real database, not in-memory)
- All endpoints respond within **500ms**

---

## Submission Format

Submit the following:

1. Your live public base URL (e.g. `http://54.23.255.34` or `http://yourdomain.com`)
2. Your GitHub repository link (must be **public**)
