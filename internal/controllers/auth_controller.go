package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"backend/internal/helpers"
	"backend/internal/models"
)

type AuthController struct {
	DB *gorm.DB
}

// getJWTSecret mendapatkan JWT secret dari environment
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // fallback untuk development
	}
	return []byte(secret)
}

// getJWTExpiresHours mendapatkan JWT expires dari environment
func getJWTExpiresHours() time.Duration {
	hoursStr := os.Getenv("JWT_EXPIRES_HOURS")
	if hoursStr == "" {
		return time.Hour * 24 // default 24 jam
	}

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return time.Hour * 24 // default jika error
	}

	return time.Hour * time.Duration(hours)
}

// NewAuthController membuat instance controller baru
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// LoginRequest struct untuk request login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login endpoint untuk authentication
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseBadRequest(c, err.Error())
		return
	}

	// Cari user berdasarkan username
	var user models.User
	if err := ac.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseUnauthorized(c, "Username tidak ditemukan")
			return
		}
		helpers.ResponseInternalServerError(c, "Database error")
		return
	}

	// Verifikasi password
	if !helpers.CheckPassword(req.Password, user.Password) {
		helpers.ResponseUnauthorized(c, "Password salah")
		return
	}

	// Generate JWT access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(getJWTExpiresHours()).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		helpers.ResponseInternalServerError(c, "Failed to generate token")
		return
	}

	// Generate refresh token
	refreshToken, err := helpers.GenerateRandomToken(32)
	if err != nil {
		helpers.ResponseInternalServerError(c, "Failed to generate refresh token")
		return
	}
	refreshExpires := time.Now().Add(7 * 24 * time.Hour) // 7 hari

	// Simpan refresh token ke DB
	rt := models.RefreshToken{
		UserID:    fmt.Sprintf("%v", user.ID),
		Token:     refreshToken,
		ExpiresAt: refreshExpires,
		Revoked:   false,
	}
	if err := ac.DB.Create(&rt).Error; err != nil {
		helpers.ResponseInternalServerError(c, "Failed to save refresh token")
		return
	}

	// Set refresh token di httpOnly cookie dengan SameSite=Lax
	durasiCookie := int(time.Until(refreshExpires).Seconds())
	if durasiCookie < 0 {
		durasiCookie = 0
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  refreshExpires,
		MaxAge:   durasiCookie,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	data := gin.H{
		"token": tokenString,
	}
	helpers.ResponseSuccess(c, "Login Berhasil", data)
}

// GetProfile returns the current logged-in user's profile based on JWT claims
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	if userID == nil || username == nil {
		helpers.ResponseUnauthorized(c, "User not authenticated")
		return
	}

	// Ambil user dari database (bisa tambahkan email/dll jika perlu)
	var user models.User
	if err := ac.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		helpers.ResponseInternalServerError(c, "User not found")
		return
	}

	data := gin.H{
		"id":       user.ID,
		"username": user.Username,
		// tambahkan field lain jika perlu
	}
	helpers.ResponseSuccess(c, "Profile fetched successfully", data)
}

// Logout endpoint (client-side logout)
func (ac *AuthController) Logout(c *gin.Context) {
	// Ambil access token dari Authorization header (untuk blacklist)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		helpers.ResponseUnauthorized(c, "Authorization header required (Bearer <token>)")
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Ambil refresh token dari cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		helpers.ResponseBadRequest(c, "refresh_token cookie required")
		return
	}

	// Blacklist access token
	hash := helpers.SHA256Hash(tokenString)
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	var exp int64 = time.Now().Add(getJWTExpiresHours()).Unix()
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if v, ok := claims["exp"].(float64); ok {
			exp = int64(v)
		}
	}
	blacklistedToken := models.BlacklistedToken{
		TokenHash: hash,
		ExpiresAt: time.Unix(exp, 0),
	}
	result := ac.DB.Create(&blacklistedToken)
	if result.Error != nil {
		log.Printf("ERROR: Gagal menyimpan token ke blacklist: %v", result.Error)
		helpers.ResponseInternalServerError(c, "Gagal menyimpan token ke blacklist")
		return
	}

	// Revoke refresh token di DB
	if err := ac.DB.Model(&models.RefreshToken{}).
		Where("token = ? AND revoked = ? AND expires_at > ?", refreshToken, false, time.Now()).
		Update("revoked", true).Error; err != nil {
		log.Printf("ERROR: Gagal revoke refresh token: %v", err)
		// Tidak gagal total, lanjutkan response
	}

	// Hapus cookie refresh_token
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	log.Printf("SUCCESS: Token berhasil di-blacklist dan refresh token di-revoke.")
	username, exists := c.Get("username")
	data := gin.H{
		"message": "Logout successful. Token blacklisted & refresh token revoked.",
	}
	if exists {
		data["username"] = username
	}
	helpers.ResponseSuccess(c, "Logout successful", data)
}

// GetUserById mengambil user berdasarkan ID
func (ac *AuthController) GetUserById(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := ac.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.ResponseNotFound(c, "User not found")
			return
		}
		helpers.ResponseInternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, "User retrieved successfully", user)
}

// RefreshToken endpoint untuk mendapatkan access token baru
func (ac *AuthController) RefreshToken(c *gin.Context) {
	// Ambil refresh token dari cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		helpers.ResponseUnauthorized(c, "Refresh token cookie required")
		return
	}

	// Cari refresh token di DB
	var rt models.RefreshToken
	if err := ac.DB.Where("token = ? AND revoked = ? AND expires_at > ?", refreshToken, false, time.Now()).First(&rt).Error; err != nil {
		helpers.ResponseUnauthorized(c, "Invalid or expired refresh token")
		return
	}

	// Cari user
	var user models.User
	if err := ac.DB.Where("id = ?", rt.UserID).First(&user).Error; err != nil {
		helpers.ResponseInternalServerError(c, "User not found")
		return
	}

	// Generate access token baru
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(getJWTExpiresHours()).Unix(),
		"iat":      time.Now().Unix(),
	})
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		helpers.ResponseInternalServerError(c, "Failed to generate token")
		return
	}

	data := gin.H{
		"token": tokenString,
		"user": gin.H{
			"username": user.Username,
		},
	}
	helpers.ResponseSuccess(c, "Token refreshed", data)
}
