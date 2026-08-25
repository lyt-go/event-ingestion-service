package idgen

import (
	"fmt"
	"time"
)

type Generator interface {
	Next() string
	Now() time.Time
}
type Atomic struct{ seq uint64 }

func (g *Atomic) Next() string   { g.seq++; return fmt.Sprintf("evt-%d", g.seq) }
func (g *Atomic) Now() time.Time { return time.Now().UTC() }
