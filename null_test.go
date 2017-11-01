package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"testing"
	"time"
)

func TestPtr(t *testing.T) {
	if PtrInt(nil).Valid {
		t.Error("PtrInt: expected Valid to equal false")
	}
	if PtrFloat64(nil).Valid {
		t.Error("PtrFloat64: expected Valid to equal false")
	}
	if PtrString(nil).Valid {
		t.Error("PtrString: expected Valid to equal false")
	}
	if PtrBool(nil).Valid {
		t.Error("PtrBool: expected Valid to equal false")
	}
	if PtrTime(nil).Valid {
		t.Error("PtrTime: expected Valid to equal false")
	}

	i := int(1)
	if n := PtrInt(&i); !n.Valid || n.Int != i {
		t.Error("PtrInt: expected Valid to equal true")
	}
	f := float64(1)
	if n := PtrFloat64(&f); !n.Valid || n.Float64 != f {
		t.Error("PtrFloat64: expected Valid to equal true")
	}
	s := "a"
	if n := PtrString(&s); !n.Valid || n.String != s {
		t.Error("PtrString: expected Valid to equal true")
	}
	b := true
	if n := PtrBool(&b); !n.Valid || n.Bool != b {
		t.Error("PtrBool: expected Valid to equal true")
	}
	tt := time.Now()
	if n := PtrTime(&tt); !n.Valid || n.Time != tt {
		t.Error("PtrTime: expected Valid to equal true")
	}
}

func TestNew(t *testing.T) {
	i := int(1)
	if n := NewInt(i); !n.Valid || n.Int != i {
		t.Error("PtrInt: expected Valid to equal true")
	}
	f := float64(1)
	if n := NewFloat64(f); !n.Valid || n.Float64 != f {
		t.Error("PtrFloat64: expected Valid to equal true")
	}
	s := "a"
	if n := NewString(s); !n.Valid || n.String != s {
		t.Error("PtrString: expected Valid to equal true")
	}
	b := true
	if n := NewBool(b); !n.Valid || n.Bool != b {
		t.Error("PtrBool: expected Valid to equal true")
	}
	tt := time.Now()
	if n := NewTime(tt); !n.Valid || n.Time != tt {
		t.Error("PtrTime: expected Valid to equal true")
	}
}

func TestScan(t *testing.T) {
	{
		i := int64(1)
		n := new(Int)
		if err := n.Scan(i); err != nil {
			t.Error(err)
		}
		if int64(n.Int) != i {
			t.Error("Int Value: mismatch")
		}
		if !n.Valid {
			t.Error("Int Value: expected valid")
		}
		n = new(Int)
		if err := n.Scan(nil); err != nil {
			t.Error(err)
		}
		if n.Valid {
			t.Error("Int Value: expected invalid")
		}
	}
	{
		i := float64(1)
		n := new(Float64)
		if err := n.Scan(i); err != nil {
			t.Error(err)
		}
		if n.Float64 != i {
			t.Error("Float64 Value: mismatch")
		}
		if !n.Valid {
			t.Error("Float64 Value: expected valid")
		}
		n = new(Float64)
		if err := n.Scan(nil); err != nil {
			t.Error(err)
		}
		if n.Valid {
			t.Error("Float64 Value: expected invalid")
		}
	}
	{
		i := string("1")
		n := new(String)
		if err := n.Scan(i); err != nil {
			t.Error(err)
		}
		if n.String != i {
			t.Error("String Value: mismatch")
		}
		if !n.Valid {
			t.Error("String Value: expected valid")
		}
		n = new(String)
		if err := n.Scan(nil); err != nil {
			t.Error(err)
		}
		if n.Valid {
			t.Error("String Value: expected invalid")
		}
	}
	{
		i := time.Now()
		n := new(Time)
		if err := n.Scan(i); err != nil {
			t.Error(err)
		}
		if n.Time != i {
			t.Error("Time Value: mismatch")
		}
		if !n.Valid {
			t.Error("Time Value: expected valid")
		}
		n = new(Time)
		if err := n.Scan(nil); err != nil {
			t.Error(err)
		}
		if n.Valid {
			t.Error("Time Value: expected invalid")
		}
	}
}

