package ore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOreKeyNil(t *testing.T) {
	k := oreKey(nil)

	if got := k; got != "n" {
		t.Errorf("got `%v`, expected `%v`", got, "n")
	}
}

func TestOreKeyEmpty(t *testing.T) {
	k := oreKey()

	if got := k; got != "" {
		t.Errorf("got `%v`, expected `%v`", got, "s")
	}
}

func TestOreKey1String(t *testing.T) {
	k := oreKey("ore")
	expect := "sore"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKey2String(t *testing.T) {
	k := oreKey("ore", "package")
	expect := "sorespackage"

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

func TestOreKey2Int(t *testing.T) {
	k := oreKey(10, 30)
	expect := "i10i30"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyStringInt(t *testing.T) {
	k := oreKey("ore", 97)
	expect := "sorei97"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKey2StringInt(t *testing.T) {
	k := oreKey("ore", 97, "di", 5)
	expect := "sorei97sdii5"

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

func TestOreKeyStringer(t *testing.T) {
	n := &c{
		Counter: 16,
	}

	assert.Panics(t, func() {
		oreKey(n)
	})
}

func TestOreKeyStruct(t *testing.T) {
	n := &simpleCounter{
		counter: 17,
	}

	assert.Panics(t, func() {
		oreKey(n)
	})
}

func TestOreKeyVarious(t *testing.T) {
	k := oreKey("firas", 16, "ore", 3.14, 1/6, -9, -1494546.452)
	expect := "sfirasi16soref640x1.91eb851eb851fp+01i0i-9f64-0x1.6ce1273b645a2p+20"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}
