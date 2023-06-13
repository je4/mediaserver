package models

import "sync"
import "emperror.dev/errors"

type CollectionsDatabase interface {
	CollectionsLoadAll(colls *Collections) error
}

func NewCollections(p *Pool, db CollectionsDatabase) (*Collections, error) {
	var collections = &Collections{
		Pool:    p,
		RWMutex: sync.RWMutex{},
		db:      db,
		cols:    map[string]*Collection{},
	}

	return collections, errors.WithStack(db.CollectionsLoadAll(collections))
}

type Collections struct {
	*Pool
	sync.RWMutex
	db   CollectionsDatabase
	cols map[string]*Collection
}

func (cols *Collections) Add(coll *Collection) {
	cols.cols[coll.Name] = coll
}

func (cols *Collections) Clear() {
	cols.cols = map[string]*Collection{}
}

/*
func (cols *Collections) New(newColl *Collection) (*Collection, error) {
	result, err := cols.db.CollectionsNew(newColl)
	return result, err
}
*/

func (cols *Collections) Get(name string) (*Collection, error) {
	cols.RLock()
	defer cols.RUnlock()
	coll, ok := cols.cols[name]
	if !ok {
		return nil, errors.Wrapf(notFound, "collection '%s'", name)
	}
	return coll, nil
}
