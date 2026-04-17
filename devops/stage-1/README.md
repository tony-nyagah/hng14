# Simple API

A minimal REST API built with Go and Gin, deployed on a Ubuntu VPS behind an Nginx reverse proxy and managed by systemd.

**Live URL:** http://ec2-52-91-20-85.compute-1.amazonaws.com

---

## Endpoints

| Method | Path | Response |
|--------|------|----------|
| GET | `/` | `{"message": "API is running"}` |
| GET | `/health` | `{"message": "healthy"}` |
| GET | `/me` | `{"name": "...", "email": "...", "github": "..."}` |

All endpoints return:
- HTTP status `200`
- `Content-Type: application/json`
- Response within 500ms

### Examples

```
GET /
{"message":"API is running"}

GET /health
{"message":"healthy"}

GET /me
{"name":"Antony Nyagah","email":"tony.m.nyagah@gmail.com","github":"https://github.com/tony-nyagah"}
```

---

## Running Locally

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)

### Steps

```bash
# 1. Clone the repo
git clone https://github.com/tony-nyagah/hng14.git
cd hng14/devops/stage-1

# 2. Download dependencies
go mod download

# 3. Run the server
go run main.go
```

The API will be available at http://localhost:8080.

---

## Deploying to a VPS

This documents exactly how the live server was set up. The server runs Ubuntu 24.04 on AWS EC2 (t2.micro).

### Prerequisites on your local machine

- Go 1.21+ (to cross-compile the binary)
- SSH access to the server with a key pair
- Nginx already installed on the server (`sudo apt install nginx`)

### 1. Build the Linux binary

Cross-compile from any OS targeting Linux amd64:

```bash
GOOS=linux GOARCH=amd64 go build -o simple-api .
```

### 2. Create the systemd service file

Save this as `simple-api.service`:

```ini
[Unit]
Description=Simple API (Go/Gin)
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/simple-api
ExecStart=/home/ubuntu/simple-api/simple-api
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

Adjust `User` and paths to match the user on your server.

### 3. Create the Nginx config file

Save this as `simple-api.nginx.conf`:

```nginx
server {
    listen 80 default_server;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 5s;
        proxy_read_timeout 10s;
    }
}
```

> If you already have other Nginx sites on port 80 with different `server_name` values, this catch-all block with `default_server` will handle all unmatched requests without disturbing them.

### 4. Transfer files to the server

```bash
# Create the app directory on the server
ssh -i /path/to/key.pem user@your-server "mkdir -p ~/simple-api"

# Copy the binary and config files
scp -i /path/to/key.pem simple-api simple-api.service simple-api.nginx.conf user@your-server:~/simple-api/
```

### 5. Configure systemd

SSH into the server and run:

```bash
chmod +x ~/simple-api/simple-api

sudo cp ~/simple-api/simple-api.service /etc/systemd/system/simple-api.service
sudo systemctl daemon-reload
sudo systemctl enable simple-api   # start on boot
sudo systemctl start simple-api

# Verify it is running
sudo systemctl status simple-api
```

### 6. Configure Nginx

```bash
sudo cp ~/simple-api/simple-api.nginx.conf /etc/nginx/sites-available/simple-api
sudo ln -s /etc/nginx/sites-available/simple-api /etc/nginx/sites-enabled/simple-api

# Test the config before reloading
sudo nginx -t
sudo systemctl reload nginx
```

### 7. Verify

```bash
curl http://your-server-ip/
curl http://your-server-ip/health
curl http://your-server-ip/me
```

---

## Project Structure

```
stage-1/
├── main.go                  # API source code
├── go.mod
├── go.sum
├── simple-api.service       # systemd unit file
└── simple-api.nginx.conf    # Nginx reverse proxy config
```

---

## How It Works

```
Internet
   |
  :80
   |
 Nginx  (reverse proxy)
   |
  :8080
   |
 Go/Gin app  (simple-api, managed by systemd)
```

The app binds only to `127.0.0.1:8080` — it is never exposed directly to the internet. Nginx accepts all public traffic on port 80 and forwards it to the app. Systemd keeps the process alive across restarts and reboots.
