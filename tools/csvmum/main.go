package main

import (
	"encoding/json"
	"fmt"

	"github.com/bridgelightcloud/bogie/pkg/csvmum"
)

func main() {
	b()
}

func b() {
	type a struct {
		One   string `json:"one" csv:"one"`
		Two   string `json:"" csv:""`
		Three string `json:"-" csv:"-"`
		four  string
	}

	_a := []a{{
		One:   "{one}",
		Two:   "{two}",
		Three: "{three}",
		four:  "{four}",
	}, {
		One:   "uno",
		Two:   "dos",
		Three: "tres",
		four:  "cuatro",
	}}

	j, err := json.MarshalIndent(_a, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", j)

	h, err := csvmum.Marshal(_a)
	if err != nil {
		panic(err)
	}
	fmt.Printf("h: %v\n", h)
}

func a() {
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
