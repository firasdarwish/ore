package main

import (
	"context"
	i "examples/benchperf/internal"
	"log"

	"github.com/firasdarwish/ore"
	"github.com/samber/do/v2"
)

func main() {
	i.RegisterToOreContainer(ore.DefaultContainer)
	a1, _ := ore.Get[*i.A](context.Background())
	log.Println(a1.ToString())
	a2, _ := ore.Get[*i.A](context.Background())
	log.Println(a2.ToString())

	i.ResetCounter()

	injector := i.BuildContainerDo()
	a3 := do.MustInvoke[*i.A](injector)
	log.Println(a3.ToString())
	a4 := do.MustInvoke[*i.A](injector)
	log.Println(a4.ToString())
}
