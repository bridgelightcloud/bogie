package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bridgelightcloud/bogie/pkg/gtfs"
	"github.com/bridgelightcloud/bogie/pkg/util"
)

func main() {
	tt := util.TrackTime("create GTFS collection")
	defer tt()

	gtfsDir := "gtfs_files"

	if _, err := os.Stat(gtfsDir); err != nil {
		log.Fatalf("Error finding %s: %s \n", gtfsDir, err.Error())
	}

	zipFiles, err := filepath.Glob(filepath.Join(gtfsDir, "*.zip"))
	if err != nil {
		log.Fatalf("Malformed file path: %s\n", err.Error())
	}

	col, err := gtfs.CreateGTFSCollection(zipFiles)
	if err != nil {
		log.Fatalf("Error creating GTFS schedule collection: %s\n", err)
		tt()
	}

	fmt.Println(gtfs.Overview(col))

	errFile, err := os.Create("gtfs_files/gtfs_errors.txt")
	if err != nil {
		log.Fatalf("Error creating error file: %s\n", err.Error())
	}
	defer errFile.Close()

	for _, e := range col {
		for _, err := range e.Errors() {
			errFile.WriteString(fmt.Sprintf("%s\n", err))
		}
	}
}
