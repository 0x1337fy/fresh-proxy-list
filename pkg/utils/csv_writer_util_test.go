package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"testing"
)

const (
	errClosedPipeMessage = "io: read/write on closed pipe"
	flushErrorMessage    = "flush error"
	closeErrorMessage    = "close error"
)

type MockCloser struct {
	closerError error
}

func (m *MockCloser) Close() error {
	return m.closerError
}

func TestNewCSVWriter(t *testing.T) {
	tests := []struct {
		name     string
		wantType interface{}
	}{
		{
			name:     "NewCSVWriter returns *CSVWriter",
			wantType: &CSVWriterUtil{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCSVWriter()
			if _, ok := got.(*CSVWriterUtil); !ok {
				t.Errorf("NewCSVWriter() = %T, want %T", got, tt.wantType)
			}
		})
	}
}

func TestCSVWriterOpen(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() *CSVWriterUtil
		fields struct {
			writer *csv.Writer
			closer io.Closer
		}
		args struct {
			w io.Writer
			c io.Closer
		}
		wantErr error
	}{
		{
			name: "Open with writer and closer",
			setup: func() *CSVWriterUtil {
				return &CSVWriterUtil{}
			},
			args: struct {
				w io.Writer
				c io.Closer
			}{
				w: &bytes.Buffer{},
				c: &MockCloser{},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := tt.setup()
			err := uc.Open(tt.args.w, tt.args.c)
			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Open() error = %v, want %v", err, tt.wantErr)
				}
			} else if (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("Open() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCSVWriterFlush(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() *CSVWriterUtil
		fields struct {
			writer *csv.Writer
		}
		wantErr error
	}{
		{
			name: "Flush writer",
			setup: func() *CSVWriterUtil {
				buffer := &bytes.Buffer{}
				writer := csv.NewWriter(buffer)
				return &CSVWriterUtil{
					writer: writer,
				}
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := tt.setup()
			uc.Flush()
		})
	}
}

func TestCSVWriterWrite(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() *CSVWriterUtil
		fields struct {
			writer *csv.Writer
		}
		args struct {
			record []string
		}
		wantErr error
	}{
		{
			name: "Write record to writer",
			setup: func() *CSVWriterUtil {
				buffer := &bytes.Buffer{}
				writer := csv.NewWriter(buffer)
				return &CSVWriterUtil{
					writer: writer,
				}
			},
			args: struct {
				record []string
			}{
				record: []string{"field1", "field2"},
			},
			wantErr: nil,
		},
		{
			name: "Write to closed pipe",
			setup: func() *CSVWriterUtil {
				return &CSVWriterUtil{
					writer: nil,
				}
			},
			args: struct {
				record []string
			}{
				record: []string{"field1", "field2"},
			},
			wantErr: errors.New(errClosedPipeMessage),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := tt.setup()
			err := uc.Write(tt.args.record)
			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Write() error = %v, want %v", err, tt.wantErr)
				}
			} else if (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("Write() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCSVWriterClose(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() *CSVWriterUtil
		fields struct {
			closer io.Closer
		}
		wantErr error
	}{
		{
			name: "Close with no error",
			setup: func() *CSVWriterUtil {
				return &CSVWriterUtil{
					closer: &MockCloser{},
				}
			},
			wantErr: nil,
		},
		{
			name: "Close with error",
			setup: func() *CSVWriterUtil {
				return &CSVWriterUtil{
					closer: &MockCloser{
						closerError: errors.New(closeErrorMessage),
					},
				}
			},
			wantErr: errors.New(closeErrorMessage),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := tt.setup()
			err := uc.Close()
			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Close() error = %v, want %v", err, tt.wantErr)
				}
			} else if (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("Close() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
