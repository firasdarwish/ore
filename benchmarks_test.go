package ore

import (
	"context"
	"testing"
)

func BenchmarkRegisterLazyFunc(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterLazyFunc[someCounter](Scoped, func(ctx context.Context) (someCounter, context.Context) {
			return &simpleCounter{}, ctx
		})
	}
}

func BenchmarkRegisterLazyCreator(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterLazyCreator[someCounter](Scoped, &simpleCounter{})
	}
}

func BenchmarkRegisterEagerSingleton(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterEagerSingleton[someCounter](&simpleCounter{})
	}
}

func BenchmarkGet(b *testing.B) {
	clearAll()

	RegisterLazyFunc[someCounter](Scoped, func(ctx context.Context) (someCounter, context.Context) {
		return &simpleCounter{}, ctx
	})

	RegisterEagerSingleton[someCounter](&simpleCounter{})

	RegisterLazyCreator[someCounter](Scoped, &simpleCounter{})

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get[someCounter](ctx)
	}
}

func BenchmarkGetList(b *testing.B) {
	clearAll()

	RegisterLazyFunc[someCounter](Scoped, func(ctx context.Context) (someCounter, context.Context) {
		return &simpleCounter{}, ctx
	})

	RegisterEagerSingleton[someCounter](&simpleCounter{})

	RegisterLazyCreator[someCounter](Scoped, &simpleCounter{})

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetList[someCounter](ctx)
	}
}
