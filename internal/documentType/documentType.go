package documentType

import "github.com/google/uuid"

const (
	Event = "event"
	Note  = "note"
	User  = "user"
)

var NameMap = map[string]uuid.UUID{
	Event: uuid.MustParse("88c2333e-2bc2-4063-b865-719c24211d2c"),
	Note:  uuid.MustParse("22322bb8-a869-43b2-9991-8466fd02a19b"),
	User:  uuid.MustParse("8d930d82-24e9-4fc1-824b-9e1253d4ee02"),
}

var IDMap map[uuid.UUID]string

func init() {
	IDMap = make(map[uuid.UUID]string, len(NameMap))
	for k, v := range NameMap {
		IDMap[v] = k
	}
}
