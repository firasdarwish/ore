package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"
	"testing"
)

func BenchmarkRegisterLazyFunc(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterLazyFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})
	}
}

func BenchmarkRegisterLazyCreator(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterLazyCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
	}
}

func BenchmarkRegisterEagerSingleton(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})
	}
}

func BenchmarkInitialGet(b *testing.B) {
	clearAll()

	RegisterLazyFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	RegisterLazyCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})

	Seal()
	Validate()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get[interfaces.SomeCounter](ctx)
	}
}

func BenchmarkGet(b *testing.B) {
	clearAll()

	RegisterLazyFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	RegisterLazyCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
	Seal()
	Validate()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ctx = Get[interfaces.SomeCounter](ctx)
	}
}

func BenchmarkInitialGetList(b *testing.B) {
	clearAll()

	RegisterLazyFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	RegisterLazyCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
	Seal()
	Validate()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetList[interfaces.SomeCounter](ctx)
	}
}

func BenchmarkGetList(b *testing.B) {
	clearAll()

	RegisterLazyFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	RegisterEagerSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	RegisterLazyCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
	Seal()
	Validate()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ctx = GetList[interfaces.SomeCounter](ctx)
	}
}
