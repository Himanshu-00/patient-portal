package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"patient-portal/controllers"
	"patient-portal/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPatientRepository for testing
type MockPatientRepository struct {
	mock.Mock
}

func (m *MockPatientRepository) Create(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) FindByID(id uint) (*models.Patient, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) Update(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPatientRepository) FindAll() ([]models.Patient, error) {
	args := m.Called()
	return args.Get(0).([]models.Patient), args.Error(1)
}

func TestCreatePatient_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockPatientRepository)
	patientCtrl := controllers.NewPatientController(mockRepo)

	router := gin.Default()
	router.POST("/patients", patientCtrl.CreatePatient)

	dob, _ := time.Parse(time.RFC3339, "1990-01-01T00:00:00Z")
	newPatient := models.Patient{
		FullName:      "John Doe",
		DateOfBirth:   dob,
		Gender:        "Male",
		ContactNumber: "123-456-7890",
	}

	// Mock expectations
	mockRepo.On("Create", &newPatient).Return(nil)

	// Test request
	reqBody, _ := json.Marshal(newPatient)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Verify
	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePatient_DoctorUpdateNotes(t *testing.T) {
	// Setup
	mockRepo := new(MockPatientRepository)
	patientCtrl := controllers.NewPatientController(mockRepo)

	router := gin.Default()
	router.PATCH("/patients/:id", patientCtrl.UpdatePatient)

	// Existing patient
	existingPatient := &models.Patient{
		ID:            1,
		FullName:      "Jane Smith",
		DateOfBirth:   time.Now(),
		Gender:        "Female",
		ContactNumber: "555-1234",
		MedicalNotes:  "Initial notes",
	}

	// Mock expectations
	mockRepo.On("FindByID", uint(1)).Return(existingPatient, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Patient")).Run(func(args mock.Arguments) {
		patient := args.Get(0).(*models.Patient)
		assert.Equal(t, "Updated medical notes", patient.MedicalNotes)
		assert.Equal(t, "Jane Smith", patient.FullName) // Should not change
	}).Return(nil)

	// Test request
	updateData := map[string]interface{}{
		"medical_notes": "Updated medical notes",
		"full_name":     "Should Not Change",
	}
	reqBody, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/patients/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Simulate doctor role in context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Set("role", "doctor")

	patientCtrl.UpdatePatient(c)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePatient_ReceptionistFullUpdate(t *testing.T) {
	// Setup
	mockRepo := new(MockPatientRepository)
	patientCtrl := controllers.NewPatientController(mockRepo)

	// Existing patient
	existingPatient := &models.Patient{
		ID:            2,
		FullName:      "Original Name",
		DateOfBirth:   time.Now(),
		Gender:        "Male",
		ContactNumber: "555-0000",
	}

	// Mock expectations
	mockRepo.On("FindByID", uint(2)).Return(existingPatient, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Patient")).Run(func(args mock.Arguments) {
		patient := args.Get(0).(*models.Patient)
		assert.Equal(t, "Updated Name", patient.FullName)
		assert.Equal(t, "555-1111", patient.ContactNumber)
	}).Return(nil)

	// Test request
	updateData := map[string]interface{}{
		"full_name":      "Updated Name",
		"contact_number": "555-1111",
	}
	reqBody, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/patients/2", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Simulate receptionist role in context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "2"}}
	c.Set("role", "receptionist")

	patientCtrl.UpdatePatient(c)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}
