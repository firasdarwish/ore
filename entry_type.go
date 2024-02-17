package ore

type ObjectType string

const (
	Singleton ObjectType = "singleton"
	Transient            = "transient"
	Scoped               = "scoped"
)
