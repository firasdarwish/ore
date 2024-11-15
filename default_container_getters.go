package ore

import (
	"context"
)

// Get Retrieves an instance based on type and key (panics if no valid implementations)
func Get[T any](ctx context.Context, key ...KeyStringer) (T, context.Context) {
	return GetFromContainer[T](DefaultContainer, ctx, key...)
}

// GetList Retrieves a list of instances based on type and key
func GetList[T any](ctx context.Context, key ...KeyStringer) ([]T, context.Context) {
	return GetListFromContainer[T](DefaultContainer, ctx, key...)
}
