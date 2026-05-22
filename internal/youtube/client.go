package youtube

import (
	"context"
	"fmt"

	"github.com/Bootcamp-IA-P6/Project_9_Backend_Team2/internal/dtos"
	"google.golang.org/api/option"
	yt "google.golang.org/api/youtube/v3"
)

// YouTubeService es nuestro conector oficial
type YouTubeService struct {
	client *yt.Service
}

// NewYouTubeService inicializa la conexión con Google usando mi API Key
func NewYouTubeService(apiKey string) (*YouTubeService, error) {
	ctx := context.Background()
	client, err := yt.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error al iniciar el cliente de youtube: %v", err)
	}
	return &YouTubeService{client: client}, nil
}

// GetVideoComments extrae los comentarios de un vídeo específico y los limpia
func (s *YouTubeService) GetVideoComments(videoID string, maxResults int64) ([]dtos.Comment, error) {
	// 1. Preparamos la petición a la API de YouTube
	call := s.client.CommentThreads.List([]string{"snippet"}).
		VideoId(videoID).
		MaxResults(maxResults).
		TextFormat("plainText") // Pedimos texto limpio sin HTML

	// 2. Ejecutamos la petición
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error al descargar los comentarios: %v", err)
	}

	// 3. Mapeamos la respuesta cruda de Google a nuestro DTO limpio
	var comentariosLimpios []dtos.Comment

	for _, item := range response.Items {
		comentarioOriginal := item.Snippet.TopLevelComment.Snippet
		
		comentariosLimpios = append(comentariosLimpios, dtos.Comment{
			Author: comentarioOriginal.AuthorDisplayName,
			Text:   comentarioOriginal.TextDisplay,
		})
	}

	return comentariosLimpios, nil
}