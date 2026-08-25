package service_test

import (
	"context"
	"eventingestion/internal/service"
	"eventingestion/internal/store"
	"eventingestion/pkg/idgen"
	"testing"
)

func TestPublishedSnapshotOwnsInputValues(t *testing.T) {
	m := store.NewMemory()
	s := service.New(m, m, m, m, &idgen.Atomic{})
	values := map[string]string{"state": "ready"}
	if _, err := s.PublishSnapshot(context.Background(), "feed", 1, values); err != nil {
		t.Fatalf("snapshot publish failed")
	}
	values["state"] = "mutated"
	got, err := s.LatestSnapshot(context.Background(), "feed")
	if err != nil {
		t.Fatalf("snapshot read failed")
	}
	if got.Values["state"] != "ready" {
		t.Fatalf("input map changed published snapshot")
	}
}
