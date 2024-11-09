package ore

type Lifetime int

// The bigger the value, the longer the lifetime
const (
	Transient Lifetime = 0
	Scoped    Lifetime = 1
	Singleton Lifetime = 2
)

func (this Lifetime) String() string {
	switch this {
	case 0:
		return "Transient"
	case 1:
		return "Scoped"
	case 2:
		return "Singleton"
	default:
		return "Unknow"
	}
}