func TestValue(t *testing.T) {
	{
		i := int(1)
		n := NewInt(i)
		v, err := n.Value()
		if err != nil {
			t.Error(err)
		}
		if int(v.(int64)) != i {
			t.Error("Int Value: mismatch")
		}
		n.Valid = false
		v, err = n.Value()
		if err != nil {
			t.Error(err)
		}
		if v != nil {
			t.Error("Int Value: expected nil")
		}
	}
	{
		i := float64(1)
		n := NewFloat64(i)
		v, err := n.Value()
		if err != nil {
			t.Error(err)
		}
		if v.(float64) != i {
			t.Error("Float64 Value: mismatch")
		}
		n.Valid = false
		v, err = n.Value()
		if err != nil {
			t.Error(err)
		}
		if v != nil {
			t.Error("Float64 Value: expected nil")
		}
	}
	{
		i := string("1")
		n := NewString(i)
		v, err := n.Value()
		if err != nil {
			t.Error(err)
		}
		if v.(string) != i {
			t.Error("String Value: mismatch")
		}
		n.Valid = false
		v, err = n.Value()
		if err != nil {
			t.Error(err)
		}
		if v != nil {
			t.Error("String Value: expected nil")
		}
	}
	{
		i := time.Now()
		n := NewTime(i)
		v, err := n.Value()
		if err != nil {
			t.Error(err)
		}
		if v.(time.Time) != i {
			t.Error("Time Value: mismatch")
		}
		n.Valid = false
		v, err = n.Value()
		if err != nil {
			t.Error(err)
		}
		if v != nil {
			t.Error("Time Value: expected nil")
		}
	}
}

func TestMarshal(t *testing.T) {
	null, err := json.Marshal(nil)
	if err != nil {
		t.Fatal(err)
	}
	{
		i := int(1)
		n := NewInt(i)
		a, err := n.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		b, err := json.Marshal(i)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(a, b) {
			t.Error("Int Marshal: mismatch")
		}
		n.Valid = false
		a, err = n.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(a, null) {
			t.Error("Int Marshal: null mismatch")
		}
	}
	{
		i := float64(1)
		n := NewFloat64(i)
		a, err := n.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		b, err := json.Marshal(i)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(a, b) {
			t.Error("Float64 Marshal: mismatch")
		}
		n.Valid = false
		a, err = n.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(a, null) {
			t.Error("Float64 Marshal: null mismatch")
		}
	}
	{
		i := string("1")
		n := NewString(i)
		a, err := n.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		b, err := json.Marshal(i)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(a, b) {
			t.Error("String Marshal: mismatch")
		}
		n.Valid = false
		a, err = n.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(a, null) {
			t.Error("String Marshal: null mismatch")
		}
	}
	{
		i := time.Now()
		n := NewTime(i)
		a, err := n.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		b, err := json.Marshal(i)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(a, b) {
			t.Error("Time Marshal: mismatch")
		}
		n.Valid = false
		a, err = n.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(a, null) {
			t.Error("Time Marshal: null mismatch")
		}
	}
}

func TestUnmarshal(t *testing.T) {
	null, err := json.Marshal(nil)
	if err != nil {
		t.Fatal(err)
	}

	{
		i := int(1)
		n := new(Int)
		b, err := json.Marshal(i)
		if err != nil {
			t.Error(err)
		}
		if err := n.UnmarshalJSON(b); err != nil {
			t.Error(err)
		}
		if n.Int != i || !n.Valid {
			t.Error("Int Unmarshal: failed", n.Int, i, n.Valid)
		}
		n = new(Int)
		if err := n.UnmarshalJSON(null); err != nil {
			t.Error(err)
		}
		if n.Valid {
			t.Error("Int Unmarshal: expected invalid")
		}
	}
	{
		i := float64(1)
		n := new(Float64)
		b, err := json.Marshal(i)
		if err != nil {
			t.Error(err)
		}
		if err := n.UnmarshalJSON(b); err != nil {
			t.Error(err)
		}
		if n.Float64 != i || !n.Valid {
			t.Error("Float64 Unmarshal: failed")
		}
		n = new(Float64)
		if err := n.UnmarshalJSON(null); err != nil {
			t.Error(err)
		}
		if n.Valid {
			t.Error("Float64 Unmarshal: expected invalid")
		}
	}
	{
		i := string(1)
		n := new(String)
		b, err := json.Marshal(i)
		if err != nil {
			t.Error(err)
		}
		if err := n.UnmarshalJSON(b); err != nil {
			t.Error(err)
		}
		if n.String != i || !n.Valid {
			t.Error("String Unmarshal: failed")
		}
		n = new(String)
		if err := n.UnmarshalJSON(null); err != nil {
			t.Error(err)
		}
		if n.Valid {
			t.Error("String Unmarshal: expected invalid")
		}
	}
	{
		i := time.Now()
		n := new(Time)
		b, err := json.Marshal(i)
		if err != nil {
			t.Error(err)
		}
		if err := n.UnmarshalJSON(b); err != nil {
			t.Error(err)
		}
		if !n.Time.Equal(i) || !n.Valid {
			t.Error("Time Unmarshal: failed")
		}
		n = new(Time)
		if err := n.UnmarshalJSON(null); err != nil {
			t.Error(err)
		}
		if n.Valid {
			t.Error("Time Unmarshal: expected invalid")
		}
	}
}

