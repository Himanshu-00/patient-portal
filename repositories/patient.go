package repositories

import (
	"patient-portal/models"

	"gorm.io/gorm"
)

type PatientRepository interface {
	Create(patient *models.Patient) error
	FindByID(id uint) (*models.Patient, error)
	Update(patient *models.Patient) error
	Delete(id uint) error
	FindAll() ([]models.Patient, error)
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) Create(patient *models.Patient) error {
	return r.db.Create(patient).Error
}

func (r *patientRepository) FindByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.First(&patient, id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

func (r *patientRepository) Delete(id uint) error {
	return r.db.Delete(&models.Patient{}, id).Error
}

func (r *patientRepository) FindAll() ([]models.Patient, error) {
	var patients []models.Patient
	if err := r.db.Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}
