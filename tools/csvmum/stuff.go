package main

import (
	"fmt"
	"reflect"
)

func getTypeFromSlice(v any) {
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
