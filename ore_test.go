package ore

import (
	"testing"

	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestSeal(t *testing.T) {
	clearAll()
	RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
	Seal()
	assert.Panics(t, func() {
		RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
	})
}

func TestIsSeal(t *testing.T) {
	clearAll()
	RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})
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

func TestTypeIdentifierNilkey(t *testing.T) {
	id1 := typeIdentifier[*A1](nilKey)
	id2 := typeIdentifier[*A1](0)
	//nilKey looks like 0 but it is not equal 0
	assert.NotEqual(t, id1, id2)

	id21 := typeIdentifier[*A1](0)
	assert.Equal(t, id2, id21)

	id3 := typeIdentifier[*A1]("a")
	id4 := typeIdentifier[*A1]("a")
	assert.Equal(t, id3, id4)
}

func TestTypeIdentifierComplexKey(t *testing.T) {
	key1 := contextKey{
		typeID: typeID{
			pointerTypeName: "toto",
			oreKey:          models.Trader{Name: "Hiep"},
		},
		containerID: 1,
		resolverID:  2,
	}
	key2 := contextKey{
		typeID: typeID{
			pointerTypeName: "toto",
			oreKey:          models.Trader{Name: "Hiep"},
		},
		containerID: 1,
		resolverID:  2,
	}

	id1 := typeIdentifier[A1](key1)
	id2 := typeIdentifier[A1](key2)

	assert.Equal(t, id1, id2)
}
