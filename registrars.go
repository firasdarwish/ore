package ore

// RegisterLazyCreator Registers a lazily initialized value using a `Creator[T]` interface
func RegisterLazyCreator[T any](lifetime Lifetime, creator Creator[T], key ...KeyStringer) {
	if creator == nil {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		lifetime:        lifetime,
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
		lifetime:          Singleton,
		singletonConcrete: &impl,
	}
	appendToContainer[T](e, key)
}

// RegisterLazyFunc Registers a lazily initialized value using an `Initializer[T]` function signature
func RegisterLazyFunc[T any](lifetime Lifetime, initializer Initializer[T], key ...KeyStringer) {
	if initializer == nil {
		panic(nilVal[T]())
	}

	e := serviceResolverImpl[T]{
		lifetime:             lifetime,
		anonymousInitializer: &initializer,
	}
	appendToContainer[T](e, key)
}
