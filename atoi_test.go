package null

import (
	"reflect"
	"testing"
)

type atoui64Test struct {
	in  string
	out uint64
	err error
}

var atoui64tests = []atoui64Test{
	{"", 0, ErrSyntax},
	{"0", 0, nil},
	{"1", 1, nil},
	{"12345", 12345, nil},
	{"012345", 12345, nil},
	{"12345x", 0, ErrSyntax},
	{"98765432100", 98765432100, nil},
	{"18446744073709551615", 1<<64 - 1, nil},
	{"18446744073709551616", 1<<64 - 1, ErrRange},
	{"18446744073709551620", 1<<64 - 1, ErrRange},
}

func init() {
	// The atoi routines return NumErrors wrapping
	// the error and the string. Convert the tables above.
	for i := range atoui64tests {
		test := &atoui64tests[i]
		if test.err != nil {
			test.err = &NumError{"ParseUint", test.in, test.err}
		}
	}
}

func TestParseUint64(t *testing.T) {
	for i := range atoui64tests {
		test := &atoui64tests[i]
		out, err := ParseUint(test.in, 64)
		if test.out != out || !reflect.DeepEqual(test.err, err) {
			t.Errorf("Atoui64(%q) = %v, %v want %v, %v",
				test.in, out, err, test.out, test.err)
		}
	}
}

func BenchmarkParseUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseUint("12345678901234", 64)
	}
}
