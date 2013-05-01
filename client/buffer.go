package client

import (
	"bytes"
	"errors"
)

// Buffer describes a buffer of information that could be printed out
// to the user. It could later be assigned to a certain GUI tab.
// Can store nicks if the buffer is for a channel
type Buffer struct {
	buf   bytes.Buffer
	nicks string
}

// Write writes a string to the buffer. Returns error if something goes wrong
func (b *Buffer) Write(str string) error {
	_, err := b.buf.WriteString(str)
	return err
}

// Flush returns the whole buffer in the format of a string.
// Returns an error if the buffer is empty.
func (b *Buffer) Flush() (string, error) {
	s := b.buf.String()
	if s == "" {
		return "", errors.New("Buffer is empty.")
	}
	b.buf.Reset()
	return s, nil
}
