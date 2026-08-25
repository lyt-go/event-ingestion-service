package store_test

import (
	"errors"
	"eventingestion/internal/model"
	"eventingestion/internal/store"
	"sync"
	"testing"
)

func TestConcurrentAcksHaveSingleWinner(t *testing.T) {
	m := store.NewMemory()
	if err := m.SaveDelivery(model.Delivery{ID: "d1", Subscriber: "worker"}); err != nil {
		t.Fatalf("delivery setup failed")
	}
	const consumers = 128
	start := make(chan struct{})
	errs := make(chan error, consumers)
	var wg sync.WaitGroup
	for i := 0; i < consumers; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; errs <- m.AckDelivery("d1") }()
	}
	close(start)
	wg.Wait()
	close(errs)
	var ok, acked int
	for err := range errs {
		if err == nil {
			ok++
		} else if errors.Is(err, model.ErrAcked) {
			acked++
		} else {
			t.Fatalf("unexpected ack error")
		}
	}
	if ok != 1 || acked != consumers-1 {
		t.Fatalf("ack results did not have one winner")
	}
}
