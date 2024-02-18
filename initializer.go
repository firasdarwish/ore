package ore

// RegisterLazyFunc Registers a lazily initialized value using an `Initializer[T]` function signature
func RegisterLazyFunc[T any](entryType RegistrationType, initializer Initializer[T], key ...KeyStringer) {
	if initializer == nil {
		panic(nilVal[T]())
	}

	e := entry[T]{
		registrationType:     entryType,
		anonymousInitializer: &initializer,
	}
	appendToContainer[T](e, key)
}
