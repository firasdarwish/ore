package internal

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/samber/do/v2"
)

var counter uint64
var countIdEnabled = false

func ResetCounter() {
	atomic.StoreUint64(&counter, 0)
}

func generateId(prefix string) string {
	if countIdEnabled {
		return fmt.Sprintf("%s-%02d", prefix, atomic.AddUint64(&counter, 1))
	}
	return prefix
}

type A struct {
	id string
	b  *B `do:""`
	c  *C `do:""`
}

func NewA(b *B, c *C) *A {
	return &A{id: generateId("A"), b: b, c: c}
}
func (this *A) ToString() string {
	return fmt.Sprintf("%s { %s, %s }", this.id, this.b.ToString(), this.c.ToString())
}

type B struct {
	id string
	d  *D `do:""`
	e  *E `do:""`
}

func NewB(d *D, e *E) *B {
	return &B{id: generateId("B"), d: d, e: e}
}
func (this *B) ToString() string {
	return fmt.Sprintf("%s { %s, %s }", this.id, this.d.ToString(), this.e.ToString())
}

type C struct {
	id string
}

func NewC() *C {
	return &C{id: generateId("C")}
}

func (this *C) ToString() string {
	return this.id
}

type D struct {
	id string
	f  *F `do:""`
	h  H  `do:""`
}

func NewD(f *F, h H) *D {
	return &D{id: generateId("D"), f: f, h: h}
}
func (this *D) ToString() string {
	return fmt.Sprintf("%s { %s, %s }", this.id, this.f.ToString(), this.h.ToString())
}

var _ do.Shutdowner = (*D)(nil)

func (this *D) Shutdown() {
	if countIdEnabled {
		log.Println("Shutdown " + this.id)
	}
}

type E struct {
	id string
	g  []G `do:""`
}

func NewE(g []G) *E {
	return &E{id: generateId("E"), g: g}
}
func (this *E) ToString() string {
	resu := this.id + "{ "
	for _, gItem := range this.g {
		resu += gItem.ToString() + ", "
	}
	resu += " }"
	return resu
}

var _ do.Shutdowner = (*E)(nil)

func (this *E) Shutdown() {
	if countIdEnabled {
		log.Println("Shutdown " + this.id)
	}
}

type F struct {
	id string
}

func NewF() *F {
	return &F{id: generateId("F")}
}
func (this *F) ToString() string {
	return this.id
}

type G interface {
	GetId() string
	ToString() string
}

var _ G = (*Ga)(nil)

type Ga struct {
	id string
}

func NewGa() *Ga {
	return &Ga{id: generateId("Ga")}
}
func (this *Ga) GetId() string {
	return this.id
}
func (this *Ga) ToString() string {
	return this.id
}

var _ do.Shutdowner = (*Ga)(nil)

func (this *Ga) Shutdown() {
	if countIdEnabled {
		log.Println("Shutdown " + this.id)
	}
}

var _ G = (*Gb)(nil)

type Gb struct {
	id string
}

func NewGb() *Gb {
	return &Gb{id: generateId("Gb")}
}
func (this *Gb) ToString() string {
	return this.id
}

func (this *Gb) GetId() string {
	return this.id
}

var _ G = (*Gc)(nil)

type Gc struct {
	id string
}

func NewGc() *Gc {
	return &Gc{id: generateId("Gc")}
}
func (this *Gc) ToString() string {
	return this.id
}

func (this *Gc) GetId() string {
	return this.id
}

var _ G = (*DGa)(nil)

type DGa struct {
	core *Ga `do:""`
	id   string
}

func NewDGa(core *Ga) *DGa {
	return &DGa{core: core, id: generateId("DGa")}
}
func (this *DGa) ToString() string {
	return fmt.Sprintf("%s { %s }", this.id, this.core.ToString())
}

func (this *DGa) GetId() string {
	return this.id
}

type H interface {
	do.Shutdowner
	ToString() string
}

type Hr struct {
	id string
}

var _ H = (*Hr)(nil)

func NewHr() *Hr {
	return &Hr{id: generateId("Hr")}
}
func (this *Hr) ToString() string {
	return this.id
}
func (this *Hr) Shutdown() {
	if countIdEnabled {
		log.Println("Shutdown " + this.id)
	}
}

type Hm struct {
	id string
}

var _ H = (*Hm)(nil)

func NewHm() *Hm {
	return &Hm{id: generateId("Hm")}
}
func (this *Hm) ToString() string {
	return this.id
}
func (this *Hm) Shutdown() {
	if countIdEnabled {
		log.Println("Shutdown " + this.id)
	}
}
