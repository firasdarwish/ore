package ore

// RegisterAlias Registers an interface type to a concrete implementation.
// Allowing you to register the concrete implementation to the default container and later get the interface from it.
func RegisterAlias[TInterface, TImpl any]() {
	registerAliasToContainer[TInterface, TImpl](DefaultContainer)
}

// RegisterAliasToContainer Registers an interface type to a concrete implementation in the given container.
// Allowing you to register the concrete implementation to the container and later get the interface from it.
func RegisterAliasToContainer[TInterface, TImpl any](con *Container) {
	registerAliasToContainer[TInterface, TImpl](con)
}
