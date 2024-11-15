package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"
	"testing"
)

func BenchmarkRegisterFunc(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
			return &models.SimpleCounter{}, ctx
		})
	}
}

func BenchmarkRegisterCreator(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
	}
}

func BenchmarkRegisterEagerSingleton(b *testing.B) {
	clearAll()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})
	}
}

func BenchmarkInitialGet(b *testing.B) {
	clearAll()

	RegisterFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})

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

	RegisterFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
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

	RegisterFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
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

	RegisterFunc[interfaces.SomeCounter](Scoped, func(ctx context.Context) (interfaces.SomeCounter, context.Context) {
		return &models.SimpleCounter{}, ctx
	})

	RegisterSingleton[interfaces.SomeCounter](&models.SimpleCounter{})

	RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
	Seal()
	Validate()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ctx = GetList[interfaces.SomeCounter](ctx)
	}
}
