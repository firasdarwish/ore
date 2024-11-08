package ore

import "time"

// concrete holds the resolved instance value and other metadata
type concrete struct {
	//the value implementation
	value any

	//invocationTime is the time when the resolver had been invoked, it is different from the "creationTime"
	//of the concrete. Eg: A calls B, then the invocationTime of A is before B, but the creationTime of A is after B
	//(because B was created before A)
	invocationTime time.Time

	//the lifetime of this concrete
	lifetime Lifetime

	//the invocation deep level, the bigger the value, the deeper it was resolved in the dependency chain
	//for example: A depends on B, B depends on C, C depends on D
	//A will have invocationLevel = 1, B = 2, C = 3, D = 4
	invocationLevel int
}
