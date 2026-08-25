package store

import (
	"eventingestion/internal/model"
	"sync"
)

type Memory struct {
	mu            sync.RWMutex
	events        map[string]model.Event
	subscriptions map[string]model.Subscription
	deliveries    map[string]model.Delivery
	snapshots     map[string]model.Snapshot
}

func NewMemory() *Memory {
	return &Memory{events: map[string]model.Event{}, subscriptions: map[string]model.Subscription{}, deliveries: map[string]model.Delivery{}, snapshots: map[string]model.Snapshot{}}
}

func (m *Memory) SaveEvent(e model.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events[e.ID] = cloneEvent(e)
	return nil
}
func (m *Memory) GetEvent(id string) (model.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.events[id]
	if !ok {
		return model.Event{}, model.ErrNotFound
	}
	return cloneEvent(e), nil
}
func (m *Memory) ListEvents(stream string) []model.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]model.Event, 0)
	for _, e := range m.events {
		if e.Stream == stream {
			out = append(out, cloneEvent(e))
		}
	}
	return out
}
func (m *Memory) SaveSubscription(s model.Subscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscriptions[s.ID] = s
	return nil
}
func (m *Memory) ListSubscriptions(stream, eventType string) []model.Subscription {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]model.Subscription, 0)
	for _, s := range m.subscriptions {
		if s.Stream == stream && (s.EventType == eventType || s.EventType == "*") {
			out = append(out, s)
		}
	}
	return out
}
func (m *Memory) SaveDelivery(d model.Delivery) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deliveries[d.ID] = cloneDelivery(d)
	return nil
}
func (m *Memory) ListPending(subscriber string) []model.Delivery {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]model.Delivery, 0)
	for _, d := range m.deliveries {
		if d.Subscriber == subscriber && !d.Acked {
			out = append(out, cloneDelivery(d))
		}
	}
	return out
}
func (m *Memory) AckDelivery(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.deliveries[id]
	if !ok {
		return model.ErrNotFound
	}
	if d.Acked {
		return model.ErrAcked
	}
	d.Acked = true
	m.deliveries[id] = d
	return nil
}
func (m *Memory) PublishSnapshot(s model.Snapshot) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.snapshots[s.Stream] = cloneSnapshot(s)
	return nil
}
func (m *Memory) LatestSnapshot(stream string) (model.Snapshot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.snapshots[stream]
	if !ok {
		return model.Snapshot{}, model.ErrNotFound
	}
	return cloneSnapshot(s), nil
}

func cloneEvent(e model.Event) model.Event          { e.Payload = cloneMap(e.Payload); return e }
func cloneDelivery(d model.Delivery) model.Delivery { d.Payload = cloneMap(d.Payload); return d }
func cloneSnapshot(s model.Snapshot) model.Snapshot { s.Values = cloneMap(s.Values); return s }
func cloneMap(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
