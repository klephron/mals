package execution

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Name  string
	Count int
	Inner innerStruct
}

type innerStruct struct {
	Value float64
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func assertEqual(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func newMem(pairs ...any) map[string]any {
	m := map[string]any{}
	for i := 0; i < len(pairs)-1; i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestDerefValue(t *testing.T) {
	t.Run("plain int", func(t *testing.T) {
		v := derefValue(reflect.ValueOf(42))
		assertEqual(t, v.Kind(), reflect.Int)
	})

	t.Run("pointer to int", func(t *testing.T) {
		n := 7
		v := derefValue(reflect.ValueOf(&n))
		assertEqual(t, v.Kind(), reflect.Int)
		assertEqual(t, int(v.Int()), 7)
	})

	t.Run("pointer to pointer to string", func(t *testing.T) {
		s := "hello"
		ps := &s
		v := derefValue(reflect.ValueOf(&ps))
		assertEqual(t, v.Kind(), reflect.String)
	})

	t.Run("interface wrapping int", func(t *testing.T) {
		var i any = 99
		v := derefValue(reflect.ValueOf(&i).Elem()) // interface Value
		assertEqual(t, v.Kind(), reflect.Int)
	})
}

func TestParseIndex(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		idx, err := parseIndex("2", 5)
		assertNoErr(t, err)
		assertEqual(t, idx, 2)
	})

	t.Run("zero", func(t *testing.T) {
		idx, err := parseIndex("0", 1)
		assertNoErr(t, err)
		assertEqual(t, idx, 0)
	})

	t.Run("last element", func(t *testing.T) {
		idx, err := parseIndex("4", 5)
		assertNoErr(t, err)
		assertEqual(t, idx, 4)
	})

	t.Run("non-integer", func(t *testing.T) {
		_, err := parseIndex("foo", 5)
		assertErr(t, err)
	})

	t.Run("negative", func(t *testing.T) {
		_, err := parseIndex("-1", 5)
		assertErr(t, err)
	})

	t.Run("equal to length", func(t *testing.T) {
		_, err := parseIndex("5", 5)
		assertErr(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := parseIndex("", 5)
		assertErr(t, err)
	})
}

func TestStructField(t *testing.T) {
	v := reflect.ValueOf(testStruct{Name: "x", Count: 3})

	t.Run("exact match", func(t *testing.T) {
		fv, err := structField(v, "Name")
		assertNoErr(t, err)
		assertEqual(t, fv.String(), "x")
	})

	t.Run("second field", func(t *testing.T) {
		fv, err := structField(v, "Count")
		assertNoErr(t, err)
		assertEqual(t, int(fv.Int()), 3)
	})

	t.Run("missing field", func(t *testing.T) {
		_, err := structField(v, "Missing")
		assertErr(t, err)
	})

	t.Run("case mismatch rejected", func(t *testing.T) {
		_, err := structField(v, "name")
		assertErr(t, err)
	})

	t.Run("case mismatch upper", func(t *testing.T) {
		_, err := structField(v, "NAME")
		assertErr(t, err)
	})
}

func TestAssignValue(t *testing.T) {
	t.Run("string to string", func(t *testing.T) {
		var s string
		dst := reflect.ValueOf(&s).Elem()
		assertNoErr(t, assignValue(dst, "hello"))
		assertEqual(t, s, "hello")
	})

	t.Run("int32 converted to int64", func(t *testing.T) {
		var n int64
		dst := reflect.ValueOf(&n).Elem()
		assertNoErr(t, assignValue(dst, int32(5)))
		assertEqual(t, n, int64(5))
	})

	t.Run("nil zeros string", func(t *testing.T) {
		s := "nonempty"
		dst := reflect.ValueOf(&s).Elem()
		assertNoErr(t, assignValue(dst, nil))
		assertEqual(t, s, "")
	})

	t.Run("nil zeros int", func(t *testing.T) {
		n := 42
		dst := reflect.ValueOf(&n).Elem()
		assertNoErr(t, assignValue(dst, nil))
		assertEqual(t, n, 0)
	})

	t.Run("nil zeros pointer", func(t *testing.T) {
		x := 1
		p := &x
		dst := reflect.ValueOf(&p).Elem()
		assertNoErr(t, assignValue(dst, nil))
		if p != nil {
			t.Fatal("expected nil pointer")
		}
	})

	t.Run("concrete into interface field", func(t *testing.T) {
		var i any
		dst := reflect.ValueOf(&i).Elem()
		assertNoErr(t, assignValue(dst, "concrete"))
		assertEqual(t, i, "concrete")
	})

	t.Run("type mismatch", func(t *testing.T) {
		var n int
		dst := reflect.ValueOf(&n).Elem()
		err := assignValue(dst, "not an int")
		assertErr(t, err)
	})
}

func TestSetMapValue(t *testing.T) {
	t.Run("string any map — set string", func(t *testing.T) {
		m := map[string]any{"k": "old"}
		mv := reflect.ValueOf(m)
		key := reflect.ValueOf("k")
		assertNoErr(t, setMapValue(mv, key, "new"))
		assertEqual(t, m["k"], "new")
	})

	t.Run("string any map — set int", func(t *testing.T) {
		m := map[string]any{}
		mv := reflect.ValueOf(m)
		key := reflect.ValueOf("n")
		assertNoErr(t, setMapValue(mv, key, 42))
		assertEqual(t, m["n"], 42)
	})

	t.Run("nil deletes key", func(t *testing.T) {
		m := map[string]any{"del": "gone"}
		mv := reflect.ValueOf(m)
		key := reflect.ValueOf("del")
		assertNoErr(t, setMapValue(mv, key, nil))
		if _, ok := m["del"]; ok {
			t.Fatal("expected key to be deleted")
		}
	})

	t.Run("concrete typed map — correct type", func(t *testing.T) {
		m := map[string]int{"x": 0}
		mv := reflect.ValueOf(m)
		key := reflect.ValueOf("x")
		assertNoErr(t, setMapValue(mv, key, 7))
		assertEqual(t, m["x"], 7)
	})

	t.Run("concrete typed map — convertible type", func(t *testing.T) {
		m := map[string]int64{"x": 0}
		mv := reflect.ValueOf(m)
		key := reflect.ValueOf("x")
		assertNoErr(t, setMapValue(mv, key, int32(9)))
		assertEqual(t, m["x"], int64(9))
	})

	t.Run("concrete typed map — wrong type", func(t *testing.T) {
		m := map[string]int{"x": 0}
		mv := reflect.ValueOf(m)
		key := reflect.ValueOf("x")
		assertErr(t, setMapValue(mv, key, "wrong"))
	})
}

func TestTraverseGet(t *testing.T) {
	t.Run("map key", func(t *testing.T) {
		m := map[string]any{"foo": "bar"}
		v, err := traverseGet(m, "foo")
		assertNoErr(t, err)
		assertEqual(t, v, "bar")
	})

	t.Run("map key missing", func(t *testing.T) {
		m := map[string]any{}
		_, err := traverseGet(m, "missing")
		assertErr(t, err)
	})

	t.Run("nested map via interface unwrap", func(t *testing.T) {
		inner := map[string]any{"z": 99}
		m := map[string]any{"inner": inner}
		v, err := traverseGet(m, "inner")
		assertNoErr(t, err)
		assertEqual(t, v, inner)
	})

	t.Run("struct field", func(t *testing.T) {
		s := testStruct{Name: "hi"}
		v, err := traverseGet(s, "Name")
		assertNoErr(t, err)
		assertEqual(t, v, "hi")
	})

	t.Run("struct field missing", func(t *testing.T) {
		s := testStruct{}
		_, err := traverseGet(s, "Nope")
		assertErr(t, err)
	})

	t.Run("slice index", func(t *testing.T) {
		sl := []any{"a", "b", "c"}
		v, err := traverseGet(sl, "1")
		assertNoErr(t, err)
		assertEqual(t, v, "b")
	})

	t.Run("slice out of range", func(t *testing.T) {
		sl := []any{"a"}
		_, err := traverseGet(sl, "5")
		assertErr(t, err)
	})

	t.Run("nil input", func(t *testing.T) {
		_, err := traverseGet(nil, "key")
		assertErr(t, err)
	})

	t.Run("scalar not traversable", func(t *testing.T) {
		_, err := traverseGet(42, "key")
		assertErr(t, err)
	})

	t.Run("pointer to map", func(t *testing.T) {
		m := map[string]any{"p": "ptr"}
		v, err := traverseGet(&m, "p")
		assertNoErr(t, err)
		assertEqual(t, v, "ptr")
	})
}

func TestTraverseSet_Map(t *testing.T) {
	t.Run("set top-level key", func(t *testing.T) {
		mem := newMem("x", 1)
		assertNoErr(t, traverseSet(reflect.ValueOf(mem), []string{"x"}, 42))
		assertEqual(t, mem["x"], 42)
	})

	t.Run("set new key", func(t *testing.T) {
		mem := map[string]any{}
		// setting a new key on map[string]any is valid
		assertNoErr(t, traverseSet(reflect.ValueOf(mem), []string{"new"}, "val"))
		assertEqual(t, mem["new"], "val")
	})

	t.Run("nil deletes key", func(t *testing.T) {
		mem := newMem("del", "gone")
		assertNoErr(t, traverseSet(reflect.ValueOf(mem), []string{"del"}, nil))
		if _, ok := mem["del"]; ok {
			t.Fatal("expected key deleted")
		}
	})

	t.Run("nested map set", func(t *testing.T) {
		inner := map[string]any{"b": "old"}
		mem := newMem("a", inner)
		assertNoErr(t, traverseSet(reflect.ValueOf(mem), []string{"a", "b"}, "new"))
		assertEqual(t, mem["a"].(map[string]any)["b"], "new")
	})

	t.Run("nested map missing intermediate key", func(t *testing.T) {
		mem := map[string]any{}
		err := traverseSet(reflect.ValueOf(mem), []string{"missing", "deep"}, "x")
		assertErr(t, err)
	})

	t.Run("slice element via map", func(t *testing.T) {
		sl := []any{"x", "y", "z"}
		mem := newMem("sl", sl)
		assertNoErr(t, traverseSet(reflect.ValueOf(mem), []string{"sl", "1"}, "Y"))
		assertEqual(t, mem["sl"].([]any)[1], "Y")
	})

	t.Run("empty path", func(t *testing.T) {
		mem := map[string]any{}
		err := traverseSet(reflect.ValueOf(mem), []string{}, "v")
		assertErr(t, err)
	})
}

func TestTraverseSet_Struct(t *testing.T) {
	t.Run("exported field", func(t *testing.T) {
		s := testStruct{Name: "before"}
		assertNoErr(t, traverseSet(reflect.ValueOf(&s).Elem(), []string{"Name"}, "after"))
		assertEqual(t, s.Name, "after")
	})

	t.Run("nested struct field", func(t *testing.T) {
		s := testStruct{Inner: innerStruct{Value: 1.0}}
		assertNoErr(t, traverseSet(reflect.ValueOf(&s).Elem(), []string{"Inner", "Value"}, 9.9))
		assertEqual(t, s.Inner.Value, 9.9)
	})

	t.Run("missing field", func(t *testing.T) {
		s := testStruct{}
		err := traverseSet(reflect.ValueOf(&s).Elem(), []string{"Ghost"}, "x")
		assertErr(t, err)
	})

	t.Run("case mismatch rejected", func(t *testing.T) {
		s := testStruct{}
		err := traverseSet(reflect.ValueOf(&s).Elem(), []string{"name"}, "x")
		assertErr(t, err)
	})
}

func TestTraverseSet_Slice(t *testing.T) {
	t.Run("set element", func(t *testing.T) {
		sl := []int{1, 2, 3}
		assertNoErr(t, traverseSet(reflect.ValueOf(sl), []string{"0"}, 99))
		assertEqual(t, sl[0], 99)
	})

	t.Run("out of range", func(t *testing.T) {
		sl := []int{1}
		err := traverseSet(reflect.ValueOf(sl), []string{"5"}, 0)
		assertErr(t, err)
	})

	t.Run("non-integer key", func(t *testing.T) {
		sl := []int{1}
		err := traverseSet(reflect.ValueOf(sl), []string{"foo"}, 0)
		assertErr(t, err)
	})
}

func TestTraverseSet_NilPointer(t *testing.T) {
	var p *map[string]any
	err := traverseSet(reflect.ValueOf(p), []string{"k"}, "v")
	assertErr(t, err)
}
