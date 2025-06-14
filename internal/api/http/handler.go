package http

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"sync"
)

type mover interface {
	Move(rightDutyCycle float64, leftDutyCycle float64) error
}

type Handler struct {
	mover  mover
	mu     sync.Mutex
	logger *slog.Logger
}

func NewHandler(mover mover, logger *slog.Logger) *Handler {
	return &Handler{
		mover:  mover,
		logger: logger.With(slog.String("where", "http handler")),
	}
}

func (h *Handler) Move(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	request := &MoveRequest{}

	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Message: err.Error()})

		return
	}

	err = h.mover.Move(*request.RightDutyCycle, *request.LeftDutyCycle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Message: "internal server error"})
		h.logger.Error("internal server error", slog.Any("error", err))

		return
	}

	c.Status(http.StatusNoContent)
}
