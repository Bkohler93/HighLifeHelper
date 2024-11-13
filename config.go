package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/bkohler93/highlifehelper/store"
	"github.com/gorilla/websocket"
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
	conns            map[string][]Connection //map[GroupName][]Connection
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

	wsUrl := os.Getenv("WS_URL")

	return &Config{
		SessionStore:     sessionStore,
		StorageToolStore: store.NewStorageToolStore(),
		tpl:              tpl,
		conns:            make(map[string][]Connection),
		WsUrl:            wsUrl,
	}
}

func (c *Config) CheckAndCloseWebSocket(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
		}
		sessionID := cookie.Value
		s, _ := c.SessionStore.GetSession(sessionID)
		groupName := s.LoginID

		conns := c.conns[groupName]
		newConns := conns[:0]

		for _, conn := range conns {
			if conn.sessionID == sessionID {
				conn.conn.Close()
			} else {
				newConns = append(newConns, conn)
			}
		}
		c.conns[groupName] = newConns
		next.ServeHTTP(w, r)
	})
}
