package models

type Storage struct {
	StorageID  int64
	Name       string
	FileBase   string
	DataDir    string
	SubItemDir string
	TempDir    string
}
