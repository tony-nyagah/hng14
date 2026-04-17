# Namelytics — Profile Intelligence Service

Namelytics is a lightweight REST API that enriches a person's name into a full demographic profile. Given a name, it concurrently queries three public inference APIs — **genderize.io**, **agify.io**, and **nationalize.io** — then assembles and stores a structured profile containing gender, predicted age, age group, and most likely nationality.

Profiles are stored in-memory and keyed by name, making repeated requests for the same name idempotent (no duplicate entries are created).

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | [Go 1.26](https://go.dev/) |
| Web framework | [Gin v1.12](https://github.com/gin-gonic/gin) |
| Unique IDs | [google/uuid](https://github.com/google/uuid) |
| Gender inference | [genderize.io](https://genderize.io/) |
| Age inference | [agify.io](https://agify.io/) |
| Nationality inference | [nationalize.io](https://nationalize.io/) |
| Storage | In-memory map (thread-safe via `sync.RWMutex`) |

---

## Running Locally

**Prerequisites:** Go 1.21+ installed.

```bash
# 1. Clone the repo and navigate to this service
git clone https://github.com/tony-nyagah/hng14.git
cd hng14/backend/stage-1

# 2. Download dependencies
go mod download

# 3. Run the server
go run main.go
```

The server listens on port **8060** by default.

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8060` | Port the HTTP server binds to |

Override the port before running:

```bash
PORT=9090 go run main.go
```

---

## Base URL

```
https://hng14.nyagah.me/backend/stage-1
```

For local development, replace this with:

```
http://localhost:8060
```

---

## API Reference

All responses follow one of two envelope shapes:

**Success**
```json
{ "status": "success", "data": <payload> }
```

**Error**
```json
{ "status": "error", "message": "<human-readable description>" }
```

---

### POST `/api/profiles`

Builds a demographic profile for the given name by querying the three inference APIs concurrently. If a profile for that name already exists, it is returned immediately (idempotent).

**Request body**

```json
{ "name": "ella" }
```

**Responses**

| Status | Meaning |
|---|---|
| `201 Created` | Profile was newly created |
| `200 OK` | Profile already existed; returning cached copy |
| `400 Bad Request` | `name` field missing or empty |
| `502 Bad Gateway` | One or more upstream inference APIs failed |

**Example — create (201)**

```bash
curl -s -X POST https://hng14.nyagah.me/backend/stage-1/api/profiles \
  -H "Content-Type: application/json" \
  -d '{"name": "ella"}'
```

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

**Example — already exists (200)**

```bash
curl -s -X POST https://hng14.nyagah.me/backend/stage-1/api/profiles \
  -H "Content-Type: application/json" \
  -d '{"name": "ella"}'
```

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

---

### GET `/api/profiles`

Returns all profiles currently held in memory.

**Example**

```bash
curl -s https://hng14.nyagah.me/backend/stage-1/api/profiles
```

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
    },
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "james",
      "gender": "male",
      "gender_probability": 0.97,
      "sample_size": 5678,
      "age": 55,
      "age_group": "adult",
      "country_id": "US",
      "created_at": "2025-01-02T12:30:00Z"
    }
  ]
}
```

> Returns an empty array `[]` when no profiles exist yet.

---

### GET `/api/profiles/:id`

Returns a single profile matched by its UUID.

**Path parameter**

| Param | Type | Description |
|---|---|---|
| `id` | `string` (UUID) | The profile's unique identifier |

**Responses**

| Status | Meaning |
|---|---|
| `200 OK` | Profile found |
| `404 Not Found` | No profile with that UUID exists |

**Example**

```bash
curl -s https://hng14.nyagah.me/backend/stage-1/api/profiles/b3f9c1e2-7d4a-4c91-9c2a-1f0a8e5b6d12
```

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

**Example — not found (404)**

```json
{
  "status": "error",
  "message": "Profile not found"
}
```

---

### DELETE `/api/profiles/:id`

Permanently removes a profile from the in-memory store.

**Path parameter**

| Param | Type | Description |
|---|---|---|
| `id` | `string` (UUID) | The profile's unique identifier |

**Responses**

| Status | Meaning |
|---|---|
| `200 OK` | Profile successfully deleted |
| `404 Not Found` | No profile with that UUID exists |

**Example**

```bash
curl -s -X DELETE https://hng14.nyagah.me/backend/stage-1/api/profiles/b3f9c1e2-7d4a-4c91-9c2a-1f0a8e5b6d12
```

```json
{
  "status": "success",
  "message": "Profile deleted successfully"
}
```

**Example — not found (404)**

```json
{
  "status": "error",
  "message": "Profile not found"
}
```

---

## Profile Schema

| Field | Type | Description |
|---|---|---|
| `id` | `string` (UUID v4) | Auto-generated unique identifier |
| `name` | `string` | The name the profile was built from |
| `gender` | `string` | `"male"` or `"female"` (empty string if indeterminate) |
| `gender_probability` | `float64` | Confidence score from genderize.io (0–1) |
| `sample_size` | `int` | Number of samples genderize.io used |
| `age` | `int` | Predicted age from agify.io |
| `age_group` | `string` | `"child"` (0–12), `"teenager"` (13–19), `"adult"` (20–59), `"senior"` (60+) |
| `country_id` | `string` | ISO 3166-1 alpha-2 country code with the highest probability |
| `created_at` | `string` (RFC 3339) | UTC timestamp of when the profile was created |

---

## Notes

- **Concurrency:** all three external API calls are made in parallel using goroutines and a `sync.WaitGroup`, keeping latency close to the slowest single call rather than the sum of all three.
- **Idempotency:** submitting the same name twice never creates a duplicate. The existing profile is returned with `200 OK` on subsequent requests.
- **Persistence:** the store is in-memory only — all profiles are lost on restart. This is intentional for this stage of the project.
