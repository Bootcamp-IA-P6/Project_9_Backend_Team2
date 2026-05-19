package dtos

// React nos enviará un JSON con la URL del video
type AnalyzeRequest struct {
	VideoURL string `json:"video_url" binding:"required"`
}

// La estructura de un solo comentario limpio
type Comment struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

// Lo que le devolveremos a React al final del proceso
type AnalyzeResponse struct {
	VideoID       string    `json:"video_id"`
	TotalComments int       `json:"total_comments"`
	ToxicityScore float64   `json:"toxicity_score"`
	Status        string    `json:"status"`
	// Aquí luego meteremos los resultados detallados
}