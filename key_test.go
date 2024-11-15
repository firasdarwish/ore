package ore

import (
	"github.com/firasdarwish/ore/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOreKeyNil(t *testing.T) {
	k := oreKey(nil)

	if got := k; got != "n" {
		t.Errorf("got `%v`, expected `%v`", got, "n")
	}
}

func TestOreKey1String(t *testing.T) {
	k := oreKey("ore")
	expect := "sore"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKey1Int(t *testing.T) {
	k := oreKey(10)
	expect := "i10"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyUint(t *testing.T) {
	var n uint
	n = 5
	k := oreKey(n)
	expect := "ui5"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyFloat32(t *testing.T) {
	var n float32
	n = 5.751
	k := oreKey(n)
	expect := "f320x1.701062p+02"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyFloat64(t *testing.T) {
	var n float64 = 5.7519
	k := oreKey(n)
	expect := "f640x1.701f212d77319p+02"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyStruct(t *testing.T) {
	n := &models.SimpleCounter{
		Counter: 17,
	}

	assert.Panics(t, func() {
		oreKey(n)
	})
}
