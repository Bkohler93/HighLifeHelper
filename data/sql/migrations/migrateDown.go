package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	db, err := sql.Open("postgres", "user=postgres.mfsffpjrfalpmwjeaaio password=H0iBNPGEjc3eXY host=aws-0-us-east-2.pooler.supabase.com port=6543 dbname=postgres")
	if err != nil {
		log.Fatalf("error opening database:%s", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("error pinging database:%s", err)
	}

	goose.SetDialect("postgres")

	if err := goose.DownTo(db, "data/sql/migrations", 0); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}
