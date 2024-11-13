package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bkohler93/highlifehelper/data/db"
	"github.com/pressly/goose/v3"
)

type StorageToolStore struct {
	qMap  map[string]*db.Queries
	dbMap map[string]*sql.DB
}

func NewStorageToolStore() *StorageToolStore {
	return &StorageToolStore{
		qMap:  map[string]*db.Queries{},
		dbMap: map[string]*sql.DB{},
	}
}

func (s *StorageToolStore) Connect(groupName string) error {
	if _, ok := s.qMap[groupName]; ok {
		return nil
	}

	database, err := sql.Open("sqlite", fmt.Sprintf(os.Getenv("DATABASE_PATH") + "/groupStorage/%s.db", groupName))
	if err != nil {
		log.Fatalf("error opening database for group %s: %v", groupName, err)
		return err
	}

	goose.SetDialect("sqlite")
	/* uncomment if you wish to migrate down */
	// if err := goose.Down(database, "data/sql/migrations/groupStorage"); err != nil {
	// 	log.Fatalf("failed to run migrations on db %s.db: %v", groupName, err)
	// }

	if err := goose.Up(database, "data/sql/migrations/groupStorage"); err != nil {
		log.Fatalf("failed to run migrations on db %s.db: %v", groupName, err)
	}

	queries := db.New(database)
	s.dbMap[groupName] = database
	s.qMap[groupName] = queries

	return nil
}

func (s *StorageToolStore) Disconnect(groupName string) error {
	if db, ok := s.dbMap[groupName]; !ok {
		return fmt.Errorf("no database found for group %s", groupName)
	} else {
		return db.Close()
	}
}

func (s *StorageToolStore) UpdateStorage(groupName string, storageID int, storageName string, clearSlabsQty, clearBlocksQty, cloudySlabsQty, cloudyBlocksQty int) error {
	return s.qMap[groupName].UpdateStorage(context.Background(), db.UpdateStorageParams{
		GroupName:   groupName,
		StorageName: storageName,
		ClearSlabQty: sql.NullInt64{
			Int64: int64(clearSlabsQty),
			Valid: true,
		},
		ClearBlockQty: sql.NullInt64{
			Int64: int64(clearBlocksQty),
			Valid: true,
		},
		CloudySlabQty: sql.NullInt64{
			Int64: int64(cloudySlabsQty),
			Valid: true,
		},
		CloudyBlockQty: sql.NullInt64{
			Int64: int64(cloudyBlocksQty),
			Valid: true,
		},
		CreatedAt: time.Now(),
		ID:        int64(storageID),
	})
}

func (s *StorageToolStore) GetGroupStorage(groupName string) ([]db.Storage, error) {
	storages, err := s.qMap[groupName].GetGroupStorages(context.Background(), groupName)
	if err != nil {
		return storages, err
	}

	return storages, nil
}

func (s *StorageToolStore) CreateStorage(groupName string) (db.Storage, error) {
	return s.qMap[groupName].CreateStorage(context.Background(), db.CreateStorageParams{
		GroupName:      groupName,
		StorageName:    "New Storage",
		ClearSlabQty:   sql.NullInt64{},
		ClearBlockQty:  sql.NullInt64{},
		CloudySlabQty:  sql.NullInt64{},
		CloudyBlockQty: sql.NullInt64{},
		CreatedAt:      time.Now(),
	})
}

func (s *StorageToolStore) DeleteStorage(groupName string, storageID int) error {
	return s.qMap[groupName].DeleteStorage(context.Background(), int64(storageID))
}
