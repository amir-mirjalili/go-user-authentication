package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DSNBuilder interface {
	BuildDSN() string
}

var DSNRegistry = make(map[string]DSNBuilder)

func RegisterDSNBuilder(dsn string, builder DSNBuilder) {
	DSNRegistry[dsn] = builder
}

func GetDSNBuilder(dsn string) (DSNBuilder, error) {
	builder, exists := DSNRegistry[dsn]
	if !exists {
		return nil, fmt.Errorf("db driver %s is not registered", dsn)
	}
	return builder, nil
}

type PostgresDSNBuilder struct{}

func (m *PostgresDSNBuilder) BuildDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		"5432",
		os.Getenv("DB_SSL_MODE"),
	)
}

func init() {
	RegisterDSNBuilder("postgres", &PostgresDSNBuilder{})
}

type Database struct {
	DB      *sql.DB
	Dialect string
}

func Connect() (*Database, error) {
	driver := "postgres"

	dsnB, err := GetDSNBuilder(driver)
	if err != nil {
		return nil, err
	}

	dsn := dsnB.BuildDSN()
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = migrate(conn)
	if err != nil {
		return nil, err
	}

	fmt.Println("Database connection established successfully")
	return &Database{DB: conn, Dialect: driver}, nil
}

func Close(conn *sql.DB) error {
	return conn.Close()
}

func migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			phone_number VARCHAR(20) UNIQUE NOT NULL,
			registered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS otps (
			phone_number VARCHAR(20) PRIMARY KEY,
			code VARCHAR(6) NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS otp_requests (
			id SERIAL PRIMARY KEY,
			phone_number VARCHAR(20) NOT NULL,
			requested_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_otp_requests_phone_time ON otp_requests(phone_number, requested_at)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
