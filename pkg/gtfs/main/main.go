package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bridgelightcloud/bogie/pkg/gtfs"
)

func main() {
	st := time.Now()

	_, err := gtfs.OpenScheduleFromFile("gtfs_files/google_transit_20240812-20250110_v05.zip")
	if err != nil {
		log.Fatal("Error validating schedule: ", err)
	}

	et := time.Now()

	fmt.Println("Time taken to validate schedule: ", et.Sub(st))
}

func printAsFormattedJSON(data any) {
	if res, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Println(string(res))
	}
}
