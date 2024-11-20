package main

import (
	"fmt"

	"github.com/bridgelightcloud/bogie/pkg/csvmum"
)

func main() {
	type Test struct {
		One   string `csv:""`
		Two   int
		Three bool `csv:"three"`
		four  float64
	}

	testData := []Test{
		{One: "one", Two: 2, Three: true, four: 4.0},
		{One: "uno", Two: 20, Three: false, four: 8.0},
		{One: "un", Two: 200, Three: true, four: 16.0},
	}

	csv, err := csvmum.Marshal(testData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("csv: %v\n", csv)

	var t []Test
	err = csvmum.Unmarshal(csv, &t)
	if err != nil {
		panic(err)
	}

	fmt.Printf("t: %v\n", t)
}
