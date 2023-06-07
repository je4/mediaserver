package database

import (
	"emperror.dev/errors"
)

var notFound = errors.New("entity not found")

func isNotFound(err error) bool {
	return errors.Cause(err) == notFound
}
