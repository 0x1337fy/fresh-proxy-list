package utils

import (
	"bytes"
	"encoding/csv"
	"io"
	"reflect"
	"testing"
)

const (
	writerBufferSize = 1024
)

type MockCloser struct {
	closerError error
}

func (m *MockCloser) Close() error {
	return m.closerError
}

func TestNewCSVWriter(t *testing.T) {
	tests := []struct {
		name string
		want interface{}
	}{
		{
			name: "NewCSVWriter returns *CSVWriter",
			want: &CSVWriterUtil{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCSVWriter()
			if _, ok := got.(*CSVWriterUtil); !ok {
				t.Errorf("NewCSVWriter() = %T, want %T", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			writer io.Writer
		}
		want *csv.Writer
	}{
		{
			name: "TestInit",
			args: struct {
				writer io.Writer
			}{
				writer: bytes.NewBuffer(make([]byte, writerBufferSize)),
			},
			want: csv.NewWriter(bytes.NewBuffer(make([]byte, writerBufferSize))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewCSVWriter()
			got := u.Init(tt.args.writer)
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Init() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}

func TestFlush(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *CSVWriterUtil
		args  struct {
			csvWriter *csv.Writer
		}
	}{
		{
			name: "TestFlush",
			setup: func() *CSVWriterUtil {
				return NewCSVWriter().(*CSVWriterUtil)
			},
			args: struct {
				csvWriter *csv.Writer
			}{
				csvWriter: csv.NewWriter(bytes.NewBuffer(make([]byte, writerBufferSize))),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.setup()
			u.Flush(tt.args.csvWriter)
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *CSVWriterUtil
		args  struct {
			writer io.Writer
			record []string
		}
		wantError error
	}{
		{
			name: "TestWrite",
			setup: func() *CSVWriterUtil {
				return NewCSVWriter().(*CSVWriterUtil)
			},
			args: struct {
				writer io.Writer
				record []string
			}{
				writer: &bytes.Buffer{},
				record: []string{"a", "b", "c"},
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.setup()
			csvWriter := u.Init(tt.args.writer)
			err := u.Write(csvWriter, tt.args.record)
			if err != nil && tt.wantError != nil {
				if err.Error() != tt.wantError.Error() {
					t.Errorf("Write() error = %v, want %v", err, tt.wantError)
				}
			} else if (err == nil && tt.wantError != nil) || (err != nil && tt.wantError == nil) {
				t.Errorf("Write() error = %v, want %v", err, tt.wantError)
			}
		})
	}
}
