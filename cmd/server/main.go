package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"backend/internal/db"
	"backend/internal/httpapi"
	"backend/internal/models"
	"backend/internal/seeders"
)

func getEnvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func ensureDirectoryExists(dirPath string) error {
	if dirPath == "" {
		return nil
	}
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0o755)
	}
	return nil
}

func mustConnectDatabase(dsnOrPath string) *gorm.DB {
	// Only ensure directory for SQLite path, skip for Postgres URLs
	lower := strings.ToLower(dsnOrPath)
	if !(strings.HasPrefix(lower, "postgres://") || strings.HasPrefix(lower, "postgresql://")) {
		if err := ensureDirectoryExists(filepath.Dir(dsnOrPath)); err != nil {
			log.Fatalf("failed to ensure data directory: %v", err)
		}
	}
	if err := db.Connect(dsnOrPath); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db.Instance()
}

func main() {
	// Load .env file dari parent directory jika ada
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: Error loading .env file from parent: %v", err)
		// fallback: coba load dari current dir
		_ = godotenv.Load()
	}

	port := getEnvOrDefault("PORT", "8001")

	// Cek env MySQL
	dbHost := getEnvOrDefault("DB_HOST", "")
	dbPort := getEnvOrDefault("DB_PORT", "")
	dbName := getEnvOrDefault("DB_NAME", "")
	dbUser := getEnvOrDefault("DB_USER", "")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "")

	var databaseURL string
	if dbHost != "" && dbPort != "" && dbName != "" && dbUser != "" {
		// Format DSN MySQL: username:password@tcp(host:port)/dbname?parseTime=true
		databaseURL = dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"
		log.Printf("Connecting to MySQL: %s", databaseURL)
	} else {
		// Prefer DATABASE_URL (can be Postgres or MySQL); fallback to SQLite path via DB_PATH
		databaseURL = getEnvOrDefault("DATABASE_URL", "")
		if databaseURL == "" {
			databaseURL = getEnvOrDefault("DB_PATH", filepath.Join(".", "data", "app.db"))
		}
		log.Printf("Using database: %s", databaseURL)
	}

	database := mustConnectDatabase(databaseURL)
	log.Printf("Database connected successfully")

	// Auto migrate models
	log.Printf("Running auto migration...")
	if err := database.AutoMigrate(&models.User{}, &models.Participant{}, &models.BlacklistedToken{}, &models.RefreshToken{}); err != nil {
		log.Printf("Migration error: %v", err)
	} else {
		log.Printf("Migration completed successfully")
	}

	// Run seeders
	log.Printf("Running user seeder...")
	if err := seeders.SeedUsers(database); err != nil {
		log.Printf("Warning: Failed to seed users: %v", err)
	} else {
		log.Printf("User seeder completed successfully")
	}

	// Set Gin mode from environment
	ginMode := getEnvOrDefault("GIN_MODE", "debug")
	gin.SetMode(ginMode)

	engine := gin.Default()

	// Setup CORS
	allowOriginsRaw := getEnvOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173,http://10.255.209.77:3000,http://10.255.209.77:5173")
	allowOrigins := strings.Split(allowOriginsRaw, ",")
	for i := range allowOrigins {
		allowOrigins[i] = strings.TrimSpace(allowOrigins[i])
	}
	log.Printf("CORS Allowed Origins: %v", allowOrigins)

	corsConfig := cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     strings.Split(getEnvOrDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
		AllowHeaders:     strings.Split(getEnvOrDefault("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization,X-Requested-With"), ","),
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	engine.Use(cors.New(corsConfig))

	// Setup API routes
	httpapi.SetupRouter(engine, database)

	// Health check endpoint
	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "Youth College Backend API",
		})
	})

	// API Info endpoint
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Youth College Backend API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"auth": gin.H{
					"login":  "POST /api/login",
					"logout": "POST /api/logout (protected)",
				},
				"participants": gin.H{
					"create": "POST /api/participants",
					"list":   "GET /api/participants (protected)",
					"get":    "GET /api/participants/:id (protected)",
					"update": "PUT /api/participants/:id (protected)",
					"delete": "DELETE /api/participants/:id (protected)",
				},
				"health": "GET /healthz",
			},
		})
	})

	log.Printf("🚀 Server starting on port %s in %s mode", port, ginMode)
	log.Printf("📖 API Documentation: http://localhost:%s/", port)
	log.Printf("❤️  Health Check: http://localhost:%s/healthz", port)

	if err := engine.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
