package service_test

import (
	"context"
	"eventingestion/internal/service"
	"eventingestion/internal/store"
	"eventingestion/pkg/idgen"
	"testing"
)

func TestEventReachesEveryMatchingSubscriber(t *testing.T) {
	m := store.NewMemory()
	s := service.New(m, m, m, m, &idgen.Atomic{})
	ctx := context.Background()
	for _, subscriber := range []string{"alpha", "beta"} {
		if _, err := s.Subscribe(ctx, subscriber, "feed", "change"); err != nil {
			t.Fatalf("subscription setup failed")
		}
	}
	if _, err := s.Receive(ctx, "feed", "change", map[string]string{"v": "1"}); err != nil {
		t.Fatalf("event receive failed")
	}
	for _, subscriber := range []string{"alpha", "beta"} {
		got, err := s.Pending(ctx, subscriber)
		if err != nil {
			t.Fatalf("pending read failed")
		}
		if len(got) != 1 {
			t.Fatalf("every matching subscriber must receive one delivery")
		}
	}
}
