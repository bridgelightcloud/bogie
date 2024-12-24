package csvmum

import (
	"reflect"
	"testing"
)

func TestParseString(t *testing.T) {
	field := reflect.Value{}
	value := "string"
	v, err := parseString(field, value)
	if err != nil {
		t.Errorf("parseString() failed: %v", err)
	}
	if v.String() != value {
		t.Errorf("parseString() failed: expected %v, got %v", value, v.String())
	}
}

func TestParseInt(t *testing.T) {
	field := reflect.Value{}
	value := "1"
	v, err := parseInt(field, value)
	if err != nil {
		t.Errorf("parseInt() failed: %v", err)
	}
	if v.Int() != 1 {
		t.Errorf("parseInt() failed: expected %v, got %v", 1, v.Int())
	}
}

func TestParseFloat(t *testing.T) {
	field := reflect.Value{}
	value := "1.1"
	v, err := parseFloat64(field, value)
	if err != nil {
		t.Errorf("parseFloat() failed: %v", err)
	}
	if v.Float() != 1.1 {
		t.Errorf("parseFloat() failed: expected %v, got %v", 1.1, v.Float())
	}
}

func TestParseBool(t *testing.T) {
	field := reflect.Value{}
	value := "true"
	v, err := parseBool(field, value)
	if err != nil {
		t.Errorf("parseBool() failed: %v", err)
	}
	if v.Bool() != true {
		t.Errorf("parseBool() failed: expected %v, got %v", true, v.Bool())
	}
}

func TestParsePointer(t *testing.T) {
	field := reflect.ValueOf(ptr(1))
	value := "1"
	v, err := parsePointer(field, value)
	if err != nil {
		t.Errorf("parsePointer() failed: %v", err)
	}
	if v.Elem().Int() != 1 {
		t.Errorf("parsePointer() failed: expected %v, got %v", 1, v.Elem().Int())
	}
}
