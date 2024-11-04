package ore

import (
	"errors"
	"fmt"
	"reflect"
)

func noValidImplementation[T any]() error {
	return fmt.Errorf("implementation not found for type: %s", reflect.TypeFor[T]())
}

func nilVal[T any]() error {
	return fmt.Errorf("nil implementation for type: %s", reflect.TypeFor[T]())
}

var alreadyBuilt = errors.New("services container is already built")
var alreadyBuiltCannotAdd = errors.New("cannot appendToContainer, services container is already built")
var nilKey = errors.New("cannot have nil keys")
