package null

import "strconv"

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
	const maxUint64 = (1<<64 - 1)

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
