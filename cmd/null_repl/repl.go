package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

func InitDatabase() {
	// "root:password@/dbname"
	db, err := sql.Open("mysql", "root@/null_test")
	if err != nil {
		Fatal(err)
	}
	defer db.Close()

	if err := InitializeNumericTable(db); err != nil {
		Fatal(err)
	}
}

func main() {
	InitDatabase()
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
