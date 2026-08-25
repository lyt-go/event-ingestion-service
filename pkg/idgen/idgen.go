package idgen

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Generator interface {
	Next() string
	Now() time.Time
}
type Atomic struct{ seq uint64 }

func (g *Atomic) Next() string   { return fmt.Sprintf("evt-%d", atomic.AddUint64(&g.seq, 1)) }
func (g *Atomic) Now() time.Time { return time.Now().UTC() }
