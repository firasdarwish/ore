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

func TestContainerName(t *testing.T) {
	if got := Name(); got != "DEFAULT" {
		t.Errorf("Name() = %v; want %v", got, "DEFAULT")
	}
}

func TestContainerSetName(t *testing.T) {
	con := NewContainer().SetName("ORE")
	if got := con.Name(); got != "ORE" {
		t.Errorf("Name() = %v; want %v", got, "ORE")
	}
}

func TestNewContainerName(t *testing.T) {
	con := NewContainer()
	if got := con.Name(); got != "" {
		t.Errorf("Name() = `%v`; want `%v`", got, "")
	}
}

func TestContainerReSetName(t *testing.T) {
	con := NewContainer().SetName("ORE")
	assert.Panics(t, func() {
		con.SetName("ORE1")
	})
}

func TestContainerID(t *testing.T) {
	if got := ContainerID(); got != 1 {
		t.Errorf("ContainerID() = %v; want 1", got)
	}
}

func TestNewContainerIDSerial(t *testing.T) {
	clearAll()
	con := NewContainer()
	if got := con.ContainerID(); got < 2 {
		t.Errorf("ContainerID() = %v; want >= 2", got)
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
