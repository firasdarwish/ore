package models

import "fmt"

type Disposer interface {
	fmt.Stringer
	Dispose()
}

var _ Disposer = (*DisposableService1)(nil)

type DisposableService1 struct {
	Name string
}

func (*DisposableService1) Dispose() {}
func (this *DisposableService1) String() string {
	return this.Name
}

var _ Disposer = (*DisposableService2)(nil)

type DisposableService2 struct {
	Name string
}

func (*DisposableService2) Dispose() {}
func (this *DisposableService2) String() string {
	return this.Name
}

var _ Disposer = (*DisposableService3)(nil)

type DisposableService3 struct {
	Name string
}

func (*DisposableService3) Dispose() {}
func (this *DisposableService3) String() string {
	return this.Name
}

var _ Disposer = (*DisposableService4)(nil)

type DisposableService4 struct {
	Name string
}

func (*DisposableService4) Dispose() {}
func (this *DisposableService4) String() string {
	return this.Name
}
