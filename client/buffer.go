package client

import (
	"bytes"
	"errors"
)

// Buffer describes a buffer of information that could be printed out
// to the user. It could later be assigned to a certain GUI tab.
type Buffer struct {
	buf bytes.Buffer
}

// Write writes a string to the buffer. Returns error if something goes wrong
func (b *Buffer) Write(str string) err {
	_, err := b.buf.WriteString(str)
	return err
}

// Flush returns the whole buffer in the format of a string.
// Returns an error if the buffer is empty.
func (b *Buffer) Flush() (string, error) {
	s := b.buf.String()
	if s == nil {
		return nil, errors.New("Buffer is empty.")
	}
	b.buffer.Reset()
	return s, nil
}
