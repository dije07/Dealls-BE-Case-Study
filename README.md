# 🧾 Payslip Management System (Go + Echo + PostgreSQL)

This project is a scalable backend API for managing payroll, attendance, overtime, and reimbursements — built with **Go**, **Echo**, and **PostgreSQL**. Designed for seamless integration with HR or admin dashboards.

---

## 📦 Features

- 🔐 JWT-based authentication (admin & employee roles)
- 🕒 Attendance submission with weekend blocking
- ⏱️ Overtime submission with hourly limits
- 💸 Reimbursement submission with audit trail
- 📄 Payslip generation with salary breakdown
- 📊 Admin summary report per payroll period
- 🧪 100% unit test coverage (handlers, services, repositories)
- 🧾 Audit logging for all critical actions
- 🌐 IP address & request ID tracking for traceability
- 🧪 100% unit test coverage (handlers, services, repositories)
- 🧾 Audit logging for all critical actions
- 🌐 IP address & request ID tracking for traceability

---

## 🛠️ Tech Stack

| Component        | Tech                         |
|------------------|------------------------------|
| Language         | Go (Golang)                  |
| Web Framework    | Echo                         |
| Database         | PostgreSQL                   |
| ORM              | GORM                         |
| Auth             | JWT Middleware               |
| Testing          | Testify + Echo Integration   |
| Auth             | JWT Middleware               |
| Testing          | Testify + Echo Integration   |

---

## 🚀 Getting Started

### 1. Clone this repo

<pre>
git clone https://github.com/dije07/Dealls-BE-Case-Study.git
git clone https://github.com/dije07/Dealls-BE-Case-Study.git
cd payslip-system
</pre>

### 2. Setup `.env`

<pre>
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=postgres
</pre>

### 3. Run PostgreSQL (example via Docker)

<pre>
docker run --name payslip-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=yourpassword \
  -e POSTGRES_DB=payslip \
  -p 5432:5432 \
  -d postgres
</pre>

### 4. Run the app

<pre>
go run main.go
</pre>

---

## 🧪 Running Tests

<pre>
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
go test ./tests/integration/...
</pre>

Includes:

- ✅ Unit tests for handlers, services, repositories
- ✅ Integration tests for all endpoints
go test ./tests/integration/...

---

## 📬 API Endpoints (Summary)

| Method | Endpoint                   | Role     | Description                 |
|--------|----------------------------|----------|-----------------------------|
| POST   | `/login`                   | All      | Login and get JWT           |
| GET    | `/me`                      | All      | Get current user profile    |
| GET    | `/me`                      | All      | Get current user profile    |
| POST   | `/api/attendance`          | Employee | Submit attendance           |
| GET    | `/api/attendance`          | Employee | View attendance history     |
| POST   | `/api/overtime`            | Employee | Submit overtime             |
| GET    | `/api/overtime`            | Employee | View overtime history       |
| POST   | `/api/reimbursement`       | Employee | Submit reimbursement        |
| GET    | `/api/reimbursement`       | Employee | View reimbursement history  |
| GET    | `/api/payslip/:id`         | Employee | View payslip                |
| GET    | `/api/payslip/:id`         | Employee | View payslip                |
| POST   | `/api/payroll-period`      | Admin    | Create payroll period       |
| POST   | `/api/run-payroll`         | Admin    | Run payroll                 |
| POST   | `/api/run-payroll`         | Admin    | Run payroll                 |
| GET    | `/api/payslip-summary/:id` | Admin    | View summary of payslips    |

---

## 📂 Folder Structure

<pre>
.
├── handlers/            # HTTP endpoint logic
├── services/            # Business logic
├── repositories/        # DB queries with interface abstraction
├── repositories/        # DB queries with interface abstraction
├── models/              # GORM models
├── middleware/          # JWT, audit log, request ID
├── middleware/          # JWT, audit log, request ID
├── routes/              # Route registration
├── seeder/              # Seed roles/users
├── tests/               # Unit & integration test cases
└── main.go              # Application entry point
</pre>

---

## 🧾 Auditing & Traceability

- All modifying actions (POST/PUT/DELETE) are recorded in the `audit_logs` table
- Each record includes:
  - Action type (e.g., "POST")
  - Entity (e.g., `/api/attendance`)
  - Entity ID (record ID)
  - PerformedBy (user ID)
  - IP address
  - Request ID (traceable across logs)

---

## 📮 Postman Collection

You can test the API using the provided Postman collection:

📁 **File:** `payslip_system.postman_collection.json`  
🔗 Or use the online collection:
[Open in Postman](https://api.postman.com/collections/11930344-397e3b12-b7ed-42a9-b0a3-b67d828673d4?access_key=PMAT-01JXBYRZ6R02HQ69G8BGA8YDBE)