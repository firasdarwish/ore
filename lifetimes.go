package ore

type Lifetime string

const (
	Singleton Lifetime = "singleton"
	Transient          = "transient"
	Scoped             = "scoped"
)
