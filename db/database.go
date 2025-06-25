package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"patient-portal/config"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {

	// Validate we have a connection string
	if cfg.DB_DSN == "" {
		return nil, fmt.Errorf("DB_DSN is empty")
	}

	log.Printf("🔧 Raw DB_DSN: %s", cfg.DB_DSN)

	// Parse the connection URL
	parsedURL, err := url.Parse(cfg.DB_DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB_DSN: %w", err)
	}

	// Extract components
	host := parsedURL.Hostname()
	port := parsedURL.Port()
	user := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()
	dbname := strings.TrimPrefix(parsedURL.Path, "/")

	// Handle CA certificate path
	caCertPath := cfg.DB_CA_CERT
	if !filepath.IsAbs(caCertPath) {
		wd, _ := os.Getwd()
		caCertPath = filepath.Join(wd, caCertPath)
	}

	// Check if CA certificate file exists
	if _, err := os.Stat(caCertPath); os.IsNotExist(err) {
		log.Printf("⚠️ Warning: CA certificate file not found at %s", caCertPath)
		log.Println("🔧 Attempting connection with sslmode=require instead")

		// Build DSN without CA certificate for Aiven (they have valid SSL certs)
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			host,
			port,
			user,
			password,
			dbname,
		)

		return connectWithDSN(dsn, user, host, port, dbname)
	}

	// Build DSN with SSL parameters and CA certificate
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-ca sslrootcert=%s",
		host,
		port,
		user,
		password,
		dbname,
		caCertPath,
	)

	return connectWithDSN(dsn, user, host, port, dbname)
}

func connectWithDSN(dsn, user, host, port, dbname string) (*gorm.DB, error) {
	log.Printf("🔧 Connecting to: %s@%s:%s/%s", user, host, port, dbname)

	// Create standard database connection
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Configure GORM with custom config for better compatibility
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // Reduce log verbosity
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: false, // Disable prepared statements to avoid parameter issues
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true, // Use simple protocol to avoid complex parameter issues
	}), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM: %w", err)
	}

	log.Println("✅ Database connection established")
	return db, nil
}
