package orm

type Entity interface {
	// CollectionName refers to PocketBase entity collection's name
	CollectionName() string
}
