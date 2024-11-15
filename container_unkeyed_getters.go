package ore

import (
	"context"
)

// GetFromContainer Retrieves an instance from the given container based on type and key (panics if no valid implementations)
func GetFromContainer[T any](con *Container, ctx context.Context) (T, context.Context) {
	return getFromContainer[T](con, ctx, nil)
}

// GetListFromContainer Retrieves a list of instances from the given container based on type and key
func GetListFromContainer[T any](con *Container, ctx context.Context) ([]T, context.Context) {
	return getListFromContainer[T](con, ctx, nil)
}
