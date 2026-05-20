package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Driver de Postgres
)

type DBClient struct {
	Conn *sql.DB
}

// NewDatabaseClient conecta Go con Supabase
func NewDatabaseClient(databaseURL string) (*DBClient, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %v", err)
	}

	// Verificar la conexión real
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("no se pudo conectar a Supabase: %v", err)
	}

	return &DBClient{Conn: db}, nil
}

//InsertarPrueba mete una fila de mentira para comprobar que tenemos permisos de escritura
func (c *DBClient) InsertarPrueba() error {
	// IMPORTANTE: En Postgres, si la tabla tiene mayúsculas, hay que ponerla entre comillas
	query := `INSERT INTO "ToxicFilterAIPredictions" (comment_id, video_id, text, is_toxic) 
			  VALUES ($1, $2, $3, $4)`
	
	// Metemos unos datos inventados
	_, err := c.Conn.Exec(query, "comentario_test_01", "VIDEO_123", "Este es un comentario de prueba desde Go", false)
	return err
}