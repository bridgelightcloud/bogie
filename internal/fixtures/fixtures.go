package fixtures

import (
	"time"

	"github.com/google/uuid"
)

var testUUIDString = "8d930d82-24e9-4fc1-824b-9e1253d4ee02"

var testTime = time.Now()

// GetTestUUID returns a UUID equal to 8d930d82-24e9-4fc1-824b-9e1253d4ee02
func GetTestUUID() uuid.UUID {
	return uuid.MustParse(testUUIDString)
}

// GetTestUUIDSlice returns a slce of size len composed of copies of
// the test UUID 8d930d82-24e9-4fc1-824b-9e1253d4ee02
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

// MaybeRefreshUUID checks if a UUID u is equal to uuid.NIL, and
// if it is, it replaces the value with the value of a new UUID
func MaybeRefreshUUID(u *uuid.UUID) {
	if *u == uuid.Nil {
		*u = uuid.New()
	}
}

func IntPtr(i int) *int {
	return &i
}
