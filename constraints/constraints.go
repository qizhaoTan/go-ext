// Package constraints defines the type constraints shared by the go-ext packages.
//
// The standard library exposes [cmp.Ordered], so this package deliberately does
// not redefine it. What is missing from the standard library is the numeric
// decomposition below, which lets generic code express "any number" without
// giving up compile-time safety.
package constraints

// Signed is a constraint permitting any signed integer type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint permitting any unsigned integer type.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint permitting any integer type.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint permitting any floating-point type.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint permitting any complex numeric type.
type Complex interface {
	~complex64 | ~complex128
}

// Number is a constraint permitting any integer or floating-point type.
//
// Complex is excluded on purpose: complex values are not ordered, so functions
// such as Max or Sum-with-comparison cannot be written against them.
type Number interface {
	Integer | Float
}
