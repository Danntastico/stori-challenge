# 🎓 Complete Deployment Learning Guide
## From Zero to Production HTTPS Deployment

**Created:** October 28, 2025  
**Project:** Stori Challenge - Full Stack Financial Tracker  
**Final URL:** https://stori.danntastico.dev/

---

## 📚 Table of Contents

1. [The Big Picture](#1-the-big-picture)
2. [Core Concepts](#2-core-concepts)
3. [Step-by-Step Process](#3-step-by-step-process)
4. [Command Reference](#4-command-reference)
5. [Troubleshooting](#5-troubleshooting)
6. [Quiz Yourself](#6-quiz-yourself)

---

## 1. The Big Picture

### 🎯 What We Built

```
User's Browser
     ↓
Domain Name (stori.danntastico.dev)
     ↓
DNS Resolution (converts name → IP address)
     ↓
EC2 Server (107.22.66.194)
     ↓
Nginx (web server + reverse proxy)
  ├─> Frontend (React app - static files)
  └─> Backend (Go API - running on :8080)
```

### 🤔 Why This Architecture?

**Question to consider:** Why not just access the server by IP address?

**Answer:** 
- IP addresses are hard to remember (107.22.66.194 vs stori.danntastico.dev)
- SSL certificates require domain names (Let's Encrypt won't issue certs for raw IPs)
- Professional appearance for portfolio/interviews
- Flexibility (can change server without changing URL)

---

## 2. Core Concepts

### 2.1 DNS (Domain Name System)

#### 🧠 What is DNS?

DNS is like the **phonebook of the internet**. It converts human-readable names into computer-readable IP addresses.

**Real-world analogy:**
- You want to call your friend "John" → you look up John's phone number in your contacts
- You want to visit "stori.danntastico.dev" → DNS looks up its IP address (107.22.66.194)

#### 📋 DNS Record Types

| Record Type | Purpose | Example |
|-------------|---------|---------|
| **A Record** | Maps domain → IPv4 address | `stori.danntastico.dev` → `107.22.66.194` |
| **AAAA Record** | Maps domain → IPv6 address | (not used in our project) |
| **CNAME** | Maps domain → another domain | `www.example.com` → `example.com` |
| **MX** | Mail server records | (for email, not used here) |

**What we configured:**
```
Type: A Record
Host: stori
Value: 107.22.66.194
Result: stori.danntastico.dev → 107.22.66.194
```

#### 🕐 DNS Propagation

**Question:** Why didn't the domain work immediately after adding the DNS record?

**Answer:** DNS changes need to **propagate** (spread) across the internet:
- Your DNS provider (Namecheap) updates their servers instantly
- Other DNS servers around the world gradually sync this information
- Takes 5 minutes to 48 hours (usually 10-30 minutes)
- Like updating a phonebook - takes time for everyone to get the new edition

**How to check DNS propagation:**
```bash
nslookup stori.danntastico.dev
# Returns the IP address if propagated
```

---

### 2.2 SSL/TLS & HTTPS

#### 🔒 What is SSL/TLS?

**SSL** (Secure Sockets Layer) / **TLS** (Transport Layer Security) encrypts data between browser and server.

**Without HTTPS (HTTP only):**
```
Browser: "Get me the dashboard"  ← Anyone can read this!
Server: "Here's the data: {...}"  ← Anyone can read this!
```

**With HTTPS (HTTP + SSL):**
```
Browser: "x9#mK2@pLq!..." ← Encrypted, unreadable
Server: "z7&nP4@rMt!..." ← Encrypted, unreadable
```

#### 🤔 Why Do We Need SSL Certificates?

1. **Encryption** - Prevents eavesdropping (man-in-the-middle attacks)
2. **Authentication** - Proves you're talking to the real server, not an impostor
3. **Trust** - Browsers show green padlock 🔒
4. **Required** - Modern web requires HTTPS (especially for `.dev` domains!)

#### 🆓 Let's Encrypt

**What is it?**
- Free, automated Certificate Authority (CA)
- Issues SSL certificates for free
- Trusted by all major browsers
- Certificates valid for 90 days (auto-renewable)

**Why 90 days?**
- Forces automation (can't forget to renew manually)
- Limits damage if private key is compromised
- Encourages best practices

#### 🤖 Certbot

**What is Certbot?**
- Automated tool for getting Let's Encrypt certificates
- Handles certificate requests, renewals, and nginx configuration

**How it works:**
```
1. You: "Certbot, get me a cert for stori.danntastico.dev"
2. Certbot → Let's Encrypt: "I want a cert for stori.danntastico.dev"
3. Let's Encrypt: "Prove you control that domain!"
4. Let's Encrypt: Creates challenge file at:
   http://stori.danntastico.dev/.well-known/acme-challenge/RANDOM_STRING
5. Certbot: Serves that file via nginx
6. Let's Encrypt: Checks if file is accessible (verifies ownership)
7. Let's Encrypt: "Proof confirmed! Here's your certificate!"
8. Certbot: Installs certificate and configures nginx
```

**This is why port 80 must be open!** Let's Encrypt needs to access that challenge file.

---

### 2.3 Nginx

#### 🌐 What is Nginx?

**Nginx** (pronounced "engine-x") is a **web server** and **reverse proxy**.

**Two main roles in our project:**

#### Role 1: Static File Server (for Frontend)

**What it does:**
```
Browser requests: https://stori.danntastico.dev/
Nginx: "Here's the index.html file!"
Browser requests: https://stori.danntastico.dev/assets/index-B6BlV7hA.js
Nginx: "Here's that JavaScript file!"
```

Like a file cabinet - stores and serves files when requested.

#### Role 2: Reverse Proxy (for Backend)

**What it does:**
```
Browser requests: https://stori.danntastico.dev/api/health
Nginx: "I'll forward this to the Go backend at localhost:8080"
Go backend: "Here's the response: {status: healthy}"
Nginx: "I'll forward this response back to the browser"
```

**Why use a reverse proxy?**

🤔 **Question:** Why not expose the Go backend directly to the internet on port 8080?

**Answers:**
1. **SSL Termination** - Nginx handles HTTPS, Go app can stay simple (HTTP)
2. **Single Entry Point** - One server, one SSL certificate for frontend + backend
3. **Security** - Backend never directly exposed, nginx filters requests
4. **Flexibility** - Can change backend port/technology without changing public URLs
5. **Performance** - Nginx is optimized for serving static files, Go handles business logic
6. **Load Balancing** - Can add multiple backend servers behind nginx (future scaling)

**Real-world analogy:**
- **Nginx** = Receptionist at a company
- **Backend servers** = Employees in offices
- Clients talk to receptionist → receptionist routes to correct employee → receptionist returns response

---

### 2.4 Ports & Firewalls

#### 🚪 What are Ports?

Think of your server as an **apartment building**:
- **IP address** = Building address (107.22.66.194)
- **Ports** = Individual apartment numbers (80, 443, 8080, etc.)
- **Services** = Residents in each apartment

**Common ports:**
```
Port 80   = HTTP  (unencrypted web traffic)
Port 443  = HTTPS (encrypted web traffic)
Port 22   = SSH   (secure remote access)
Port 8080 = Our Go backend (internal only)
```

#### 🛡️ Firewalls & Security Groups

**AWS Security Groups** = Virtual firewall for your EC2 instance

**Rules we configured:**
```
Port 22 (SSH)   → Only your IP can connect (remote admin)
Port 80 (HTTP)  → Open to world (Let's Encrypt verification + redirects)
Port 443 (HTTPS) → Open to world (public web traffic)
Port 8080       → NOT open (backend stays internal!)
```

🔒 **Security principle:** Only open ports that MUST be public. Keep everything else internal.

---

## 3. Step-by-Step Process

### Phase 1: Domain & DNS Setup

#### Step 1.1: Purchase Domain

**Command:** (Done via web interface)
- Provider: Namecheap.com
- Domain: `danntastico.dev`
- Cost: $6.98/year

**Why `.dev`?**
- Professional for developers
- HSTS preloaded (forces HTTPS by default)
- Shows technical awareness in interviews

#### Step 1.2: Configure DNS

**Where:** Namecheap Dashboard → Advanced DNS

**What we added:**
```
Type:  A Record
Host:  stori
Value: 107.22.66.194
TTL:   Automatic (300 seconds = 5 minutes)
```

**What this means:**
- Full domain becomes: `stori.danntastico.dev`
- Points to EC2 instance IP
- TTL = Time To Live (how long DNS servers cache this info)

**Verification command:**
```bash
nslookup stori.danntastico.dev
```

**Expected output:**
```
Name:   stori.danntastico.dev
Address: 107.22.66.194
```

---

### Phase 2: Prepare Application Files

#### Step 2.1: Update Frontend Environment

**File:** `frontend/.env`

**Command:**
```bash
cd /home/misterfancybg/Documents/stori/stori-challenge/frontend
echo "VITE_API_BASE_URL=https://stori.danntastico.dev/api" > .env
```

**Why?**
- Frontend needs to know where to find the API
- Changed from IP address (`http://107.22.66.194:8080/api`) to domain
- Changed from HTTP to HTTPS

**How Vite uses this:**
```javascript
// In your frontend code:
axios.get(`${import.meta.env.VITE_API_BASE_URL}/health`)
// Becomes:
axios.get('https://stori.danntastico.dev/api/health')
```

#### Step 2.2: Build Frontend

**Command:**
```bash
cd /home/misterfancybg/Documents/stori/stori-challenge/frontend
npm run build
```

**What happens:**
1. Vite reads all React/TypeScript files
2. Bundles them into optimized JavaScript
3. Processes CSS (TailwindCSS → regular CSS)
4. Minifies everything (removes whitespace, shortens variable names)
5. Outputs to `dist/` folder:
   ```
   dist/
   ├── index.html              (entry point)
   ├── assets/
   │   ├── index-B6BlV7hA.js  (bundled JavaScript)
   │   └── index-D8wwfAmX.css (bundled CSS)
   └── vite.svg
   ```

**Why build?**
- React/TS can't run directly in browsers
- Optimization for production (smaller files = faster loading)
- All imports resolved, code split

---

### Phase 3: EC2 Deployment

#### Step 3.1: Create Frontend Directory

**Command:**
```bash
ssh -i ~/Documents/stori/stori-expenses-backend-key.pem ec2-user@107.22.66.194 \
  "sudo mkdir -p /var/www/stori && sudo chown ec2-user:ec2-user /var/www/stori"
```

**Breaking it down:**
- `ssh` = Secure Shell (remote connection to server)
- `-i ~/Documents/stori/stori-expenses-backend-key.pem` = Use this private key for authentication
- `ec2-user@107.22.66.194` = Connect as user `ec2-user` to IP `107.22.66.194`
- `"..."` = Commands to run on the remote server

**Inside the quotes:**
- `sudo` = Run as administrator (root)
- `mkdir -p /var/www/stori` = Create directory (and parent dirs if needed)
- `sudo chown ec2-user:ec2-user /var/www/stori` = Change owner to ec2-user (so we can upload files without sudo)

**Why `/var/www/`?**
- Standard Linux location for web files
- Convention that all sysadmins recognize

#### Step 3.2: Upload Frontend Files

**Command:**
```bash
cd /home/misterfancybg/Documents/stori/stori-challenge/frontend
scp -i ~/Documents/stori/stori-expenses-backend-key.pem -r dist/* ec2-user@107.22.66.194:/var/www/stori/
```

**Breaking it down:**
- `scp` = Secure Copy (copy files over SSH)
- `-i ...` = Use this private key
- `-r` = Recursive (copy directories and their contents)
- `dist/*` = Source: all files in dist/ folder
- `ec2-user@107.22.66.194:/var/www/stori/` = Destination: remote server path

**What happens:**
```
Local:                        Remote:
dist/index.html        →      /var/www/stori/index.html
dist/assets/index.js   →      /var/www/stori/assets/index.js
dist/assets/index.css  →      /var/www/stori/assets/index.css
```

---

### Phase 4: Nginx Configuration

#### Step 4.1: Create Nginx Config File

**File:** `nginx-stori.conf`

**Content explanation:**

```nginx
server {
    listen 80;
    server_name stori.danntastico.dev;
```
**Translation:** "Listen on port 80 (HTTP) for requests to stori.danntastico.dev"

```nginx
    root /var/www/stori;
    index index.html;
```
**Translation:** "Look for files in /var/www/stori/, default file is index.html"

```nginx
    # Frontend - serve React app
    location / {
        try_files $uri $uri/ /index.html;
    }
```
**Translation:** "For any request:
1. Try to find exact file (`$uri`)
2. If not found, try as directory (`$uri/`)
3. If still not found, serve index.html (React Router handles the rest)"

**Why the fallback to index.html?**
- React uses client-side routing
- URLs like `/dashboard` don't exist as actual files on server
- Nginx serves index.html → React router reads URL → shows correct page

```nginx
    # Backend API - proxy to Go service
    location /api/ {
        proxy_pass http://localhost:8080/api/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
```

**Translation:** "For any request starting with /api/:
- Forward it to localhost:8080/api/
- Pass along important headers (so backend knows original client info)"

**Headers explained:**
- `Host` = Original domain requested (stori.danntastico.dev)
- `X-Real-IP` = Client's actual IP address
- `X-Forwarded-For` = Chain of proxies (for logging)
- `X-Forwarded-Proto` = Original protocol (http or https)

#### Step 4.2: Upload & Install Config

**Upload:**
```bash
scp -i ~/Documents/stori/stori-expenses-backend-key.pem \
  /home/misterfancybg/Documents/stori/stori-challenge/nginx-stori.conf \
  ec2-user@107.22.66.194:/tmp/nginx-stori.conf
```

**Move to correct location:**
```bash
ssh -i ~/Documents/stori/stori-expenses-backend-key.pem ec2-user@107.22.66.194 \
  "sudo mv /tmp/nginx-stori.conf /etc/nginx/conf.d/stori.conf"
```

**Why `/etc/nginx/conf.d/`?**
- Nginx automatically loads all `.conf` files from this directory
- Keeps configurations organized (one file per site)
- Main nginx.conf includes: `include /etc/nginx/conf.d/*.conf;`

#### Step 4.3: Test & Reload

**Test configuration syntax:**
```bash
ssh -i ~/Documents/stori/stori-expenses-backend-key.pem ec2-user@107.22.66.194 \
  "sudo nginx -t"
```

**Output:**
```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

**Why test first?**
- Bad configuration could crash nginx
- Testing prevents downtime

**Reload nginx:**
```bash
ssh -i ~/Documents/stori/stori-expenses-backend-key.pem ec2-user@107.22.66.194 \
  "sudo systemctl reload nginx"
```

**`reload` vs `restart`:**
- `reload` = Gracefully reload config without dropping connections
- `restart` = Stop everything, then start again (brief downtime)

---

### Phase 5: SSL Certificate Setup

#### Step 5.1: Open Port 80 (Firewall)

**Why?**
Let's Encrypt needs to verify domain ownership by accessing:
```
http://stori.danntastico.dev/.well-known/acme-challenge/RANDOM_TOKEN
```

If port 80 is closed, verification fails!

**Find security group:**
```bash
aws ec2 describe-instances \
  --filters "Name=ip-address,Values=107.22.66.194" \
  --query 'Reservations[0].Instances[0].SecurityGroups[*].GroupId' \
  --output text
```

**Output:** `sg-0c62599075fef806c`

**Open port 80:**
```bash
aws ec2 authorize-security-group-ingress \
  --group-id sg-0c62599075fef806c \
  --protocol tcp \
  --port 80 \
  --cidr 0.0.0.0/0
```

**Breaking it down:**
- `authorize-security-group-ingress` = Add inbound rule
- `--group-id` = Which security group to modify
- `--protocol tcp --port 80` = Allow TCP traffic on port 80
- `--cidr 0.0.0.0/0` = From anywhere (entire internet)

**Test it works:**
```bash
curl -I http://stori.danntastico.dev/
```

**Expected:** `HTTP/1.1 200 OK` (nginx is answering on port 80)

#### Step 5.2: Install Certbot

**Command:**
```bash
ssh -i ~/Documents/stori/stori-expenses-backend-key.pem ec2-user@107.22.66.194 \
  "sudo dnf install -y certbot python3-certbot-nginx"
```

**What gets installed:**
- `certbot` = Main program
- `python3-certbot-nginx` = Nginx plugin (auto-configures nginx for HTTPS)

**Package manager `dnf`:**
- Amazon Linux 2023 uses `dnf` (Fedora/RedHat family)
- Similar to `apt` (Debian/Ubuntu) or `brew` (macOS)
- Manages software installation, updates, dependencies

#### Step 5.3: Get SSL Certificate

**Command:**
```bash
ssh -i ~/Documents/stori/stori-expenses-backend-key.pem ec2-user@107.22.66.194 \
  "sudo certbot --nginx -d stori.danntastico.dev --non-interactive --agree-tos --email danntastico@gmail.com --redirect"
```

**Breaking it down:**
- `certbot` = Run the certbot program
- `--nginx` = Use nginx plugin (auto-configure)
- `-d stori.danntastico.dev` = Domain to get certificate for
- `--non-interactive` = Don't ask questions (assume defaults)
- `--agree-tos` = Agree to Let's Encrypt Terms of Service
- `--email danntastico@gmail.com` = Email for renewal notifications
- `--redirect` = Automatically redirect HTTP → HTTPS

**What certbot does (step by step):**

1. **Registers account with Let's Encrypt**
   ```
   "Hi Let's Encrypt, I'm danntastico@gmail.com"
   ```

2. **Requests certificate**
   ```
   "I want a certificate for stori.danntastico.dev"
   ```

3. **Let's Encrypt issues challenge**
   ```
   "Put this file at: http://stori.danntastico.dev/.well-known/acme-challenge/RANDOM"
   ```

4. **Certbot creates challenge file**
   - Modifies nginx config temporarily
   - Serves the challenge file

5. **Let's Encrypt verifies**
   ```
   Let's Encrypt → http://stori.danntastico.dev/.well-known/acme-challenge/RANDOM
   "File found! Domain ownership verified!"
   ```

6. **Let's Encrypt issues certificate**
   ```
   Files saved to:
   - /etc/letsencrypt/live/stori.danntastico.dev/fullchain.pem (certificate)
   - /etc/letsencrypt/live/stori.danntastico.dev/privkey.pem (private key)
   ```

7. **Certbot configures nginx**
   - Modifies `/etc/nginx/conf.d/stori.conf`
   - Adds SSL configuration
   - Sets up HTTP → HTTPS redirect

**Final nginx config (auto-generated by certbot):**
```nginx
server {
    listen 80;
    server_name stori.danntastico.dev;
    
    # Redirect all HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    server_name stori.danntastico.dev;
    
    # SSL Certificate files
    ssl_certificate /etc/letsencrypt/live/stori.danntastico.dev/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/stori.danntastico.dev/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;
    
    # Rest of configuration (root, locations, etc.)
    root /var/www/stori;
    location / { ... }
    location /api/ { ... }
}
```

#### Step 5.4: Auto-Renewal

**Certbot automatically sets up renewal!**

**How?**
```bash
# Systemd timer (like cron)
systemctl list-timers | grep certbot
```

**What happens:**
- Timer runs twice daily
- Checks if certificates expire in < 30 days
- If yes, automatically renews
- Reloads nginx to apply new cert

**Manual renewal test:**
```bash
sudo certbot renew --dry-run
```

---

### Phase 6: Update Backend CORS

#### Why?

CORS = Cross-Origin Resource Sharing

**Security rule:** Browsers block requests between different origins unless explicitly allowed.

**Origin = Protocol + Domain + Port**

Examples:
- `https://stori.danntastico.dev` (origin 1)
- `https://stori.danntastico.dev:443` (same as above, port 443 is default for HTTPS)
- `http://localhost:8080` (origin 2 - different!)

**Our situation:**
- Frontend runs at: `https://stori.danntastico.dev`
- Backend runs at: `https://stori.danntastico.dev/api/` (same origin!)

But backend needs to know which origins to trust!

#### Step 6.1: Update Environment File

**Command:**
```bash
ssh -i ~/Documents/stori/stori-expenses-backend-key.pem ec2-user@107.22.66.194 \
  "echo 'PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:5173,https://stori.danntastico.dev
LOG_LEVEL=info
OPENAI_API_KEY=' > /home/ec2-user/.stori-backend.env"
```

**What this does:**
- Creates environment file with allowed origins
- `http://localhost:5173` = For local development
- `https://stori.danntastico.dev` = For production

**How backend uses this:**
```go
// In backend code (middleware/cors.go)
allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
// If request comes from allowed origin, add header:
w.Header().Set("Access-Control-Allow-Origin", origin)
```

#### Step 6.2: Restart Backend

**Command:**
```bash
ssh -i ~/Documents/stori/stori-expenses-backend-key.pem ec2-user@107.22.66.194 \
  "sudo systemctl restart stori-backend"
```

**What happens:**
1. Systemd stops the current Go process
2. Reads new environment variables from `.stori-backend.env`
3. Starts new Go process with new settings

---

## 4. Command Reference

### 4.1 SSH Commands

**Basic SSH connection:**
```bash
ssh -i /path/to/key.pem username@hostname
```

**Run command on remote server:**
```bash
ssh -i /path/to/key.pem username@hostname "command to run"
```

**Examples:**
```bash
# Check if file exists
ssh -i key.pem ec2-user@107.22.66.194 "ls -la /var/www/stori"

# View logs
ssh -i key.pem ec2-user@107.22.66.194 "sudo journalctl -u stori-backend -n 50"

# Check nginx status
ssh -i key.pem ec2-user@107.22.66.194 "sudo systemctl status nginx"
```

---

### 4.2 SCP Commands

**Upload file:**
```bash
scp -i /path/to/key.pem /local/file username@hostname:/remote/path/
```

**Upload directory (recursive):**
```bash
scp -i /path/to/key.pem -r /local/dir/* username@hostname:/remote/path/
```

**Download file:**
```bash
scp -i /path/to/key.pem username@hostname:/remote/file /local/path/
```

**Examples:**
```bash
# Upload frontend build
scp -i key.pem -r dist/* ec2-user@107.22.66.194:/var/www/stori/

# Download logs
scp -i key.pem ec2-user@107.22.66.194:/var/log/nginx/error.log ./
```

---

### 4.3 Nginx Commands

**Test configuration:**
```bash
sudo nginx -t
```

**Reload configuration (no downtime):**
```bash
sudo systemctl reload nginx
```

**Restart nginx:**
```bash
sudo systemctl restart nginx
```

**Check status:**
```bash
sudo systemctl status nginx
```

**View logs:**
```bash
# Error log
sudo tail -f /var/log/nginx/error.log

# Access log
sudo tail -f /var/log/nginx/access.log
```

---

### 4.4 Certbot Commands

**Get new certificate:**
```bash
sudo certbot --nginx -d yourdomain.com
```

**List all certificates:**
```bash
sudo certbot certificates
```

**Renew all certificates:**
```bash
sudo certbot renew
```

**Test renewal (dry run):**
```bash
sudo certbot renew --dry-run
```

**Revoke certificate:**
```bash
sudo certbot revoke --cert-path /etc/letsencrypt/live/domain/fullchain.pem
```

---

### 4.5 Systemd Commands

**Start service:**
```bash
sudo systemctl start service-name
```

**Stop service:**
```bash
sudo systemctl stop service-name
```

**Restart service:**
```bash
sudo systemctl restart service-name
```

**Check status:**
```bash
sudo systemctl status service-name
```

**Enable (start on boot):**
```bash
sudo systemctl enable service-name
```

**View logs:**
```bash
sudo journalctl -u service-name
# Last 50 lines
sudo journalctl -u service-name -n 50
# Follow (real-time)
sudo journalctl -u service-name -f
```

---

### 4.6 DNS Commands

**Check DNS resolution:**
```bash
nslookup domain.com
```

**Alternative DNS check:**
```bash
dig domain.com
```

**Check from specific DNS server:**
```bash
nslookup domain.com 8.8.8.8  # Google's DNS
```

---

### 4.7 Network/Testing Commands

**Test HTTP connection:**
```bash
curl http://domain.com
```

**Test HTTPS connection:**
```bash
curl https://domain.com
```

**Show response headers only:**
```bash
curl -I https://domain.com
```

**Test specific endpoint:**
```bash
curl https://stori.danntastico.dev/api/health
```

**Test with verbose output:**
```bash
curl -v https://domain.com
```

**Check open ports:**
```bash
# From remote server
sudo netstat -tlnp | grep nginx
```

---

## 5. Troubleshooting

### 5.1 Common Issues & Solutions

#### Issue 1: DNS Not Resolving

**Symptom:**
```bash
nslookup stori.danntastico.dev
# Server can't find stori.danntastico.dev: NXDOMAIN
```

**Solutions:**
1. **Wait for propagation** (5-30 minutes)
2. **Check DNS record in registrar**
   - Verify Host = `stori`
   - Verify Value = correct IP
   - Verify Type = `A Record`
3. **Clear local DNS cache:**
   ```bash
   # Linux
   sudo systemd-resolve --flush-caches
   
   # macOS
   sudo dscacheutil -flushcache
   ```
4. **Test with Google DNS:**
   ```bash
   nslookup stori.danntastico.dev 8.8.8.8
   ```

---

#### Issue 2: Port Not Accessible

**Symptom:**
```bash
curl http://stori.danntastico.dev
# curl: (28) Connection timeout
```

**Solutions:**
1. **Check security group** (AWS Console → EC2 → Security Groups)
   - Verify port 80 is open (0.0.0.0/0)
   - Verify port 443 is open (0.0.0.0/0)

2. **Check nginx is running:**
   ```bash
   ssh -i key.pem ec2-user@IP "sudo systemctl status nginx"
   ```

3. **Check nginx is listening:**
   ```bash
   ssh -i key.pem ec2-user@IP "sudo netstat -tlnp | grep nginx"
   # Should show: 0.0.0.0:80 and 0.0.0.0:443
   ```

---

#### Issue 3: SSL Certificate Fails

**Symptom:**
```
Certbot failed to authenticate some domains
Detail: Timeout during connect (likely firewall problem)
```

**Solutions:**
1. **Ensure port 80 is open** (see Issue 2)
2. **Verify DNS is working:**
   ```bash
   nslookup stori.danntastico.dev
   ```
3. **Check nginx can serve challenge files:**
   ```bash
   curl http://stori.danntastico.dev/.well-known/acme-challenge/test
   ```
4. **Review certbot logs:**
   ```bash
   ssh -i key.pem ec2-user@IP "sudo cat /var/log/letsencrypt/letsencrypt.log"
   ```

---

#### Issue 4: Mixed Content Errors

**Symptom:**
```
Mixed Content: The page at 'https://...' was loaded over HTTPS,
but requested an insecure resource 'http://...'
```

**Solutions:**
1. **Check frontend .env file:**
   ```bash
   cat frontend/.env
   # Should be: VITE_API_BASE_URL=https://domain.com/api
   # NOT: http://...
   ```

2. **Rebuild frontend after changing .env:**
   ```bash
   cd frontend
   npm run build
   ```

3. **Ensure ALL resources use HTTPS**

---

#### Issue 5: CORS Errors

**Symptom:**
```
Access to XMLHttpRequest at 'https://api...' from origin 'https://frontend...'
has been blocked by CORS policy
```

**Solutions:**
1. **Check backend CORS configuration:**
   ```bash
   ssh -i key.pem ec2-user@IP "cat /home/ec2-user/.stori-backend.env"
   # Verify CORS_ALLOWED_ORIGINS includes your frontend domain
   ```

2. **Restart backend after changes:**
   ```bash
   ssh -i key.pem ec2-user@IP "sudo systemctl restart stori-backend"
   ```

3. **Check browser console** for exact error message

---

#### Issue 6: Nginx Configuration Errors

**Symptom:**
```bash
sudo nginx -t
# nginx: configuration file /etc/nginx/nginx.conf test failed
```

**Solutions:**
1. **Read error message carefully** - it tells you line number
2. **Common mistakes:**
   - Missing semicolon `;`
   - Unmatched braces `{ }`
   - Wrong file paths
3. **Restore previous config if needed:**
   ```bash
   sudo cp /etc/nginx/conf.d/stori.conf.bak /etc/nginx/conf.d/stori.conf
   ```

---

#### Issue 7: Backend Not Responding

**Symptom:**
```bash
curl https://stori.danntastico.dev/api/health
# 502 Bad Gateway
```

**Solutions:**
1. **Check if backend is running:**
   ```bash
   ssh -i key.pem ec2-user@IP "sudo systemctl status stori-backend"
   ```

2. **Check backend logs:**
   ```bash
   ssh -i key.pem ec2-user@IP "sudo journalctl -u stori-backend -n 50"
   ```

3. **Verify backend is listening on correct port:**
   ```bash
   ssh -i key.pem ec2-user@IP "sudo netstat -tlnp | grep 8080"
   # Should show Go process listening on 127.0.0.1:8080
   ```

4. **Restart backend:**
   ```bash
   ssh -i key.pem ec2-user@IP "sudo systemctl restart stori-backend"
   ```

---

## 6. Quiz Yourself

Test your understanding! Try to answer these without looking back:

### Basic Concepts

1. **What does DNS do?**
   <details>
   <summary>Show answer</summary>
   Converts human-readable domain names (stori.danntastico.dev) into IP addresses (107.22.66.194) that computers use to communicate.
   </details>

2. **What's the difference between HTTP and HTTPS?**
   <details>
   <summary>Show answer</summary>
   HTTPS = HTTP + SSL/TLS encryption. HTTPS encrypts data in transit, preventing eavesdropping and tampering. HTTPS also authenticates the server.
   </details>

3. **What is a reverse proxy?**
   <details>
   <summary>Show answer</summary>
   A server (like nginx) that sits in front of backend servers, accepting client requests and forwarding them to the appropriate backend server. It's "reverse" because it forwards requests to servers, not clients.
   </details>

4. **Why do we need port 80 open for Let's Encrypt?**
   <details>
   <summary>Show answer</summary>
   Let's Encrypt verifies domain ownership by accessing a challenge file at http://domain/.well-known/acme-challenge/TOKEN. If port 80 is closed, Let's Encrypt can't reach this file and verification fails.
   </details>

### Commands

5. **What does this command do?**
   ```bash
   sudo systemctl reload nginx
   ```
   <details>
   <summary>Show answer</summary>
   Reloads nginx configuration without stopping the service (no downtime). Nginx gracefully applies new configuration while keeping existing connections alive.
   </details>

6. **What's the difference between these two?**
   ```bash
   ssh user@host "command"
   scp file user@host:/path/
   ```
   <details>
   <summary>Show answer</summary>
   - `ssh` runs a command on the remote server
   - `scp` copies files to/from the remote server
   Both use SSH protocol for secure communication.
   </details>

### Architecture

7. **Draw the request flow for: `https://stori.danntastico.dev/api/health`**
   <details>
   <summary>Show answer</summary>
   ```
   1. Browser → DNS lookup → Gets IP 107.22.66.194
   2. Browser → HTTPS request to 107.22.66.194:443
   3. Nginx → Terminates SSL, sees path /api/health
   4. Nginx → Proxies to localhost:8080/api/health
   5. Go backend → Processes request, returns JSON
   6. Nginx → Forwards response to browser
   7. Browser → Displays response
   ```
   </details>

8. **Why do we use nginx in front of the Go backend instead of exposing Go directly?**
   <details>
   <summary>Show answer</summary>
   Multiple reasons:
   - SSL termination (nginx handles HTTPS, Go stays simple)
   - Single entry point (one certificate for frontend + backend)
   - Security (backend never directly exposed)
   - Performance (nginx optimized for static files)
   - Flexibility (can change backend without changing public URLs)
   </details>

### Troubleshooting

9. **You just updated the frontend code and ran `npm run build`. Users still see the old version. What's wrong?**
   <details>
   <summary>Show answer</summary>
   You forgot to upload the new `dist/` files to the server!
   
   ```bash
   cd frontend
   scp -i key.pem -r dist/* ec2-user@IP:/var/www/stori/
   ```
   
   Nginx serves files from `/var/www/stori/`, so you must update those files after each build.
   </details>

10. **You changed backend CORS settings but still getting CORS errors. What did you forget?**
    <details>
    <summary>Show answer</summary>
    Restart the backend service!
    
    ```bash
    ssh -i key.pem ec2-user@IP "sudo systemctl restart stori-backend"
    ```
    
    Environment variables are read when the process starts, not dynamically.
    </details>

---

## 7. Key Takeaways

### 🎯 Essential Concepts

1. **DNS is the phonebook of the internet** - converts names to IPs
2. **HTTPS = HTTP + SSL/TLS** - encrypts data, authenticates servers
3. **Let's Encrypt provides free SSL certificates** - automated with Certbot
4. **Nginx is a web server + reverse proxy** - serves files and forwards requests
5. **Ports are like apartment numbers** - different services on different ports
6. **Security groups are firewalls** - control which ports are accessible

### 🛠️ Essential Commands

```bash
# DNS
nslookup domain.com

# Test website
curl -I https://domain.com

# SSH
ssh -i key.pem user@host "command"

# Copy files
scp -i key.pem file user@host:/path/

# Nginx
sudo nginx -t              # Test config
sudo systemctl reload nginx # Apply config

# Certbot
sudo certbot --nginx -d domain.com

# Systemd
sudo systemctl status service
sudo systemctl restart service
```

### 📊 Our Final Architecture

```
Internet
    ↓
DNS Resolution (stori.danntastico.dev → 107.22.66.194)
    ↓
AWS Security Group (firewall)
    ├─ Port 80  ✅ (HTTP, Let's Encrypt)
    ├─ Port 443 ✅ (HTTPS, public traffic)
    └─ Port 22  ✅ (SSH, admin only)
    ↓
EC2 Instance (Amazon Linux 2023)
    ↓
Nginx (web server + reverse proxy)
    ├─ SSL termination (HTTPS → HTTP internally)
    ├─ Serve frontend: /var/www/stori/ (React static files)
    └─ Proxy backend: localhost:8080/api/ (Go service)
        ↓
    Go Backend (systemd service)
        └─ Port 8080 (internal only)
```

### 🔐 Security Best Practices

1. **Use HTTPS everywhere** - Never send sensitive data over HTTP
2. **Don't expose backend directly** - Use nginx as a gateway
3. **Open only necessary ports** - Close everything else
4. **Use Let's Encrypt** - Free, trusted certificates
5. **Enable auto-renewal** - Don't let certificates expire
6. **Configure CORS properly** - Only allow trusted origins
7. **Use environment variables** - Never hardcode secrets in code

---

## 8. Next Steps for Learning

### 📚 Deep Dive Topics

If you want to learn more about specific areas:

1. **DNS Deep Dive**
   - How DNS resolution works (recursive queries)
   - DNS caching and TTL
   - Different record types (A, AAAA, CNAME, MX, TXT)
   - DNS security (DNSSEC)

2. **SSL/TLS Deep Dive**
   - Public key cryptography (RSA, ECDSA)
   - Certificate chains and trust
   - TLS handshake process
   - Perfect forward secrecy

3. **Nginx Deep Dive**
   - Load balancing strategies
   - Caching strategies
   - Rate limiting
   - Security hardening

4. **Linux System Administration**
   - Systemd in depth
   - Log management (journald, syslog)
   - User permissions and security
   - Process management

5. **AWS Deep Dive**
   - EC2 instance types and pricing
   - Load Balancers (ALB, NLB)
   - Auto Scaling
   - CloudWatch monitoring

### 🛠️ Practice Projects

To reinforce your learning:

1. **Deploy a second subdomain**
   - Create `blog.danntastico.dev`
   - Deploy a simple static site
   - Practice DNS + SSL setup again

2. **Add a database**
   - Install PostgreSQL on EC2
   - Update Go backend to use real database
   - Handle database backups

3. **Set up monitoring**
   - Install Prometheus + Grafana
   - Monitor nginx, Go backend
   - Set up alerts

4. **Implement CI/CD**
   - GitHub Actions workflow
   - Auto-deploy on git push
   - Run tests before deploying

---

## 9. Glossary

**A Record** - DNS record mapping domain name to IPv4 address

**Certificate Authority (CA)** - Organization that issues SSL certificates

**CORS** - Cross-Origin Resource Sharing, browser security policy

**DNS** - Domain Name System, internet's phonebook

**Firewall** - Security system controlling network traffic

**HTTPS** - HTTP Secure, encrypted HTTP using SSL/TLS

**Let's Encrypt** - Free, automated Certificate Authority

**Nginx** - Web server and reverse proxy software

**Port** - Virtual communication endpoint on a server

**Propagation** - Time for DNS changes to spread globally

**Reverse Proxy** - Server forwarding requests to backend servers

**Security Group** - AWS virtual firewall for EC2 instances

**SSL/TLS** - Encryption protocols for secure communication

**Subdomain** - Prefix to main domain (stori.danntastico.dev)

**Systemd** - Linux system and service manager

**TTL** - Time To Live, how long data is cached

---

## 10. Resources for Further Learning

### Official Documentation

- **Nginx**: https://nginx.org/en/docs/
- **Let's Encrypt**: https://letsencrypt.org/docs/
- **Certbot**: https://certbot.eff.org/
- **AWS EC2**: https://docs.aws.amazon.com/ec2/

### Free Courses

- **Linux Journey**: https://linuxjourney.com/
- **Nginx Fundamentals**: YouTube - "Nginx Crash Course"
- **DNS Explained**: Cloudflare Learning Center

### Tools

- **DNS Propagation Checker**: https://dnschecker.org/
- **SSL Checker**: https://www.ssllabs.com/ssltest/
- **HTTP Request Tester**: https://reqbin.com/

---

**📝 Final Note:**

This deployment involved many moving pieces! Don't worry if you don't understand everything perfectly on first read. The key is:

1. **Understand the big picture** (architecture diagram)
2. **Know where to find this guide** (for specific commands)
3. **Practice by doing** (deploy another project)
4. **Ask "why?"** before every command

The more you practice, the more natural it becomes! 🚀

---

**Written for:** Danilo's Stori Challenge  
**Date:** October 28, 2025  
**Status:** Production Deployment Complete ✅

