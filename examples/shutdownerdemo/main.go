package main

import (
	"context"
	"log"
	"sync"

	"github.com/firasdarwish/ore"
)

type shutdowner interface {
	Shutdown() error
}

type disposer interface {
	Dispose(ctx context.Context) error
}

type myGlobalRepo struct {
}

var _ shutdowner = (*myGlobalRepo)(nil)

func (*myGlobalRepo) Shutdown() error {
	log.Println("shutdown globalRepo")
	return nil
}

type myScopedRepo struct {
}

var _ disposer = (*myScopedRepo)(nil)

func (*myScopedRepo) Dispose(ctx context.Context) error {
	log.Println("dispose scopedRepo")
	return nil
}

func (*myScopedRepo) New(ctx context.Context) (*myScopedRepo, context.Context) {
	return &myScopedRepo{}, ctx
}

func main() {
	ore.RegisterEagerSingleton[*myGlobalRepo](&myGlobalRepo{})
	ore.RegisterLazyCreator(ore.Scoped, &myScopedRepo{})

	ore.Validate()

	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context is canceled when main exits

	//start a go routine that will clean up resources when the context is canceled
	go func() {
		<-ctx.Done() // Wait for the context to be canceled
		// Perform your cleanup tasks here
		disposables := ore.GetResolvedScopedInstances[disposer](ctx)
		for _, d := range disposables {
			_ = d.Dispose(ctx)
		}
		wg.Done()
	}()

	//invoke the scoped service
	_, ctx = ore.Get[*myScopedRepo](ctx)

	//cancel the context will trigger the cleanup
	cancel()
	wg.Wait() // Wait for the goroutine to finish

	shutdownables := ore.GetResolvedSingletons[shutdowner]()
	for _, s := range shutdownables {
		_ = s.Shutdown()
	}
}
