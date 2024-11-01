package ore

import (
	"context"
	"testing"
)

func TestGetList(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

		counters, _ := GetList[someCounter](context.Background())

		if got := len(counters); got != 1 {
			t.Errorf("got %v, expected %v", got, 1)
		}
	}
}

func TestGetListShouldNotPanicIfNoImplementations(t *testing.T) {
	clearAll()
	services, _ := GetList[someCounter](context.Background())
	if len(services) != 0 {
		t.Errorf("got %v, expected %v", len(services), 0)
	}
}

func TestGetListKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		key := "somekeyhere"

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{}, key)
		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{}, key)
		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{}, key)
		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{}, "Firas")

		counters, _ := GetList[someCounter](context.Background(), key)
		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}
