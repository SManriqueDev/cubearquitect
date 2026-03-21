package orchestrator

import (
	"log"
	"sync"
	"time"
)

type Event struct {
	Type         string `json:"type"`
	DeploymentID string `json:"deployment_id"`
	NodeID       string `json:"node_id,omitempty"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message,omitempty"`
	LevelIdx     int    `json:"level_idx,omitempty"`
	Timestamp    int64  `json:"timestamp"`
}

type EventHub struct {
	subscriptions map[string][]chan *Event
	mu            sync.RWMutex
	broadcastCh   chan *Event
}

func NewEventHub() *EventHub {
	hub := &EventHub{
		subscriptions: make(map[string][]chan *Event),
		broadcastCh:   make(chan *Event, 100),
	}
	go hub.broadcast()
	return hub
}

func (h *EventHub) Subscribe(deploymentID string) chan *Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan *Event, 10)
	h.subscriptions[deploymentID] = append(h.subscriptions[deploymentID], ch)
	log.Printf("[EventHub] New subscriber for deployment %s (total: %d)", deploymentID, len(h.subscriptions[deploymentID]))

	return ch
}

func (h *EventHub) Unsubscribe(deploymentID string, ch chan *Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	channels := h.subscriptions[deploymentID]
	for i, c := range channels {
		if c == ch {
			h.subscriptions[deploymentID] = append(channels[:i], channels[i+1:]...)
			close(ch)
			log.Printf("[EventHub] Unsubscribed from deployment %s (remaining: %d)", deploymentID, len(h.subscriptions[deploymentID]))
			break
		}
	}

	if len(h.subscriptions[deploymentID]) == 0 {
		delete(h.subscriptions, deploymentID)
	}
}

func (h *EventHub) Publish(event *Event) {
	h.broadcastCh <- event
}

func (h *EventHub) broadcast() {
	for event := range h.broadcastCh {
		h.mu.RLock()
		channels := h.subscriptions[event.DeploymentID]
		h.mu.RUnlock()

		for _, ch := range channels {
			select {
			case ch <- event:
			default:
				log.Printf("[EventHub] Channel full for deployment %s, dropping event", event.DeploymentID)
			}
		}
	}
}

func EventFromNodeStatus(deploymentID string, nodeID string, status *NodeStatus) *Event {
	return &Event{
		Type:         "node_update",
		DeploymentID: deploymentID,
		NodeID:       nodeID,
		Status:       status.Status,
		Message:      status.Error,
		Timestamp:    status.Timestamp,
	}
}

func EventLevelStart(deploymentID string, levelIdx int) *Event {
	return &Event{
		Type:         "level_start",
		DeploymentID: deploymentID,
		LevelIdx:     levelIdx,
		Timestamp:    time.Now().UnixMilli(),
	}
}

func EventLevelComplete(deploymentID string, levelIdx int) *Event {
	return &Event{
		Type:         "level_complete",
		DeploymentID: deploymentID,
		LevelIdx:     levelIdx,
		Timestamp:    time.Now().UnixMilli(),
	}
}

func EventError(deploymentID string, message string) *Event {
	return &Event{
		Type:         "error",
		DeploymentID: deploymentID,
		Message:      message,
		Timestamp:    time.Now().UnixMilli(),
	}
}
