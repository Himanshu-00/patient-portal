package main

import (
	"log"
	"patient-portal/config"
	"patient-portal/controllers"
	"patient-portal/db"
	_ "patient-portal/docs"
	"patient-portal/models"
	"patient-portal/repositories"
	"patient-portal/routes"
	"patient-portal/utils"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title Hospital Management System API
// @version 1.0
// @description API for hospital management system with receptionist and doctor portals
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	} else {
		log.Println("✅ Loaded .env file")
	}
	// Load configuration
	cfg := config.LoadConfig()
	// Verify we have a DB_DSN
	if cfg.DB_DSN == "" {
		log.Fatal("🚨 DB_DSN is not set. Please add it to your .env file or environment variables")
	}
	log.Printf("🔧 Loaded configuration:")
	log.Printf("  DB_DSN: %s", maskPassword(cfg.DB_DSN))
	log.Printf("  DB_CA_CERT: %s", cfg.DB_CA_CERT)

	// Connect to database
	dbConn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := dbConn.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		log.Println("🔧 Configured database connection pool")
	}

	// Set GORM logger to reduce verbosity and avoid parameter mismatch issues
	dbConn.Logger = logger.Default.LogMode(logger.Error)

	// Auto migrate models with better error handling
	log.Println("🔄 Starting database migration...")

	// Try to migrate User model first
	if err := dbConn.AutoMigrate(&models.User{}); err != nil {
		log.Printf("⚠️ Warning: Failed to migrate User model: %v", err)
		log.Println("🔧 Attempting to create User table manually...")
		if err := createUserTableManually(dbConn); err != nil {
			log.Printf("❌ Failed to create User table manually: %v", err)
		}
	}

	// Migrate Patient model
	if err := dbConn.AutoMigrate(&models.Patient{}); err != nil {
		log.Printf("⚠️ Warning: Failed to migrate Patient model: %v", err)
		log.Println("🔧 Attempting to create Patient table manually...")
		if err := createPatientTableManually(dbConn); err != nil {
			log.Printf("❌ Failed to create Patient table manually: %v", err)
		}
	}

	log.Println("✅ Database migration completed")

	// Create sample data
	createSampleData(dbConn)

	// Initialize repositories
	patientRepo := repositories.NewPatientRepository(dbConn)

	// Initialize controllers
	authCtrl := controllers.NewAuthController(dbConn)
	patientCtrl := controllers.NewPatientController(patientRepo)

	// Setup routes
	router := routes.SetupRouter(authCtrl, patientCtrl)

	// Start server
	log.Println("🚀 Server starting on :8080")
	log.Println("📖 Swagger documentation available at: http://localhost:8080/swagger/index.html")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// createUserTableManually creates the users table manually if AutoMigrate fails
func createUserTableManually(db *gorm.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL
	);`

	return db.Exec(sql).Error
}

// createPatientTableManually creates the patients table manually if AutoMigrate fails
func createPatientTableManually(db *gorm.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS patients (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		date_of_birth TIMESTAMP NOT NULL,
		gender VARCHAR(50),
		contact_number VARCHAR(50),
		medical_notes TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	return db.Exec(sql).Error
}

// maskPassword hides password in connection string
func maskPassword(dsn string) string {
	if strings.Contains(dsn, "password=") {
		return regexp.MustCompile(`password=[^ ]+`).ReplaceAllString(dsn, "password=***")
	}
	if strings.Contains(dsn, "@") {
		parts := strings.SplitN(dsn, "@", 2)
		if len(parts) > 0 {
			auth := strings.SplitN(parts[0], "://", 2)
			if len(auth) > 1 {
				creds := strings.SplitN(auth[1], ":", 2)
				if len(creds) > 1 {
					return auth[0] + "://" + creds[0] + ":***@" + parts[1]
				}
			}
		}
	}
	return dsn
}

func createSampleData(db *gorm.DB) {
	log.Println("🔧 Creating sample data...")

	// Create sample users
	users := []models.User{
		{
			Username: "reception@hospital.com",
			Password: utils.MustHash("Pass123!"),
			Role:     "receptionist",
		},
		{
			Username: "doctor@hospital.com",
			Password: utils.MustHash("DoctorPass1!"),
			Role:     "doctor",
		},
	}

	for _, user := range users {
		var count int64
		db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
		if count == 0 {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("⚠️ Warning: Failed to create user %s: %v", user.Username, err)
			} else {
				log.Printf("✅ Created user: %s", user.Username)
			}
		}
	}

	// Create sample patient
	var patientCount int64
	db.Model(&models.Patient{}).Count(&patientCount)
	if patientCount == 0 {
		samplePatient := &models.Patient{
			FullName:      "John Doe",
			DateOfBirth:   time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC),
			Gender:        "Male",
			ContactNumber: "+1 555-1234",
			MedicalNotes:  "Annual checkup scheduled",
		}
		if err := db.Create(samplePatient).Error; err != nil {
			log.Printf("⚠️ Warning: Failed to create sample patient: %v", err)
		} else {
			log.Println("✅ Created sample patient: John Doe")
		}
	}

	log.Println("✅ Sample data creation completed")
}
