package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

// DBO est le Database Object central, similaire à un EntityManager en Java
type DBO struct {
	conn *sql.DB
}

// Config contient les paramètres de connexion
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDBO crée une nouvelle instance du DBO et ouvre la connexion
func NewDBO(cfg Config) (*DBO, error) {
	dsn := fmt.Sprintf(
		// URI de connexion PostgreSQL
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erreur ouverture connexion: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("erreur ping base de données: %w", err)
	}

	log.Println("Connexion à PostgreSQL établie")
	return &DBO{conn: conn}, nil
}

// NewDBOFromEnv crée un DBO depuis les variables d'environnement
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

// Close ferme la connexion à la base de données
func (d *DBO) Close() error {
	return d.conn.Close()
}

// Méthodes SQL de base, à utiliser dans les repositories

// QueryRows exécute une requête SELECT et retourne les lignes
func (d *DBO) QueryRows(query string, args ...any) (*sql.Rows, error) {
	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("QueryRows error: %w", err)
	}
	return rows, nil
}

// QueryRow exécute une requête SELECT et retourne une seule ligne
func (d *DBO) QueryRow(query string, args ...any) *sql.Row {
	return d.conn.QueryRow(query, args...)
}

// Exec exécute une requête INSERT / UPDATE / DELETE
// Retourne le nombre de lignes affectées
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

// ExecReturning exécute une requête INSERT ... RETURNING id (ou autre colonne)
func (d *DBO) ExecReturning(query string, args ...any) *sql.Row {
	return d.conn.QueryRow(query, args...)
}

// Gestion des transactions :
// BeginTx démarre une transaction et retourne un TxDBO
func (d *DBO) BeginTx() (*TxDBO, error) {
	tx, err := d.conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("BeginTx error: %w", err)
	}
	return &TxDBO{tx: tx}, nil
}

// WithTx exécute une fonction dans une transaction, commit ou rollback automatique
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
		return fmt.Errorf("transaction annulée: %w", err)
	}
	return tx.Commit()
}

// SeedIfEmpty exécute le script SQL si la base est vide (table auteurs absente ou vide)
func (d *DBO) SeedIfEmpty(scriptPath string) error {
	var tableCount int
	err := d.conn.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'auteurs'
	`).Scan(&tableCount)
	if err != nil {
		return fmt.Errorf("erreur vérification tables: %w", err)
	}

	if tableCount > 0 {
		var rowCount int
		if err := d.conn.QueryRow("SELECT COUNT(*) FROM auteurs").Scan(&rowCount); err != nil {
			return fmt.Errorf("erreur vérification données: %w", err)
		}
		if rowCount > 0 {
			log.Println("Base de données déjà peuplée, aucun chargement nécessaire")
			return nil
		}
	}

	script, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("erreur lecture script SQL %s: %w", scriptPath, err)
	}

	stmts := splitSQLStatements(string(script))
	for i, stmt := range stmts {
		if _, err := d.conn.Exec(stmt); err != nil {
			preview := stmt
			if len(preview) > 80 {
				preview = preview[:80] + "..."
			}
			return fmt.Errorf("erreur instruction #%d (%s): %w", i+1, preview, err)
		}
	}

	log.Println("Base de données peuplée avec succès")
	return nil
}

// splitSQLStatements découpe un script SQL en instructions individuelles.
// Gère les chaînes entre guillemets simples (y compris les '' échappés) et les commentaires --.
// Cela évite les split naïfs sur les points-virgules qui pourraient être présents dans les données.
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
			// commentaire ligne : sauter jusqu'à la fin de ligne
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