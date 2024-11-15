package ore

import (
	"context"
	"fmt"
	"testing"

	m "github.com/firasdarwish/ore/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestModuleIsolation(t *testing.T) {
	for _, lifetime := range types {
		t.Run(fmt.Sprintf("Module isolation %s", lifetime), func(t *testing.T) {
			con1 := NewContainer()
			RegisterFuncToContainer(con1, lifetime, func(ctx context.Context) (*m.Trader, context.Context) {
				return &m.Trader{Name: "John"}, ctx
			})

			con2 := NewContainer()
			RegisterFuncToContainer(con2, lifetime, func(ctx context.Context) (*m.Trader, context.Context) {
				return &m.Trader{Name: "Mary"}, ctx
			})

			ctx := context.Background()

			trader1, ctx := GetFromContainer[*m.Trader](con1, ctx)
			assert.Equal(t, "John", trader1.Name)

			trader2, _ := GetFromContainer[*m.Trader](con2, ctx)
			assert.Equal(t, "Mary", trader2.Name)
		})
	}
}
