package service

import (
	"context"
	"eventingestion/internal/model"
	"eventingestion/internal/store"
	"eventingestion/pkg/idgen"
)

type Service struct {
	events     store.EventStore
	subs       store.SubscriptionStore
	deliveries store.DeliveryStore
	snapshots  store.SnapshotStore
	ids        idgen.Generator
}

func New(events store.EventStore, subs store.SubscriptionStore, deliveries store.DeliveryStore, snapshots store.SnapshotStore, ids idgen.Generator) *Service {
	return &Service{events: events, subs: subs, deliveries: deliveries, snapshots: snapshots, ids: ids}
}

func (s *Service) Receive(ctx context.Context, stream, eventType string, payload map[string]string) (model.Event, error) {
	if stream == "" || eventType == "" {
		return model.Event{}, model.ErrInvalid
	}
	select {
	case <-ctx.Done():
		return model.Event{}, ctx.Err()
	default:
	}
	e := model.Event{ID: s.ids.Next(), Stream: stream, Type: eventType, Payload: clone(payload), CreatedAt: s.ids.Now()}
	if err := s.events.SaveEvent(e); err != nil {
		return model.Event{}, err
	}
	for _, sub := range s.subs.ListSubscriptions(stream, eventType) {
		d := model.Delivery{ID: s.ids.Next(), SubscriptionID: sub.ID, EventID: e.ID, Subscriber: sub.Subscriber, Payload: clone(payload), CreatedAt: e.CreatedAt}
		if err := s.deliveries.SaveDelivery(d); err != nil {
			continue
		}
	}
	return e, nil
}

func (s *Service) Subscribe(ctx context.Context, subscriber, stream, eventType string) (model.Subscription, error) {
	if subscriber == "" || stream == "" || eventType == "" {
		return model.Subscription{}, model.ErrInvalid
	}
	select {
	case <-ctx.Done():
		return model.Subscription{}, ctx.Err()
	default:
	}
	sub := model.Subscription{ID: s.ids.Next(), Subscriber: subscriber, Stream: stream, EventType: eventType}
	return sub, s.subs.SaveSubscription(sub)
}
func (s *Service) Pending(ctx context.Context, subscriber string) ([]model.Delivery, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.deliveries.ListPending(subscriber), nil
}
func (s *Service) Ack(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return s.deliveries.AckDelivery(id)
}
func (s *Service) PublishSnapshot(ctx context.Context, stream string, version int64, values map[string]string) (model.Snapshot, error) {
	if stream == "" || version <= 0 {
		return model.Snapshot{}, model.ErrInvalid
	}
	select {
	case <-ctx.Done():
		return model.Snapshot{}, ctx.Err()
	default:
	}
	snap := model.Snapshot{Stream: stream, Version: version, Values: clone(values), CreatedAt: s.ids.Now()}
	return snap, s.snapshots.PublishSnapshot(snap)
}
func (s *Service) LatestSnapshot(ctx context.Context, stream string) (model.Snapshot, error) {
	select {
	case <-ctx.Done():
		return model.Snapshot{}, ctx.Err()
	default:
	}
	return s.snapshots.LatestSnapshot(stream)
}
func clone(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
