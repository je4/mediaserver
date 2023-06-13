package models

import (
	"emperror.dev/errors"
	"time"
)

func NewPool(
	collectionsDatabase CollectionsDatabase,
	storagesDatabase StoragesDatabase,
	itemsDatabase ItemsDatabase,
	itemsCacheSize int,
	itemsCacheExpiration time.Duration,
) (*Pool, error) {
	var err error
	p := &Pool{}

	p.Storages, err = NewStorages(p, storagesDatabase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create storages")
	}
	p.Collections, err = NewCollections(p, collectionsDatabase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create collections")
	}
	p.Items, err = NewItems(p, itemsDatabase, itemsCacheSize, itemsCacheExpiration)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create items")
	}
	return p, nil
}

type Pool struct {
	Storages    *Storages
	Collections *Collections
	Items       *Items
}
