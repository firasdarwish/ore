package ore

import (
	"testing"
)

func TestBuild(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)

	RegisterLazyCreator[someCounter](Scoped, &simpleCounter{})
	Build()
	RegisterLazyCreator[someCounter](Scoped, &simpleCounter{})
}

type A1 struct{}
type A2 struct{}

func TestTypeIdentifier(t *testing.T) {
	id1 := typeIdentifier[*A1]([]KeyStringer{})
	id2 := typeIdentifier[*A2]([]KeyStringer{})
	if id1 == id2 {
		t.Errorf("got the same identifier value %v, expected different values", id1)
	}

	id3 := typeIdentifier[*A1]([]KeyStringer{"a", "b"})
	id4 := typeIdentifier[*A1]([]KeyStringer{"a", "b"})
	if id3 != id4 {
		t.Errorf("got %v, expected %v", id3, id4)
	}
}
