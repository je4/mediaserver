package models

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/bluele/gcache"
	"time"
)

type ItemsDatabase interface {
	LoadItem(collection string, signature string) (*Item, error)
}

func NewItems(p *Pool, db ItemsDatabase, cacheSize int, cacheExpiration time.Duration) (*Items, error) {
	loaderFunc := func(i any) (any, error) {
		collsig, ok := i.(string)
		if !ok {
			return nil, errors.Errorf("invalid key %v", i)
		}
		item, err := db.LoadItem(collsig, "")
		if err != nil {
			return nil, errors.Wrapf(err, "cannot load item '%s'", collsig)
		}
		return item, nil
	}
	items := &Items{
		Pool:  p,
		cache: gcache.New(cacheSize).ARC().LoaderFunc(loaderFunc).Expiration(cacheExpiration).Build(),
	}
	return items, nil
}

type Items struct {
	*Pool
	cache       gcache.Cache
	collections *Collections
}

func (items *Items) Load(collection, signature string) (*Item, error) {
	collsig := fmt.Sprintf("%s/%s", collection, signature)
	itemAny, err := items.cache.Get(collsig)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot load item '%s'", collsig)
	}
	item, ok := itemAny.(*Item)
	if !ok {
		return nil, errors.Errorf("got invalid cache entry for item '%s' - %v", collsig, itemAny)

	}
	return item, nil
}
