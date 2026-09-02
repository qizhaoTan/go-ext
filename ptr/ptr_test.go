package ptr_test

import (
	"testing"

	"github.com/qizhaoTan/go-ext/ptr"
)

func TestToAndFrom(t *testing.T) {
	p := ptr.To(42)
	if p == nil || *p != 42 {
		t.Fatalf("To(42) = %v, want pointer to 42", p)
	}
	if got := ptr.From(p); got != 42 {
		t.Errorf("From() = %d, want 42", got)
	}
}

func TestToCopiesTheValue(t *testing.T) {
	v := 1
	p := ptr.To(v)
	v = 2
	if *p != 1 {
		t.Errorf("To did not copy: *p = %d, want 1", *p)
	}
}

func TestFromNilYieldsZero(t *testing.T) {
	if got := ptr.From[int](nil); got != 0 {
		t.Errorf("From(nil) = %d, want 0", got)
	}
	if got := ptr.From[string](nil); got != "" {
		t.Errorf("From(nil) = %q, want empty", got)
	}
}

func TestFromOr(t *testing.T) {
	if got := ptr.FromOr(ptr.To("x"), "fallback"); got != "x" {
		t.Errorf("FromOr(non-nil) = %q, want %q", got, "x")
	}
	if got := ptr.FromOr(nil, "fallback"); got != "fallback" {
		t.Errorf("FromOr(nil) = %q, want %q", got, "fallback")
	}
}

func TestToNilable(t *testing.T) {
	if got := ptr.ToNilable(""); got != nil {
		t.Errorf("ToNilable(zero) = %v, want nil", got)
	}
	got := ptr.ToNilable("x")
	if got == nil || *got != "x" {
		t.Errorf("ToNilable(non-zero) = %v, want pointer to x", got)
	}
}

func TestIsNil(t *testing.T) {
	if !ptr.IsNil[int](nil) {
		t.Error("IsNil(nil) = false, want true")
	}
	if ptr.IsNil(ptr.To(0)) {
		t.Error("IsNil(pointer to zero) = true, want false")
	}
}
