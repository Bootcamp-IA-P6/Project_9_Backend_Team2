package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Bootcamp-IA-P6/Project_9_Backend_Team2/internal/ai"
	"github.com/Bootcamp-IA-P6/Project_9_Backend_Team2/internal/dtos"
	"github.com/Bootcamp-IA-P6/Project_9_Backend_Team2/internal/youtube"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Función para extraer el ID del vídeo
func extractVideoID(videoURL string) string {
	u, err := url.Parse(videoURL)
	if err != nil {
		return ""
	}
	return u.Query().Get("v")
}

func main() {
	// 1. Cargar entorno
	godotenv.Load()
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	aiURL := os.Getenv("AI_MODEL_URL")

	if apiKey == "" || aiURL == "" {
		log.Fatal("Faltan variables en el archivo .env (YOUTUBE_API_KEY o AI_MODEL_URL)")
	}

	// 2. Inicializar servicios
	ytService, _ := youtube.NewYouTubeService(apiKey)
	aiService := ai.NewAIService(aiURL)

	// 3. Levantar servidor
	r := gin.Default()

	// Configuración CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 4. EL ENDPOINT MÁGICO
	r.POST("/api/v1/analyze", func(c *gin.Context) {
		var req dtos.AnalyzeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}

		videoID := extractVideoID(req.VideoURL)
		
		// A. Extraer comentarios (pedimos 10 para no saturar tu IA en pruebas)
		comments, err := ytService.GetVideoComments(videoID, 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo en YouTube"})
			return
		}

		// B. Analizar cada comentario con la IA
		var resultadosFinales []gin.H
		comentariosToxicos := 0

		for _, comment := range comments {
			// Enviar a la IA
			analisis, err := aiService.EvaluateComment(comment.Text)
			if err != nil {
				log.Printf("Aviso: Fallo al analizar un comentario: %v", err)
				continue // Si uno falla, saltamos al siguiente
			}

			if analisis.EsToxico {
				comentariosToxicos++
			}

			// Guardar el resultado combinado
			resultadosFinales = append(resultadosFinales, gin.H{
				"autor":      comment.Author,
				"texto":      comment.Text,
				"evaluacion": analisis,
			})
		}

		// C. Calcular métricas finales
		totalValidos := len(resultadosFinales)
		porcentajeToxicidad := 0.0
		if totalValidos > 0 {
			porcentajeToxicidad = (float64(comentariosToxicos) / float64(totalValidos)) * 100
		}

		// D. Devolver la respuesta a React
		c.JSON(http.StatusOK, gin.H{
			"video_id":             videoID,
			"comentarios_totales":  totalValidos,
			"comentarios_toxicos":  comentariosToxicos,
			"porcentaje_toxicidad": porcentajeToxicidad,
			"detalle":              resultadosFinales,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}