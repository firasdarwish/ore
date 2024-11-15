package ore

import (
	"errors"
	"fmt"
	"reflect"
)

func noValidImplementation[T any]() error {
	return fmt.Errorf("implementation not found for type: %s", reflect.TypeFor[T]())
}

func invalidKeyType(t reflect.Type) error {
	return fmt.Errorf("cannot use type: `%s` as a key", t)
}

func nilVal[T any]() error {
	return fmt.Errorf("nil implementation for type: %s", reflect.TypeFor[T]())
}

func lifetimeMisalignment(resolver resolverMetadata, depResolver resolverMetadata) error {
	return fmt.Errorf("detected lifetime misalignment: %s depends on %s", resolver, depResolver)
}

func cyclicDependency(resolver resolverMetadata) error {
	return fmt.Errorf("detected cyclic dependency where: %s depends on itself", resolver)
}

func placeholderValueNotProvided(resolver resolverMetadata) error {
	return fmt.Errorf("no value has been provided for this placeholder: %s", resolver)
}

func typeAlreadyRegistered(typeID typeID) error {
	return fmt.Errorf("the type '%s' has already been registered (as a Resolver or as a Placeholder). Cannot override it with other Placeholder", typeID)
}

var alreadyBuilt = errors.New("services container is already sealed")
var alreadyBuiltCannotAdd = errors.New("cannot register new resolvers, container is sealed")
var nilKey = errors.New("cannot use nil key")
