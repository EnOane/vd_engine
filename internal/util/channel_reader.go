package util

import (
	"bytes"
	"io"
)

// ChannelReader - ридер для чтения данных из канала
type ChannelReader struct {
	ch     <-chan []byte
	buffer *bytes.Buffer
	done   bool
}

// NewChannelReader создает новый ChannelReader
func NewChannelReader(ch <-chan []byte) *ChannelReader {
	return &ChannelReader{
		ch:     ch,
		buffer: bytes.NewBuffer(nil),
	}
}

func (r *ChannelReader) Read(p []byte) (n int, err error) {
	if r.buffer.Len() == 0 {
		if r.done {
			return 0, io.EOF
		}

		data, ok := <-r.ch
		if !ok {
			r.done = true
			return 0, io.EOF
		}

		_, err := r.buffer.Write(data)
		if err != nil {
			return 0, err
		}
	}

	return r.buffer.Read(p)
}
