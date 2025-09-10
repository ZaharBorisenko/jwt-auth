# JWT Auth API

A secure and scalable JWT-based authentication API built with pure Go (without frameworks). Features user registration, login, email verification, password management, and admin functionality.

## 🛠️ Built With

### Backend
- **Go** (Golang) - Pure Go without frameworks
- **net/http** - Standard HTTP package
- **JWT** - JSON Web Tokens for authentication
- **Redis** - Token blacklisting and caching
- **PostgreSQL** - Database storage
- **SQLx** - Database access

### Authentication & Security
- **JWT Tokens** - Stateless authentication
- **BCrypt** - Password hashing
- **Redis Blacklist** - Token invalidation
- **Rate Limiting** - Protection against brute force
- **CORS** - Cross-origin resource sharing

### Deployment & Containerization
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **golang-migrate** - Database migrations

### Documentation & Tools
- **Swagger/OpenAPI** - API documentation
- **Go Validator** - Request validation
- **Godotenv** - Environment configuration
- **UUID** - Unique identifier generation

## 📖 API Documentation

Interactive Swagger documentation is available at:  
http://localhost:8080/swagger/index.html

![alt text]("C:\Users\Захар\Desktop\swag.jpg")
## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)
```bash
# Clone the repository
git clone <your-repo-url>
cd jwt-auth

# Copy environment file
cp .env.example .env
# Edit .env with your preferences

# Start all services (PostgreSQL, Redis, Migrations, and Go API)
docker-compose up -d

# View logs
docker-compose logs -f app
```

### Option 2: Manual Setup
```bash
# Install dependencies
go mod download

# Install Swagger CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
swag init -g main.go

# Set up environment variables
cp .env.example .env
# Edit .env with your database and Redis settings

# Run server
go run main.go
```

## 🐳 Docker Configuration

Your project uses a multi-stage Docker build with:

- **Builder stage** - Compiles the application and includes golang-migrate
- **Final stage** - Lightweight Alpine-based production image
- **Health checks** - For database connectivity
- **Volume mounts** - For persistent storage

### Services:
- **app** - Go API application
- **db** - PostgreSQL database
- **migrate** - Database migrations
- **redis** - Redis caching and sessions

## 🔧 Environment Variables

Create a `.env` file based on `.env.example`:

```env
# Server
SERVER_PORT=8080

# Database
DB_HOST=db
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=password
DB_NAME=jwt_auth

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT
JWT_SECRET=your_super_secret_jwt_key_change_in_production

# For local development without Docker:
# DB_HOST=localhost
# REDIS_HOST=localhost
```

## 📋 API Endpoints

### Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/register` | Register new user | No |
| POST | `/login` | User login | No |
| POST | `/logout` | User logout | Yes |
| POST | `/verify-email` | Verify email | No |
| POST | `/resend-verification` | Resend verification code | No |
| POST | `/change-password` | Change password | Yes |

### Users
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/profile/{id}` | Get user profile | Yes |
| PUT | `/profile/{id}` | Update user profile | Yes |

### Admin
| Method | Endpoint | Description | Auth Required | Admin Only |
|--------|----------|-------------|---------------|------------|
| GET | `/admin/users` | Get all users | Yes | Yes |
| DELETE | `/admin/user/{id}` | Delete user | Yes | Yes |
| GET | `/admin/blacklist` | Get blacklisted tokens | Yes | Yes |

## 🐳 Docker Commands

### Basic Commands
```bash
# Start all services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f
docker-compose logs -f app
docker-compose logs -f db
docker-compose logs -f redis

# Restart specific service
docker-compose restart app

# Rebuild and start
docker-compose up -d --build

# Remove volumes (warning: deletes data)
docker-compose down -v
```

### Development Commands
```bash
# Run database migrations manually
docker-compose run migrate

# Access PostgreSQL database
docker-compose exec db psql -U admin -d jwt_auth

# Access Redis CLI
docker-compose exec redis redis-cli

# Run tests inside container
docker-compose exec app go test ./...
```

### Service Ports:
- **API**: http://localhost:8080
- **PostgreSQL**: localhost:5433 → container:5432
- **Redis**: localhost:6380 → container:6379

## 📊 Example Requests

### Register User
```http
POST /register
Content-Type: application/json

{
  "userName": "johndoe",
  "firstName": "John",
  "lastName": "Doe",
  "email": "john.doe@example.com",
  "password": "password123"
}
```

### Login
```http
POST /login
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

### Change Password
```http
POST /change-password
Authorization: Bearer <your_token>
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "oldPassword": "oldPassword123",
  "newPassword": "newPassword456"
}
```

## 🏗️ Project Structure

```
jwt-auth/
├── handlers/          # HTTP handlers
├── models/           # Data structures
├── storage/          # Database and Redis
│   ├── repositories/ # Data access layer
│   └── services/     # Business logic
├── middleware/       # Auth, rate limiting, etc.
├── helpers/          # Utilities
├── route/           # Router setup
├── database/
│   └── migration/    # Database migrations
├── docs/            # Swagger documentation
├── Dockerfile       # Multi-stage Docker build
├── docker-compose.yml # Complete environment
├── .env.example     # Environment template
├── go.mod          # Go dependencies
├── main.go         # Application entry point
└── README.md       # This file
```

## 🔧 Development

### With Docker Compose (Recommended)
```bash
# Start development environment
docker-compose up -d

# Monitor application logs
docker-compose logs -f app

# Run tests
docker-compose exec app go test ./...

# Apply new migrations
docker-compose run migrate
```

### Without Docker
```bash
# Ensure PostgreSQL and Redis are running locally
go run main.go
```

## 🚀 Production Deployment

For production, consider:

1. **Use environment-specific .env files**
2. **Enable PostgreSQL SSL connections**
3. **Use Redis with password authentication**
4. **Set up reverse proxy (nginx/caddy)**
5. **Configure proper monitoring and logging**
6. **Use Docker secrets for sensitive data**

## 📝 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📞 Support

If you have any questions or issues, please open an issue on GitHub.

---

**Security Note**: Always use strong secrets in production and never commit sensitive data to version control!