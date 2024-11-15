package ore

import (
	"container/list"
	"context"
	"fmt"
	"time"
)

type (
	Initializer[T any] func(ctx context.Context) (T, context.Context)
)

type serviceResolver interface {
	fmt.Stringer
	resolveService(ctn *Container, ctx context.Context) (*concrete, context.Context)

	//getInvokedSingleton returns the invoked singleton value, or false if the resolver is not a singleton or has not been invoked
	getInvokedSingleton() (con *concrete, isInvokedSingleton bool)

	//isPlaceHolder returns true if this resolver is a placeholder.
	//A placeholder is a special Scoped resolver that doesn't have Creator or Factory (a.k.a anonymousInitializer) function
	isPlaceHolder() bool

	//providePlaceHolderDefaultValue provides a default value for a placeholder for validation
	providePlaceHolderDefaultValue(ctn *Container, ctx context.Context) context.Context

	// isScopedValueResolved returns true if this resolver is a scoped resolver and the scoped value has been already resolved.
	// in case this resolver is a placeholder, then it returns true if the placeholder value has been provided.
	isScopedValueResolved(ctx context.Context) bool
}

type resolverMetadata struct {
	id       contextKey
	lifetime Lifetime
}

type serviceResolverImpl[T any] struct {
	resolverMetadata
	anonymousInitializer *Initializer[T]
	creatorInstance      Creator[T]
	singletonConcrete    *concrete
}

// resolversStack is a stack of [resolverMetadata], similar to a call stack describing How a resolver has
// to call other resolvers to resolve its dependencies.
// Before a resolver creates a new concrete value it would be registered (pushed) to the stack.
// Once the concrete is resolved (with help of other resolvers), then it would be removed (popped) from the stack.
//
// While a Resolver forms a tree with other dependent resolvers.
//
// Example:
//
//	A calls B and C; B calls D; C calls E.
//
// then resolversStack holds a "path" in the tree from the root to one of the bottom.
//
// Example:
//
//	A -> B -> D or A -> C -> E
//
// The resolversStack is stored in the context. Analyze the stack will help to
//
//   - (1) detect the invocation level
//   - (2) detect cyclic dependencies
//   - (3) detect lifetime misalignment (when a service of longer lifetime depends on a service of shorter lifetime)
type resolversStack = *list.List

// make sure that the `serviceResolverImpl` struct implements the `serviceResolver` interface
var _ serviceResolver = serviceResolverImpl[any]{}

func (this serviceResolverImpl[T]) resolveService(ctn *Container, ctx context.Context) (*concrete, context.Context) {
	// try get concrete implementation
	if this.lifetime == Singleton && this.singletonConcrete != nil {
		return this.singletonConcrete, ctx
	}

	// get the currentStack from the context
	var currentStack resolversStack
	if !ctn.DisableValidation {
		untypedCurrentStack := ctx.Value(contextKeyResolversStack)
		if untypedCurrentStack != nil {
			currentStack = untypedCurrentStack.(resolversStack)
		}
		validateLifetime(currentStack, this.resolverMetadata)
	}

	// try to get concrete from context scope
	if this.lifetime == Scoped {
		scopedConcrete, ok := ctx.Value(this.id).(*concrete)
		if ok {
			return scopedConcrete, ctx
		}
	}

	if this.isPlaceHolder() {
		panic(placeholderValueNotProvided(this.resolverMetadata))
	}

	// this resolver is about to create a new concrete value, we have to put it to the resolversStack until the creation done

	var marker *list.Element
	if !ctn.DisableValidation {
		if currentStack == nil {
			currentStack = list.New()
			ctx = context.WithValue(ctx, contextKeyResolversStack, currentStack)
		}
		// push the current resolver to the resolversStack
		marker = pushToStack(currentStack, this.resolverMetadata)
	}
	var concreteValue T
	invocationTime := time.Now()

	// first, try to make concrete implementation from `anonymousInitializer`
	// if nil, try the concrete implementation `Creator`
	if this.anonymousInitializer != nil {
		concreteValue, ctx = (*this.anonymousInitializer)(ctx)
	} else {
		concreteValue, ctx = this.creatorInstance.New(ctx)
	}

	invocationLevel := 0
	if !ctn.DisableValidation {
		invocationLevel = currentStack.Len()

		//the concreteValue is created, we must pop the current resolvers from the stack
		//so that future resolvers won't link to it
		currentStack.Remove(marker)
	}

	con := &concrete{
		value:           concreteValue,
		lifetime:        this.lifetime,
		invocationTime:  invocationTime,
		invocationLevel: invocationLevel,
	}

	// if scoped, attach to the current context
	if this.lifetime == Scoped {
		ctx = addScopedConcreteToContext(ctx, this.id, con)
	}

	// if was lazily-created, then attach the newly-created concrete implementation
	// to the service resolver
	if this.lifetime == Singleton {
		this.singletonConcrete = con
		replaceResolver(ctn, this)
		return con, ctx
	}

	return con, ctx
}

