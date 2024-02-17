package ore

type RegistrationType string

const (
	Singleton RegistrationType = "singleton"
	Transient                  = "transient"
	Scoped                     = "scoped"
)
