package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
)

func main() {
	count := os.Args[1]

	if count == "" {
		count = "1"
	}

	c, err := strconv.Atoi(count)

	if err != nil {
		fmt.Println("Usage: uuidgen [count]")
	}

	for i := 0; i < c; i++ {
		uuid := uuid.New()
		fmt.Println(uuid.String())
		// var arr [16]byte = uuid
		// fmt.Printf("%#v\n\n", arr)
	}
}
