package ore

type Lifetime string

const (
	Singleton Lifetime = "singleton"
	Transient Lifetime = "transient"
	Scoped    Lifetime = "scoped"
)
