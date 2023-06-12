package models

import "sync"
import "emperror.dev/errors"

type CollectionsDatabase interface {
	LoadAll(colls *Collections) error
}

func NewCollections(db CollectionsDatabase) (*Collections, error) {
	var collections = &Collections{
		RWMutex:     sync.RWMutex{},
		db:          db,
		collections: map[string]*Collection{},
	}

	return collections, errors.WithStack(db.LoadAll(collections))
}

type Collections struct {
	sync.RWMutex
	db          CollectionsDatabase
	collections map[string]*Collection
}

func (colls *Collections) Add(coll *Collection) {
	colls.collections[coll.Name] = coll
}

func (colls *Collections) Clear() {
	colls.collections = map[string]*Collection{}
}

func (colls *Collections) Get(name string) (*Collection, error) {
	colls.RLock()
	defer colls.RUnlock()
	coll, ok := colls.collections[name]
	if !ok {
		return nil, errors.Wrapf(notFound, "collection '%s'", name)
	}
	return coll, nil
}
