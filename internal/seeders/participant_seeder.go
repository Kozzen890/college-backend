package seeders

import (
	"backend/internal/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func SeedParticipants(db *gorm.DB) error {

	participants := []models.Participant{}
	for i := 1; i <= 15; i++ {
		participants = append(participants, models.Participant{
			Name:      "Participant " + strconv.Itoa(i),
			Place:     "Semarang",
			BirthDate: time.Date(2000, time.Month((i%12)+1), (i%28)+1, 0, 0, 0, 0, time.UTC),
			Kampus:    "Kampus " + strconv.Itoa(i),
			Jurusan:   "Jurusan " + strconv.Itoa(i),
			Angkatan:  strconv.Itoa(2020 + (i % 5)),
			Phone:     "0812345678" + strconv.Itoa(i),
		})
	}
	return db.Create(&participants).Error
}
