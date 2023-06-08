package models

import "time"

type Item struct {
	ItemID        int64
	CollectionID  int64
	Signature     string
	Urn           string
	Type          string
	SubType       string
	ObjectType    string
	ParentID      int64
	Mimetype      string
	Error         string
	SHA512        string
	Metadata      string
	CreationDate  time.Time
	LastModified  time.Time
	Disbled       bool
	Public        bool
	PublicActions []string
	Status        string
}
