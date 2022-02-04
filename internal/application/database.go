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
		configs["DB_USER"],
		configs["DB_PASSWORD"],
		configs["DB_HOST"],
		configs["DB_PORT"],
		configs["DB_NAME"],
	)
}
