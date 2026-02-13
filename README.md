# Maulana Laundry - Backend API

This is the backend service for Maulana Laundry, powered by **Go Fiber** and **GORM**.

## 🏗 Architecture
We follow a modular approach inspired by Clean Architecture:
- `cmd/`: Application entry point.
- `handlers/`: Request/Response handling logic.
- `models/`: Database entities and schemas.
- `routes/`: API endpoint definitions.
- `middleware/`: Auth (JWT), Logger, and Security.
- `database/`: Connection setup and migrations.

## 🛠 Tech Stack
- **Framework**: [Fiber v2](https://gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Security**: JWT & Bcrypt

## ⚙️ Configuration
Environment variables are managed via `.env`. See `.env.example` for required fields.

| Variable | Description | Default |
| :--- | :--- | :--- |
| `SERVER_PORT` | Port for the API server | `8080` |
| `DB_HOST` | PostgreSQL Host | `localhost` |
| `JWT_SECRET` | Secret key for JWT signing | `Required` |

## 🚀 Running Locally

1. Install dependencies:
   ```bash
   go mod download
   ```
2. Run migrations and start server:
   ```bash
   go run cmd/main.go
   ```

## 🔐 Auth Flow
1. **Login**: `POST /api/auth/login` returns a JWT.
2. **Authorization**: Include `Authorization: Bearer <token>` in headers for protected routes.
