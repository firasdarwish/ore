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

var _ Disposer = (*DisposableService5)(nil)

type DisposableService5 struct {
	Name string
}

func (*DisposableService5) Dispose() {}
func (this *DisposableService5) String() string {
	return this.Name
}

func FindIndexOf(disposables []Disposer, name string) int {
	for i, disposable := range disposables {
		if disposable.String() == name {
			return i
		}
	}
	return -1
}
