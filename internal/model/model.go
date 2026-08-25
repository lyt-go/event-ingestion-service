package model

import "time"

type Event struct {
	ID        string            `json:"id"`
	Stream    string            `json:"stream"`
	Type      string            `json:"type"`
	Payload   map[string]string `json:"payload"`
	CreatedAt time.Time         `json:"created_at"`
}

type Subscription struct {
	ID         string `json:"id"`
	Subscriber string `json:"subscriber"`
	Stream     string `json:"stream"`
	EventType  string `json:"event_type"`
}

type Delivery struct {
	ID             string            `json:"id"`
	SubscriptionID string            `json:"subscription_id"`
	EventID        string            `json:"event_id"`
	Subscriber     string            `json:"subscriber"`
	Payload        map[string]string `json:"payload"`
	Acked          bool              `json:"acked"`
	CreatedAt      time.Time         `json:"created_at"`
}

type Snapshot struct {
	Stream    string            `json:"stream"`
	Version   int64             `json:"version"`
	Values    map[string]string `json:"values"`
	CreatedAt time.Time         `json:"created_at"`
}
