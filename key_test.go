package ore

import (
	"testing"
)

func TestOreKeyNil(t *testing.T) {
	k := oreKey(nil)

	if got := k; got != "" {
		t.Errorf("got `%v`, expected `%v`", got, "")
	}
}

func TestOreKeyEmpty(t *testing.T) {
	k := oreKey([]KeyStringer{})

	if got := k; got != "" {
		t.Errorf("got `%v`, expected `%v`", got, "")
	}
}

func TestOreKey1String(t *testing.T) {
	k := oreKey([]KeyStringer{"ore"})
	expect := "ore"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKey2String(t *testing.T) {
	k := oreKey([]KeyStringer{"ore", "package"})
	expect := "orepackage"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKey1Int(t *testing.T) {
	k := oreKey([]KeyStringer{10})
	expect := "10"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKey2Int(t *testing.T) {
	k := oreKey([]KeyStringer{10, 30})
	expect := "1030"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyStringInt(t *testing.T) {
	k := oreKey([]KeyStringer{"ore", 97})
	expect := "ore97"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKey2StringInt(t *testing.T) {
	k := oreKey([]KeyStringer{"ore", 97, "di", 5})
	expect := "ore97di5"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyUint(t *testing.T) {
	var n uint
	n = 5
	k := oreKey([]KeyStringer{n})
	expect := "5"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyFloat32(t *testing.T) {
	var n float32
	n = 5.751
	k := oreKey([]KeyStringer{n})
	expect := "0x1.701062p+02"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyFloat64(t *testing.T) {
	var n float64
	n = 5.7519
	k := oreKey([]KeyStringer{n})
	expect := "0x1.701f212d77319p+02"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyStringer(t *testing.T) {
	n := &c{
		Counter: 16,
	}

	k := oreKey([]KeyStringer{n})
	expect := "Counter is: 16"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyStruct(t *testing.T) {
	n := &simpleCounter{
		counter: 17,
	}

	k := oreKey([]KeyStringer{n})
	expect := "&{17}"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}

func TestOreKeyVarious(t *testing.T) {
	k := oreKey([]KeyStringer{"firas", 16, "ore", 3.14, 1 / 6, -9, -1494546.452, &simpleCounter{
		counter: 17,
	}, &c{
		Counter: 18,
	}})
	expect := "firas16ore0x1.91eb851eb851fp+010-9-0x1.6ce1273b645a2p+20&{17}Counter is: 18"

	if got := k; got != expect {
		t.Errorf("got `%v`, expected `%v`", got, expect)
	}
}
