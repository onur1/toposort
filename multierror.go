package toposort

import (
	"errors"
	"fmt"
)

// MultiError stores multiple errors.
//
// Borrowed from the AppEngine SDK.
type MultiError []error

func (m MultiError) Is(target error) bool {
	for _, e := range m {
		if errors.Is(e, target) {
			return true
		}
	}
	return false
}

func (m MultiError) Error() string {
	s, n := "", 0
	for _, e := range m {
		if e != nil {
			if n == 0 {
				s = e.Error()
			}
			n++
		}
	}
	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, n-1)
}
