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
    umbral := 0.85 
    
    if val, err := strconv.ParseFloat(umbralStr, 64); err == nil {
        umbral = val
    }

    esEtiquetaToxica := aiResp.EtiquetaModelo == "Toxic" || aiResp.EtiquetaModelo == "NEGATIVE"

    if esEtiquetaToxica {
        if aiResp.ScoreConfianza >= umbral {
            // Es verdaderamente tóxico
            aiResp.EsToxico = true
        } else {
            // Falso positivo detectado: El umbral lo salva
            aiResp.EsToxico = false
            aiResp.EtiquetaModelo = "Safe" // 👈 Sobrescribimos la etiqueta para "engañar" al frontend
        }
    } else {
        // Era sano desde el principio
        aiResp.EsToxico = false
        aiResp.EtiquetaModelo = "Safe" // 👈 Aseguramos que el texto siempre sea "Safe"
    }
    // ==========================================

    return aiResp, nil
}