package ore

import (
	"context"
)

var (
	DefaultContainer = NewContainer()

	//contextKeysRepositoryID is a special context key. The value of this key is the collection of other context keys stored in the context.
	contextKeysRepositoryID specialContextKey = "__ORE_CTX_KEYS_REPO"
	//contextKeyResolversStack is a special context key. The value of this key is the [ResolversStack].
	contextKeyResolversStack specialContextKey = "__ORE_DEP_STACK"

	//placeholderResolverID is a special resolverID of every "placeholder". "placeholder" is a special resolver
	//describing a "promise" for a concrete value, which will be provided in runtime.
	placeholderResolverID = -1
)

var types = []Lifetime{Singleton, Transient, Scoped}

type contextKeysRepository = []contextKey

type KeyStringer any

type Creator[T any] interface {
	New(ctx context.Context) (T, context.Context)
}

func init() {
	DefaultContainer.SetName("DEFAULT")
}

// Generates a unique identifier for a service resolver based on type and key(s)
func getTypeID(pointerTypeName pointerTypeName, key KeyStringer) typeID {
	validateOreKeyType(key)
	return typeID{pointerTypeName, key}
}

// Generates a unique identifier for a service resolver based on type and key(s)
func typeIdentifier[T any](key KeyStringer) typeID {
	return getTypeID(getPointerTypeName[T](), key)
}

// Appends a service resolver to the container with type and key
func addResolver[T any](this *Container, resolver serviceResolverImpl[T], key KeyStringer) {
	if this.isSealed {
		panic(alreadyBuiltCannotAdd)
	}

	typeID := typeIdentifier[T](key)

	this.lock.Lock()
	defer this.lock.Unlock()

	resolverID := len(this.resolvers[typeID])
	if resolver.isPlaceholder() {
		if resolverID > 0 {
			panic(typeAlreadyRegistered(typeID))
		}
		resolverID = placeholderResolverID
	}

	resolver.id = contextKey{
		typeID:      typeID,
		containerID: this.containerID,
		resolverID:  resolverID,
	}
	this.resolvers[typeID] = append(this.resolvers[typeID], resolver)
}

func replaceResolver[T any](this *Container, resolver serviceResolverImpl[T]) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.resolvers[resolver.id.typeID][resolver.id.resolverID] = resolver
}

func addAliases[TInterface, TImpl any](this *Container) {
	originalType := getPointerTypeName[TImpl]()
	aliasType := getPointerTypeName[TInterface]()
	if originalType == aliasType {
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, ot := range this.aliases[aliasType] {
		if ot == originalType {
			return //already registered
		}
	}
	this.aliases[aliasType] = append(this.aliases[aliasType], originalType)
}

// Seal puts the DEFAULT container into read-only mode, preventing any further registrations.
func Seal() {
	DefaultContainer.Seal()
}

// IsSealed checks whether the DEFAULT container is sealed (in readonly mode)
func IsSealed() bool {
	return DefaultContainer.IsSealed()
}

// Validate invokes all registered resolvers. It panics if any of them fails.
// It is recommended to call this function on application start, or in the CI/CD test pipeline
// The objective is to panic early when the container is bad configured. For eg:
//
//   - (1) Missing dependency (forget to register certain resolvers)
//   - (2) cyclic dependency
//   - (3) lifetime misalignment (a longer lifetime service depends on a shorter one).
func Validate() {
	DefaultContainer.Validate()
}

func ContainerID() int32 {
	return DefaultContainer.containerID
}

func Name() string {
	return DefaultContainer.name
}
