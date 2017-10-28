package null

import (
	"errors"
	"strconv"
)

// ErrRange indicates that a value is out of range for the target type.
var ErrRange = errors.New("value out of range")

// ErrSyntax indicates that a value does not have the right syntax for the target type.
var ErrSyntax = errors.New("invalid syntax")

// A NumError records a failed conversion.
type NumError struct {
	Func string // the failing function (ParseBool, ParseInt, ParseUint, ParseFloat)
	Num  string // the input
	Err  error  // the reason the conversion failed (ErrRange, ErrSyntax)
}

func (e *NumError) Error() string {
	// TODO: Remove use of strconv
	return "strconv." + e.Func + ": " + "parsing " + strconv.Quote(e.Num) + ": " + e.Err.Error()
}

const intSize = 32 << (^uint(0) >> 63)

// IntSize is the size in bits of an int or uint value.
const IntSize = intSize

const maxUint64 = (1<<64 - 1)

func ParseUint(s string, bitSize int) (uint64, error) {
	const base = 10
	const cutoff = maxUint64/10 + 1

	if len(s) < 1 {
		return 0, &NumError{"ParseUint", s, ErrSyntax}
	}

	if bitSize == 0 {
		bitSize = int(IntSize)
	}
	maxVal := uint64(1<<uint(bitSize) - 1)

	var n uint64
	var err error
	for i := 0; i < len(s); i++ {
		v := uint64(s[i]) - '0'
		if v > 9 {
			n = 0
			err = ErrSyntax
			goto Error
		}
		if n >= cutoff {
			// n*base overflows
			n = maxUint64
			err = ErrRange
			goto Error
		}
		n *= base

		n1 := n + v
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = maxUint64
			err = ErrRange
			goto Error
		}
		n = n1
	}

	return n, nil

Error:
	return n, &NumError{"ParseUint", s, err}
}

/*
func ParseUint(s string, bitSize int) (uint64, error) {
	const base = 10
	const cutoff = maxUint64/10 + 1

	var i int
	var n uint64
	var maxVal uint64
	var err error

	if bitSize == 0 {
		bitSize = int(IntSize)
	}

	if len(s) < 1 {
		err = ErrSyntax
		goto Error
	}

	maxVal = 1<<uint(bitSize) - 1

	for ; i < len(s); i++ {
		var v byte
		d := s[i]
		if '0' <= s[i] && s[i] <= '9' {
			v = d - '0'
		} else {
			n = 0
			err = ErrSyntax
			goto Error
		}
		if n >= cutoff {
			// n*base overflows
			n = maxUint64
			err = ErrRange
			goto Error
		}
		n *= uint64(base)

		n1 := n + uint64(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = maxUint64
			err = ErrRange
			goto Error
		}
		n = n1
	}

	return n, nil

Error:
	return n, &NumError{"ParseUint", s, err}
}
*/
