// Package ptr converts between values and pointers.
//
// It exists for the common case of optional fields in JSON payloads and API
// clients, where a *string distinguishes "absent" from "empty" but is awkward
// to construct inline.
package ptr

// To returns a pointer to a copy of v.
//
// It is the inline form of the two-line temporary-variable dance Go otherwise
// requires, which cannot be avoided because &"literal" is not valid Go.
func To[T any](v T) *T {
	return &v
}

// From returns the value p points to, or the zero value of T if p is nil.
func From[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// FromOr returns the value p points to, or fallback if p is nil.
func FromOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}

// ToNilable returns nil if v is the zero value of T, and a pointer to v
// otherwise. It is the usual way to omit empty fields from an outgoing payload.
func ToNilable[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

// IsNil reports whether p is nil.
func IsNil[T any](p *T) bool {
	return p == nil
}
