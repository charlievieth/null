package null

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"testing"
)

type convertFloatTest struct {
	in      interface{}
	out     float64
	bitSize int
	err     error // TODO: use an error type
}

var convertFloatTests = []convertFloatTest{
	{-1, -1, 64, nil},
	{0, 0, 64, nil},
	{1, 1, 64, nil},
	{float64(math.SmallestNonzeroFloat64), math.SmallestNonzeroFloat64, 64, nil},
	{float64(math.MaxFloat64), math.MaxFloat64, 64, nil},
	{float64(math.MaxFloat32), math.MaxFloat32, 32, nil},
	{float64(math.SmallestNonzeroFloat32), math.SmallestNonzeroFloat32, 32, nil},

	// Error
	{float64(math.MaxFloat64), math.MaxFloat64, 32, strconv.ErrRange},
	{float64(math.SmallestNonzeroFloat64), math.SmallestNonzeroFloat64, 32, strconv.ErrRange},
}

func init() {
	for _, test := range convertFloatTests {
		if s, ok := formatNumber(test.in); ok {
			convertFloatTests = append(convertFloatTests, convertFloatTest{
				s, test.out, test.bitSize, test.err,
			})
			convertFloatTests = append(convertFloatTests, convertFloatTest{
				[]byte(s), test.out, test.bitSize, test.err,
			})
		}
	}
}

func TestConvertFloat(t *testing.T) {
	for _, test := range convertFloatTests {
		out, err := convertFloat(test.in, test.bitSize)
		ok := test.err == nil
		switch {
		case ok != (err == nil):
			t.Errorf("convertFloat(%v - %d - %T) = %v, %v want %v, %v",
				test.in, test.bitSize, test.in, out, err, test.out, test.err)
		case err == nil && test.out != out:
			// TODO (CEV): Do we actually want to check the return value?
			t.Errorf("convertFloat(%v - %d - %T) = %v, %v want %v, %v",
				test.in, test.bitSize, test.in, out, err, test.out, test.err)
		}
	}
}

var re = regexp.MustCompile

// syntactic checks on form of marshaled floating point numbers.
var badFloatREs = []*regexp.Regexp{
	re(`p`),                     // no binary exponential notation
	re(`^\+`),                   // no leading + sign
	re(`^-?0[^.]`),              // no unnecessary leading zeros
	re(`^-?\.`),                 // leading zero required before decimal point
	re(`\.(e|$)`),               // no trailing decimal
	re(`\.[0-9]+0(e|$)`),        // no trailing zero in fraction
	re(`^-?(0|[0-9]{2,})\..*e`), // exponential notation must have normalized mantissa
	re(`e[0-9]`),                // positive exponent must be signed
	re(`e[+-]0`),                // exponent must not have leading zeros
	re(`e-[1-6]$`),              // not tiny enough for exponential notation
	re(`e+(.|1.|20)$`),          // not big enough for exponential notation
	re(`^-?0\.0000000`),         // too tiny, should use exponential notation
	re(`^-?[0-9]{22}`),          // too big, should use exponential notation
	re(`[1-9][0-9]{16}[1-9]`),   // too many significant digits in integer
	re(`[1-9][0-9.]{17}[1-9]`),  // too many significant digits in decimal
	// below here for float32 only
	re(`[1-9][0-9]{8}[1-9]`),  // too many significant digits in integer
	re(`[1-9][0-9.]{9}[1-9]`), // too many significant digits in decimal
}

func TestMarshalFloat(t *testing.T) {
	t.Parallel()
	nfail := 0
	test := func(f float64, bits int) {
		vf := interface{}(Float64{Valid: true, Float64: f})
		if bits == 32 {
			vf = interface{}(Float32{Valid: true, Float32: float32(f)})
		}
		bout, err := json.Marshal(vf)
		if err != nil {
			t.Errorf("Marshal(%T(%g)): %v", vf, vf, err)
			nfail++
			return
		}
		out := string(bout)

		// result must convert back to the same float
		g, err := strconv.ParseFloat(out, bits)
		if err != nil {
			t.Errorf("Marshal(%T(%g)) = %q, cannot parse back: %v", vf, vf, out, err)
			nfail++
			return
		}
		if f != g || fmt.Sprint(f) != fmt.Sprint(g) { // fmt.Sprint handles ±0
			t.Errorf("Marshal(%T(%g)) = %q (is %g, not %g)", vf, vf, out, float32(g), vf)
			nfail++
			return
		}

		bad := badFloatREs
		if bits == 64 {
			bad = bad[:len(bad)-2]
		}
		for _, re := range bad {
			if re.MatchString(out) {
				t.Errorf("Marshal(%T(%g)) = %q, must not match /%s/", vf, vf, out, re)
				nfail++
				return
			}
		}
	}

	var (
		bigger  = math.Inf(+1)
		smaller = math.Inf(-1)
	)

	var digits = "1.2345678901234567890123"
	for i := len(digits); i >= 2; i-- {
		for exp := -30; exp <= 30; exp++ {
			for _, sign := range "+-" {
				for bits := 32; bits <= 64; bits += 32 {
					s := fmt.Sprintf("%c%se%d", sign, digits[:i], exp)
					f, err := strconv.ParseFloat(s, bits)
					if err != nil {
						log.Fatal(err)
					}
					next := math.Nextafter
					if bits == 32 {
						next = func(g, h float64) float64 {
							return float64(math.Nextafter32(float32(g), float32(h)))
						}
					}
					test(f, bits)
					test(next(f, bigger), bits)
					test(next(f, smaller), bits)
					if nfail > 50 {
						t.Fatalf("stopping test early")
					}
				}
			}
		}
	}
	test(0, 64)
	test(math.Copysign(0, -1), 64)
	test(0, 32)
	test(math.Copysign(0, -1), 32)
}
