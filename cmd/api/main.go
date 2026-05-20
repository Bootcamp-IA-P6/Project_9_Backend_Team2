package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Bootcamp-IA-P6/Project_9_Backend_Team2/internal/ai"
	"github.com/Bootcamp-IA-P6/Project_9_Backend_Team2/internal/dtos"
	"github.com/Bootcamp-IA-P6/Project_9_Backend_Team2/internal/youtube"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
	// IMPORT CORREGIDO CON LA RUTA REAL DE TU PROYECTO
	"github.com/Bootcamp-IA-P6/Project_9_Backend_Team2/internal/database"
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
	dbURL := os.Getenv("DATABASE_URL") // Cargamos la URL de Supabase
	fmt.Println("URL que está leyendo Go:", dbURL)

	if apiKey == "" || aiURL == "" || dbURL == "" {
		log.Fatal("Faltan variables en el archivo .env (YOUTUBE_API_KEY, AI_MODEL_URL o DATABASE_URL)")
	}

	// 2. Inicializar base de datos y hacer prueba de inserción
	dbClient, err := database.NewDatabaseClient(dbURL)
	if err != nil {
		log.Fatal("❌ Error conectando a Supabase:", err)
	}
	defer dbClient.Conn.Close()
	fmt.Println("✅ ¡Conexión a Supabase perfecta!")

	err = dbClient.InsertarPrueba()
	if err != nil {
		log.Println("⚠️ Aviso: Fallo al insertar la fila de prueba:", err)
	} else {
		fmt.Println("💾 ¡Fila insertada! Ve a mirar la tabla en la web de Supabase.")
	}

	// 3. Inicializar servicios externos
	ytService, _ := youtube.NewYouTubeService(apiKey)
	aiService := ai.NewAIService(aiURL)

	// 4. Levantar servidor
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

	// 5. EL ENDPOINT MÁGICO
	r.POST("/api/v1/analyze", func(c *gin.Context) {
		var req dtos.AnalyzeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}

		videoID := extractVideoID(req.VideoURL)
		
		// A. Extraer comentarios
		comments, err := ytService.GetVideoComments(videoID, 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo en YouTube"})
			return
		}

		// B. Analizar cada comentario con la IA
		var resultadosFinales []gin.H
		comentariosToxicos := 0

		for _, comment := range comments {
			analisis, err := aiService.EvaluateComment(comment.Text)
			if err != nil {
				log.Printf("Aviso: Fallo al analizar un comentario: %v", err)
				continue 
			}

			if analisis.EsToxico {
				comentariosToxicos++
			}

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
