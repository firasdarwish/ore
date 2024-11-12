package main

import (
	"context"
	i "examples/benchperf/internal"
	"testing"

	"github.com/firasdarwish/ore"
	"github.com/samber/do/v2"
)

var container = i.BuildContainerOre(false)
var unsafeContainer = i.BuildContainerOre(true)
var injector = i.BuildContainerDo()
var ctx = context.Background()

func Benchmark_Ore(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, ctx = ore.GetFromContainer[*i.A](container, ctx)
	}
}

func Benchmark_OreNoValidation(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, ctx = ore.GetFromContainer[*i.A](unsafeContainer, ctx)
	}
}

func Benchmark_SamberDo(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = do.MustInvoke[*i.A](injector)
	}
}
