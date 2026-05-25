package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os" // Necesario para leer la variable de entorno
    "strconv" // Necesario para convertir texto a número
)

// AIService gestiona la conexión con tu modelo de ML
type AIService struct {
	ModelURL string
}

// NewAIService inicializa el cliente con la URL de tu .env
func NewAIService(modelURL string) *AIService {
	return &AIService{ModelURL: modelURL}
}

// AIRequest es el JSON que le enviamos a la IA ({"texto": "..."})
type AIRequest struct {
	Texto string `json:"texto"`
}

// AIResponse es el JSON que esperamos recibir de la IA
type AIResponse struct {
	TextoAnalizado string  `json:"texto_analizado"`
	EtiquetaModelo string  `json:"etiqueta_modelo"`
	ScoreConfianza float64 `json:"score_confianza"`
	EsToxico       bool    `json:"es_toxico"`
}

// EvaluateComment envía un texto individual a tu API y devuelve el resultado
func (s *AIService) EvaluateComment(text string) (AIResponse, error) {
    // 1. Preparamos el JSON de envío
    reqBody := AIRequest{Texto: text}
    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return AIResponse{}, fmt.Errorf("error empaquetando JSON: %v", err)
    }

    // 2. Hacemos la petición POST a tu API de Python
    resp, err := http.Post(s.ModelURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return AIResponse{}, fmt.Errorf("error contactando con la IA: %v", err)
    }
    defer resp.Body.Close()

   // 3. Leemos y transformamos la respuesta
    body, _ := io.ReadAll(resp.Body)
    var aiResp AIResponse
    if err := json.Unmarshal(body, &aiResp); err != nil {
        return AIResponse{}, fmt.Errorf("error decodificando respuesta de la IA: %v", err)
    }

    // ==========================================
    // 4. APLICAR EL UMBRAL DE DECISIÓN (THRESHOLD)
    // ==========================================
    umbralStr := os.Getenv("UMBRAL_TOXICIDAD")
    umbral := 0.85 // Nuestro umbral estricto para evitar falsos positivos
    
    if val, err := strconv.ParseFloat(umbralStr, 64); err == nil {
        umbral = val
    }

    // 1º Comprobamos si la API de Python clasificó el texto como tóxico
    // (Usamos exactamente las mismas etiquetas que hay en la api de Hugging Face )
    esEtiquetaToxica := aiResp.EtiquetaModelo == "Toxic" || aiResp.EtiquetaModelo == "NEGATIVE"

    if esEtiquetaToxica {
        // 2º Si es tóxico, le pasamos NUESTRO filtro. 
        // ¿Está el modelo lo suficientemente seguro (>= 0.85)?
        aiResp.EsToxico = aiResp.ScoreConfianza >= umbral
    } else {
        // 3º Si la etiqueta es sana (ej: "POSITIVE"), nos aseguramos de que sea false
        aiResp.EsToxico = false
    }
    // ==========================================

    return aiResp, nil
}
