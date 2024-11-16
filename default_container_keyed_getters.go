package ore

import (
	"context"
)

// GetKeyed Retrieves an instance based on type and key (panics if no valid implementations)
func GetKeyed[T any, K comparable](ctx context.Context, key K) (T, context.Context) {
	return getFromContainer[T](DefaultContainer, ctx, key)
}

// GetKeyedList Retrieves a list of instances based on type and key
func GetKeyedList[T any, K comparable](ctx context.Context, key K) ([]T, context.Context) {
	return getListFromContainer[T](DefaultContainer, ctx, key)
}
