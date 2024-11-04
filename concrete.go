package ore

import "time"

// concrete holds the resolved instance value and other metadata
type concrete struct {
	// the value implementation
	value any
	// the creation time
	createdAt time.Time
	// the lifetime of this concrete
	lifetime Lifetime
}
