# CSVMUM
CSV Marshal/Unmarshal

CSVMUM can convert a slice or map of structs into a slice of slices of strings, which can then be written with `csv.Write`

Example

```go
package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/bridgelightcloud/bogie/pkg/csvmum"
)

func main() {
	type Test struct {
		One   string
		Two   int
		Three bool

		four float64
	}

	tt := []Test{
		{One: "one", Two: 1},
		{One: "two", Two: 2, Three: true},
		{One: "three", four: 4.0},
	}

	out, err := csvmum.Marshal(tt)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("%d headers\n", len(out[0]))

	csv.NewWriter(os.Stdout).WriteAll(out)
}
```

Output
```
3 headers
One,Two,Three
one,1,false
two,2,true
three,0,false
```