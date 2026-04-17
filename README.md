# HNG14

My solutions to the [HNG14](https://hng.tech) internship tasks, organized by track and stage.

## Stack

- **Backend** — Go, deployed on Hetzner via Docker + Traefik
- **Frontend** — Vue 3 + Vite + Tailwind CSS v4 + DaisyUI, deployed on Vercel

## Backend

| Stage | Description | Live | Test Endpoint |
| ----- | ----------- | ---- | ------------- |
| 0 | Gender classification API | [hng14.nyagah.me/api/classify](https://hng14.nyagah.me/api/classify?name=john) | `GET /api/classify?name=john` |
| 1 | Profile enrichment API (gender, age, nationality) | [hng14.nyagah.me/api/profiles](https://hng14.nyagah.me/api/profiles) | `POST /api/profiles` |

## Frontend

| Stage | Description | Live | Source |
| ----- | ----------- | ---- | ------ |
| 0 | Neubrutalist todo item card | [hng14-taskcard.vercel.app](https://hng14-taskcard.vercel.app/) | [frontend/stage-0/task-card](https://github.com/tony-nyagah/hng14/tree/main/frontend/stage-0/task-card) |

## DevOps

| Stage | Description | Live Base URL | Repo Link |
| ----- | ----------- | ------------- | --------- |
| 0 | Linux server setup | [https://devops.nyagah.me/](https://devops.nyagah.me/) | [devops/stage-0](https://github.com/tony-nyagah/hng14/tree/main/devops/stage-0) |
| 1 | Simple API deployment | [http://ec2-52-91-20-85.compute-1.amazonaws.com](http://ec2-52-91-20-85.compute-1.amazonaws.com) | [devops/stage-1](https://github.com/tony-nyagah/hng14/tree/main/devops/stage-1) |

## Repository structure

```
backend/
  stage-0/   Go API — gender classification
  stage-1/   Go API — profile enrichment (gender, age, nationality)
frontend/
  stage-0/
    task-card/   Vue 3 todo card component
devops/
  stage-0/   Linux server setup
  stage-1/   Simple API deployment
```
