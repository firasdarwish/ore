package ore

// concrete holds the resolved instance value and other metadata
type concrete struct {
	//the value implementation
	value any

	//invocationOrder is the order when the resolver had been invoked, it is the opposite of the "creationTime"
	//of the concrete. Eg: A calls B, then the invocationOrder is [A -> B], but the creationTime/order is [B -> A]
	//(because B was created before A)
	invocationOrder uint32

	//the lifetime of this concrete
	lifetime Lifetime

	//the invocation depth level, the bigger the value, the deeper it was resolved in the dependency chain
	//for example: A depends on B, B depends on C, C depends on D
	//A will have invocationLevel = 1, B = 2, C = 3, D = 4
	invocationLevel int
}
