package main

import (
	"database/sql"
	"html/template"
	"log"
	"os"

	"github.com/bkohler93/highlifehelper/store"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	baseHTML     = "static/html/"
	basePartials = "static/html/partials/"
	baseViews    = "static/html/views/"
)

type Config struct {
	SessionStore     *store.SessionStore
	StorageToolStore *store.StorageToolStore
	tpl              *template.Template
	conns            map[string][]Connection
	WsUrl            string
}

type Connection struct {
	conn      *websocket.Conn
	sessionID string
}

func NewConfig(db *sql.DB) *Config {
	sessionStore := store.NewSessionStore(db)
	tpl, err := template.ParseFiles(baseHTML+"index.html", basePartials+"registerUser.html", basePartials+"cookRecipe.html", basePartials+"cookSuppliesNeeded.html", basePartials+"sidebar.html", basePartials+"storagePropertyCard.html", baseViews+"cookTool.html", baseViews+"storageTool.html", baseViews+"truckTool.html")
	if err != nil {
		log.Fatalf("error parsing templates:%v", err)
	}

	godotenv.Load()
	wsUrl := os.Getenv("WS_URL")

	return &Config{
		SessionStore:     sessionStore,
		StorageToolStore: store.NewStorageToolStore(),
		tpl:              tpl,
		conns:            make(map[string][]Connection),
		WsUrl:            wsUrl,
	}
}
