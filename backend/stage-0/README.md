# Stage 0 — Name Gender Classification API

Calls the [Genderize API](https://genderize.io) and returns a structured classification result.

## Submission Details

- API Base URL: `https://hng14.nyagah.me/backend/stage-0`
- GitHub Repository: `https://github.com/tony-nyagah/hng14`

Live test endpoint:

```bash
curl "https://hng14.nyagah.me/backend/stage-0/api/classify?name=john"
```

## Endpoint

```
GET /api/classify?name={name}
```

### Success Response
```json
{
  "status": "success",
  "data": {
    "name": "john",
    "gender": "male",
    "probability": 0.99,
    "sample_size": 1234,
    "is_confident": true,
    "processed_at": "2026-04-01T12:00:00Z"
  }
}
```

### Error Responses
| Status | Reason |
|--------|--------|
| 400 | Missing or empty `name` parameter |
| 422 | No prediction available for the name |
| 502 | Upstream API unreachable |
| 500 | Internal server error |

## Run Locally

```bash
go run main.go
curl "http://localhost:8059/api/classify?name=john"
```

## Run with Docker

```bash
docker build -t stage-0 .
docker run -p 8059:8059 stage-0
```
