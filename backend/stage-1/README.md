# Namelytics

A REST API that takes a name and returns a demographic profile by querying genderize.io, agify.io, and nationalize.io concurrently.

**Base URL:** `https://hng14.nyagah.me/backend/stage-1`

## Run Locally

```bash
go run main.go
# Runs on :8060 by default. Override with PORT=xxxx
```

## Endpoints

### `POST /api/profiles`
Create a profile from a name. Returns the existing profile if the name was already submitted.

### `GET /api/profiles`
Returns all stored profiles.

### `GET /api/profiles/:id`
Returns a single profile by UUID. Returns `404` if not found.

### `DELETE /api/profiles/:id`
Deletes a profile by UUID. Returns `404` if not found.

## Build



```bash

docker build -t namelytics .

docker run -p 8060:8060 namelytics

```
