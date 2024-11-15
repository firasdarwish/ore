package ore

import (
	"context"
)

// GetKeyedFromContainer Retrieves an instance from the given container based on type and key (panics if no valid implementations)
func GetKeyedFromContainer[T any](con *Container, ctx context.Context, key KeyStringer) (T, context.Context) {
	if key == nil {
		panic(nilKey)
	}
	return getFromContainer[T](con, ctx, key)
}

// GetKeyedListFromContainer Retrieves a list of instances from the given container based on type and key
func GetKeyedListFromContainer[T any](con *Container, ctx context.Context, key KeyStringer) ([]T, context.Context) {
	if key == nil {
		panic(nilKey)
	}
	return getListFromContainer[T](con, ctx, key)
}
