// Package null implements nullable types for passing values between JSON
// and databases.
package null

// N.B.: Be mindful of which method receivers are pointers and which are
// values.

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var nullLiteral = []byte("null")

// A Int is a nullable int that can be scanned into and from databases,
// and marshaled into and from JSON.
type Int struct {
	Int   int
	Valid bool
}

// NewInt, returns a new valid Int.
func NewInt(i int) Int {
	return Int{
		Int:   i,
		Valid: true,
	}
}

// PtrInt, returns a new Int from a pointer.
func PtrInt(i *int) Int {
	if i == nil {
		return Int{Valid: false}
	}
	return Int{
		Int:   *i,
		Valid: true,
	}
}

func (i *Int) Scan(value interface{}) error {
	if value == nil {
		i.Int, i.Valid = 0, false
		return nil
	}
	n, err := convertInt(value, strconv.IntSize)
	if err != nil {
		i.Int, i.Valid = 0, false
		return err
	}
	i.Int = int(n)
	i.Valid = true
	return nil
}

// Value, returns the database driver value of Int i.
func (i Int) Value() (driver.Value, error) {
	if i.Valid {
		return int64(i.Int), nil
	}
	return nil, nil
}

// MarshalJSON, marshals Int i into JSON.
func (i Int) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return strconv.AppendInt(nil, int64(i.Int), 10), nil
	}
	return nullLiteral, nil
}

// MarshalJSON, unmarshals JSON data into Int i.
func (i *Int) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		i.Int, i.Valid = 0, false
		return nil
	}
	n, err := parseInt(unquote(data), strconv.IntSize)
	if err == nil {
		i.Int = int(n)
	}
	i.Valid = (err == nil)
	return err
}

// Ptr, returns the value of Int i as a pointer.
func (i Int) Ptr() *int {
	if !i.Valid {
		return nil
	}
	n := i.Int
	return &n
}

// A Float64 is a nullable float64 that can be scanned into and from databases,
// and marshaled into and from JSON.
type Float64 struct {
	Float64 float64
	Valid   bool
}

// NewFloat64, returns a new valid Float64 with value f.
func NewFloat64(f float64) Float64 {
	return Float64{
		Float64: f,
		Valid:   true,
	}
}

// PtrFloat64, returns a new Float64 from a pointer.
func PtrFloat64(f *float64) Float64 {
	if f == nil {
		return Float64{Valid: false}
	}
	return Float64{
		Float64: *f,
		Valid:   true,
	}
}

// Scan, scans a database value into Float f.
func (f *Float64) Scan(value interface{}) error {
	if value == nil {
		f.Float64, f.Valid = 0, false
		return nil
	}
	ff, err := convertFloat(value, 64)
	if err != nil {
		f.Float64, f.Valid = 0, false
		return err
	}
	f.Float64, f.Valid = ff, err == nil
	return nil
}

// Value, returns the database driver value of Float64 f.
func (f Float64) Value() (driver.Value, error) {
	if f.Valid {
		return f.Float64, nil
	}
	return nil, nil
}

// MarshalJSON, marshals Float64 f into JSON.
func (f Float64) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return encodeFloat(f.Float64, 64)
	}
	return nullLiteral, nil
}

// UnmarshalJSON, unmarshals JSON data into Flaot64 f.
func (f *Float64) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		f.Float64, f.Valid = 0, false
		return nil
	}
	f.Float64, err = strconv.ParseFloat(string(unquote(data)), 64)
	f.Valid = (err == nil)
	return err
}

// Ptr, returns the value of Float64 f as a pointer.
func (f Float64) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	n := f.Float64
	return &n
}

// A Float32 is a nullable float32 that can be scanned into and from databases,
// and marshaled into and from JSON.
type Float32 struct {
	Float32 float32
	Valid   bool
}

// NewFloat32, returns a new valid Float32 with value f.
func NewFloat32(f float32) Float32 {
	return Float32{
		Float32: f,
		Valid:   true,
	}
}

// PtrFloat32, returns a new Float32 from a pointer.
func PtrFloat32(f *float32) Float32 {
	if f == nil {
		return Float32{Valid: false}
	}
	return Float32{
		Float32: *f,
		Valid:   true,
	}
}

// Scan, scans a database value into Float f.
func (f *Float32) Scan(value interface{}) error {
	if value == nil {
		f.Float32, f.Valid = 0, false
		return nil
	}
	ff, err := convertFloat(value, 32)
	if err != nil {
		f.Float32, f.Valid = 0, false
		return err
	}
	f.Float32, f.Valid = float32(ff), err == nil
	return nil
}

// Value, returns the database driver value of Float32 f.
func (f Float32) Value() (driver.Value, error) {
	if f.Valid {
		return float64(f.Float32), nil
	}
	return nil, nil
}

// MarshalJSON, marshals Float32 f into JSON.
func (f Float32) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return encodeFloat(float64(f.Float32), 32)
	}
	return nullLiteral, nil
}

// UnmarshalJSON, unmarshals JSON data into Flaot64 f.
func (f *Float32) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		f.Float32, f.Valid = 0, false
		return nil
	}
	var ff float64
	ff, err = strconv.ParseFloat(string(unquote(data)), 64)
	f.Valid = (err == nil)
	f.Float32 = float32(ff)
	return err
}

// Ptr, returns the value of Float32 f as a pointer.
func (f Float32) Ptr() *float32 {
	if !f.Valid {
		return nil
	}
	n := f.Float32
	return &n
}

// A String is a nullable string that can be scanned into and from databases,
// and marshaled into and from JSON.
type String struct {
	String string
	Valid  bool
}

