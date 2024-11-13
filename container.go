package ore

import (
	"context"
	"sync"
	"sync/atomic"
)

type Container struct {
	//DisableValidation is false by default, Set to true to skip validation.
	// Use case: you called the [Validate] function (either in the test pipeline or on application startup).
	// So you are confident that your registrations are good:
	//
	//   - no missing dependencies
	//   - no circular dependencies
	//   - no lifetime misalignment (a longer lifetime service depends on a shorter one).
	//
	// You don't need Ore to validate over and over again each time it creates a new concrete.
	// It's a waste of resource especially when you will need Ore to create a million of transient concretes
	// and any "pico" seconds or memory allocation matter for you.
	//
	// In this case, you can set DisableValidation = true.
	//
	// This config would impact also the [GetResolvedSingletons] and the [GetResolvedScopedInstances] functions,
	// the returning order would be no longer guaranteed.
	DisableValidation bool
	containerID       int32
	lock              *sync.RWMutex
	isSealed          bool
	resolvers         map[typeID][]serviceResolver

	//map interface type to the implementations type
	aliases map[pointerTypeName][]pointerTypeName
}

var lastContainerID atomic.Int32

func NewContainer() *Container {
	return &Container{
		containerID: lastContainerID.Add(1),
		lock:        &sync.RWMutex{},
		isSealed:    false,
		resolvers:   map[typeID][]serviceResolver{},
		aliases:     map[pointerTypeName][]pointerTypeName{},
	}
}

// Validate invokes all registered resolvers. It panics if any of them fails.
// It is recommended to call this function on application start, or in the CI/CD test pipeline
// The objective is to panic early when the container is bad configured. For eg:
//
//   - (1) Missing dependency (forget to register certain resolvers)
//   - (2) cyclic dependency
//   - (3) lifetime misalignment (a longer lifetime service depends on a shorter one).
func (this *Container) Validate() {
	if this.DisableValidation {
		panic("Validation is disabled")
	}
	ctx := context.Background()

	//provide default value for all placeHolders
	for _, resolvers := range this.resolvers {
		for _, resolver := range resolvers {
			if resolver.isPlaceHolder() {
				ctx = resolver.providePlaceHolderDefaultValue(this, ctx)
			}
		}
	}

	//invoke all resolver to detect potential registration problem
	for _, resolvers := range this.resolvers {
		for _, resolver := range resolvers {
			_, ctx = resolver.resolveService(this, ctx)
		}
	}
}

// Seal puts the container into read-only mode, preventing any further registrations.
func (this *Container) Seal() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.isSealed {
		panic(alreadyBuilt)
	}

	this.isSealed = true
}

// IsSealed checks whether the container is sealed (in readonly mode)
func (this *Container) IsSealed() bool {
	return this.isSealed
}
