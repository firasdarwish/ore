package main

import (
	"context"
	i "examples/benchperf/internal"
	"testing"

	"github.com/firasdarwish/ore"
	"github.com/samber/do/v2"
)

// func Benchmark_Ore_NoValidation(b *testing.B) {
// 	i.BuildContainerOre()
// 	ore.DisableValidation = true
// 	ctx := context.Background()
// 	b.ResetTimer()
// 	for n := 0; n < b.N; n++ {
// 		_, ctx = ore.Get[*i.A](ctx)
// 	}
// }

var _ = i.BuildContainerOre()
var injector = i.BuildContainerDo()
var ctx = context.Background()

func Benchmark_Ore(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, ctx = ore.Get[*i.A](ctx)
	}
}

func Benchmark_SamberDo(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = do.MustInvoke[*i.A](injector)
	}
}
