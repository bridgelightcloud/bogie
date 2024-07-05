package ddbmarshal

import (
	"testing"
	"time"
)

func TestMarshalAttributeValue(t *testing.T) {
	t.Parallel()

	// Test cases
	tests := []struct {
		name string
		data any
	}{
		{
			name: "Test case 1",
			data: nil,
		},
		{
			name: "Test case 2",
			data: "string",
		},
		{
			name: "Test case 3",
			data: 12345,
		},
		{
			name: "Test case 4",
			data: []int{1, 2, 3},
		},
		{
			name: "Test case 5",
			data: struct{ Name string }{Name: "John"},
		},
		{
			name: "Test case 6",
			data: UnixTimestamp(time.Now()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MarshalAttributeValue(tt.data)
		})
	}
}
