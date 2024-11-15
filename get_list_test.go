package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"
	"testing"
)

func TestGetList(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{})

		counters, _ := GetList[interfaces.SomeCounter](context.Background())

		if got := len(counters); got != 1 {
			t.Errorf("got %v, expected %v", got, 1)
		}
	}
}

func TestGetListShouldNotPanicIfNoImplementations(t *testing.T) {
	clearAll()
	services, _ := GetList[interfaces.SomeCounter](context.Background())
	if len(services) != 0 {
		t.Errorf("got %v, expected %v", len(services), 0)
	}
}

func TestGetListKeyed(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		key := "somekeyhere"

		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, key)
		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, key)
		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, key)
		RegisterCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, "Firas")

		counters, _ := GetList[interfaces.SomeCounter](context.Background(), key)
		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}
