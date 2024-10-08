package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bridgelightcloud/bogie/internal/db"
	"github.com/bridgelightcloud/bogie/internal/models"
	"github.com/google/uuid"
)

func gen(name string) {
	file, err := os.Open("/workspaces/bogie/tools/fill/" + name + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var units []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		units = append(units, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	evts := []models.Event{}
	for _, unit := range units {
		count := 1 + strings.Count(unit, "-")
		for i := 0; i < count; i++ {
			evt := models.Event{
				Id:     uuid.New(),
				Type:   db.DocTypeEvent,
				Agency: strings.ToUpper(name),
				UnitID: strings.TrimRightFunc(unit, func(r rune) bool {
					return r == '-'
				}),
			}
			evts = append(evts, evt)
		}
	}
	out, err := json.Marshal(evts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(out))
}

func main() {
	gen("muni")
	gen("bart")
}
