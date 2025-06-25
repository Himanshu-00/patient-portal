package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
	Role     string // "receptionist" or "doctor"
}

type Patient struct {
	ID            uint      `gorm:"primaryKey"`
	FullName      string    `json:"full_name"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	Gender        string    `json:"gender"`
	ContactNumber string    `json:"contact_number"`
	MedicalNotes  string    `json:"medical_notes"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
