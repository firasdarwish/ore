package ore

import "testing"

func TestBuild(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)

	RegisterLazyCreator[Counter](Scoped, &simpleCounter{})
	Build()
	RegisterLazyCreator[Counter](Scoped, &simpleCounter{})
}
