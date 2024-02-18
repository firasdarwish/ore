package ore

func isNil[T comparable](impl T) bool {
	var mock T
	return impl == mock
}

// RegisterEagerSingleton Registers an eagerly instantiated singleton value
func RegisterEagerSingleton[T comparable](impl T, key ...KeyStringer) {
	if isNil[T](impl) {
		panic(nilVal[T]())
	}

	e := entry[T]{
		registrationType: Singleton,
		concrete:         &impl,
	}
	appendToContainer[T](e, key)
}
