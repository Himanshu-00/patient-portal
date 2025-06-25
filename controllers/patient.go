package controllers

import (
	"net/http"
	"patient-portal/models"
	"patient-portal/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PatientController struct {
	repo repositories.PatientRepository
}

func NewPatientController(repo repositories.PatientRepository) *PatientController {
	return &PatientController{repo: repo}
}

// GetAllPatients godoc
// @Summary Get all patients
// @Description Retrieve all patients
// @Tags Patients
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Patient
// @Failure 500 {object} ErrorResponse
// @Router /api/patients [get]
func (c *PatientController) GetAllPatients(ctx *gin.Context) {
	patients, err := c.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}
	ctx.JSON(http.StatusOK, patients)
}

// GetPatientByID godoc
// @Summary Get patient by ID
// @Description Get a specific patient by their ID
// @Tags Patients
// @Security BearerAuth
// @Produce json
// @Param id path int true "Patient ID"
// @Success 200 {object} models.Patient
// @Failure 404 {object} ErrorResponse
// @Router /api/patients/{id} [get]
func (c *PatientController) GetPatientByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	patient, err := c.repo.FindByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}
	ctx.JSON(http.StatusOK, patient)
}

// CreatePatient godoc
// @Summary Create a new patient
// @Description Create a new patient (Receptionist only)
// @Tags Patients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param patient body models.Patient true "Patient data"
// @Success 201 {object} models.Patient
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/patients [post]
func (c *PatientController) CreatePatient(ctx *gin.Context) {
	var patient models.Patient
	if err := ctx.ShouldBindJSON(&patient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.repo.Create(&patient); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient"})
		return
	}

	ctx.JSON(http.StatusCreated, patient)
}

// UpdatePatient godoc
// @Summary Update patient information
// @Description Update patient data (full update for receptionist, medical notes only for doctor)
// @Tags Patients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param patient body models.Patient true "Patient data"
// @Success 200 {object} models.Patient
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/patients/{id} [patch]
func (c *PatientController) UpdatePatient(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	patient, err := c.repo.FindByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Doctors can only update medical notes
	if ctx.GetString("role") == "doctor" {
		var updateData struct {
			MedicalNotes string `json:"medical_notes"`
		}
		if err := ctx.ShouldBindJSON(&updateData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		patient.MedicalNotes = updateData.MedicalNotes
	} else {
		if err := ctx.ShouldBindJSON(patient); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := c.repo.Update(patient); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
		return
	}

	ctx.JSON(http.StatusOK, patient)
}

// DeletePatient godoc
// @Summary Delete a patient
// @Description Delete a patient (Receptionist only)
// @Tags Patients
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Router /api/patients/{id} [delete]
func (c *PatientController) DeletePatient(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete patient"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// UpdateMedicalNotes godoc
// @Summary Update patient medical notes
// @Description Update medical notes for a patient (Doctor only)
// @Tags Patients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param notes body object{medical_notes=string} true "Medical notes data"
// @Success 200 {object} models.Patient
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/patients/medical-notes/{id} [put]
func (c *PatientController) UpdateMedicalNotes(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	var updateData struct {
		MedicalNotes string `json:"medical_notes"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patient, err := c.repo.FindByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	patient.MedicalNotes = updateData.MedicalNotes
	if err := c.repo.Update(patient); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update medical notes"})
		return
	}

	ctx.JSON(http.StatusOK, patient)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
