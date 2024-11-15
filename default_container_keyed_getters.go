package ore

import (
	"context"
)

// GetKeyed Retrieves an instance based on type and key (panics if no valid implementations)
func GetKeyed[T any](ctx context.Context, key KeyStringer) (T, context.Context) {
	if key == nil {
		panic(nilKey)
	}
	return getFromContainer[T](DefaultContainer, ctx, key)
}

// GetKeyedList Retrieves a list of instances based on type and key
func GetKeyedList[T any](ctx context.Context, key KeyStringer) ([]T, context.Context) {
	if key == nil {
		panic(nilKey)
	}
	return getListFromContainer[T](DefaultContainer, ctx, key)
}
