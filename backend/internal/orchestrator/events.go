package orchestrator

import (
	"log"
	"sync"
	"time"
)

const maxBufferedEvents = 100

type Event struct {
	Type         string   `json:"type"`
	DeploymentID string   `json:"deployment_id"`
	NodeID       string   `json:"node_id,omitempty"`
	NodeIDs      []string `json:"node_ids,omitempty"`
	Status       string   `json:"status,omitempty"`
	Message      string   `json:"message,omitempty"`
	LevelIdx     int      `json:"level_idx,omitempty"`
	Timestamp    int64    `json:"timestamp"`
}

type EventHub struct {
	subscriptions  map[string][]chan *Event
	bufferedEvents map[string][]*Event
	mu             sync.RWMutex
	broadcastCh    chan *Event
}

func NewEventHub() *EventHub {
	hub := &EventHub{
		subscriptions:  make(map[string][]chan *Event),
		bufferedEvents: make(map[string][]*Event),
		broadcastCh:    make(chan *Event, 100),
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
			// Do not close the channel here to avoid racing with concurrent broadcasts.
			log.Printf("[EventHub] Unsubscribed from deployment %s (remaining: %d)", deploymentID, len(h.subscriptions[deploymentID]))
			break
		}
	}

	if len(h.subscriptions[deploymentID]) == 0 {
		delete(h.subscriptions, deploymentID)
	}
}

func (h *EventHub) GetBufferedEventsAndClear(deploymentID string) []*Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	events := h.bufferedEvents[deploymentID]
	delete(h.bufferedEvents, deploymentID)

	if len(events) > 0 {
		log.Printf("[EventHub] Flushed %d buffered events for deployment %s", len(events), deploymentID)
	}

	return events
}

// SubscribeAndDrain atomically registers a new subscriber for the given
// deployment and returns (and clears) any buffered events under the same lock.
// This prevents a race where events published between a separate
// GetBufferedEventsAndClear and Subscribe call would be missed.
func (h *EventHub) SubscribeAndDrain(deploymentID string) (chan *Event, []*Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan *Event, 10)
	h.subscriptions[deploymentID] = append(h.subscriptions[deploymentID], ch)
	log.Printf("[EventHub] New subscriber for deployment %s (total: %d)", deploymentID, len(h.subscriptions[deploymentID]))

	events := h.bufferedEvents[deploymentID]
	delete(h.bufferedEvents, deploymentID)

	if len(events) > 0 {
		log.Printf("[EventHub] Flushed %d buffered events for deployment %s", len(events), deploymentID)
	}

	return ch, events
}

func (h *EventHub) Publish(event *Event) {
	h.mu.Lock()
	hasSubscribers := len(h.subscriptions[event.DeploymentID]) > 0

	// Si no hay subscribers, guardar en buffer para subscribers tardíos
	if !hasSubscribers {
		if h.bufferedEvents[event.DeploymentID] == nil {
			h.bufferedEvents[event.DeploymentID] = []*Event{}
		}
		if len(h.bufferedEvents[event.DeploymentID]) < maxBufferedEvents {
			h.bufferedEvents[event.DeploymentID] = append(
				h.bufferedEvents[event.DeploymentID],
				event,
			)
			log.Printf("[EventHub] Buffered event %s for deployment %s (buffer size: %d)",
				event.Type, event.DeploymentID, len(h.bufferedEvents[event.DeploymentID]))
		}
	}
	h.mu.Unlock()

	// Broadcast normal a subscribers existentes
	select {
	case h.broadcastCh <- event:
	default:
		log.Printf("[EventHub] Broadcast channel full, dropping event for deployment %s", event.DeploymentID)
	}
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
		Timestamp:    time.Now().UnixMilli(),
	}
}

func EventLevelStart(deploymentID string, levelIdx int, nodeIDs []string) *Event {
	return &Event{
		Type:         "level_start",
		DeploymentID: deploymentID,
		LevelIdx:     levelIdx,
		NodeIDs:      nodeIDs,
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
