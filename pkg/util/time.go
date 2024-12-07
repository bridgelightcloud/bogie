package util

import (
	"fmt"
	"time"
)

type TimeTracker func()

func TrackTime(message string)  TimeTracker {
	st := time.Now()
	return func() {
		et := time.Now()
		fmt.Printf("Time taken to %s: %s\n", message, et.Sub(st).String())
	}
}