func TestTimeNow(t *testing.T) {
	n := &Time{}
	n.Now()
	if time.Since(n.Time) > time.Microsecond {
		t.Fatal("Err")
	}
}

// Test that Int marshals negative, zero and positive values.
func TestMarshalInt(t *testing.T) {
	n := &Int{Valid: true}
	for i := -10; i < 10; i++ {
		n.Int = i
		a, err := n.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		b, err := json.Marshal(n.Int)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(a, b) {
			t.Fatal(i)
		}
	}
}

func TestString(t *testing.T) {
	// Test that invalid UTF-8 is coerced to valid UTF-8,
	// matching the behavior of JSON Unmarshal.
	invalidUTF := []string{
		"\"hello\xffworld\"",
		"\"hello\xc2\xc2world\"",
		"\"hello\xc2\xffworld\"",
		"\"hello\\ud800world\"",
		"\"hello\\ud800\\ud800world\"",
		"\"hello\\ud800\\ud800world\"",
		"\"hello\xed\xa0\x80\xed\xb0\x80world\"",
	}
	for _, s := range invalidUTF {
		n := NewString(s)
		a, err := n.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		b, err := json.Marshal(n.String)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(a, b) {
			t.Errorf("TestString: failed to coerce: %s", s)
		}
	}
}

// Scan

func BenchmarkIntScan_Int64(b *testing.B) {
	var v Int
	for i := 0; i < b.N; i++ {
		if err := v.Scan(int64(123456)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIntScan_Bytes(b *testing.B) {
	var s interface{} = []byte("123456")
	var v Int
	for i := 0; i < b.N; i++ {
		if err := v.Scan(s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIntScan_Base(b *testing.B) {
	var s interface{} = []byte("123456")
	var v sql.NullInt64
	for i := 0; i < b.N; i++ {
		if err := v.Scan(s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFloat64Scan(b *testing.B) {
	var v Float64
	for i := 0; i < b.N; i++ {
		v.Scan(float64(123.456))
	}
}

func BenchmarkBoolScan(b *testing.B) {
	var v Bool
	for i := 0; i < b.N; i++ {
		v.Scan(true)
	}
}

func BenchmarkStringScan(b *testing.B) {
	var v String
	for i := 0; i < b.N; i++ {
		v.Scan("value string")
	}
}

func BenchmarkTimeScan(b *testing.B) {
	var v Time
	t := time.Now()
	for i := 0; i < b.N; i++ {
		v.Scan(t)
	}
}

// MarshalJSON

func BenchmarkIntMarshalJSON(b *testing.B) {
	v := Int{123456, true}
	for i := 0; i < b.N; i++ {
		v.MarshalJSON()
	}
}

func BenchmarkFloat64MarshalJSON(b *testing.B) {
	v := Float64{123.456, true}
	for i := 0; i < b.N; i++ {
		v.MarshalJSON()
	}
}

func BenchmarkBoolMarshalJSON(b *testing.B) {
	v := Bool{true, true}
	for i := 0; i < b.N; i++ {
		v.MarshalJSON()
	}
}

func BenchmarkStringMarshalJSON(b *testing.B) {
	v := String{"value string", true}
	for i := 0; i < b.N; i++ {
		v.MarshalJSON()
	}
}

func BenchmarkTimeMarshalJSON(b *testing.B) {
	v := Time{time.Now(), true}
	for i := 0; i < b.N; i++ {
		v.MarshalJSON()
	}
}

// UnmarshalJSON

func BenchmarkIntUnmarshalJSON(b *testing.B) {
	data := []byte("123456")
	var v Int
	for i := 0; i < b.N; i++ {
		v.UnmarshalJSON(data)
	}
}

func BenchmarkFloat64UnmarshalJSON(b *testing.B) {
	data := []byte("123.456")
	var v Float64
	for i := 0; i < b.N; i++ {
		v.UnmarshalJSON(data)
	}
}

func BenchmarkBoolUnmarshalJSON(b *testing.B) {
	data := []byte("true")
	var v Bool
	for i := 0; i < b.N; i++ {
		v.UnmarshalJSON(data)
	}
}

func BenchmarkStringUnmarshalJSON(b *testing.B) {
	data := []byte(`"value string"`)
	var v String
	for i := 0; i < b.N; i++ {
		v.UnmarshalJSON(data)
	}
}

func BenchmarkTimeUnmarshalJSON(b *testing.B) {
	data := []byte(`"2017-10-20T23:04:56.123456Z"`)
	var v Time
	for i := 0; i < b.N; i++ {
		v.UnmarshalJSON(data)
	}
}
