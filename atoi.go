package null

import (
	"fmt"
	"strconv"
)

const (
	maxUint64 = (1<<64 - 1)
	maxInt64  = 1<<63 - 1
)

func parseInt(s []byte, bitSize int) (int64, error) {
	if bitSize == 0 {
		bitSize = int(strconv.IntSize)
	}
	if len(s) == 0 {
		return 0, &strconv.NumError{"ParseInt", string(s), strconv.ErrSyntax}
	}

	// Pick off leading sign.
	s0 := s
	neg := false
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	// Convert unsigned and check range.
	un, err := parseUint(s, bitSize)
	if err != nil && err.(*strconv.NumError).Err != strconv.ErrRange {
		err.(*strconv.NumError).Func = "ParseInt"
		err.(*strconv.NumError).Num = string(s0)
		return 0, err
	}
	cutoff := uint64(1 << uint(bitSize-1))
	if !neg && un >= cutoff {
		return int64(cutoff - 1), &strconv.NumError{"ParseInt", string(s0), strconv.ErrRange}
	}
	if neg && un > cutoff {
		return -int64(cutoff), &strconv.NumError{"ParseInt", string(s0), strconv.ErrRange}
	}
	n := int64(un)
	if neg {
		n = -n
	}

	return n, nil
}

func parseUint(s []byte, bitSize int) (uint64, error) {
	if len(s) < 1 {
		return 0, &strconv.NumError{"ParseUint", string(s), strconv.ErrSyntax}
	}
	if bitSize == 0 {
		bitSize = int(strconv.IntSize)
	}
	maxVal := uint64(1<<uint(bitSize) - 1)

	var n uint64
	var err error
	for i := 0; i < len(s); i++ {
		v := uint64(s[i]) - '0'
		if v > 9 {
			n = 0
			err = strconv.ErrSyntax
			goto Error
		}
		if n >= maxUint64/10+1 {
			// n*10 overflows
			n = maxUint64
			err = strconv.ErrRange
			goto Error
		}
		n *= 10

		n1 := n + v
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = maxUint64
			err = strconv.ErrRange
			goto Error
		}
		n = n1
	}

	return n, nil

Error:
	return n, &strconv.NumError{"ParseUint", string(s), err}
}

func convertInt(value interface{}, bitSize int) (int64, error) {
	var n int64
	var err error
	if bitSize == 0 {
		bitSize = strconv.IntSize
	}
	cutoff := uint64(1 << uint(bitSize-1))

	switch v := value.(type) {
	// Conformant types
	case int64:
		n = v
	case string:
		n, err = strconv.ParseInt(v, 10, bitSize)
	case []byte:
		n, err = parseInt(v, bitSize)

	// Accept other numeric types
	case int:
		n = int64(v)
	case int8:
		n = int64(v)
	case int16:
		n = int64(v)
	case int32:
		n = int64(v)
	case uint8:
		n = int64(v)
	case uint16:
		n = int64(v)
	case uint32:
		n = int64(v)
	case uint:
		if v <= maxInt64 {
			n = int64(v)
		} else {
			n = int64(cutoff - 1)
			err = strconv.ErrRange
		}
	case uint64:
		if v <= maxInt64 {
			n = int64(v)
		} else {
			n = int64(cutoff - 1)
			err = strconv.ErrRange
		}
	default:
		err = fmt.Errorf("unsupported Scan, storing driver.Value type %T into type int64", value)
	}

	// TODO: Match sql.convertAssign error message
	//
	if err != nil {
		// Special case for uint conversions
		if err == strconv.ErrRange {
			err = &strconv.NumError{"ParseInt", strconv.FormatInt(n, 10), strconv.ErrRange}
		}
		return n, err
	}

	if n >= 0 && uint64(n) >= cutoff {
		n = int64(cutoff - 1)
		err = &strconv.NumError{"ParseInt", strconv.FormatInt(n, 10), strconv.ErrRange}
	}
	if n < 0 && uint64(-n) > cutoff {
		n = -int64(cutoff)
		err = &strconv.NumError{"ParseInt", strconv.FormatInt(n, 10), strconv.ErrRange}
	}

	return n, err
}

func convertUint(value interface{}, bitSize int) (uint64, error) {
	var n uint64
	var err error
	if bitSize == 0 {
		bitSize = strconv.IntSize
	}
	cutoff := uint64(1<<uint(bitSize) - 1)

	switch v := value.(type) {
	// Conformant types
	case int64:
		if v < 0 {
			goto ErrOverflow
		}
		n = uint64(v)
	case string:
		n, err = strconv.ParseUint(v, 10, bitSize)
	case []byte:
		n, err = parseUint(v, bitSize)

	// Accept other numeric types
	case uint:
		n = uint64(v)
	case uint8:
		n = uint64(v)
	case uint16:
		n = uint64(v)
	case uint32:
		n = uint64(v)
	case uint64:
		n = v
	case int:
		if v < 0 {
			goto ErrOverflow
		}
		n = uint64(v)
	case int8:
		if v < 0 {
			goto ErrOverflow
		}
		n = uint64(v)
	case int16:
		if v < 0 {
			goto ErrOverflow
		}
		n = uint64(v)
	case int32:
		if v < 0 {
			goto ErrOverflow
		}
		n = uint64(v)
	default:
		err = fmt.Errorf("unsupported Scan, storing driver.Value type %T into type int64", value)
	}

	// TODO: Match sql.convertAssign error message
	//
	if err != nil {
		// Special case for uint conversions
		if err == strconv.ErrRange {
			err = &strconv.NumError{"ParseInt", strconv.FormatUint(n, 10), strconv.ErrRange}
		}
	} else if n > cutoff {
		n = maxUint64
		err = &strconv.NumError{"ParseInt", strconv.FormatUint(n, 10), strconv.ErrRange}
	}
	return n, err

ErrOverflow:
	return 0, &strconv.NumError{"ParseInt", strconv.FormatUint(n, 10), strconv.ErrRange}
}
