# Stage 2 DevOps Task: Containerize & Ship a Microservices Application

| Field    | Details                                                      |
|----------|--------------------------------------------------------------|
| Deadline | 4/24/2026 1:59am                                             |
| Task     | Containerize & Ship a Microservices Application              |

---

## Overview

In Stage 1, you built and deployed a personal API on a live server using Nginx. In Stage 2, you will take a provided multi-service application and make it production-ready through containerization and a full CI/CD pipeline.

The application is provided — it has **intentional bugs**. Finding them, fixing them, and documenting every single one is a graded part of this assessment.

**Starter repo:** https://github.com/chukwukelu2023/hng14-stage2-devops

---

## The Application

A job processing system made up of four services:

| Service  | Stack          | Role                                              |
|----------|----------------|---------------------------------------------------|
| Frontend | Node.js        | Users submit and track jobs                       |
| API      | Python/FastAPI | Creates jobs and serves status updates            |
| Worker   | Python         | Picks up and processes jobs from the queue        |
| Redis    | Redis          | Shared queue between the API and worker           |

---

## What You Must Do

### 1. Fix the Application

- Read through **all source files** before touching any infrastructure
- Find and fix all bugs — misconfigurations, bad practices, missing production requirements
- Document **every single issue** in `FIXES.md`:
  - File name
  - Line number
  - What the problem was
  - What you changed
- Vague entries will not receive marks

### 2. Containerize It

**Dockerfiles** (one per service):

- Use multi-stage builds where appropriate — final image must not contain build tools or dev dependencies
- All services must run as a **non-root user**
- Each Dockerfile must include a working `HEALTHCHECK` instruction
- No secrets, `.env` files, or credentials may be copied into any image

**`docker-compose.yml`:**

- All services communicate over a **named internal network**
- Redis must **not** be exposed on the host machine
- Services must only start after their dependencies are confirmed **healthy** (not just started)
- All configuration via environment variables — nothing hardcoded in the Compose file
- Include **CPU and memory limits** for every service

### 3. Build the CI/CD Pipeline

Implement a GitHub Actions pipeline on `ubuntu-latest` (free tier only — no self-hosted runners, no paid services).

Stages must run in **strict order** — a failure in any stage blocks all subsequent stages:

```
lint → test → build → security scan → integration test → deploy
```

| Stage              | Requirements                                                                                                  |
|--------------------|---------------------------------------------------------------------------------------------------------------|
| **Lint**           | Python (`flake8`), JavaScript (`eslint`), all Dockerfiles (`hadolint`)                                        |
| **Test**           | ≥ 3 unit tests for the API using `pytest` with Redis mocked; generate and upload coverage report as artifact  |
| **Build**          | Build all 3 images, tag with git SHA + `latest`, push to a local Docker registry (service container in job)   |
| **Security scan**  | Scan all images with Trivy; fail on any `CRITICAL` finding; upload results as SARIF artifact                  |
| **Integration test** | Bring full stack up in runner, submit a job via the frontend, poll until complete, assert final status, tear down cleanly regardless of outcome |
| **Deploy**         | Runs on `main` pushes only; scripted rolling update — new container must pass health check before old one stops; abort and leave old container running if health check fails within 60s |

### 4. Document It

- **`README.md`** — how to bring the full stack up on a clean machine from scratch: prerequisites, all commands, and what a successful startup looks like
- **`FIXES.md`** — every bug found: file, line number, what it was, how you fixed it
- **`.env.example`** — committed with placeholder values for every required variable

---

## Evaluation Criteria / Acceptance Criteria

- Your fork must be **public** — private repos will not be graded
- Everything must be committed — nothing should only exist locally
- `.env` must **never** appear in the repo or git history — this will cost you heavily
- No hardcoded secrets, passwords, or tokens anywhere (YAML, Python, JavaScript, or git history all count)
- No cloud accounts required — runs entirely on your machine and GitHub's free tier
- Do **not** open a pull request to the starter repo

---

## Submission Format

When your repository is ready, head to **#stage-2-devops** and use `/submit`. You will be asked for:

1. Your GitHub username
2. The full URL of your forked repository

**You have 3 attempts.** Read your report carefully before resubmitting — fix what failed, then try again.

**Points:** 100 | **Pass Mark:** 75
