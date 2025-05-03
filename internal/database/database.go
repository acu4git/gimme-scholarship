package database

import (
	"fmt"
	"log"
	"os"

	"github.com/acu4git/gimme-scholarship/internal/domain/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
)

const (
	tableCourses = "courses"
)

var (
	host     string
	port     string
	username string
	password string
	dbname   string
)

type Database struct {
	sess *dbr.Session
}

func NewDatabase() (*Database, error) {
	if host = os.Getenv("DB_HOST"); host == "" {
		host = "localhost"
	}
	if port = os.Getenv("DB_PORT"); port == "" {
		port = "3308"
	}
	if username = os.Getenv("DB_USER"); username == "" {
		username = "root"
	}
	if password = os.Getenv("DB_PASSWORD"); password == "" {
		password = "root"
	}
	if dbname = os.Getenv("DB_NAME"); dbname == "" {
		dbname = "gimme_scholarship"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbname)
	conn, err := dbr.Open("mysql", dsn, nil)
	if err != nil {
		log.Println("failed at dbr.Open()")
		return nil, err
	}

	sess := conn.NewSession(nil)

	return &Database{sess: sess}, nil
}

func (db *Database) CreateUser(user model.User) error {
	return nil
}
