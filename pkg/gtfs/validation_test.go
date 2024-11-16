package gtfs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateColor(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value        string
		expectedErr error
	}{{
		value:        "000000",
		expectedErr: nil,
	}, {
		value:        "FFFFFF",
		expectedErr: nil,
	}, {
		value:        "123456",
		expectedErr: nil,
	}, {
		value:        "ABCDEF",
		expectedErr: nil,
	}, {
		value:        "abc123",
		expectedErr: nil,
	}, {
		value:        "abC14D",
		expectedErr: nil,
	}, {
		value:        "1234567",
		expectedErr: fmt.Errorf("invalid color: 1234567"),
	}, {
		value:        "ABCDEF1",
		expectedErr: fmt.Errorf("invalid color: ABCDEF1"),
	}, {
		value:        "12345",
		expectedErr: fmt.Errorf("invalid color: 12345"),
	}, {
		value:        "ABCDE",
		expectedErr: fmt.Errorf("invalid color: ABCDE"),
	}, {
		value:        "12345G",
		expectedErr: fmt.Errorf("invalid color: 12345G"),
	}, {
		value:        "ABCDEG",
		expectedErr: fmt.Errorf("invalid color: ABCDEG"),
	}, {
		value:        "",
		expectedErr: fmt.Errorf("invalid color: "),
	}, {
		value:        " 04FE2B",
		expectedErr: fmt.Errorf("invalid color:  04FE2B"),
	}, {
		value:        "#A5FF32",
		expectedErr: fmt.Errorf("invalid color: #A5FF32"),
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			_, err := parseColor(tc.value)

			assert.Equal(tc.expectedErr, err)
		})
	}
}
