package db

import (
	"database/sql"
	"fmt"
	"pengoe/internal/utils"
	"sync"

	_ "github.com/libsql/libsql-client-go/libsql"
)

/*
dbManager is a struct that manages the database connection.
*/
type dbManager struct {
	db   *sql.DB
	once sync.Once
}

/*
NewDBManager is a function that returns a new DBManager.
*/
func NewDBManager() *dbManager {
	return &dbManager{}
}

/*
Global Manager.
*/
var Manager *dbManager = NewDBManager()

/*
GetDB is a function that returns a database connection.
*/
func (manager *dbManager) GetDB() (*sql.DB, error) {
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
func (manager *dbManager) connect() error {
	url := utils.Env.DBUrl
	token := utils.Env.DBToken

	connectionStr := fmt.Sprintf("%s?authToken=%s", url, token)

	db, dbErr := sql.Open("libsql", connectionStr)
	if dbErr != nil {
		return dbErr
	}

	manager.db = db
	return nil
}
