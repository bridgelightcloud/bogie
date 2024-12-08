# CSVMUM
CSV Marshal/Unmarshal

CSVMUM can convert a slice or map of structs into a slice of slices of strings, which can then be written with `csv.Write`

## Marshal

Example

```go
func marshal() {
	csvm, err := csvmum.NewMarshaler[testD](os.Stdout)
	if err != nil {
		panic(err)
	}

	csvm.Marshal(testD{One: "uno", Two: "dos"})
	csvm.Marshal(testD{One: "1", Two: "2"})
	csvm.Flush()
}
```

Output
```
td: {uno dos}
td: {1 2}
```

## Unmarshal

Example

```go
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
```

Output
```
one,two
uno,dos
1,2
```