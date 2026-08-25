package store

import "eventingestion/internal/model"

type EventStore interface {
	SaveEvent(model.Event) error
	GetEvent(string) (model.Event, error)
	ListEvents(string) []model.Event
}

type SubscriptionStore interface {
	SaveSubscription(model.Subscription) error
	ListSubscriptions(string, string) []model.Subscription
}

type DeliveryStore interface {
	SaveDelivery(model.Delivery) error
	ListPending(string) []model.Delivery
	AckDelivery(string) error
}

type SnapshotStore interface {
	PublishSnapshot(model.Snapshot) error
	LatestSnapshot(string) (model.Snapshot, error)
}
