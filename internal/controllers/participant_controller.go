package controllers

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/internal/forms"
	"backend/internal/helpers"
	"backend/internal/models"
)

type ParticipantController struct {
	DB *gorm.DB
}

// NewParticipantController membuat instance controller baru
func NewParticipantController(db *gorm.DB) *ParticipantController {
	return &ParticipantController{DB: db}
}

// CreateParticipant membuat participant baru
func (pc *ParticipantController) CreateParticipant(c *gin.Context) {
	var form forms.ParticipantForm
	if err := c.ShouldBindJSON(&form); err != nil {
		helpers.ResponseBadRequest(c, err.Error())
		return
	}
	// Konversi string ke time.Time
	var birthDate time.Time
	if form.BirthDate != "" {
		var err error
		birthDate, err = time.Parse("2006-01-02", form.BirthDate)
		if err != nil {
			helpers.ResponseBadRequest(c, "Format birth_date harus YYYY-MM-DD")
			return
		}
	}

	participant := models.Participant{
		Name:      form.Name,
		Place:     form.Place,
		BirthDate: birthDate,
		Kampus:    form.Kampus,
		Jurusan:   form.Jurusan,
		Angkatan:  form.Angkatan,
		Phone:     form.Phone,
	}

	if err := pc.DB.Create(&participant).Error; err != nil {
		helpers.ResponseInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseCreated(c, "Participant created successfully", participant)
}

// GetAllParticipants mengambil semua participants dengan pagination, search, dan sorting
func (pc *ParticipantController) GetAllParticipants(c *gin.Context) {
	// Query parameters untuk pagination, search, dan sorting
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Convert to int
	pageInt := 1
	limitInt := 10

	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageInt = p
	}

	if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
		limitInt = l
	}

	// Validate sort parameters
	validSortFields := map[string]bool{
		"id":         true,
		"name":       true,
		"place":      true,
		"birth_date": true,
		"kampus":     true,
		"jurusan":    true,
		"angkatan":   true,
		"phone":      true,
		"created_at": true,
		"updated_at": true,
	}

	validSortOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	if !validSortOrders[sortOrder] {
		sortOrder = "desc"
	}

	offset := (pageInt - 1) * limitInt

	// Query dengan search, sorting, dan pagination
	query := pc.DB.Model(&models.Participant{})

	if search != "" {
		query = query.Where("name LIKE ?",
			"%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated data dengan sorting
	var participants []models.Participant
	orderClause := sortBy + " " + sortOrder
	if err := query.Order(orderClause).Offset(offset).Limit(limitInt).Find(&participants).Error; err != nil {
		helpers.ResponseInternalServerError(c, err.Error())
		return
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(total) / float64(limitInt)))
	hasNext := pageInt < totalPages
	hasPrev := pageInt > 1

	data := gin.H{
		"participants": participants,
		"pagination": gin.H{
			"current_page": pageInt,
			"per_page":     limitInt,
			"total_items":  total,
			"total_pages":  totalPages,
			"has_next":     hasNext,
			"has_prev":     hasPrev,
		},
		"filters": gin.H{
			"search":     search,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	}

	helpers.ResponseSuccess(c, "Participants retrieved successfully", data)
}

// GetParticipant mengambil data participant berdasarkan ID
func (pc *ParticipantController) GetParticipant(c *gin.Context) {
	id := c.Param("id")
	var participant models.Participant

	if err := pc.DB.Where("id = ?", id).First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseNotFound(c, "Participant not found")
			return
		}
		helpers.ResponseInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, "Participant retrieved successfully", participant)
}

// UpdateParticipant mengupdate data participant
func (pc *ParticipantController) UpdateParticipant(c *gin.Context) {
	id := c.Param("id")
	var participant models.Participant

	// Cek apakah participant ada
	if err := pc.DB.Where("id = ?", id).First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseNotFound(c, "Participant not found")
			return
		}
		helpers.ResponseInternalServerError(c, err.Error())
		return
	}

	// Bind data update pakai form agar birth_date bisa string
	var form forms.ParticipantForm
	if err := c.ShouldBindJSON(&form); err != nil {
		helpers.ResponseBadRequest(c, err.Error())
		return
	}

	// Konversi birth_date string ke time.Time
	var birthDate time.Time
	if form.BirthDate != "" {
		var err error
		birthDate, err = time.Parse("2006-01-02", form.BirthDate)
		if err != nil {
			helpers.ResponseBadRequest(c, "Format birth_date harus YYYY-MM-DD")
			return
		}
	}

	// Update data
	participant.Name = form.Name
	participant.Place = form.Place
	participant.BirthDate = birthDate
	participant.Kampus = form.Kampus
	participant.Jurusan = form.Jurusan
	participant.Angkatan = form.Angkatan
	participant.Phone = form.Phone

	if err := pc.DB.Save(&participant).Error; err != nil {
		helpers.ResponseInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, "Participant updated successfully", participant)
}

// DeleteParticipant menghapus participant
func (pc *ParticipantController) DeleteParticipant(c *gin.Context) {
	id := c.Param("id")
	var participant models.Participant

	// Cek apakah participant ada
	if err := pc.DB.Where("id = ?", id).First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseNotFound(c, "Participant not found")
			return
		}
		helpers.ResponseInternalServerError(c, err.Error())
		return
	}

	// Hapus participant
	if err := pc.DB.Delete(&participant).Error; err != nil {
		helpers.ResponseInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, "Participant deleted successfully", nil)
}

// CountParticipant returns the total number of participants
func (pc *ParticipantController) CountParticipant(c *gin.Context) {
	var count int64
	if err := pc.DB.Model(&models.Participant{}).Count(&count).Error; err != nil {
		helpers.ResponseInternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, "Total participants counted successfully", gin.H{"total": count})
}
