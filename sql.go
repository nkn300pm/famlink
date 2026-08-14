package main

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	
)

func GetCon() (*sql.DB, error) {
	
	db, err := sql.Open("mysql", "famlink:balance_shaft@/famlink?parseTime=true&loc=Local")
	return db, err
}

func isDuplicateError(err error) bool {
	// MySQL error code 1062 = duplicate entry
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
