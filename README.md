# FME Skills Database Backend

This repository houses the backend codebase for the Federal Ministry of Education's (FME) Skills Database platform. It provides RESTful APIs for managing students, artisans, employers, MDAs (Ministries, Departments, and Agencies), STCs (Skill Training Centers), and various other platform resources.

## 🏗 Architecture & Tech Stack
The backend is built with extensibility and performance in mind using **Go**.

- **Framework**: [Gin](https://github.com/gin-gonic/gin) for blazing-fast HTTP routing.
- **Database**: PostgreSQL (using `pgvector` Docker image).
- **ORM**: [GORM](https://gorm.io/) for database models and interactions.
- **Authentication**: JWT (JSON Web Tokens) with OTP support.
- **Email/Notifications**: Postmark.
- **Architecture Pattern**: Domain-Driven/Modular Architecture. Code related to specific modules (e.g., Course, Student, User) is encapsulated logically under the `internal/` directory.

## 📂 Project Structure

The project follows standard Go layout conventions alongside domain-driven module grouping:

```text
fme_backend/
├── cmd/
│   ├── server/           # Main application entrypoint (main.go)
│   └── migrations/       # Database migration scripts and initializers
├── internal/
│   ├── config/           # Environment variables and DB connection setup
│   ├── middlewares/      # Gin middlewares (Auth, Logging, etc.)
│   ├── models/           # Shared database models
│   ├── utilities/        # Shared helper functions (hashing, email, etc.)
│   ├── user/             # User domain (models, controllers, schemas)
│   ├── student/          # Student domain
│   ├── mdas/             # MDA domain
│   ├── course/           # Courses domain
│   ├── employers/        # Employers domain
│   ├── jobs/             # Job postings domain
│   └── artisans/         # Artisans domain
├── docker-compose.yml    # Docker Compose orchestration
├── dockerfile            # Multi-stage Docker builder
└── go.mod / go.sum       # Go module dependencies
```

## 🗄 Database Models

The database models are represented as Go structs and are managed via GORM. Because the project utilizes a modular design, models are located within their respective domain folders:

- **`internal/user/models.go`**: Core `User` model, handling authentication, roles, and OTPs.
- **`internal/student/models.go`**: Main entity models for students, their demographics, and associated fields.
- **Other modules** (`internal/employers/`, `internal/artisans/`, etc.): Contains structural definitions and schemas extending data specific to those platform actors.

All standard models integrate `gorm.Model` (providing `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`). Migrations are managed centrally within `cmd/migrations/migration.go` to ensure systematic synchronization with the PostgreSQL database.

## 🚀 Getting Started locally

### Prerequisites
- [Go 1.21+](https://golang.org/)
- [Docker](https://www.docker.com/) and Docker Compose (highly recommended for seamless DB setup)

### Environment Variables
Create a `.env` file in the root directory. You can utilize the following values referenced by `docker-compose.yml` for basic local development:
```env
DATABASE_URL=postgres://fme_user:secure_password_123@localhost:5435/fme_db_backend?sslmode=disable
ENV_TYPE=testing
HASH_SECRET=b21c91b2ab9f6218ce94ee072c43298cb28a
FME_PASSWORD=Pass123*
FME_EMAIL=fme.testing@gmail.com
PM_SERVER_TOKEN=c209741a-e135-4f9c-8591-9a7f56347144
PM_ACCT_TOKEN=178aab19-a495-4196-9703-b3dd9d2cdb7f
HOME_MAIL=skillsdbproject@coderina.org
```
*(Notice the Postgres port `5435` in the DATABASE_URL. Update your host/port logic corresponding to how you run the backend).*

### Option A: Running with Docker (Recommended)
You can spin up the entire application (Database + Backend Server) easily with Docker. The backend container will automatically run the migrations and start the server.

1. Build and run the containers:
   ```bash
   docker-compose up --build
   ```
2. The server will be accessible at `http://localhost:80`. (Mapped from the container's internal `8080`).

### Option B: Running Natively
If you prefer running the Go application natively but want to use Docker solely for the database:

1. Target only the database utilizing Docker-compose:
   ```bash
   docker-compose up -d db
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run Database Migrations:
   ```bash
   go run cmd/migrations/migration.go
   ```
4. Start the Application:
   ```bash
   go run cmd/server/main.go
   ```

The application's APIs will now be active on `http://localhost:8080`.
