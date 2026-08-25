package handler

import (
	"context"
	"encoding/json"
	"errors"
	"eventingestion/internal/model"
	"eventingestion/internal/service"
	"net/http"
	"strings"
)

type Server struct{ svc *service.Service }

func New(svc *service.Service) http.Handler {
	s := &Server{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/events", s.events)
	mux.HandleFunc("/v1/subscriptions", s.subscriptions)
	mux.HandleFunc("/v1/deliveries", s.deliveries)
	mux.HandleFunc("/v1/deliveries/", s.ack)
	mux.HandleFunc("/v1/snapshots", s.snapshots)
	mux.HandleFunc("/v1/snapshots/", s.latest)
	return mux
}
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
func decode(r *http.Request, v any) error { return json.NewDecoder(r.Body).Decode(v) }
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var in struct {
		Stream  string            `json:"stream"`
		Type    string            `json:"type"`
		Payload map[string]string `json:"payload"`
	}
	if err := decode(r, &in); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	e, err := s.svc.Receive(r.Context(), in.Stream, in.Type, in.Payload)
	if err != nil {
		statusError(w, err)
		return
	}
	writeJSON(w, 201, e)
}
func (s *Server) subscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var in struct{ Subscriber, Stream, EventType string }
	if err := decode(r, &in); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	sub, err := s.svc.Subscribe(r.Context(), in.Subscriber, in.Stream, in.EventType)
	if err != nil {
		statusError(w, err)
		return
	}
	writeJSON(w, 201, sub)
}
func (s *Server) deliveries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, 200, []model.Delivery{})
}
func (s *Server) ack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/ack") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v1/deliveries/"), "/ack")
	if err := s.svc.Ack(r.Context(), id); err != nil {
		statusError(w, err)
		return
	}
	writeJSON(w, 200, map[string]string{"status": "acked"})
}
func (s *Server) snapshots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var in struct {
		Stream  string            `json:"stream"`
		Version int64             `json:"version"`
		Values  map[string]string `json:"values"`
	}
	if err := decode(r, &in); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	snap, err := s.svc.PublishSnapshot(r.Context(), in.Stream, in.Version, in.Values)
	if err != nil {
		statusError(w, err)
		return
	}
	writeJSON(w, 201, snap)
}
func (s *Server) latest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	stream := strings.TrimPrefix(r.URL.Path, "/v1/snapshots/")
	snap, err := s.svc.LatestSnapshot(r.Context(), stream)
	if err != nil {
		statusError(w, err)
		return
	}
	writeJSON(w, 200, snap)
}
func statusError(w http.ResponseWriter, err error) {
	code := 500
	if errors.Is(err, model.ErrInvalid) {
		code = 400
	}
	if errors.Is(err, model.ErrNotFound) {
		code = 404
	}
	if errors.Is(err, model.ErrAcked) {
		code = 409
	}
	if errors.Is(err, context.Canceled) {
		code = 499
	}
	writeJSON(w, code, map[string]string{"error": err.Error()})
}
