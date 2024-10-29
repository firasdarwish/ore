package ore

import (
	"context"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})

		c, _ := Get[someCounter](context.Background())

		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}

func TestGetLatestByDefault(t *testing.T) {
	for _, registrationType := range types {
		clearAll()

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{})
		c, _ := Get[someCounter](context.Background())
		c.AddOne()
		c.AddOne()

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter2{})
		c, _ = Get[someCounter](context.Background())
		c.AddOne()
		c.AddOne()
		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 4 {
			t.Errorf("got %v, expected %v", got, 4)
		}
	}
}

func TestGetPanicIfNoImplementations(t *testing.T) {
	clearAll()
	defer mustHavePanicked(t)
	Get[someCounter](context.Background())
}

func TestGetKeyed(t *testing.T) {
	for i, registrationType := range types {
		clearAll()

		key := fmt.Sprintf("keynum: %v", i)

		RegisterLazyCreator[someCounter](registrationType, &simpleCounter{}, key)

		c, _ := Get[someCounter](context.Background(), key)

		c.AddOne()
		c.AddOne()

		if got := c.GetCount(); got != 2 {
			t.Errorf("got %v, expected %v", got, 2)
		}
	}
}
