package infrastructure

import (
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/hub"
	"log"
)

type HubEventPublisher struct {
	hub *hub.Hub
}

func NewHubEventPublisher(h *hub.Hub) *HubEventPublisher {
	return &HubEventPublisher{hub: h}
}

func (p *HubEventPublisher) Publish(event any) {
	log.Println("Hub received broadcast")
	switch e := event.(type) {

	case domain.MessageSentEvent:
		p.hub.Broadcast <- hub.MessageEvent{
			Message:    e.Message,
			Recipients: e.Recipients,
		}
	}
}
