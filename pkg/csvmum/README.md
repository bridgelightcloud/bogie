# CSVMUM
CSV Marshal/Unmarshal

CSVMUM is a CSV marshaler/unmarshaler. 

## Marshal

Example

```go
func marshal() {
	type person struct {
		Name string `csv:"name"`
		Age  int    `csv:"age"`
	}

	csvm, err := csvmum.NewMarshaler[person](os.Stdout)
	if err != nil {
		panic(err)
	}

	csvm.Marshal(person{Name: "Seanny Phoenix", Age: 38})
	csvm.Marshal(person{Name: "Somebody", Age: 27})
	csvm.Flush()
}
```

Output
```
name,age
Seanny Phoenix,38
Somebody,27
```

## Unmarshal

Flush must be called for the csv writer to write to the underlying writer.

Example

```go
func unmarshal() {
	type person struct {
		Name string `csv:"name"`
		Age  int    `csv:"age"`
	}

	r := bytes.NewBuffer([]byte("name,age\nNobody,0\nSpot,2\n"))
	csvu, err := csvmum.NewUnmarshaler[person](r)
	if err != nil {
		panic(err)
	}

	pp := []person{}
	for {
		var p person
		err = csvu.Unmarshal(&p)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		pp = append(pp, p)
	}

	fmt.Println(pp)
}
```

Output
```
[{Nobody 0} {Spot 2}]
```

## Tags

The `csv` tag can be used in structs to define the column name for a field. As with `json.Marshal` and `json.Unmarshal`, a field tagged with a hyphen (`-`) will be ignored.

Example

```go
func tags() {
	type tagged struct {
		AsIs       string  // marshaled as "AsIs"
		Renamed    float64 `csv:"renamed"` // marshaled as "renamed"
		unexported int     // not marshaled
		Ignored    bool    `csv:"-"` // not marshaled
	}

	taggedData := tagged{
		AsIs:       "as is",
		Renamed:    27.72,
		unexported: 2,
		Ignored:    true,
	}

	csvm, err := csvmum.NewMarshaler[tagged](os.Stdout)
	if err != nil {
		panic(err)
	}

	csvm.Marshal(taggedData)
	csvm.Flush()
}
```

Output

```
AsIs,renamed
as is,27.72
```