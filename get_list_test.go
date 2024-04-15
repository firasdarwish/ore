package ore

import (
	"context"
	"testing"
)

func TestGetList(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[Counter](registrationType, &simpleCounter{})

		counters, _ := GetList[Counter](context.Background())

		if got := len(counters); got != 1 {
			t.Errorf("got %v, expected %v", got, 1)
		}
	}
}

func TestGetListShouldNotPanicIfNoImplementations(t *testing.T) {
	clearAll()
	services, _ := GetList[Counter](context.Background())
	if len(services) != 0 {
		t.Errorf("got %v, expected %v", len(services), 0)
	}
}

func TestGetListKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		key := "somekeyhere"

		RegisterLazyCreator[Counter](registrationType, &simpleCounter{}, key)
		RegisterLazyCreator[Counter](registrationType, &simpleCounter{}, key)
		RegisterLazyCreator[Counter](registrationType, &simpleCounter{}, key)
		RegisterLazyCreator[Counter](registrationType, &simpleCounter{}, "Firas")

		counters, _ := GetList[Counter](context.Background(), key)
		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}
