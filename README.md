# 🚀 PotaFlow  
### A lightweight workflow automation engine — Connect triggers and actions to automate anything.

PotaFlow is a full-stack workflow automation platform inspired by tools like Zapier, n8n, and Make.  
Built primarily in **Golang**, it showcases backend engineering, worker pools, scheduling, API integrations, and a clean minimal UI.

This project was created as a capstone for a full-stack / backend engineering curriculum, with the goal of demonstrating real-world software engineering skills and providing a foundation for freelance automation work.

---

## ✨ Features

### 🟩 Workflow Automation
Create workflows composed of:

**Trigger → Actions…**

Examples:
- Webhook → Slack message  
- Cron schedule → Email summary  
- File upload → Process → Store → Notify  
- HTTP request → Google Sheets append  

### 🟦 Triggers
- **Webhook Trigger** — fire workflows from external systems  
- **Cron Trigger** — run workflows on schedules (hourly, daily, etc.)  
- *(More coming soon…)*

### 🟨 Actions
- **Slack** — send messages to channels  
- **Email** — via SendGrid or SMTP  
- **HTTP Action** — send POST/GET requests  
- **Google Sheets** — append rows  
- **Custom Logic** — run your own handlers  
- *(Extensible by design)*

### 🧵 Concurrency & Worker Pool
PotaFlow uses a custom Go worker pool to execute actions concurrently:
- Configurable worker count  
- Automatic retry with backoff  
- Run logging  
- Dead-letter queue (optional)

### 🔐 Authentication
- JWT access tokens  
- Refresh token rotation  
- Argon2 password hashing  
- Role-based route protection

### 📊 Logs & Monitoring
- Workflow run history  
- Step-by-step action logs  
- Error reporting and retry statuses  

### 🖥️ UI (Frontend)
A small, clean interface for:
- Managing workflows  
- Viewing logs  
- Editing triggers/actions  
- Viewing real-time run results  

---

## 🏗️ Architecture Overview

PotaFlow is split into two services:

### **API Service (`cmd/api/`)**
Responsible for:
- REST endpoints  
- Trigger registration  
- Workflow CRUD  
- Authentication  
- Webhook handling  
- UI static file hosting  

### **Worker Service (`cmd/worker/`)**
Responsible for:
- Executing queued workflow runs  
- Processing triggers  
- Running action steps  
- Retry + backoff logic  
- Logging results  

### **Tech Stack**
- **Golang** (primary language)
- **PostgreSQL** (database)
- **sqlc** (compile-time SQL queries)
- **Chi** (HTTP router)
- **Viper** (configuration)
- **Zerolog** (structured logging)
- **Docker / Docker Compose**
- **React or HTMX** (minimal frontend)

---

## 📂 Project Structure

```
mini-zapier/
├── cmd/
│   ├── api/
│   └── worker/
├── internal/
│   ├── auth/
│   ├── workflows/
│   ├── workerpool/
│   ├── database/
│   ├── http/
│   ├── config/
│   └── integrations/
├── ui/
├── migrations/
├── deployments/
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites
- Go 1.22+
- Docker + Docker Compose
- PostgreSQL (or run via Compose)
- Node.js (if building frontend manually)

---

### 1. Clone the repo
```bash
git clone https://github.com/<your-username>/pota-flow.git
cd pota-flow
```

### 2. Copy environment file
```bash
cp .env.example .env
```

Fill in:
- Database URL  
- JWT secret  
- API tokens (Slack, SendGrid, Google, etc.)

---

### 3. Start database (Docker)
```bash
docker compose up -d db
```

### 4. Run migrations
```bash
make migrate
```

### 5. Run API server
```bash
make api
```

### 6. Run worker service
```bash
make worker
```

Or run the entire stack with:

```bash
docker compose up --build
```

---

## 🧪 Running Tests

```bash
make test
```

---

## 📈 Roadmap

- [ ] Webhook secret verification  
- [ ] Visual workflow builder  
- [ ] OAuth integrations (Google/Slack)  
- [ ] Workflow templates  
- [ ] Prometheus metrics dashboard  
- [ ] Redis queue option (Asynq)  
- [ ] Plugin system for custom actions  
- [ ] Multi-tenant organizations  

---

## 🤝 Contributing
Contributions are welcome — open an issue or submit a PR.

---

## 📜 License
MIT License — free for personal and commercial use.

---

## 👤 Author

**Cory Gleason**  
Creator of PotaFlow  
Backend Engineer | Automation Enthusiast  

---

## ⭐ Support

If this project inspires or helps you, please star the repo.  
It means a lot and increases visibility. ⭐
