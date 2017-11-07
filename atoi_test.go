package null

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

// Integer limit values.
const (
	MaxInt8   = 1<<7 - 1
	MinInt8   = -1 << 7
	MaxInt16  = 1<<15 - 1
	MinInt16  = -1 << 15
	MaxInt32  = 1<<31 - 1
	MinInt32  = -1 << 31
	MaxInt64  = 1<<63 - 1
	MinInt64  = -1 << 63
	MaxUint8  = 1<<8 - 1
	MaxUint16 = 1<<16 - 1
	MaxUint32 = 1<<32 - 1
	MaxUint64 = 1<<64 - 1

	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
)

type atoui64Test struct {
	in  string
	out uint64
	err error
}

var atoui64tests = []atoui64Test{
	{"", 0, strconv.ErrSyntax},
	{"0", 0, nil},
	{"1", 1, nil},
	{"12345", 12345, nil},
	{"012345", 12345, nil},
	{"12345x", 0, strconv.ErrSyntax},
	{"98765432100", 98765432100, nil},
	{"18446744073709551615", 1<<64 - 1, nil},
	{"18446744073709551616", 1<<64 - 1, strconv.ErrRange},
	{"18446744073709551620", 1<<64 - 1, strconv.ErrRange},
}

type atoi64Test struct {
	in  string
	out int64
	err error
}

var atoi64tests = []atoi64Test{
	{"", 0, strconv.ErrSyntax},
	{"0", 0, nil},
	{"-0", 0, nil},
	{"1", 1, nil},
	{"-1", -1, nil},
	{"12345", 12345, nil},
	{"-12345", -12345, nil},
	{"012345", 12345, nil},
	{"-012345", -12345, nil},
	{"98765432100", 98765432100, nil},
	{"-98765432100", -98765432100, nil},
	{"9223372036854775807", 1<<63 - 1, nil},
	{"-9223372036854775807", -(1<<63 - 1), nil},
	{"9223372036854775808", 1<<63 - 1, strconv.ErrRange},
	{"-9223372036854775808", -1 << 63, nil},
	{"9223372036854775809", 1<<63 - 1, strconv.ErrRange},
	{"-9223372036854775809", -1 << 63, strconv.ErrRange},
}

type atoui32Test struct {
	in  string
	out uint32
	err error
}

var atoui32tests = []atoui32Test{
	{"", 0, strconv.ErrSyntax},
	{"0", 0, nil},
	{"1", 1, nil},
	{"12345", 12345, nil},
	{"012345", 12345, nil},
	{"12345x", 0, strconv.ErrSyntax},
	{"987654321", 987654321, nil},
	{"4294967295", 1<<32 - 1, nil},
	{"4294967296", 1<<32 - 1, strconv.ErrRange},
}

type atoi32Test struct {
	in  string
	out int32
	err error
}

var atoi32tests = []atoi32Test{
	{"", 0, strconv.ErrSyntax},
	{"0", 0, nil},
	{"-0", 0, nil},
	{"1", 1, nil},
	{"-1", -1, nil},
	{"12345", 12345, nil},
	{"-12345", -12345, nil},
	{"012345", 12345, nil},
	{"-012345", -12345, nil},
	{"12345x", 0, strconv.ErrSyntax},
	{"-12345x", 0, strconv.ErrSyntax},
	{"987654321", 987654321, nil},
	{"-987654321", -987654321, nil},
	{"2147483647", 1<<31 - 1, nil},
	{"-2147483647", -(1<<31 - 1), nil},
	{"2147483648", 1<<31 - 1, strconv.ErrRange},
	{"-2147483648", -1 << 31, nil},
	{"2147483649", 1<<31 - 1, strconv.ErrRange},
	{"-2147483649", -1 << 31, strconv.ErrRange},
}

type convertIntTest struct {
	in      interface{}
	out     int64
	bitSize int
	err     error // TODO: use an error type
}

