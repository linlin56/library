package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

// DBO is the central Database Object, similar to an EntityManager in Java
type DBO struct {
	conn *sql.DB
}

// Config holds the connection settings
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDBO creates a new DBO instance and opens the connection
func NewDBO(cfg Config) (*DBO, error) {
	dsn := fmt.Sprintf(
		// PostgreSQL connection URI
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return &DBO{conn: conn}, nil
}

// NewDBOFromEnv creates a DBO from environment variables
func NewDBOFromEnv() (*DBO, error) {
	cfg := Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "mydb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
	return NewDBO(cfg)
}

// Close closes the database connection
func (d *DBO) Close() error {
	return d.conn.Close()
}

// Basic SQL methods, to be used in the repositories

// QueryRows runs a SELECT query and returns the rows
func (d *DBO) QueryRows(query string, args ...any) (*sql.Rows, error) {
	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("QueryRows error: %w", err)
	}
	return rows, nil
}

// QueryRow runs a SELECT query and returns a single row
func (d *DBO) QueryRow(query string, args ...any) *sql.Row {
	return d.conn.QueryRow(query, args...)
}

// Exec runs an INSERT / UPDATE / DELETE query
// Returns the number of affected rows
func (d *DBO) Exec(query string, args ...any) (int64, error) {
	result, err := d.conn.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("Exec error: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("RowsAffected error: %w", err)
	}
	return rowsAffected, nil
}

// ExecReturning runs an INSERT ... RETURNING id (or another column) query
func (d *DBO) ExecReturning(query string, args ...any) *sql.Row {
	return d.conn.QueryRow(query, args...)
}

// Transaction management:
// BeginTx starts a transaction and returns a TxDBO
func (d *DBO) BeginTx() (*TxDBO, error) {
	tx, err := d.conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("BeginTx error: %w", err)
	}
	return &TxDBO{tx: tx}, nil
}

// WithTx runs a function inside a transaction, with automatic commit or rollback
func (d *DBO) WithTx(fn func(tx *TxDBO) error) error {
	tx, err := d.BeginTx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("transaction rolled back: %w", err)
	}
	return tx.Commit()
}

// SeedIfEmpty runs the SQL script if the database is empty (authors table missing or empty)
func (d *DBO) SeedIfEmpty(scriptPath string) error {
	var tableCount int
	err := d.conn.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'authors'
	`).Scan(&tableCount)
	if err != nil {
		return fmt.Errorf("error checking tables: %w", err)
	}

	if tableCount > 0 {
		var rowCount int
		if err := d.conn.QueryRow("SELECT COUNT(*) FROM authors").Scan(&rowCount); err != nil {
			return fmt.Errorf("error checking data: %w", err)
		}
		if rowCount > 0 {
			log.Println("Database already seeded, nothing to load")
			return nil
		}
	}

	script, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("error reading SQL script %s: %w", scriptPath, err)
	}

	stmts := splitSQLStatements(string(script))
	for i, stmt := range stmts {
		if _, err := d.conn.Exec(stmt); err != nil {
			preview := stmt
			if len(preview) > 80 {
				preview = preview[:80] + "..."
			}
			return fmt.Errorf("error in statement #%d (%s): %w", i+1, preview, err)
		}
	}

	log.Println("Database seeded successfully")
	return nil
}

// splitSQLStatements splits a SQL script into individual statements.
// Handles single-quoted strings (including escaped '') and -- comments.
// This avoids naive splitting on semicolons that may appear in the data.
func splitSQLStatements(script string) []string {
	var stmts []string
	var cur strings.Builder
	inString := false
	n := len(script)

	for i := 0; i < n; i++ {
		ch := script[i]

		if inString {
			cur.WriteByte(ch)
			if ch == '\'' {
				if i+1 < n && script[i+1] == '\'' {
					cur.WriteByte(script[i+1])
					i++
				} else {
					inString = false
				}
			}
			continue
		}

		switch {
		case ch == '\'':
			inString = true
			cur.WriteByte(ch)
		case ch == '-' && i+1 < n && script[i+1] == '-':
			// line comment: skip to the end of the line
			for i < n && script[i] != '\n' {
				i++
			}
		case ch == ';':
			if stmt := strings.TrimSpace(cur.String()); stmt != "" {
				stmts = append(stmts, stmt)
			}
			cur.Reset()
		default:
			cur.WriteByte(ch)
		}
	}

	if stmt := strings.TrimSpace(cur.String()); stmt != "" {
		stmts = append(stmts, stmt)
	}

	return stmts
}

// --- Helper ---

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}