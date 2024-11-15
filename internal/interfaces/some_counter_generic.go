package interfaces

type Numeric interface {
	uint
}

type SomeCounterGeneric[T Numeric] interface {
	Add(number T)
	GetCount() T
}