var convertIntTests = []convertIntTest{
	// Valid
	{0, 0, 0, nil},
	{1, 1, 0, nil},
	{-1, -1, 0, nil},

	// Accepted Types
	{int8(1), 1, 0, nil},
	{int16(1), 1, 0, nil},
	{int32(1), 1, 0, nil},
	{int64(1), 1, 0, nil},
	{int(1), 1, 0, nil},
	{uint8(1), 1, 0, nil},
	{uint16(1), 1, 0, nil},
	{uint32(1), 1, 0, nil},
	{uint64(1), 1, 0, nil},
	{uint(1), 1, 0, nil},

	// Limits
	{int64(MaxInt), int64(MaxInt), 0, nil},
	{int8(MaxInt8), MaxInt8, 8, nil},
	{int16(MaxInt16), MaxInt16, 16, nil},
	{int32(MaxInt32), MaxInt32, 32, nil},
	{int64(MaxInt64), int64(MaxInt64), 64, nil},

	{int8(MinInt8), MinInt8, 8, nil},
	{int16(MinInt16), MinInt16, 16, nil},
	{int32(MinInt32), MinInt32, 32, nil},
	{int64(MinInt64), int64(MinInt64), 64, nil},

	// Make sure 64-bit values are handled correctly
	{int64(MaxInt64), int64(MaxInt64), 64, nil},
	{int64(MinInt64), int64(MinInt64), 64, nil},

	// Zero base (target int size)
	{int64(MaxInt8), MaxInt8, 0, nil},
	{int64(MaxInt16), MaxInt16, 0, nil},
	{int64(MaxInt32), MaxInt32, 0, nil},

	// Error
	{uint64(MaxInt8) + 1, MaxInt8, 8, strconv.ErrRange},
	{uint64(MaxInt16) + 1, MaxInt16, 16, strconv.ErrRange},
	{uint64(MaxInt32) + 1, MaxInt32, 32, strconv.ErrRange},
	{uint64(MaxInt64) + 1, int64(MaxInt64), 64, strconv.ErrRange},
	{uint64(MaxInt) + 1, int64(MaxInt), 0, strconv.ErrRange},

	{MinInt8 - 1, MinInt8, 8, strconv.ErrRange},
	{MinInt16 - 1, MinInt16, 16, strconv.ErrRange},
	{int64(MinInt32 - 1), MinInt32, 32, strconv.ErrRange},

	{"9223372036854775809", MaxInt64, 0, strconv.ErrRange},
	{"-9223372036854775809", MinInt64, 0, strconv.ErrRange},

	{new(int), 0, 0, errors.New("fixme")},
	{nil, 0, 0, errors.New("fixme")},
	{true, 0, 0, errors.New("fixme")},
	{"true", 0, 0, errors.New("fixme")},
	{"0x12345", 0, 0, errors.New("fixme")},
}

type convertUintTest struct {
	in      interface{}
	out     uint64
	bitSize int
	err     error // TODO: use an error type
}

var convertUintTests = []convertUintTest{
	// Valid
	{0, 0, 0, nil},
	{1, 1, 0, nil},

	// Accepted Types
	{int8(1), 1, 0, nil},
	{int16(1), 1, 0, nil},
	{int32(1), 1, 0, nil},
	{int64(1), 1, 0, nil},
	{int(1), 1, 0, nil},
	{uint8(1), 1, 0, nil},
	{uint16(1), 1, 0, nil},
	{uint32(1), 1, 0, nil},
	{uint64(1), 1, 0, nil},
	{uint(1), 1, 0, nil},

	// Limits
	{uint64(MaxInt), uint64(MaxInt), 0, nil},
	{uint8(MaxUint8), MaxUint8, 8, nil},
	{uint16(MaxUint16), MaxUint16, 16, nil},
	{uint32(MaxUint32), MaxUint32, 32, nil},
	{uint64(MaxUint64), MaxUint64, 64, nil},

	// Zero base
	{uint64(MaxUint), uint64(MaxUint), 0, nil},
	{uint8(MaxUint8), MaxUint8, 0, nil},
	{uint16(MaxUint16), MaxUint16, 0, nil},
	{uint32(MaxUint32), MaxUint32, 0, nil},

	// 64-bit
	{uint64(MaxInt64), MaxInt64, 64, nil},

	// Error cases

	{uint64(MaxUint8) + 1, MaxUint64, 8, strconv.ErrRange},
	{uint64(MaxUint16) + 1, MaxUint64, 16, strconv.ErrRange},
	{uint64(MaxUint32) + 1, MaxUint64, 32, strconv.ErrRange},

	{int8(MinInt8), 0, 64, strconv.ErrRange},
	{int16(MinInt16), 0, 64, strconv.ErrRange},
	{int32(MinInt32), 0, 64, strconv.ErrRange},
	{int64(MinInt64), 0, 64, strconv.ErrRange},

	{"18446744073709551620", MaxUint64, 0, strconv.ErrRange},
	{"-1", 0, 0, strconv.ErrSyntax},

	{new(int), 0, 0, errors.New("fixme")},
	{nil, 0, 0, errors.New("fixme")},
	{true, 0, 0, errors.New("fixme")},
	{"true", 0, 0, errors.New("fixme")},
	{"0x12345", 0, 0, errors.New("fixme")},
}

