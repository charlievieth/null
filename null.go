// Package null implements nullable types for passing values between JSON
// and databases.
package null

// N.B.: Be mindful of which method receivers are pointers and which are
// values.

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

var nullLiteral = []byte("null")

func convertInt(value interface{}, bitSize int) (int64, error) {
	const maxInt64 = 1<<63 - 1

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
	const maxUint64 = (1<<64 - 1)

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
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int), nil
}

// MarshalJSON, marshals Int i into JSON.
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return nullLiteral, nil
	}
	return strconv.AppendInt([]byte{}, int64(i.Int), 10), nil
}

// MarshalJSON, unmarshals JSON data into Int i.
func (i *Int) UnmarshalJSON(data []byte) (err error) {
	if null(data) {
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
	var n sql.NullFloat64
	err := n.Scan(value)
	f.Float64, f.Valid = n.Float64, n.Valid
	return err
}

// Value, returns the database driver value of Float64 f.
func (f Float64) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.Float64, nil
}

// MarshalJSON, marshals Float64 f into JSON.
func (f Float64) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return nullLiteral, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'g', -1, 64)), nil
}

// UnmarshalJSON, unmarshals JSON data into Flaot64 f.
func (f *Float64) UnmarshalJSON(data []byte) (err error) {
	if null(data) {
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
func (s *String) Scan(value interface{}) error {
	var n sql.NullString
	err := n.Scan(value)
	s.String, s.Valid = n.String, n.Valid
	return err
}

// Value, returns the database driver value of String s.
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}

// MarshalJSON, marshals String s into JSON.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return nullLiteral, nil
	}
	return marshalString(s.String)
}

// UnmarshalJSON, unmarshals JSON data into String s.
func (s *String) UnmarshalJSON(data []byte) (err error) {
	if null(data) {
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
	if !b.Valid {
		return nil, nil
	}
	return b.Bool, nil
}

// MarshalJSON, marshals Bool b into JSON.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return nullLiteral, nil
	}
	if b.Bool {
		return []byte(`"true"`), nil
	}
	return []byte(`"false"`), nil
}

// UnmarshalJSON, unmarshals JSON data into Bool b.
func (b *Bool) UnmarshalJSON(data []byte) (err error) {
	if null(data) {
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
func (t *Time) Scan(value interface{}) error {
	var n mysql.NullTime
	err := n.Scan(value)
	t.Time, t.Valid = n.Time, n.Valid
	return err
}

// Value, returns the database driver value of Time t.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// MarshalJSON implements the json.Marshaler interface. The time is a quoted
// string in RFC 3339 format, with sub-second precision added if present.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return nullLiteral, nil
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface. The time is expected
// to be a quoted string in RFC 3339 format.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if null(data) {
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

// null, returns if data is a null JSON value.
func null(data []byte) bool {
	return bytes.Equal([]byte("null"), data)
}

// unquote, returns the form of JSON value b, for use by Int, Float64 and Bool.
// For JSON string values use unmarshalString.
func unquote(b []byte) []byte {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return b
	}
	return b[1 : len(b)-1]
}
