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

// GuardarResumenVideo guarda las métricas globales en la tabla Padre
func (c *DBClient) GuardarResumenVideo(videoID string, totalComments int, promedioToxicidad float64, alertasCriticas int, healthScore float64) error {
	// IMPORTANTE: Al igual que antes, usamos comillas dobles porque "VideoSummary" lleva mayúsculas
	query := `INSERT INTO "VideoSummary" (video_id, total_comments, promedio_toxicidad, alertas_criticas, health_score) 
			  VALUES ($1, $2, $3, $4, $5)`
	
	_, err := c.Conn.Exec(query, videoID, totalComments, promedioToxicidad, alertasCriticas, healthScore)
	return err
}

// GuardarPrediccionComentario guarda los comentarios individuales en la tabla Hijo
func (c *DBClient) GuardarPrediccionComentario(videoID string, commentID string, text string, isToxic bool) error {
	query := `INSERT INTO "ToxicFilterAIPredictions" (comment_id, video_id, text, is_toxic) 
			  VALUES ($1, $2, $3, $4)`
	
	_, err := c.Conn.Exec(query, commentID, videoID, text, isToxic)
	return err
}