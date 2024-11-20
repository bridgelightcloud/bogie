package csvmum

import (
	"fmt"
	"reflect"
)

func GetTypeFromSlice(v any) {
	// expect v ro be a pointer to a slice
	p := reflect.ValueOf(v)
	fmt.Printf("p: %v\n\n", p)

	pt := p.Type()
	pte := pt.Elem()
	ptee := pte.Elem()
	fmt.Printf("pt: %v\n", pt)
	fmt.Printf("pte: %v\n", pte)
	fmt.Printf("ptee: %v\n\n", ptee)

	pe := p.Elem()
	pet := pe.Type()
	pete := pet.Elem()
	fmt.Printf("pe: %v\n", pe)
	fmt.Printf("pet: %v\n", pet)
	fmt.Printf("pete: %v\n", pete)
}
