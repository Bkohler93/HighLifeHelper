package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pressly/goose/v3"
)

const (
	port = ":8080"
)

func main() {
	db, err := sql.Open("sqlite", "./data/databases/session/session.db")
	if err != nil {
		log.Fatalf("error opening database:%s", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("error pinging database:%s", err)
	}

	goose.SetDialect("sqlite")

	/* uncomment if you wish to migrate down */
	// if err := goose.Down(db, "data/sql/migrations/session"); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	if err := goose.Up(db, "data/sql/migrations/session"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	cfg := NewConfig(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Method("GET", "/", Handler(cfg.HomeHandler))

	r.Method("POST", "/login", Handler(cfg.LoginHandler))
	r.Method("GET", "/register", Handler(cfg.GetRegisterHandler))
	r.Method("POST", "/register", Handler(cfg.PostRegisterHandler))

	r.Method("GET", "/storage-tool", Handler(cfg.StorageToolHandler))
	r.Method("GET", "/storage-add-card", Handler(cfg.StorageAddCardHandler))
	r.Method("DELETE", "/storage-delete-card", Handler(cfg.StorageDeleteCardHandler))
	r.Method("POST", "/storage-tool-update", Handler(cfg.StorageToolUpdateHandler))
	r.Method("GET", "/ws-storage-tool", Handler(cfg.StorageToolHandleWs))

	r.Method("GET", "/cook-tool", Handler(cfg.CookToolHandler))
	r.Method("POST", "/cook-calculate", Handler(cfg.CookCalculateHandler))
	r.Method("DELETE", "/cook-delete-recipe", Handler(cfg.CookDeleteRecipeHandler))
	r.Method("POST", "/cook-add-recipe", Handler(cfg.CookAddRecipeHandler))

	r.Method("GET", "/truck-tool", Handler(cfg.TruckToolHandler))

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
