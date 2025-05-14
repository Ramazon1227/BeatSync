# BeatSync

BeatSync is a Go backend service for saving, retrieving, and analyzing user heart rate data, with a focus on Heart Rate Variability (HRV) analytics. The service provides RESTful APIs for user management and heart data operations, leverages InfluxDB for time-series storage, and is fully containerized for easy deployment.

---

## Features
- **User Management:** Register, login, logout, update profile, and delete users with JWT-based authentication.
- **Heart Rate Data:** Store and retrieve heart rate and PPG sensor data per user/device.
- **HRV Analysis:** Analyze heart rate variability with metrics like SDNN, RMSSD, LF/HF ratio, and more.
- **API Documentation:** Interactive Swagger docs auto-generated from code annotations.
- **Dockerized:** Runs with Docker Compose, including InfluxDB as the backing store.

---

## Architecture
- **Language:** Go (Gin framework)
- **Database:** InfluxDB (time-series)
- **API:** RESTful, documented with Swagger
- **Structure:**
  - `api/` – API entrypoint, handlers, middleware
  - `models/` – Data models for users, HRV, auth
  - `storage/` – InfluxDB integration and storage logic
  - `config/` – Configuration and environment management
  - `cmd/` – Application entrypoint (`main.go`)
  - `compose.yaml` – Docker Compose for local/dev

---

## Getting Started

### Prerequisites
- [Go](https://golang.org/) 1.20+
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Setup
1. **Clone the repo:**
   ```bash
   git clone https://github.com/Ramazon1227/BeatSync.git
   cd BeatSync
   ```
2. **Configure environment:**
   - Copy `.env.example` to `.env` and fill in required values (InfluxDB, JWT secret, etc).
3. **Run with Docker Compose:**
   ```bash
   docker-compose -f compose.yaml up --build
   ```
   The API will be available at `http://localhost:8080` by default.

### Local Development
- **Build binary:** `make build`
- **Run locally:** `make run`
- **Generate Swagger docs:** `make swag-init`

---

## API Documentation
- Interactive docs available at `/swagger/index.html` when the service is running.
- Example endpoints:
  - `POST /v1/auth/register` – Register user
  - `POST /v1/auth/login` – Login and receive JWT
  - `GET /v1/profile/{user_id}` – Get user profile
  - `POST /v1/hrv/analyze` – Analyze HRV data

---

## Environment Variables
See `.env.example` for all variables. Key settings:
- `INFLUX_URL`, `INFLUX_TOKEN`, `INFLUX_ORG`, `INFLUX_BUCKET`
- `SECRET_KEY` (JWT)
- `SERVICE_HOST`, `HTTP_PORT`, etc.

---

## Development & Contribution
- Use the Makefile for common tasks (build, run, Docker image management).
- PRs and issues welcome!

---

## License
Apache 2.0

## Contact
- API Support: support@swagger.io
- Project by [Ramazon1227](https://github.com/Ramazon1227)
