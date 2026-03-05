package ore

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

func registerCreatorToContainer[T any, K comparable](con *Container, lifetime Lifetime, creator Creator[T], key K) {
	if creator == nil {
		panic(nilVal[T]())
	}

	var once *sync.Once
	if lifetime == Singleton {
		once = &sync.Once{}
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: lifetime,
		},
		creatorInstance: creator,
		singletonOnce:   once,
	}
	addResolver[T](con, e, key)
}

func registerSingletonToContainer[T any, K comparable](con *Container, impl T, key K) {
	var mock any
	mock = impl

	if mock == nil {
		panic(nilVal[T]())
	}

	v := reflect.ValueOf(impl)
	kind := v.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface ||
		kind == reflect.Slice || kind == reflect.Map ||
		kind == reflect.Chan || kind == reflect.Func {
		if v.IsNil() {
			panic(nilVal[T]())
		}
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: Singleton,
		},
		singletonConcrete: &concrete{
			value:          impl,
			lifetime:       Singleton,
			invocationTime: time.Now(),
		},
	}
	addResolver[T](con, e, key)
}

func registerFuncToContainer[T any, K comparable](con *Container, lifetime Lifetime, initializer Initializer[T], key K) {
	if initializer == nil {
		panic(nilVal[T]())
	}

	var once *sync.Once
	if lifetime == Singleton {
		once = &sync.Once{}
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: lifetime,
		},
		anonymousInitializer: &initializer,
		singletonOnce:        once,
	}
	addResolver[T](con, e, key)
}

func registerAliasToContainer[TInterface, TImpl any](con *Container) {
	interfaceType := reflect.TypeFor[TInterface]()
	implType := reflect.TypeFor[TImpl]()

	if !implType.Implements(interfaceType) {
		panic(fmt.Errorf("%s does not implements %s", implType, interfaceType))
	}

	addAliases[TInterface, TImpl](con)
}

func registerPlaceholderToContainer[T any, K comparable](con *Container, key K) {
	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: Scoped,
		},
	}
	addResolver[T](con, e, key)
}

func provideScopedValueToContainer[T any, K comparable](con *Container, ctx context.Context, value T, key K) context.Context {
	concreteValue := &concrete{
		value:           value,
		lifetime:        Scoped,
		invocationTime:  time.Now(),
		invocationLevel: 0,
	}
	id := contextKey{
		containerID: con.containerID,
		typeID:      typeIdentifier[T](key),
		resolverID:  placeholderResolverID,
	}
	return addScopedConcreteToContext(ctx, id, concreteValue)
}
