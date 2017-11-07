package null

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func encodeFloat(f float64, bits int) ([]byte, error) {

	// match behaviour of json.floatEncoder.encode()
	if math.IsInf(f, 0) || math.IsNaN(f) {
		panic(errors.New("null: unsupported floating point value: " +
			strconv.FormatFloat(f, 'g', -1, bits)))
	}

	// Convert as if by ES6 number to string conversion.
	// This matches most other JSON generators.
	// See golang.org/issue/6384 and golang.org/issue/14135.
	// Like fmt %g, but the exponent cutoffs are different
	// and exponents themselves are not padded to two digits.
	abs := math.Abs(f)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}
	b := strconv.AppendFloat(nil, f, fmt, -1, bits)
	if fmt == 'e' {
		// clean up e-09 to e-9
		n := len(b)
		if n >= 4 && b[n-4] == 'e' && b[n-3] == '-' && b[n-2] == '0' {
			b[n-2] = b[n-1]
			b = b[:n-1]
		}
	}

	// TODO: support quoting floats

	return b, nil
}

func convertFloat(value interface{}, bitSize int) (float64, error) {
	var f float64
	var err error

	switch v := value.(type) {
	// Conformant types
	case float64:
		f = v
	case int64:
		f = float64(v)
	case string:
		f, err = strconv.ParseFloat(v, bitSize)
	case []byte:
		f, err = strconv.ParseFloat(string(v), bitSize)

	// Accept other numeric types
	case float32:
		f = float64(v)
	case int:
		f = float64(v)
	case int8:
		f = float64(v)
	case int16:
		f = float64(v)
	case int32:
		f = float64(v)
	case uint8:
		f = float64(v)
	case uint16:
		f = float64(v)
	case uint32:
		f = float64(v)
	case uint:
		f = float64(v)
	case uint64:
		f = float64(v)
	default:
		err = fmt.Errorf("unsupported Scan, storing driver.Value type %T into type float64", value)
	}
	if err == nil && bitSize == 32 {
		if bitSize == 32 {
			switch {
			case f > math.MaxFloat32:
				f = math.MaxFloat32
				err = strconv.ErrRange
			case f < math.SmallestNonzeroFloat32:
				f = math.SmallestNonzeroFloat32
				err = strconv.ErrRange
			}
		}
	}
	// TODO: Match sql.convertAssign error message
	//
	if err != nil {
		// Special case for uint conversions
		if err == strconv.ErrRange {
			err = &strconv.NumError{
				Func: "ParseFloat",
				Num:  strconv.FormatFloat(f, 'g', -1, bitSize),
				Err:  strconv.ErrRange,
			}
		}
		return f, err
	}

	return f, nil
}
