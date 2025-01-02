package csvmum

import (
	"bytes"
	"fmt"
)

type shortWriter struct{}

func (e shortWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

type closeReaderWriter struct {
	bytes.Buffer
	closed bool
}

func (c *closeReaderWriter) Close() error {
	c.closed = true
	return nil
}

func (c closeReaderWriter) Write(p []byte) (n int, err error) {
	if c.closed {
		return 0, fmt.Errorf("closed")
	}

	return c.Buffer.Write(p)
}

func (c closeReaderWriter) Read(p []byte) (n int, err error) {
	if c.closed {
		return 0, fmt.Errorf("closed")
	}

	return c.Buffer.Read(p)
}

type customMarshalAndUnmarshal struct {
	One string
}

func (c customMarshalAndUnmarshal) MarshalText() ([]byte, error) {
	if len(c.One) == 0 {
		return nil, fmt.Errorf("invalid text: %s", c.One)
	}
	return []byte(fmt.Sprintf("~%s~", c.One)), nil
}

func (c *customMarshalAndUnmarshal) UnmarshalText(text []byte) error {
	if len(text) < 2 {
		return fmt.Errorf("invalid text: %s", string(text))
	}
	c.One = string(text[1 : len(text)-1])
	return nil
}

func ptr[T any](v T) *T {
	return &v
}

func nilPtr[T any]() *T {
	return nil
}
