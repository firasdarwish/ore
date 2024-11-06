package ore

import (
	"fmt"
	"reflect"
	"time"
)

// RegisterLazyCreator Registers a lazily initialized value using a `Creator[T]` interface
func RegisterLazyCreator[T any](lifetime Lifetime, creator Creator[T], key ...KeyStringer) {
	if creator == nil {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: lifetime,
		},
		creatorInstance: creator,
	}
	appendToContainer[T](e, key)
}

// RegisterEagerSingleton Registers an eagerly instantiated singleton value
func RegisterEagerSingleton[T comparable](impl T, key ...KeyStringer) {
	if isNil[T](impl) {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: Singleton,
		},
		singletonConcrete: &concrete{
			value:     impl,
			lifetime:  Singleton,
			createdAt: time.Now(),
		},
	}
	appendToContainer[T](e, key)
}

// RegisterLazyFunc Registers a lazily initialized value using an `Initializer[T]` function signature
func RegisterLazyFunc[T any](lifetime Lifetime, initializer Initializer[T], key ...KeyStringer) {
	if initializer == nil {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: lifetime,
		},
		anonymousInitializer: &initializer,
	}
	appendToContainer[T](e, key)
}

// RegisterAlias Registers an interface type to a concrete implementation.
// Allowing you to register the concrete implementation to the container and later get the interface from it.
func RegisterAlias[TInterface, TImpl any]() {
	interfaceType := reflect.TypeFor[TInterface]()
	implType := reflect.TypeFor[TImpl]()

	if !implType.Implements(interfaceType) {
		panic(fmt.Errorf("%s does not implements %s", implType, interfaceType))
	}

	appendToAliases[TInterface, TImpl]()
}
