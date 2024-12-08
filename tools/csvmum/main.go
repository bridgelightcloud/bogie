package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/bridgelightcloud/bogie/pkg/csvmum"
)

func main() {
	marshal()
	unmarshal()
}

func marshal() {
	csvm, err := csvmum.NewMarshaler[testD](os.Stdout)
	if err != nil {
		panic(err)
	}

	csvm.Marshal(testD{One: "uno", Two: "dos"})
	csvm.Marshal(testD{One: "1", Two: "2"})
	csvm.Flush()
}

func unmarshal() {
	r := bytes.NewBuffer([]byte("one,two\nuno,dos\n1,2\n"))
	csvu, err := csvmum.NewUnmarshaler[testD](r)
	if err != nil {
		panic(err)
	}

	tds := []testD{}
	for {
		var td testD
		err = csvu.Unmarshal(&td)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		tds = append(tds, td)
	}

	for _, td := range tds {
		fmt.Printf("td: %v\n", td)
	}
}

func e() {
	r := bytes.NewBuffer([]byte("one,two\nuno,dos\n1,2\n"))
	csvu, err := csvmum.NewUnmarshaler[testD](r)
	if err != nil {
		panic(err)
	}

	tds := []testD{}
	for {
		var td testD
		err = csvu.Unmarshal(&td)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		tds = append(tds, td)
	}

	for _, td := range tds {
		fmt.Printf("td: %v\n", td)
	}
}

type testD struct {
	One string `csv:"one"`
	Two string `csv:"two"`
}

func d() {
	csvm, err := csvmum.NewMarshaler[testD](os.Stdout)
	if err != nil {
		panic(err)
	}

	csvm.Marshal(testD{One: "uno", Two: "dos"})
	csvm.Marshal(testD{One: "1", Two: "2"})
	csvm.Flush()
}

// type thing struct {
// 	Date gtfs.Date `csv:"date"`
// 	Time gtfs.Time `csv:"time"`

// 	Heh string
// }

// func c() {
// 	t := []thing{
// 		{
// 			Date: gtfs.Date{Time: time.Time{}},
// 			Time: gtfs.Time{Time: time.Time{}},
// 			Heh:  "heh",
// 		},
// 		{Date: gtfs.Date{Time: time.Date(2024, 11, 26, 14, 14, 0, 0, time.UTC)}, Heh: "heh"},
// 		{Time: gtfs.Time{Time: time.Date(0, 0, 0, 14, 15, 23, 0, time.UTC)}},
// 	}

// 	out, err := csvmum.Marshal(t)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("out: %v\n", out)

// 	td := [][]string{
// 		{"date", "time", "Heh"},
// 		{"20241126", "14:14:00", "heh"},
// 		{"19860922", "14:15:23", ""},
// 	}

// 	var t2 []thing
// 	err = csvmum.Unmarshal(td, &t2)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("t2: %v\n", t2)

// }

// func b() {
// 	type a struct {
// 		One   string `json:"one" csv:"one"`
// 		Two   string `json:"" csv:""`
// 		Three string `json:"-" csv:"-"`
// 		four  string
// 	}

// 	_a := []a{{
// 		One:   "{one}",
// 		Two:   "{two}",
// 		Three: "{three}",
// 		four:  "{four}",
// 	}, {
// 		One:   "uno",
// 		Two:   "dos",
// 		Three: "tres",
// 		four:  "cuatro",
// 	}}

// 	j, err := json.MarshalIndent(_a, "", "  ")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Printf("%s\n", j)

// 	h, err := csvmum.Marshal(_a)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("h: %v\n", h)
// }

// func a() {
// 	type Test struct {
// 		One   string `csv:""`
// 		Two   int
// 		Three bool `csv:"three"`
// 		four  float64
// 	}

// 	testData := []Test{
// 		{One: "one", Two: 2, Three: true, four: 4.0},
// 		{One: "uno", Two: 20, Three: false, four: 8.0},
// 		{One: "un", Two: 200, Three: true, four: 16.0},
// 	}

// 	csv, err := csvmum.Marshal(testData)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("csv: %v\n", csv)

// 	var t []Test
// 	err = csvmum.Unmarshal(csv, &t)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("t: %v\n", t)
// }