func formatNumber(v interface{}) (string, bool) {
	switch n := v.(type) {
	case int8:
		return strconv.FormatInt(int64(n), 10), true
	case int16:
		return strconv.FormatInt(int64(n), 10), true
	case int32:
		return strconv.FormatInt(int64(n), 10), true
	case int64:
		return strconv.FormatInt(int64(n), 10), true
	case int:
		return strconv.FormatInt(int64(n), 10), true
	case uint8:
		return strconv.FormatUint(uint64(n), 10), true
	case uint16:
		return strconv.FormatUint(uint64(n), 10), true
	case uint32:
		return strconv.FormatUint(uint64(n), 10), true
	case uint64:
		return strconv.FormatUint(uint64(n), 10), true
	case uint:
		return strconv.FormatUint(uint64(n), 10), true
	case float32:
		return strconv.FormatFloat(float64(n), 'g', -1, 32), true
	case float64:
		return strconv.FormatFloat(n, 'g', -1, 64), true
	}
	return "", false
}

func init() {
	// The atoi routines return NumErrors wrapping
	// the error and the string. Convert the tables above.
	for i := range atoui64tests {
		test := &atoui64tests[i]
		if test.err != nil {
			test.err = &strconv.NumError{"ParseUint", test.in, test.err}
		}
	}
	for i := range atoi64tests {
		test := &atoi64tests[i]
		if test.err != nil {
			test.err = &strconv.NumError{"ParseInt", test.in, test.err}
		}
	}
	for i := range atoui32tests {
		test := &atoui32tests[i]
		if test.err != nil {
			test.err = &strconv.NumError{"ParseUint", test.in, test.err}
		}
	}
	for i := range atoi32tests {
		test := &atoi32tests[i]
		if test.err != nil {
			test.err = &strconv.NumError{"ParseInt", test.in, test.err}
		}
	}

	// Create string and []byte versions of numeric tests
	for _, test := range convertIntTests {
		if s, ok := formatNumber(test.in); ok {
			convertIntTests = append(convertIntTests, convertIntTest{
				s, test.out, test.bitSize, test.err,
			})
			convertIntTests = append(convertIntTests, convertIntTest{
				[]byte(s), test.out, test.bitSize, test.err,
			})
		}
	}
	for _, test := range convertUintTests {
		if s, ok := formatNumber(test.in); ok {
			convertUintTests = append(convertUintTests, convertUintTest{
				s, test.out, test.bitSize, test.err,
			})
			convertUintTests = append(convertUintTests, convertUintTest{
				[]byte(s), test.out, test.bitSize, test.err,
			})
		}
	}
}

func TestParseUint64(t *testing.T) {
	for i := range atoui64tests {
		test := &atoui64tests[i]
		out, err := parseUint([]byte(test.in), 64)
		if test.out != out || !reflect.DeepEqual(test.err, err) {
			t.Errorf("Atoui64(%q) = %v, %v want %v, %v",
				test.in, out, err, test.out, test.err)
		}
	}
}

func TestParseInt64(t *testing.T) {
	for i := range atoi64tests {
		test := &atoi64tests[i]
		out, err := parseInt([]byte(test.in), 64)
		if test.out != out || !reflect.DeepEqual(test.err, err) {
			t.Errorf("Atoi64(%q) = %v, %v want %v, %v",
				test.in, out, err, test.out, test.err)
		}
	}
}

