package db

import "github.com/google/uuid"

const (
	EventDoc = "event"
	UserDoc  = "user"
)

var NameMap = map[string]uuid.UUID{
	EventDoc: uuid.MustParse("88c2333e-2bc2-4063-b865-719c24211d2c"),
	UserDoc:  uuid.MustParse("8d930d82-24e9-4fc1-824b-9e1253d4ee02"),
}

var IDMap map[uuid.UUID]string

func init() {
	IDMap = make(map[uuid.UUID]string, len(NameMap))
	for k, v := range NameMap {
		IDMap[v] = k
	}
}
