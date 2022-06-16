package models

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type dbConfig struct {
	host     string
	port     uint
	database string
	username string
	password string
}

var _dbConfig dbConfig = dbConfig{
	host:     "127.0.0.1",
	port:     3306,
	database: "gen-bike",
	username: "root",
}

func DBReadEnv() {
	if host := os.Getenv("MYSQL_HOST"); host != "" {
		_dbConfig.host = host
	}

	if port := os.Getenv("MYSQL_HOST"); port != "" {
		nerPort, err := strconv.Atoi(port)
		if err == nil {
			_dbConfig.port = uint(nerPort)
		}
	}

	if dbname := os.Getenv("MYSQL_DATABASE"); dbname != "" {
		_dbConfig.database = dbname
	}

	if user := os.Getenv("MYSQL_USERNAME"); user != "" {
		_dbConfig.username = user
	}

	if pass := os.Getenv("MYSQL_PASSWORD"); pass != "" {
		_dbConfig.password = pass
	}
}

func CreateDBConnection() (*sql.DB, error) {
	password := _dbConfig.password
	if password != "" {
		password = ":" + password
	}
	dataSrc := fmt.Sprintf("%s%s@tcp(%s:%d)/%s", _dbConfig.username, password, _dbConfig.host, _dbConfig.port, _dbConfig.database)
	return sql.Open("mysql", dataSrc)
}
