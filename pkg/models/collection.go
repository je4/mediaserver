package models

type Collection struct {
	CollectionID    int64
	EstateID        int64
	Name            string
	Description     string
	SignaturePrefix string
	StorageID       int64
	JWTKey          string
	Secret          string
	Public          any
}
