package execution

import (
	"fmt"
	"reflect"
	"strconv"
)

func traverseGet(current any, key string) (any, error) {
	if current == nil {
		return nil, fmt.Errorf("nil value")
	}

	v := reflect.ValueOf(current)
	v = derefValue(v)

	switch v.Kind() {
	case reflect.Map:
		mv := v.MapIndex(reflect.ValueOf(key))
		if !mv.IsValid() {
			return nil, fmt.Errorf("key %q not found", key)
		}
		// Unwrap interface so callers receive the concrete value.
		if mv.Kind() == reflect.Interface {
			mv = mv.Elem()
		}
		return mv.Interface(), nil

	case reflect.Struct:
		fv, err := structField(v, key)
		if err != nil {
			return nil, err
		}
		return fv.Interface(), nil

	case reflect.Slice, reflect.Array:
		idx, err := parseIndex(key, v.Len())
		if err != nil {
			return nil, err
		}
		return v.Index(idx).Interface(), nil

	default:
		return nil, fmt.Errorf("cannot traverse %s", v.Kind())
	}
}

func traverseSet(v reflect.Value, segments []string, value any) error {
	// Unwrap pointers and interfaces, but do not chase nil pointers.
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return fmt.Errorf("nil %s", v.Kind())
		}
		v = v.Elem()
	}

	if len(segments) == 0 {
		return fmt.Errorf("empty path")
	}

	key := segments[0]
	terminal := len(segments) == 1

	switch v.Kind() {
	case reflect.Map:
		mapKey := reflect.ValueOf(key)
		if err := ensureMapKeyType(v, &mapKey); err != nil {
			return err
		}

		if terminal {
			return setMapValue(v, mapKey, value)
		}

		// Get the existing child, make it mutable, recurse, write back.
		mv := v.MapIndex(mapKey)
		if !mv.IsValid() {
			return fmt.Errorf("key %q not found", key)
		}
		if mv.Kind() == reflect.Interface {
			mv = mv.Elem()
		}
		child := reflect.New(mv.Type()).Elem()
		child.Set(mv)

		if err := traverseSet(child, segments[1:], value); err != nil {
			return err
		}

		// Write the modified child back. Map elements are not addressable
		// so we must re-set the map index.
		v.SetMapIndex(mapKey, child)
		return nil

	case reflect.Struct:
		fv, err := structField(v, key)
		if err != nil {
			return err
		}
		if terminal {
			if !fv.CanSet() {
				return fmt.Errorf("field %q is not settable (unexported?)", key)
			}
			return assignValue(fv, value)
		}
		return traverseSet(fv, segments[1:], value)

	case reflect.Slice, reflect.Array:
		idx, err := parseIndex(key, v.Len())
		if err != nil {
			return err
		}
		ev := v.Index(idx)
		if terminal {
			return assignValue(ev, value)
		}
		return traverseSet(ev, segments[1:], value)

	default:
		return fmt.Errorf("cannot traverse %s", v.Kind())
	}
}

func assignValue(dst reflect.Value, value any) error {
	if !dst.IsValid() {
		return fmt.Errorf("invalid destination")
	}
	if !dst.CanSet() {
		return fmt.Errorf("destination is not settable")
	}

	if value == nil {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	src := reflect.ValueOf(value)

	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return nil
	}

	if dst.Kind() == reflect.Interface && src.Type().Implements(dst.Type()) {
		dst.Set(src)
		return nil
	}

	if src.Type().ConvertibleTo(dst.Type()) {
		dst.Set(src.Convert(dst.Type()))
		return nil
	}

	return fmt.Errorf("cannot assign %s to %s", src.Type(), dst.Type())
}

func setMapValue(m reflect.Value, key reflect.Value, value any) error {
	if value == nil {
		m.SetMapIndex(key, reflect.Value{})
		return nil
	}

	src := reflect.ValueOf(value)
	elemType := m.Type().Elem()

	if elemType.Kind() == reflect.Interface {
		if !src.Type().Implements(elemType) {
			return fmt.Errorf("type %s does not implement map element interface %s", src.Type(), elemType)
		}
		m.SetMapIndex(key, src)
		return nil
	}

	if src.Type().AssignableTo(elemType) {
		m.SetMapIndex(key, src)
		return nil
	}
	if src.Type().ConvertibleTo(elemType) {
		m.SetMapIndex(key, src.Convert(elemType))
		return nil
	}

	return fmt.Errorf("cannot assign %s to map element type %s", src.Type(), elemType)
}

func structField(v reflect.Value, key string) (reflect.Value, error) {
	if fv := v.FieldByName(key); fv.IsValid() {
		return fv, nil
	}
	return reflect.Value{}, fmt.Errorf("field %q not found", key)
}

func ensureMapKeyType(m reflect.Value, mapKey *reflect.Value) error {
	want := m.Type().Key()
	if mapKey.Type().AssignableTo(want) {
		return nil
	}
	if mapKey.Type().ConvertibleTo(want) {
		*mapKey = mapKey.Convert(want)
		return nil
	}
	return fmt.Errorf("map key type %s is not compatible with %s", mapKey.Type(), want)
}

func parseIndex(key string, length int) (int, error) {
	idx, err := strconv.Atoi(key)
	if err != nil {
		return 0, fmt.Errorf("expected integer index, got %q", key)
	}
	if idx < 0 || idx >= length {
		return 0, fmt.Errorf("index %d out of range [0, %d)", idx, length)
	}
	return idx, nil
}

func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}
