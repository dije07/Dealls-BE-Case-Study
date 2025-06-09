🧾 Payslip Management System (Go + Echo + PostgreSQL)

This project is a scalable backend API for managing payroll, attendance, overtime, and reimbursements — built with **Go**, **Echo**, and **PostgreSQL**. Designed for easy integration with HR or admin dashboards.

---

📦 Features

- 🔐 JWT-based authentication (admin & employee roles)
- 🕒 Attendance submission with weekend blocking
- ⏱️ Overtime submission with hourly limits
- 💸 Reimbursement submission with audit trail
- 📄 Payslip generation with salary breakdown
- 📊 Admin summary report per payroll period
- 🧪 100% unit test coverage (handlers, services, repos)

---

🛠️ Tech Stack

| Component        | Tech                         |
|------------------|------------------------------|
| Language         | Go (Golang)                  |
| Web Framework    | Echo                         |
| Database         | PostgreSQL                   |
| ORM              | GORM                         |
| JWT              | Secure Auth (Middleware)     |
| Testing          | Testify + SQLite in-memory   |

---

🚀 Getting Started

1. Clone this repo

git clone https://github.com/dije07/Dealls-BE-Case-Study.git

cd payslip-system

2. Setup .env

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=postgres

3. Run PostgreSQL (e.g. via Docker)

docker run --name payslip-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=yourpassword \
  -e POSTGRES_DB=payslip \
  -p 5432:5432 \
  -d postgres

4. Run the app

go run main.go

---

🧪 Running Tests

go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

---

📬 API Endpoints (Summary)

| Method | Endpoint                   | Role     | Description                 |
| ------ | -------------------------- | -------- | --------------------------- |
| POST   | `/login`                   | All      | Login and get JWT           |
| GET    | `/me`                      | All      | Current user profile        |
| POST   | `/api/attendance`          | Employee | Submit attendance           |
| GET    | `/api/attendance`          | Employee | View attendance history     |
| POST   | `/api/overtime`            | Employee | Submit overtime             |
| GET    | `/api/overtime`            | Employee | View overtime history       |
| POST   | `/api/reimbursement`       | Employee | Submit reimbursement        |
| GET    | `/api/reimbursement`       | Employee | View reimbursement history  |
| GET    | `/api/payslip/:id`         | Employee | View generated payslip      |
| POST   | `/api/payroll-period`      | Admin    | Create payroll period       |
| POST   | `/api/run-payroll`         | Admin    | Run payroll for that period |
| GET    | `/api/payslip-summary/:id` | Admin    | View summary of payslips    |

---

📂 Folder Structure

.
├── handlers/            # HTTP endpoint logic
├── services/            # Business logic
├── repositories/        # DB queries
├── models/              # GORM models
├── middleware/          # JWT, logging, role checks
├── routes/              # Route registration
├── seeder/              # Seed roles/users
└── main.go              # App entry point
