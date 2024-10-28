package ore

import (
	"context"
	"testing"
)

func BenchmarkRegisterLazyFunc(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterLazyFunc[Counter](Scoped, func(ctx context.Context) (Counter, context.Context) {
			return &simpleCounter{}, ctx
		})
	}
}

func BenchmarkRegisterLazyCreator(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterLazyCreator[Counter](Scoped, &simpleCounter{})
	}
}

func BenchmarkRegisterEagerSingleton(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterEagerSingleton[Counter](&simpleCounter{})
	}
}

func BenchmarkGet(b *testing.B) {
	clearAll()

	RegisterLazyFunc[Counter](Scoped, func(ctx context.Context) (Counter, context.Context) {
		return &simpleCounter{}, ctx
	})

	RegisterEagerSingleton[Counter](&simpleCounter{})

	RegisterLazyCreator[Counter](Scoped, &simpleCounter{})

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get[Counter](ctx)
	}
}

func BenchmarkGetList(b *testing.B) {
	clearAll()

	RegisterLazyFunc[Counter](Scoped, func(ctx context.Context) (Counter, context.Context) {
		return &simpleCounter{}, ctx
	})

	RegisterEagerSingleton[Counter](&simpleCounter{})

	RegisterLazyCreator[Counter](Scoped, &simpleCounter{})

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetList[Counter](ctx)
	}
}