// NewString, returns a new valid String with value s.
func NewString(s string) String {
	return String{
		String: s,
		Valid:  true,
	}
}

// PtrString, returns a new String from a pointer.
func PtrString(s *string) String {
	if s == nil {
		return String{Valid: false}
	}
	return String{
		String: *s,
		Valid:  true,
	}
}

// Scan, scans value into String s.
func (s *String) Scan(value interface{}) (err error) {
	if value == nil {
		s.String, s.Valid = "", false
		return nil
	}
	switch v := value.(type) {
	case string:
		s.String, s.Valid = v, true
	case []byte:
		s.String, s.Valid = string(v), true
	default:
		var n sql.NullString
		err = n.Scan(value)
		s.String, s.Valid = n.String, n.Valid
	}
	return err
}

// Value, returns the database driver value of String s.
func (s String) Value() (driver.Value, error) {
	if s.Valid {
		return s.String, nil
	}
	return nil, nil
}

// MarshalJSON, marshals String s into JSON.
func (s String) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return marshalString(s.String)
	}
	return nullLiteral, nil
}

// UnmarshalJSON, unmarshals JSON data into String s.
func (s *String) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		s.String, s.Valid = "", false
		return nil
	}
	s.String, err = unmarshalString(data)
	s.Valid = (err == nil)
	return err
}

// Ptr, returns the value of String s as a pointer.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	n := s.String
	return &n
}

// A Bool is a nullable bool that can be scanned into and from databases,
// and marshaled into and from JSON.
type Bool struct {
	Bool  bool
	Valid bool
}

// NewBool, returns a new valid Bool with value b.
func NewBool(b bool) Bool {
	return Bool{
		Bool:  b,
		Valid: true,
	}
}

// PtrBool, returns a new Bool from a pointer.
func PtrBool(b *bool) Bool {
	if b == nil {
		return Bool{Valid: false}
	}
	return Bool{
		Bool:  *b,
		Valid: true,
	}
}

// Scan, scans a database value into Bool b.
func (b *Bool) Scan(value interface{}) error {
	var n sql.NullBool
	err := n.Scan(value)
	b.Bool, b.Valid = n.Bool, n.Valid
	return err
}

// Value, returns the database driver value of Bool b.
func (b Bool) Value() (driver.Value, error) {
	if b.Valid {
		return b.Bool, nil
	}
	return nil, nil
}

// MarshalJSON, marshals Bool b into JSON.
func (b Bool) MarshalJSON() ([]byte, error) {
	if b.Valid {
		if b.Bool {
			return []byte(`"true"`), nil
		}
		return []byte(`"false"`), nil
	}
	return nullLiteral, nil
}

// UnmarshalJSON, unmarshals JSON data into Bool b.
func (b *Bool) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		b.Bool, b.Valid = false, false
		return nil
	}
	data = unquote(data)
	switch s := string(data); s {
	case "true":
		b.Bool, b.Valid = true, true
	case "false":
		b.Bool, b.Valid = false, true
	default:
		err = errors.New("null: cannot unmarshal '" + s + "' into type Bool")
	}
	return err
}

// Ptr, returns the value of Bool t as a pointer.
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	n := b.Bool
	return &n
}

// A Time is a nullable time.Time that can be scanned into and from databases,
// and marshaled into and from JSON.
type Time struct {
	Time  time.Time
	Valid bool
}

// NewTime, returns a new valid Time with time t.
func NewTime(t time.Time) Time {
	return Time{
		Time:  t,
		Valid: true,
	}
}

// PtrTime, returns a new Time from a pointer.
func PtrTime(t *time.Time) Time {
	if t == nil {
		return Time{Valid: false}
	}
	return Time{
		Time:  *t,
		Valid: true,
	}
}

// Scan, scans a database value into Time t.
func (t *Time) Scan(value interface{}) (err error) {
	if value == nil {
		t.Time, t.Valid = time.Time{}, false
		return
	}
	switch v := value.(type) {
	case time.Time:
		t.Time, t.Valid = v, true
	case []byte:
		t.Time, err = parseTime(string(v))
	case string:
		t.Time, err = parseTime(v)
	default:
		err = fmt.Errorf("Can't convert %T to time.Time", value)
	}
	t.Valid = (err == nil)
	return
}

// Time layout used for MySQL and SQLite
const TimeLayout = "2006-01-02 15:04:05.999999"

func parseTime(str string) (t time.Time, err error) {
	const base = "0000-00-00 00:00:00.0000000"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26:
		if str != base[:len(str)] {
			t, err = time.Parse(TimeLayout[:len(str)], str)
		}
	default:
		err = fmt.Errorf("invalid time string: %s", str)
	}
	return
}

// Value, returns the database driver value of Time t.
func (t Time) Value() (driver.Value, error) {
	if t.Valid {
		return t.Time, nil
	}
	return nil, nil
}

// MarshalJSON implements the json.Marshaler interface. The time is a quoted
// string in RFC 3339 format, with sub-second precision added if present.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return t.Time.MarshalJSON()
	}
	return nullLiteral, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. The time is expected
// to be a quoted string in RFC 3339 format.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}
	err = t.Time.UnmarshalJSON(data)
	t.Valid = (err == nil)
	return err
}

// Ptr, returns the value of Time t as a pointer.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	n := t.Time
	return &n
}

// Now, sets t's time to now.
func (t *Time) Now() {
	t.Time, t.Valid = time.Now(), true
}

// unquote, returns the form of JSON value b, for use by Int, Float64 and Bool.
// For JSON string values use unmarshalString.
func unquote(b []byte) []byte {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return b
	}
	return b[1 : len(b)-1]
}
