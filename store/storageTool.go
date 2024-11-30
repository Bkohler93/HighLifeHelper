package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/bkohler93/highlifehelper/data/db"
)

type StorageToolStore struct {
	db *db.Queries
	// qMap  map[string]*db.Queries
	// dbMap map[string]*sql.DB
}

func NewStorageToolStore(database *sql.DB) *StorageToolStore {
	q := db.New(database)
	return &StorageToolStore{db: q}
}

func (s *StorageToolStore) UpdateStorage(groupName string, storageID int, storageName string, clearSlabsQty, clearBlocksQty, cloudySlabsQty, cloudyBlocksQty int) error {
	return s.db.UpdateStorage(context.Background(), db.UpdateStorageParams{
		GroupName:   groupName,
		StorageName: storageName,
		ClearSlabQty: sql.NullInt32{
			Int32: int32(clearSlabsQty),
			Valid: true,
		},
		ClearBlockQty: sql.NullInt32{
			Int32: int32(clearBlocksQty),
			Valid: true,
		},
		CloudySlabQty: sql.NullInt32{
			Int32: int32(cloudySlabsQty),
			Valid: true,
		},
		CloudyBlockQty: sql.NullInt32{
			Int32: int32(cloudyBlocksQty),
			Valid: true,
		},
		CreatedAt: time.Now(),
		ID:        int32(storageID),
	})
}

func (s *StorageToolStore) GetGroupStorage(groupName string) ([]db.Storage, error) {
	storages, err := s.db.GetGroupStorages(context.Background(), groupName)
	if err != nil {
		return storages, err
	}

	return storages, nil
}

func (s *StorageToolStore) CreateStorage(groupName string) (db.Storage, error) {
	return s.db.CreateStorage(context.Background(), db.CreateStorageParams{
		GroupName:   groupName,
		StorageName: "New Storage",
		ClearSlabQty: sql.NullInt32{
			Int32: 0,
			Valid: true,
		},
		ClearBlockQty: sql.NullInt32{
			Int32: 0,
			Valid: true,
		},
		CloudySlabQty: sql.NullInt32{
			Int32: 0,
			Valid: true,
		},
		CloudyBlockQty: sql.NullInt32{
			Int32: 0,
			Valid: true,
		},
		CreatedAt: time.Now(),
	})
}

func (s *StorageToolStore) DeleteStorage(groupName string, storageID int) error {
	return s.db.DeleteStorage(context.Background(), int32(storageID))
}
