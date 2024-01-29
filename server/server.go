package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/koenno/aiapi/aiservice"

	"github.com/gin-gonic/gin"
)

type Answerer interface {
	Ask(ctx context.Context, question string) (string, error)
}

type Server struct {
	answerer Answerer
}

func New(answerer Answerer) *Server {
	return &Server{
		answerer: answerer,
	}
}

func (s *Server) Routes(r gin.IRoutes) {
	r.POST("/answer", s.Answer)
}

func (s *Server) Answer(gc *gin.Context) {
	var questionDesc struct {
		Question string `json:"question"`
	}
	if err := gc.BindJSON(&questionDesc); err != nil {
		log.Printf("ERR: failed to bing json: %v", err)
		gc.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	answer, err := s.answerer.Ask(ctx, questionDesc.Question)
	if errors.Is(err, aiservice.ErrModerator) || errors.Is(err, aiservice.ErrTooLong) || errors.Is(err, aiservice.ErrTooShort) {
		log.Printf("ERR: invalid request: %v", err)
		gc.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, aiservice.ErrEngine) {
		log.Printf("ERR: internal failure: %v", err)
		gc.JSON(500, gin.H{"error": "internal error"})
		return
	}

	log.Printf("INF: queston: '%s', answer: '%s'", questionDesc.Question, answer)
	gc.JSON(http.StatusOK, gin.H{"reply": answer})
}
