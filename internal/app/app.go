package app

import (
	"eventingestion/internal/service"
	"eventingestion/internal/store"
	"eventingestion/pkg/idgen"
)

func NewService() *service.Service {
	m := store.NewMemory()
	return service.New(m, m, m, m, &idgen.Atomic{})
}
