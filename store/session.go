package store

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/bkohler93/highlifehelper/data/db"
	"golang.org/x/crypto/bcrypt"
)

type SessionStore struct {
	db *db.Queries
}

type Session struct {
	SessionID string
	LoginID   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewSessionStore(database *sql.DB) *SessionStore {
	q := db.New(database)
	return &SessionStore{db: q}
}

func (s *SessionStore) ValidateSession(sessionID string) bool {
	session, err := s.db.GetSession(context.Background(), sessionID)
	if err != nil {
		log.Println(err)
		return false
	}

	if time.Now().After(session.ExpiresAt) {
		return false
	}

	return true
}

func (s *SessionStore) CreateSession(uuid string, loginID string, createdAt time.Time, expiresAt time.Time) error {
	_, err := s.db.CreateSession(context.Background(), db.CreateSessionParams{
		Uuid:      uuid,
		LoginID:   loginID,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	})
	return err
}

func (s *SessionStore) GetUser(userID string) (db.User, error) {
	return s.db.GetUser(context.Background(), userID)
}

func (s *SessionStore) CreateUser(loginID, pwHash string) (db.User, error) {
	return s.db.CreateUser(context.Background(), db.CreateUserParams{
		LoginID:   loginID,
		PwHash:    pwHash,
		CreatedAt: time.Now(),
	})
}

func (s *SessionStore) Login(loginID, pw string) bool {
	user, err := s.GetUser(loginID)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(pw))
	return err == nil
}

func (s *SessionStore) GetSession(sessionID string) (db.Session, error) {
	return s.db.GetSession(context.Background(), sessionID)
}