func TestParseUint(t *testing.T) {
	switch strconv.IntSize {
	case 32:
		for i := range atoui32tests {
			test := &atoui32tests[i]
			out, err := parseUint([]byte(test.in), 0)
			if test.out != uint32(out) || !reflect.DeepEqual(test.err, err) {
				t.Errorf("Atoui(%q) = %v, %v want %v, %v",
					test.in, out, err, test.out, test.err)
			}
		}
	case 64:
		for i := range atoui64tests {
			test := &atoui64tests[i]
			out, err := parseUint([]byte(test.in), 0)
			if test.out != uint64(out) || !reflect.DeepEqual(test.err, err) {
				t.Errorf("Atoui(%q) = %v, %v want %v, %v",
					test.in, out, err, test.out, test.err)
			}
		}
	}
}

func TestParseInt(t *testing.T) {
	switch strconv.IntSize {
	case 32:
		for i := range atoi32tests {
			test := &atoi32tests[i]
			out, err := parseInt([]byte(test.in), 0)
			if test.out != int32(out) || !reflect.DeepEqual(test.err, err) {
				t.Errorf("Atoi(%q) = %v, %v want %v, %v",
					test.in, out, err, test.out, test.err)
			}
		}
	case 64:
		for i := range atoi64tests {
			test := &atoi64tests[i]
			out, err := parseInt([]byte(test.in), 0)
			if test.out != int64(out) || !reflect.DeepEqual(test.err, err) {
				t.Errorf("Atoi(%q) = %v, %v want %v, %v",
					test.in, out, err, test.out, test.err)
			}
		}
	}
}

func TestConvertInt(t *testing.T) {
	for _, test := range convertIntTests {
		out, err := convertInt(test.in, test.bitSize)
		ok := test.err == nil
		if test.out != out || ok != (err == nil) {
			t.Errorf("convertInt(%v - %d) = %v, %v want %v, %v",
				test.in, test.bitSize, out, err, test.out, test.err)
		}
	}
}

func TestConvertUint(t *testing.T) {
	for _, test := range convertUintTests {
		out, err := convertUint(test.in, test.bitSize)
		ok := test.err == nil
		if test.out != out || ok != (err == nil) {
			t.Errorf("convertUint(%v - %d - %T) = %v, %v want %v, %v",
				test.in, test.bitSize, test.in, out, err, test.out, test.err)
		}
	}
}

// Parse benchmarks

func BenchmarkParseInt(b *testing.B) {
	s := []byte("12345678")
	for i := 0; i < b.N; i++ {
		parseInt(s, 0)
	}
}

func BenchmarkParseInt_Base(b *testing.B) {
	s := []byte("12345678")
	for i := 0; i < b.N; i++ {
		strconv.ParseInt(string(s), 10, 0)
	}
}

func BenchmarkParseIntNeg(b *testing.B) {
	s := []byte("-12345678")
	for i := 0; i < b.N; i++ {
		parseInt(s, 0)
	}
}

func BenchmarkParseIntNeg_Base(b *testing.B) {
	s := []byte("-12345678")
	for i := 0; i < b.N; i++ {
		strconv.ParseInt(string(s), 10, 0)
	}
}

func BenchmarkParseUint64(b *testing.B) {
	s := []byte("12345678901234")
	for i := 0; i < b.N; i++ {
		parseUint(s, 64)
	}
}

func BenchmarkParseUint64_Base(b *testing.B) {
	s := []byte("12345678901234")
	for i := 0; i < b.N; i++ {
		strconv.ParseUint(string(s), 10, 64)
	}
}

// Numeric Conversions

func BenchmarkConvertInt_Int64(b *testing.B) {
	const v int64 = 12345678901234
	for i := 0; i < b.N; i++ {
		if _, err := convertInt(v, 64); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertInt_String(b *testing.B) {
	const v = "12345678901234"
	for i := 0; i < b.N; i++ {
		if _, err := convertInt(v, 64); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertUint_Int64(b *testing.B) {
	const v int64 = 12345678901234
	for i := 0; i < b.N; i++ {
		if _, err := convertUint(v, 64); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertUint_String(b *testing.B) {
	const v = "12345678901234"
	for i := 0; i < b.N; i++ {
		if _, err := convertUint(v, 64); err != nil {
			b.Fatal(err)
		}
	}
}
