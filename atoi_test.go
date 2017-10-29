package null

import (
	"reflect"
	"strconv"
	"testing"
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
