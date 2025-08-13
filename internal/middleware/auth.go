package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"backend/internal/helpers"
	"backend/internal/models"
)

// getJWTSecret mendapatkan JWT secret dari environment
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // fallback untuk development
	}
	return []byte(secret)
}

// AuthMiddleware middleware untuk validasi JWT token (stateless)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization atau cookie
		tokenString := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Coba ambil dari cookie
			cookieToken, err := c.Cookie("admin_token")
			if err == nil && cookieToken != "" {
				tokenString = cookieToken
			}
		}
		if tokenString == "" {
			helpers.ResponseUnauthorized(c, "Token required (Authorization header or cookie)")
			c.Abort()
			return
		}
		// Cek blacklist
		hash := helpers.SHA256Hash(tokenString)
		log.Printf("Checking token in blacklist. Hash: %s", hash)

		db := c.MustGet("db").(*gorm.DB)
		var blacklisted models.BlacklistedToken
		if err := db.Where("token_hash = ? AND expires_at > ?", hash, time.Now()).First(&blacklisted).Error; err == nil {
			log.Printf("Token found in blacklist! Hash: %s, ExpiresAt: %v", blacklisted.TokenHash, blacklisted.ExpiresAt)
			helpers.ResponseUnauthorized(c, "Token is blacklisted. Please login again.")
			c.Abort()
			return
		} else {
			log.Printf("Token not in blacklist. Error: %v", err)
		}

		// Parse dan validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			helpers.ResponseUnauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Set user info ke context untuk digunakan di controller
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
		}

		c.Next()
	}
}
