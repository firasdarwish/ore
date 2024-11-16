package ore

import (
	"context"
)

// Get Retrieves an instance based on type and key (panics if no valid implementations)
func Get[T any](ctx context.Context) (T, context.Context) {
	return getFromContainer[T](DefaultContainer, ctx, nilKey)
}

// GetList Retrieves a list of instances based on type and key
func GetList[T any](ctx context.Context) ([]T, context.Context) {
	return getListFromContainer[T](DefaultContainer, ctx, nilKey)
}
