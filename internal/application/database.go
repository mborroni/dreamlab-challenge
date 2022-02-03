package application

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func buildDBConnections() {
	rdb, err := sql.Open("postgres", dsn("user", "pwd", "host"))
	if err != nil {
		panic(err)
	}
	db = rdb
}

func dsn(user, password, host string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost",
		"5432",
		"postgres",
		"password",
		"ip2location",
	)
}