// pushToStack appends the given resolver to the Back of the given resolversStack.
// `marker.previous` refers to the calling (parent) resolver
func pushToStack(stack resolversStack, currentResolver resolverMetadata) (marker *list.Element) {
	if stack.Len() != 0 {
		//detect cyclic dependencies
		for e := stack.Back(); e != nil; e = e.Prev() {
			if e.Value.(resolverMetadata).id == currentResolver.id {
				panic(cyclicDependency(currentResolver))
			}
		}
	}
	marker = stack.PushBack(currentResolver) // `marker.previous` refers to the calling (parent) resolver
	return marker
}

// validateLifetime checks the lifetime of the current resolver must not shorter
// than the lifetime of the caller (a.k.a parent) resolver
func validateLifetime(stack resolversStack, currentResolver resolverMetadata) {
	if stack == nil || stack.Len() == 0 {
		return
	}
	//detect lifetime misalignment
	lastElem := stack.Back()
	lastResolver := lastElem.Value.(resolverMetadata)
	if lastResolver.lifetime > currentResolver.lifetime {
		panic(lifetimeMisalignment(lastResolver, currentResolver))
	}
}

func (this serviceResolverImpl[T]) getInvokedSingleton() (con *concrete, isInvokedSingleton bool) {
	if this.lifetime == Singleton && this.singletonConcrete != nil {
		return this.singletonConcrete, true
	}
	return nil, false
}

func addToContextKeysRepository(ctx context.Context, newContextKey contextKey) context.Context {
	repository, ok := ctx.Value(contextKeysRepositoryID).(contextKeysRepository)
	if ok {
		repository = append(repository, newContextKey)
	} else {
		repository = contextKeysRepository{newContextKey}
	}
	return context.WithValue(ctx, contextKeysRepositoryID, repository)
}

func (this resolverMetadata) String() string {
	return fmt.Sprintf("Resolver(%s, type={%s}, key='%s')", this.lifetime, getUnderlyingTypeName(this.id.pointerTypeName), this.id.oreKey)
}

func (this serviceResolverImpl[T]) isPlaceHolder() bool {
	return this.lifetime == Scoped && this.anonymousInitializer == nil && this.creatorInstance == nil
}

func (this serviceResolverImpl[T]) providePlaceHolderDefaultValue(ctn *Container, ctx context.Context) context.Context {
	defaultValue := *new(T)
	concreteValue := &concrete{
		value:           defaultValue,
		lifetime:        Scoped,
		invocationTime:  time.Now(),
		invocationLevel: 0,
	}
	return addScopedConcreteToContext(ctx, this.id, concreteValue)
}

// isScopedValueResolved returns true if the scoped value has been already resolved.
// we need this to know if the placeholder value has been provided?
func (this serviceResolverImpl[T]) isScopedValueResolved(ctx context.Context) bool {
	return ctx.Value(this.id) != nil
}

func addScopedConcreteToContext(ctx context.Context, id contextKey, concrete *concrete) context.Context {
	ctx = context.WithValue(ctx, id, concrete)
	ctx = addToContextKeysRepository(ctx, id)
	return ctx
}

// func toString(resolversStack resolversStack) string {
// 	var sb string
// 	for e := resolversStack.Front(); e != nil; e = e.Next() {
// 		sb = fmt.Sprintf("%s%s\n", sb, e.Value.(resolverMetadata).String())
// 	}
// 	return sb
// }
