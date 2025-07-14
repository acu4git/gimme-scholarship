package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
)

func init() {
	user, password, host, port, dbname := "root", "root", "127.0.0.1", "3309", "gimme_scholarship_test"
	if os.Getenv("FAKE_DB_USER") != "" {
		user = os.Getenv("FAKE_DB_USER")
	}
	if os.Getenv("FAKE_DB_PASSWORD") != "" {
		password = os.Getenv("FAKE_DB_PASSWORD")
	}
	if os.Getenv("FAKE_DB_HOST") != "" {
		host = os.Getenv("FAKE_DB_HOST")
	}
	if os.Getenv("FAKE_DB_PORT") != "" {
		port = os.Getenv("FAKE_DB_PORT")
	}
	if os.Getenv("FAKE_DB_NAME") != "" {
		dbname = os.Getenv("FAKE_DB_NAME")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		user,
		password,
		host,
		port,
		dbname,
	)
	txdb.Register("txdb", "mysql", dsn)
}

type FakeDatabase struct {
	*Database
}

func NewFakeDatabase() *FakeDatabase {
	fdb := &FakeDatabase{}
	fdb.connect()
	return fdb
}

func (f *FakeDatabase) connect() {
	db, err := sql.Open("txdb", strconv.FormatInt(time.Now().UnixNano(), 10))
	if err != nil {
		panic(err)
	}
	conn := &dbr.Connection{
		DB:            db,
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.MySQL,
	}
	sess := conn.NewSession(nil)
	if err = sess.Ping(); err != nil {
		panic(err)
	}
	f.Database = &Database{sess: sess}
}

func (f *FakeDatabase) TruncateTables() {
	if err := f.sess.Close(); err != nil {
		panic(err)
	}
}
