package forms

// ParticipantForm untuk validasi input participant
// Tambahkan tag binding sesuai kebutuhan validasi

type ParticipantForm struct {
	Name      string `json:"name" binding:"required"`
	Place     string `json:"place" binding:"required,min=2,max=255"`
	BirthDate string `json:"birth_date" binding:"required"`
	Kampus    string `json:"kampus" binding:"required"`
	Jurusan   string `json:"jurusan" binding:"required"`
	Angkatan  string `json:"angkatan" binding:"required"`
	Phone     string `json:"phone" binding:"required,min=8,max=20"`
}
