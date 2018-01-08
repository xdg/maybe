package maybe

import "fmt"

// AoAoS implements the Maybe monad for a 2D slice of strings.  An AoAoS is
// considered 'valid' or 'invalid' depending on whether it contains a 2D slice
// of strings or an error value.
type AoAoS struct {
	just [][]string
	err  error
}

// NewAoAoS constructs an AoAoS from a given 2D slice of strings or error. If
// e is not nil, returns ErrAoAoS(e), otherwise returns JustAoAoS(s)
func NewAoAoS(s [][]string, e error) AoAoS {
	if e != nil {
		return ErrAoAoS(e)
	}
	return JustAoAoS(s)
}

// JustAoAoS constructs a valid AoAoS from a given 2D slice of stringss.
func JustAoAoS(s [][]string) AoAoS {
	return AoAoS{just: s}
}

// ErrAoAoS constructs an invalid AoAoS from a given error.
func ErrAoAoS(e error) AoAoS {
	return AoAoS{err: e}
}

// IsErr returns true for an invalid AoAoS.
func (m AoAoS) IsErr() bool {
	return m.err != nil
}

// Bind applies a function that takes a 2D slice of strings and returns an AoAoS.
func (m AoAoS) Bind(f func(s [][]string) AoAoS) AoAoS {
	if m.err != nil {
		return m
	}

	return f(m.just)
}

// Join applies a function that takes a 2D slice of strings and returns an AoS.
func (m AoAoS) Join(f func(s [][]string) AoS) AoS {
	if m.err != nil {
		return ErrAoS(m.err)
	}

	return f(m.just)
}

// Map applies a function to each element of a valid AoAoS and returns a new
// AoAoS.  If the AoAoS is invalid or if any function returns an invalid AoS,
// Map returns an invalid AoAoS.
func (m AoAoS) Map(f func(xs []string) AoS) AoAoS {
	if m.err != nil {
		return m
	}

	new := make([][]string, len(m.just))
	for i, v := range m.just {
		strs, err := f(v).Unbox()
		if err != nil {
			return ErrAoAoS(err)
		}
		new[i] = strs
	}

	return JustAoAoS(new)
}

// ToInt applies a function that takes a string and returns an I.  If the
// AoAoS is invalid or if any function returns an invalid I, ToInt returns an
// invalid AoAoI.  Note: unlike Map, this is a deep conversion of individual
// elements of the 2D slice of strings.
func (m AoAoS) ToInt(f func(s string) I) AoAoI {
	if m.err != nil {
		return ErrAoAoI(m.err)
	}

	new := make([][]int, len(m.just))
	for i, xs := range m.just {
		new[i] = make([]int, len(xs))
		for j, v := range xs {
			num, err := f(v).Unbox()
			if err != nil {
				return ErrAoAoI(err)
			}
			new[i][j] = num
		}
	}

	return JustAoAoI(new)
}

// String returns a string representation, mostly useful for debugging.
func (m AoAoS) String() string {
	if m.err != nil {
		return fmt.Sprintf("Err %v", m.err)
	}
	return fmt.Sprintf("Just %v", m.just)
}

// Unbox returns the underlying 2D slice of strings value or error.
func (m AoAoS) Unbox() ([][]string, error) {
	return m.just, m.err
}
