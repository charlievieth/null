package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// SMALLINT[(M)] [UNSIGNED] [ZEROFILL]
// MEDIUMINT[(M)] [UNSIGNED] [ZEROFILL]
// INTEGER[(M)] [UNSIGNED] [ZEROFILL]
// BIGINT[(M)] [UNSIGNED] [ZEROFILL]

const CreateTableStmt = `
	CREATE TABLE IF NOT EXISTS null_test (
		id           INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
	    test_int8    TINYINT,
	    test_int16   SMALLINT,
	    test_int32   INTEGER,
	    test_int64   BIGINT,
	    test_uint8   TINYINT UNSIGNED,
	    test_uint16  SMALLINT UNSIGNED,
	    test_uint32  INTEGER UNSIGNED,
	    test_uint64  BIGINT UNSIGNED
	);`

func InitializeNumericTable(db *sql.DB) error {
	const InsertStmt = `
	INSERT INTO null_test (test_int8, test_int16, test_int32, test_int64, test_uint8, test_uint16, test_uint32, test_uint64)
	VALUES
	    (1, 1, 1, 1, 1, 1, 1, 1),                                                                    -- Simple
	    (127, 32767, 2147483647, 9223372036854775807, 255, 65535, 4294967295, 18446744073709551615), -- Max
	    (-128, -32768, -2147483648, -9223372036854775808, 0, 0, 0, 0),                               -- Min
	    (NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);                                            -- NULL`

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	exit := func(err error) error {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec("DROP TABLE IF EXISTS null_test;"); err != nil {
		return exit(err)
	}
	if _, err := tx.Exec(CreateTableStmt); err != nil {
		return exit(err)
	}
	if _, err := tx.Exec(InsertStmt); err != nil {
		return exit(err)
	}
	return tx.Commit()
}

func CreateDatabase() error {
	db, err := sql.Open("mysql", "root@/")
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS null_test"); err != nil {
		return err
	}
	return nil
}

func InitDatabase() (*sql.DB, error) {
	// "root:password@/dbname"
	db, err := sql.Open("mysql", "root@/null_test")
	if err != nil {
		return nil, err
	}

	if err := InitializeNumericTable(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

type NullUint64 struct {
	Uint16 uint16
	Valid  bool
}

func (n *NullUint64) Scan(value interface{}) error {
	fmt.Printf("%T -- %v\n", value, value)
	return nil
}

func TestScanValue() {
	db, err := InitDatabase()
	if err != nil {
		Fatal(err)
	}
	defer db.Close()

	var nu NullUint64
	// err = db.QueryRow("SELECT test_uint32 FROM null_test LIMIT 1;").Scan(&nu)2
	err = db.QueryRow("SELECT test_uint32 FROM null_test ORDER BY test_uint32 DESC LIMIT 1;").Scan(&nu)
	if err != nil {
		Fatal(err)
	}
	fmt.Println("Ok")
}

// Parse Int

// ErrRange indicates that a value is out of range for the target type.
var ErrRange = errors.New("value out of range")

// ErrSyntax indicates that a value does not have the right syntax for the target type.
var ErrSyntax = errors.New("invalid syntax")

// A NumError records a failed conversion.
type NumError struct {
	Func string // the failing function (ParseBool, ParseInt, ParseUint, ParseFloat)
	Num  string // the input
	Err  error  // the reason the conversion failed (ErrRange, ErrSyntax)
}

func (e *NumError) Error() string {
	// TODO: Remove use of strconv
	return "strconv." + e.Func + ": " + "parsing " + strconv.Quote(e.Num) + ": " + e.Err.Error()
}

const intSize = 32 << (^uint(0) >> 63)

// IntSize is the size in bits of an int or uint value.
const IntSize = intSize

const maxUint64 = (1<<64 - 1)

func ParseUint(s string, bitSize int) (uint64, error) {
	const base = 10
	const cutoff = maxUint64/10 + 1

	var i int
	var n uint64
	var maxVal uint64
	var err error

	if bitSize == 0 {
		bitSize = int(IntSize)
	}

	if len(s) < 1 {
		err = ErrSyntax
		goto Error
	}

	maxVal = 1<<uint(bitSize) - 1

	for ; i < len(s); i++ {
		var v byte
		d := s[i]
		if '0' <= s[i] && s[i] <= '9' {
			v = d - '0'
		} else {
			n = 0
			err = ErrSyntax
			goto Error
		}
		if n >= cutoff {
			// n*base overflows
			n = maxUint64
			err = ErrRange
			goto Error
		}
		n *= uint64(base)

		n1 := n + uint64(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = maxUint64
			err = ErrRange
			goto Error
		}
		n = n1
	}

	return n, nil

Error:
	return n, &NumError{"ParseUint", s, err}
}

func main() {
	const x = 32 << (^uint(0) >> 63)
	a := ^uint32(0)
	b := (^uint32(0) >> 63)
	c := 32 << (^uint32(0) >> 63)
	fmt.Println("A:", a)
	fmt.Println("B:", b)
	fmt.Println("C:", c)
	fmt.Println(IntSize)
	fmt.Println(uint64(1 << uint(strconv.IntSize-1)))
	fmt.Println(uint64(1 << uint(64-1)))
}

func Fatal(err interface{}) {
	if err == nil {
		return
	}
	errMsg := "Error"
	if _, file, line, _ := runtime.Caller(1); file != "" {
		errMsg = fmt.Sprintf("Error (%s:#%d)", filepath.Base(file), line)
	}
	switch e := err.(type) {
	case string, error, fmt.Stringer:
		fmt.Fprintf(os.Stderr, "%s: %s\n", errMsg, e)
	default:
		fmt.Fprintf(os.Stderr, "%s: %#v\n", errMsg, e)
	}
	os.Exit(1)
}

/*
type NumericEntry struct {
	Int8   *uint64 `db:"test_int8"`
	Int16  *uint64 `db:"test_int16"`
	Int32  *uint64 `db:"test_int32"`
	Int64  *uint64 `db:"test_int64"`
	Uint8  *uint64 `db:"test_uint8"`
	Uint16 *uint64 `db:"test_uint16"`
	Uint32 *uint64 `db:"test_uint32"`
	Uint64 *uint64 `db:"test_uint64"`
}

func uintPtr(value interface{}) *uint64 {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case int16:
		u := uint64(v)
		return &u
	case int32:
		u := uint64(v)
		return &u
	case int64:
		u := uint64(v)
		return &u
	case uint16:
		u := uint64(v)
		return &u
	case uint32:
		u := uint64(v)
		return &u
	case uint64:
		u := uint64(v)
		return &u
	default:
		panic(fmt.Sprintf("invalid numeric type %T: %v", value, value))
	}
}

var entries = []NumericEntry{}
*/
