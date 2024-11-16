package ore

import (
	"context"
)

// GetKeyedFromContainer Retrieves an instance from the given container based on type and key (panics if no valid implementations)
func GetKeyedFromContainer[T any, K comparable](con *Container, ctx context.Context, key K) (T, context.Context) {
	return getFromContainer[T](con, ctx, key)
}

// GetKeyedListFromContainer Retrieves a list of instances from the given container based on type and key
func GetKeyedListFromContainer[T any, K comparable](con *Container, ctx context.Context, key K) ([]T, context.Context) {
	return getListFromContainer[T](con, ctx, key)
}
