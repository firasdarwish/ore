package ore

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

func registerCreatorToContainer[T any](con *Container, lifetime Lifetime, creator Creator[T], key KeyStringer) {
	if creator == nil {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: lifetime,
		},
		creatorInstance: creator,
	}
	addResolver[T](con, e, key)
}

func registerSingletonToContainer[T any](con *Container, impl T, key KeyStringer) {
	var mock any
	mock = impl

	if mock == nil {
		panic(nilVal[T]())
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

func registerFuncToContainer[T any](con *Container, lifetime Lifetime, initializer Initializer[T], key KeyStringer) {
	if initializer == nil {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: lifetime,
		},
		anonymousInitializer: &initializer,
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

func registerPlaceholderToContainer[T any](con *Container, key KeyStringer) {
	e := serviceResolverImpl[T]{
		resolverMetadata: resolverMetadata{
			lifetime: Scoped,
		},
	}
	addResolver[T](con, e, key)
}

func provideScopedValueToContainer[T any](con *Container, ctx context.Context, value T, key KeyStringer) context.Context {
	concreteValue := &concrete{
		value:           value,
		lifetime:        Scoped,
		invocationTime:  time.Now(),
		invocationLevel: 0,
	}
	id := contextKey{
		containerID: con.containerID,
		typeID:      typeIdentifier[T](key),
		resolverID:  placeHolderResolverID,
	}
	return addScopedConcreteToContext(ctx, id, concreteValue)
}
