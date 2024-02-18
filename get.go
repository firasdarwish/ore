package ore

import (
	"context"
	"fmt"
)

// Get Retrieves an instance based on type and key (throws panic for missing or invalid implementations)
func Get[T any](ctx context.Context, key ...KeyStringer) (T, context.Context) {
	tId := typeId[T](key)

	lock.RLock()
	o, ok := container[tId]
	lock.RUnlock()

	if !ok {
		panic(noValidImplementation[T]())
	}

	oLen := len(o)

	if oLen == 0 {
		panic(noValidImplementation[T]())
	}

	index := oLen - 1

	i, okInst := o[index].(entry[T])
	if !okInst {
		panic(noValidImplementation[T]())
	}

	ctxValId := fmt.Sprintln(tId, index)

	con, ctx := i.load(ctx, ctxValId)
	if i.registrationType == Singleton {
		replace[T](tId, index, i)
	}

	return con, ctx
}
