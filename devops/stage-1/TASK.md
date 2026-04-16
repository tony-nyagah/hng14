# Stage 1 DevOps Task Brief: Build & Deploy a Personal API

| Field    | Details                        |
|----------|--------------------------------|
| Deadline | 4/18/2026 1:59am               |
| Task     | Build & Deploy a Personal API  |

---

## Overview

In Stage 0, you provisioned a Linux server and configured Nginx. In Stage 1, you will write a small API yourself and deploy it. As a DevOps engineer, understanding what you're deploying matters — this task gives you just enough backend exposure to know how an API is structured and what a running service looks like from the inside.

The backend is intentionally minimal. Most of your effort should go into deployment. You may use any language or framework you're comfortable with (Node.js/Express, Python/FastAPI or Flask, PHP/Laravel, Go, etc.).

---

## What You Must Do

### 1. Build the API

Write an API with the following three endpoints:

**GET /** — returns the following JSON response exactly:
```json
{
  "message": "API is running"
}
```

**GET /health** — returns the following JSON response exactly:
```json
{
  "message": "healthy"
}
```

**GET /me** — returns the following JSON response exactly:
```json
{
  "name": "Your Full Name",
  "email": "you@example.com",
  "github": "https://github.com/yourusername"
}
```

> All three endpoints must return `Content-Type: application/json`, an HTTP status code of `200`, and respond within **500ms**.

---

### 2. Deploy It

Build your API, test it locally, then deploy it publicly using a VPS with an Nginx Reverse Proxy:

- Provision a cloud server (you can reuse your Stage 0 server if you want)
- Run your application on a non-public port
- Configure Nginx to reverse proxy public traffic to your app
- The service must be persistently running — it should not need a manual restart before a reviewer tests it

---

### 3. Document It

Push your code to a public GitHub repository with a README that includes:

- What the project is and how to run it locally
- The three endpoints and their expected responses
- Your live deployment URL

---

## Evaluation Criteria

- All endpoints must return `Content-Type: application/json` — no HTML, no plain text
- All endpoints must return HTTP status `200`
- The `/me` endpoint must contain your real details: name, email, and a valid GitHub profile link
- The app must run on a local port with Nginx proxying public traffic to it — do not expose your app port directly
- The service must stay up on its own — use `systemd`, `pm2`, `supervisor`, or equivalent to keep it alive
- All endpoints must respond within **500ms** — slow responses fail the bot check
- Your GitHub repository must be **public** — private repos will not be reviewed

---

## Submission Format

To submit, go to the `#…` channel and use `/submit` (always remember to send). Submit the following:

1. Your live public base URL (e.g. `http://54.23.255.34` or `http://yourdomain.com`)
2. Your GitHub repository link

---

| Field           | Details |
|-----------------|---------|
| Submission Link | —       |
| Points          | —       |
| Pass Mark       | —       |

> May the forces be with you, Cool Keeds! 🚀
