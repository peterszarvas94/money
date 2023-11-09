package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"

	_ "github.com/libsql/libsql-client-go/libsql"
)

/*
DBManager is a struct that manages the database connection.
*/
type DBManager struct {
	db   *sql.DB
	once sync.Once
}

/*
NewDBManager is a function that returns a new DBManager.
*/
func NewDBManager() *DBManager {
	return &DBManager{}
}

/*
GetDB is a function that returns a database connection.
*/
func (manager *DBManager) GetDB() (*sql.DB, error) {
	var err error
	manager.once.Do(func() {
		err = manager.connect()
	})
	if err != nil {
		return nil, err
	}
	return manager.db, nil
}

/*
connect is a function that opens a connection to the database.
It uses the DB_URL and DB_TOKEN environment variables (Turso db).
*/
func (manager *DBManager) connect() error {
	url, urlFound := os.LookupEnv("DB_URL")
	if !urlFound {
		return errors.New("DB_URL not found")
	}

	token, token_found := os.LookupEnv("DB_TOKEN")
	if !token_found {
		return errors.New("DB_TOKEN not found")
	}

	connectionStr := fmt.Sprintf("%s?authToken=%s", url, token)

	db, dbErr := sql.Open("libsql", connectionStr)
	if dbErr != nil {
		return dbErr
	}

	manager.db = db
	return nil
}
