package ore

import (
	"context"
	"testing"
)

func TestGetWithAlias(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyFunc(registrationType, func(ctx context.Context) (*simpleCounterUint, context.Context) {
			return &simpleCounterUint{}, ctx
		})
		RegisterAlias[someCounterGeneric[uint], *simpleCounterUint]()

		c, _ := Get[someCounterGeneric[uint]](context.Background())

		c.Add(1)
		c.Add(1)

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestGetListWithAlias(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		for i := 0; i < 3; i++ {
			RegisterLazyFunc(registrationType, func(ctx context.Context) (*simpleCounterUint, context.Context) {
				return &simpleCounterUint{}, ctx
			})
		}

		RegisterAlias[someCounterGeneric[uint], *simpleCounterUint]()

		counters, _ := GetList[someCounterGeneric[uint]](context.Background())
		if got := len(counters); got != 3 {
			t.Errorf("got %v, expected %v", got, 3)
		}

		c := counters[1]
		c.Add(1)
		c.Add(1)

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

var _ someCounterGeneric[uint] = (*simpleCounterUint)(nil)

type simpleCounterUint struct {
	counter uint
}

func (this *simpleCounterUint) Add(number uint) {
	this.counter += number
}

func (this *simpleCounterUint) GetCount() uint {
	return this.counter
}
