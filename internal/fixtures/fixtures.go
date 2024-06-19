package fixtures

import (
	"time"

	"github.com/google/uuid"
)

var testTime = time.Now()

func GetTestUUID() uuid.UUID {
	return uuid.MustParse("8d930d82-24e9-4fc1-824b-9e1253d4ee02")
}

func GetTestUUIDSlice(len int) []uuid.UUID {
	uuids := make([]uuid.UUID, len)
	for i := 0; i < len; i++ {
		uuids[i] = GetTestUUID()
	}
	return uuids
}

func GetTestTime() time.Time {
	return testTime
}

func GetTestTimePtr() *time.Time {
	return &testTime
}