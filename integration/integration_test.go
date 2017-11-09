package integration

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	DatabaseName = fmt.Sprintf("null_test_%d", time.Now().Unix())

	// Format: [username[:password]@][protocol[(address)]]
	//
	// If the database is provided (/...) use it instead
	// of creating a new one for tests.
	MySQLDSN = os.Getenv("NULL_TEST_MYSQL_DSN")
)

func ParseDSN(dsn string) (user, db string, err error) {
	a := strings.Split(strings.TrimSuffix(dsn, "/"), "/")
	switch len(a) {
	case 1:
		user = a[0]
	case 2:
		user = a[0]
		db = a[1]
	default:
		err = errors.New("invalid DSN: " + dsn)
	}
	return
}

type MySQL struct {
	DatabaseName string
	DSN          string
	CleanupDB    string // remove db on test exit
	db           *sql.DB
}

func NewMySQL(dsn string) (*MySQL, error) {
	user, db, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	_, _ = user, db
	return nil, nil
}

func (m *MySQL) CreateDatabase() error {
	db, err := sql.Open("mysql", "root@/")
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec("CREATE DATABASE " + m.DatabaseName); err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	// setup
	code := m.Run()
	// teardown
	os.Exit(code)
}
