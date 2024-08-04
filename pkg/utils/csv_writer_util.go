package utils

import (
	"encoding/csv"
	"io"
	"sync"
)

type CSVWriterUtil struct {
	writer *csv.Writer
	closer io.Closer
	mu     sync.RWMutex
}

type CSVWriterUtilInterface interface {
	Open(w io.Writer, c io.Closer) error
	Flush()
	Write(record []string) error
	Close() error
}

func NewCSVWriter() CSVWriterUtilInterface {
	return &CSVWriterUtil{
		mu: sync.RWMutex{},
	}
}

func (u *CSVWriterUtil) Open(w io.Writer, c io.Closer) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.writer = csv.NewWriter(w)
	u.closer = c
	return nil
}

func (u *CSVWriterUtil) Flush() {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.writer != nil {
		u.writer.Flush()
	}
}

func (u *CSVWriterUtil) Write(record []string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.writer == nil {
		return io.ErrClosedPipe
	}
	return u.writer.Write(record)
}

func (u *CSVWriterUtil) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closer != nil {
		if err := u.closer.Close(); err != nil {
			return err
		}
		u.closer = nil
		u.writer = nil
	}
	return nil
}
