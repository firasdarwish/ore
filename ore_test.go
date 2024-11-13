package ore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeal(t *testing.T) {
	clearAll()
	RegisterLazyCreator[someCounter](Scoped, &simpleCounter{})
	Seal()
	assert.Panics(t, func() {
		RegisterLazyCreator[someCounter](Scoped, &simpleCounter{})
	})
}

func TestIsSeal(t *testing.T) {
	clearAll()
	RegisterLazyCreator[someCounter](Scoped, &simpleCounter{})
	if got := IsSealed(); got != false {
		t.Errorf("IsSealed() = %v; want %v", got, false)
	}

	Seal()

	if got := IsSealed(); got != true {
		t.Errorf("IsSealed() = %v; want %v", got, true)
	}
}

type A1 struct{}
type A2 struct{}

func TestTypeIdentifier(t *testing.T) {
	id1 := typeIdentifier[*A1]()
	id11 := typeIdentifier[*A1]()
	id2 := typeIdentifier[*A2]()
	assert.NotEqual(t, id1, id2)
	assert.Equal(t, id1, id11)

	id3 := typeIdentifier[*A1]("a", "b")
	id4 := typeIdentifier[*A1]("a", "b")
	assert.Equal(t, id3, id4)
}
