package printer

import (
	"context"
)

// Printer .
type Printer interface {
	Print(ctx context.Context, relation map[interface{}]map[interface{}][]string) (string, error)
}
