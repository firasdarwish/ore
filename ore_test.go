package ore

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NotEqual(t, id1, id2)

	id3 := typeIdentifier[*A1]([]KeyStringer{"a", "b"})
	id4 := typeIdentifier[*A1]([]KeyStringer{"a", "b"})
	assert.Equal(t, id3, id4)
}
