# DevOps Stage 0 — Linux Server Setup & Nginx Configuration

## Submission Details

- **Live Domain:** `https://devops.nyagah.me`
- **HNG Username:** `nyagah`

Test endpoints:

```bash
curl https://devops.nyagah.me/
curl https://devops.nyagah.me/api
```

---

## What Was Done

### 1. Server

Provisioned an Ubuntu 24.04 EC2 instance on AWS (`t2.micro`, `us-east-1`).

### 2. User Setup

Created a non-root user `hngdevops` with sudo privileges and key-based SSH access:

```bash
useradd -m -s /bin/bash hngdevops
# copied authorized_keys from ubuntu user
# configured passwordless sudo for /usr/sbin/sshd and /usr/sbin/ufw
echo 'hngdevops ALL=(root) NOPASSWD:/usr/sbin/sshd,/usr/sbin/ufw' > /etc/sudoers.d/hngdevops
```

### 3. SSH Hardening

Disabled root login and password-based authentication — key-based only:

```
PermitRootLogin no
PasswordAuthentication no
```

### 4. Firewall (UFW)

Configured UFW to allow only ports 22, 80, and 443:

```bash
ufw default deny incoming
ufw allow 22
ufw allow 80
ufw allow 443
ufw enable
```

### 5. Nginx

Installed Nginx and configured two locations:

- `GET /` — serves a static HTML page with `nyagah` as visible text
- `GET /api` — returns a JSON response with `Content-Type: application/json`

```json
{
  "message": "HNGI14 Stage 0",
  "track": "DevOps",
  "username": "nyagah"
}
```

### 6. SSL

Obtained a valid Let's Encrypt certificate via Certbot:

```bash
certbot --nginx -d devops.nyagah.me --redirect
```

HTTP requests are redirected to HTTPS with a `301 Moved Permanently`.
