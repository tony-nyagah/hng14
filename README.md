# HNG14

My solutions to the [HNG14](https://hng.tech) internship tasks, organized by track and stage.

## Stack

- **Backend** — Go, deployed on Hetzner via Docker + Traefik
- **Frontend** — Vue 3 + Vite + Tailwind CSS v4 + DaisyUI, deployed on Vercel

## Backend

| Stage | Description | Live | Test Endpoint |
| ----- | ----------- | ---- | ------------- |
| 0 | Name classification API (gender, age, nationality) | [hng14.nyagah.me/backend/stage-0](https://hng14.nyagah.me/backend/stage-0) | [/api/classify?name=john](https://hng14.nyagah.me/backend/stage-0/api/classify?name=john) |

## Frontend

| Stage | Description | Live | Source |
| ----- | ----------- | ---- | ------ |
| 0 | Neubrutalist todo item card | [hng14-taskcard.vercel.app](https://hng14-taskcard.vercel.app/) | [frontend/stage-0/task-card](https://github.com/tony-nyagah/hng14/tree/main/frontend/stage-0/task-card) |

## Repository structure

```
backend/
  stage-0/   Go API — name classification
frontend/
  stage-0/
    task-card/   Vue 3 todo card component
devops/
  stage-0/   Linux server setup
```
