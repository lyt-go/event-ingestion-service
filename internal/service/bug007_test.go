package service_test

import (
	"context"
	"errors"
	"eventingestion/internal/model"
	"eventingestion/internal/service"
	"testing"
	"time"
)

type cancelSubStore struct{ cancel context.CancelFunc }

func (s *cancelSubStore) SaveSubscription(model.Subscription) error           { s.cancel(); return nil }
func (*cancelSubStore) ListSubscriptions(string, string) []model.Subscription { return nil }

type fixedIDs struct{ n int }

func (g *fixedIDs) Next() string   { g.n++; return "id" }
func (g *fixedIDs) Now() time.Time { return time.Unix(1, 0) }
func TestSubscribeReturnsCancellationAfterStoreWrite(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	subs := &cancelSubStore{cancel: cancel}
	s := service.New(&emptyEvents{}, subs, &emptyDeliveries{}, &emptySnapshots{}, &fixedIDs{})
	_, err := s.Subscribe(ctx, "worker", "feed", "change")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want cancellation, got %v", err)
	}
}

type emptyEvents struct{}

func (*emptyEvents) SaveEvent(model.Event) error          { return nil }
func (*emptyEvents) GetEvent(string) (model.Event, error) { return model.Event{}, model.ErrNotFound }
func (*emptyEvents) ListEvents(string) []model.Event      { return nil }

type emptyDeliveries struct{}

func (*emptyDeliveries) SaveDelivery(model.Delivery) error   { return nil }
func (*emptyDeliveries) ListPending(string) []model.Delivery { return nil }
func (*emptyDeliveries) AckDelivery(string) error            { return nil }

type emptySnapshots struct{}

func (*emptySnapshots) PublishSnapshot(model.Snapshot) error { return nil }
func (*emptySnapshots) LatestSnapshot(string) (model.Snapshot, error) {
	return model.Snapshot{}, model.ErrNotFound
}
