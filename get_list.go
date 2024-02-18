package ore

import (
	"context"
	"fmt"
)

// GetList Retrieves a list of instances based on type and key (throws panic for missing or invalid implementations)
func GetList[T any](ctx context.Context, key ...KeyStringer) ([]T, context.Context) {
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

	arrr := make([]T, oLen)

	for index := 0; index < oLen; index++ {
		i, okInst := o[index].(entry[T])
		if !okInst {
			panic(noValidImplementation[T]())
		}

		ctxValId := fmt.Sprintln(tId, index)
		connn, ctxxx := i.load(ctx, ctxValId)

		if i.registrationType == Singleton {
			replace[T](tId, index, i)
		}

		arrr[index] = connn
		ctx = ctxxx
	}

	return arrr, ctx
}
