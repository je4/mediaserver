package models

import (
	"emperror.dev/errors"
	"sync"
)

type StoragesDatabase interface {
	StoragesLoadAll(stors *Storages) error
}

func NewStorages(p *Pool, db StoragesDatabase) (*Storages, error) {
	storages := &Storages{
		Pool:     p,
		RWMutex:  sync.RWMutex{},
		db:       db,
		storages: map[string]*Storage{},
	}
	return storages, db.StoragesLoadAll(storages)
}

type Storages struct {
	*Pool
	sync.RWMutex
	db       StoragesDatabase
	storages map[string]*Storage
}

func (stors *Storages) Add(stor *Storage) {
	stors.storages[stor.Name] = stor
}

func (stors *Storages) Clear() {
	stors.storages = map[string]*Storage{}
}

func (stors *Storages) Get(name string) (*Storage, error) {
	stors.RLock()
	defer stors.RUnlock()
	stor, ok := stors.storages[name]
	if !ok {
		return nil, errors.Wrapf(notFound, "collection '%s'", name)
	}
	return stor, nil
}
