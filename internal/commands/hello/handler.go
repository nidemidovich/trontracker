package hello

import (
	"github.com/nidemidovich/trontracker/internal/commands"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(c commands.Context) error {
	return c.Send("Hello!")
}
