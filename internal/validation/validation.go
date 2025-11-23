package validation

import "reflect"

func IsNil(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)

	// Check if the Kind is one that can logically be nil.
	// Array and Struct are excluded because they are value types.
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Func, reflect.Map, reflect.Slice, reflect.Chan, reflect.UnsafePointer:
		// Check if the underlying value is nil.
		return v.IsNil()
	default:
		// Value types (struct, int, string, etc.) are never nil once they are inside an interface.
		return false
	}
}
