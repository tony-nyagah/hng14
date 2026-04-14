## HNG14

Solutions to the HNG14 internship tasks. I have categorized them by track and stage.

## Stage 0 Submission Details

- API Base URL: `https://hng14.nyagah.me/backend/stage-0`
- GitHub Repository: `https://github.com/tony-nyagah/hng14`

Test endpoint:

```bash
curl "https://hng14.nyagah.me/backend/stage-0/api/classify?name=john"
```

## Reset deployment (fresh start)

Use the helper script below when you want to take everything down and start clean:

```bash
./scripts/fresh-start.sh
```

If you also want to prune dangling Docker resources (images/containers/networks), run:

```bash
./scripts/fresh-start.sh --prune
```

What this does:
- `docker compose down -v --remove-orphans`
- `docker compose build --no-cache`
- `docker compose up -d --force-recreate`

## Run fresh start on the server

SSH into your Hetzner server and run from the repo path:

```bash
cd /root/hng14
git pull origin main
./scripts/fresh-start.sh
```

With prune:

```bash
cd /root/hng14
git pull origin main
./scripts/fresh-start.sh --prune
```
