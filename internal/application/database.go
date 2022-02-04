package application

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func buildDBConnections() {
	rdb, err := sql.Open("postgres", dsn())
	if err != nil {
		panic(err)
	}
	db = rdb
}

func dsn() string {
	fmt.Print(os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DATABASE"),
		os.Getenv("PORT"),
		os.Getenv("DBNAME"))

	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DATABASE"),
		os.Getenv("PORT"),
		os.Getenv("DBNAME"),
	)
}
