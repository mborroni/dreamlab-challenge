package application

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func buildDBConnections() {
	rdb, err := sql.Open("postgres", dsn())
	if err != nil {
		panic(err)
	}
	db = rdb
}

func dsn() string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		configs["user"],
		configs["password"],
		configs["database"],
		configs["port"],
		configs["dbname"],
	)
}
