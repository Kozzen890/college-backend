package seeders

import (
	"backend/internal/helpers"
	"backend/internal/models"
	"log"

	"gorm.io/gorm"
)

// SeedUsers menambahkan data user default ke database
func SeedUsers(db *gorm.DB) error {
	// Cek apakah sudah ada user
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Printf("Users already exist (%d users found), skipping seeder", count)
		return nil // Skip seeding jika sudah ada user
	}

	log.Printf("No users found, running user seeder...")

	// Data user default
	users := []struct {
		Username string
		Password string
	}{
		{
			Username: "admin.youth.college",
			Password: "youth-college2025",
		},
	}

	// Insert users dengan password yang di-hash
	for _, userData := range users {
		hashedPassword, err := helpers.HashPassword(userData.Password)
		if err != nil {
			return err
		}

		user := models.User{
			Username: userData.Username,
			Password: hashedPassword,
		}

		if err := db.Create(&user).Error; err != nil {
			return err
		}
		log.Printf("Created user: %s", userData.Username)
	}

	log.Printf("User seeder completed successfully")
	return nil
}
