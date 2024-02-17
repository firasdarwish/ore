package ore

import "context"

type Creator[T any] interface {
	New(ctx context.Context) T
}
