package ore

// RegisterEagerSingleton Registers an eagerly instantiated singleton value
func RegisterEagerSingleton[T any](impl T, key ...KeyStringer) {
	if &impl == nil {
		panic(nilVal[T]())
	}

	e := entry[T]{
		registrationType: Singleton,
		concrete:         &impl,
	}
	appendToContainer[T](e, key)
}
