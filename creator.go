package ore

// RegisterLazyCreator Registers a lazily initialized value using a `Creator[T]` interface
func RegisterLazyCreator[T any](entryType RegistrationType, creator Creator[T], key ...KeyStringer) {
	if creator == nil {
		panic(nilVal[T]())
	}

	e := entry[T]{
		registrationType: entryType,
		creatorInstance:  creator,
	}
	appendToContainer[T](e, key)
}
