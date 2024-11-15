package ore

import (
	"context"
	"github.com/firasdarwish/ore/internal/interfaces"
	"github.com/firasdarwish/ore/internal/models"
	"testing"
)

func TestGetKeyedListSingleton(t *testing.T) {
	clearAll()

	key := "ore"
	RegisterKeyedCreator[interfaces.SomeCounter](Singleton, &models.SimpleCounter{}, key)

	counters, _ := GetKeyedList[interfaces.SomeCounter](context.Background(), key)
	counters, _ = GetKeyedList[interfaces.SomeCounter](context.Background(), key)
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}

	if got := len(counters); got != 1 {
		t.Errorf("got %v, expected %v", got, 1)
	} else if got := counters[0].GetCount(); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}

	counters, _ = GetKeyedListFromContainer[interfaces.SomeCounter](DefaultContainer, context.Background(), key)
	counters, _ = GetKeyedListFromContainer[interfaces.SomeCounter](DefaultContainer, context.Background(), key)
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}

	if got := len(counters); got != 1 {
		t.Errorf("got %v, expected %v", got, 1)
	} else if got := counters[0].GetCount(); got != 4 {
		t.Errorf("got %v, expected %v", got, 4)
	}
}

func TestGetListSingleton(t *testing.T) {
	clearAll()

	RegisterCreator[interfaces.SomeCounter](Singleton, &models.SimpleCounter{})

	counters, _ := GetList[interfaces.SomeCounter](context.Background())
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}
	counters, _ = GetList[interfaces.SomeCounter](context.Background())
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}

	if got := len(counters); got != 1 {
		t.Errorf("got %v, expected %v", got, 1)
	} else if got := counters[0].GetCount(); got != 4 {
		t.Errorf("got %v, expected %v", got, 4)
	}

	counters, _ = GetListFromContainer[interfaces.SomeCounter](DefaultContainer, context.Background())
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}
	counters, _ = GetListFromContainer[interfaces.SomeCounter](DefaultContainer, context.Background())
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}

	if got := len(counters); got != 1 {
		t.Errorf("got %v, expected %v", got, 1)
	} else if got := counters[0].GetCount(); got != 8 {
		t.Errorf("got %v, expected %v", got, 8)
	}
}

func TestGetListScoped(t *testing.T) {
	clearAll()

	RegisterCreator[interfaces.SomeCounter](Scoped, &models.SimpleCounter{})

	ctx := context.Background()
	counters, ctx := GetList[interfaces.SomeCounter](ctx)
	counters, ctx = GetList[interfaces.SomeCounter](ctx)
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}

	if got := len(counters); got != 1 {
		t.Errorf("got %v, expected %v", got, 1)
	} else if got := counters[0].GetCount(); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}

	counters, ctx = GetListFromContainer[interfaces.SomeCounter](DefaultContainer, ctx)
	counters, ctx = GetListFromContainer[interfaces.SomeCounter](DefaultContainer, ctx)
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}

	if got := len(counters); got != 1 {
		t.Errorf("got %v, expected %v", got, 1)
	} else if got := counters[0].GetCount(); got != 4 {
		t.Errorf("got %v, expected %v", got, 4)
	}
}

func TestGetListTransient(t *testing.T) {
	clearAll()

	RegisterCreator[interfaces.SomeCounter](Transient, &models.SimpleCounter{})

	ctx := context.Background()
	counters, ctx := GetList[interfaces.SomeCounter](ctx)
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}

	counters, ctx = GetList[interfaces.SomeCounter](ctx)
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}

	if got := len(counters); got != 1 {
		t.Errorf("got %v, expected %v", got, 1)
	} else if got := counters[0].GetCount(); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
	}

	counters, ctx = GetListFromContainer[interfaces.SomeCounter](DefaultContainer, ctx)
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}
	counters, ctx = GetListFromContainer[interfaces.SomeCounter](DefaultContainer, ctx)
	if len(counters) > 0 {
		counters[0].AddOne()
		counters[0].AddOne()
	}

	if got := len(counters); got != 1 {
		t.Errorf("got %v, expected %v", got, 1)
	} else if got := counters[0].GetCount(); got != 2 {
		t.Errorf("got %v, expected %v", got, 2)
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

		RegisterKeyedCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, key)
		RegisterKeyedCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, key)
		RegisterKeyedCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, key)
		RegisterKeyedCreator[interfaces.SomeCounter](registrationType, &models.SimpleCounter{}, "Firas")

		counters, _ := GetKeyedList[interfaces.SomeCounter](context.Background(), key)
		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}
	}
}
