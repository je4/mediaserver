package models

import (
	"emperror.dev/errors"
	"sync"
)

type StoragesDatabase interface {
	LoadAll(stors *Storages) error
}

func NewStorages(db StoragesDatabase) (*Storages, error) {
	return &Storages{
		RWMutex:  sync.RWMutex{},
		db:       db,
		storages: map[string]*Storage{},
	}, nil
}

type Storages struct {
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
