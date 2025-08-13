# 🎓 Youth College Backend API

REST API untuk sistem pendaftaran peserta Youth College menggunakan Go, Gin Framework, GORM, dan MySQL.

## 🚀 Features

- ✅ **Authentication & Authorization** (JWT)
- ✅ **CRUD Operations** untuk Participants
- ✅ **Pagination & Search**
- ✅ **Input Validation**
- ✅ **CORS Support**
- ✅ **Environment Configuration**
- ✅ **Database Migration & Seeding**
- ✅ **Health Check Endpoint**
- ✅ **Structured Response Format**

## 📋 Prerequisites

- Go 1.21+
- MySQL 8.0+
- Git

## 🛠️ Installation

### 1. Clone Repository

```bash
git clone <repository-url>
cd youth-college-backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Setup

```bash
cp .env.example .env
```

Edit `.env` file:

```env
# Database Configuration
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=college_backend
DB_USER=root
DB_PASSWORD=root

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-for-production-change-this
JWT_EXPIRES_HOURS=24

# Server Configuration
PORT=8001
GIN_MODE=release

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization,X-Requested-With
```

### 4. Database Setup

```bash
# Buat database MySQL
mysql -u root -p
CREATE DATABASE college_backend;
```

### 5. Run Application

```bash
# Development
go run ./cmd/server

# Production Build
go build -o college-backend ./cmd/server
./college-backend
```

## 📚 API Documentation

### Base URL

```
http://localhost:8001
```

### Authentication

#### Login

```http
POST /api/login
Content-Type: application/json

{
  "username": "admin.youth.college",
  "password": "youth-college2025"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin.youth.college"
    }
  }
}
```

#### Logout

```http
POST /api/logout
Authorization: Bearer <token>
```

### Participants

#### Create Participant (Public)

```http
POST /api/participants
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "08123456789",
  "place": "Jakarta"
}
```

#### Get All Participants (Protected)

```http
GET /api/participants?page=1&limit=10&search=john
Authorization: Bearer <token>
```

**Response:**

```json
{
  "success": true,
  "message": "Participants retrieved successfully",
  "data": {
    "participants": [...],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_items": 25,
      "total_pages": 3,
      "has_next": true,
      "has_prev": false
    },
    "search": "john"
  }
}
```

#### Get Participant by ID (Protected)

```http
GET /api/participants/{id}
Authorization: Bearer <token>
```

#### Update Participant (Protected)

```http
PUT /api/participants/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "John Doe Updated",
  "email": "john.updated@example.com",
  "phone": "08987654321",
  "place": "Bandung"
}
```

#### Delete Participant (Protected)

```http
DELETE /api/participants/{id}
Authorization: Bearer <token>
```

### Health Check

```http
GET /healthz
```

### API Info

```http
GET /
```

## 🔒 Security

- **JWT Authentication** dengan secret key dari environment
- **Input Validation** untuk semua endpoints
- **CORS Configuration** untuk frontend integration
- **Password Hashing** menggunakan bcrypt

## 📊 Database Schema

### Users Table

```sql
CREATE TABLE users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Participants Table

```sql
CREATE TABLE participants (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255),
  phone VARCHAR(20) NOT NULL,
  place VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 🚀 Deployment

### Docker (Recommended)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o college-backend ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/college-backend .
COPY --from=builder /app/.env .
CMD ["./college-backend"]
```

### Production Checklist

- [ ] Change `JWT_SECRET` to secure random string
- [ ] Set `GIN_MODE=release`
- [ ] Configure proper CORS origins
- [ ] Use environment variables for all configs
- [ ] Setup SSL/HTTPS
- [ ] Configure rate limiting
- [ ] Setup logging and monitoring
- [ ] Database connection pooling
- [ ] Backup strategy

## 🧪 Testing

```bash
# Run tests
go test ./...

# Test coverage
go test -cover ./...
```

## 🤝 Contributing

1. Fork the project
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📝 License

This project is licensed under the MIT License.

## 👥 Authors

- **Youth College Team** - _Initial work_

## 🙏 Acknowledgments

- Gin Web Framework
- GORM ORM
- JWT-Go Library
- MySQL Database

### Konfigurasi

Variabel lingkungan:

- `PORT` (default `8080`)
- `DATABASE_URL` (prioritas utama)
  - Postgres: `postgres://USER:PASSWORD@HOST:5432/DBNAME?sslmode=require`
  - MySQL (tanpa skema atau dengan skema):
    - `USER:PASSWORD@tcp(HOST:3306)/DBNAME?parseTime=true&loc=Local`
    - `mysql://USER:PASSWORD@HOST:3306/DBNAME?parseTime=true&loc=Local`
- `DB_PATH` (fallback SQLite, default `./data/app.db`) — digunakan jika `DATABASE_URL` tidak di-set.

Contoh `.env` (lokal, gunakan `direnv` atau export manual):

```bash
export PORT=8080
# Pilih salah satu contoh:
# Postgres
# export DATABASE_URL="postgres://USER:PASSWORD@HOST:5432/DBNAME?sslmode=require"
# MySQL
# export DATABASE_URL="USER:PASSWORD@tcp(HOST:3306)/DBNAME?parseTime=true&loc=Local"
```
