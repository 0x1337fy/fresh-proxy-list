package repository

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"
)

type mockWriter struct {
	err error
}

func (m *mockWriter) Write(p []byte) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return len(p), nil
}

func (m *mockWriter) Close() error {
	return nil
}

func TestNewFileRepository(t *testing.T) {
	mockMkdirAll := func(path string, perm fs.FileMode) error {
		if path == "" {
			return errors.New("path cannot be empty")
		}
		return nil
	}
	mockCreate := func(name string) (io.Writer, error) {
		if name == "" {
			return nil, errors.New("file name cannot be empty")
		}
		return &bytes.Buffer{}, nil
	}
	repo := NewFileRepository(mockMkdirAll, mockCreate)

	if repo == nil {
		t.Errorf("Expected NewFileRepository to return a non-nil FileRepositoryInterface")
	}

	fileRepo, ok := repo.(*FileRepository)
	if !ok {
		t.Errorf("Expected repo to be of type *FileRepository")
	}

	if fileRepo.mkdirAll == nil {
		t.Fatal("expected mkdirAll to be set")
	}

	if fileRepo.create == nil {
		t.Fatal("expected create to be set")
	}

	if fileRepo.modePerm != fs.ModePerm {
		t.Errorf("Expected modePerm to be fs.ModePerm, got %v", fileRepo.modePerm)
	}
}

func TestSaveFile(t *testing.T) {
	type fields struct {
		mkdirAll func(path string, perm os.FileMode) error
		create   func(name string) (io.Writer, error)
	}

	type args struct {
		path   string
		data   interface{}
		format string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Create Directory Error",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return errors.New("error creating directory")
				},
				create: func(name string) (io.Writer, error) {
					return nil, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   strings.Join(testIPs, "\n"),
				format: "txt",
			},
			want:    "",
			wantErr: fmt.Errorf(testErrCreateDirectory, testFilePath, "error creating directory"),
		},
		{
			name: "Create File Error",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return nil, errors.New("error creating file")
				},
			},
			args: args{
				path:   testFilePath,
				data:   testIPs,
				format: "txt",
			},
			want:    "",
			wantErr: fmt.Errorf(testErrCreateFile, testFilePath, "error creating file"),
		},
		{
			name: "Unsupported Format",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   testIPs,
				format: "unsupported",
			},
			want:    "",
			wantErr: fmt.Errorf(testUnsupportedFormat, "unsupported"),
		},
		{
			name: "Write TXT Success",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   []string{testIP1},
				format: "txt",
			},
			want:    testIP1,
			wantErr: nil,
		},
		{
			name: "Write TXT Error",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &mockWriter{
						err: errors.New(testErrWriting),
					}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   []string{testIP1},
				format: "txt",
			},
			want:    "",
			wantErr: fmt.Errorf(testErrWritingTXT, testErrWriting),
		},
		{
			name: "Encode JSON",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   string(testProxiesToString),
				format: "json",
			},
			want:    string(testProxiesToString),
			wantErr: nil,
		},
		{
			name: "Encode JSON Error",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &mockWriter{
						err: errors.New(testErrWriting),
					}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   testProxies,
				format: "json",
			},
			want:    "",
			wantErr: fmt.Errorf(testErrEncode, "JSON", testErrWriting),
		},
		{
			name: "Encode CSV With String Data",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   testIPs,
				format: "csv",
			},
			want:    string(testIPsToText) + "\n",
			wantErr: nil,
		},
		{
			name: "Encode CSV With Proxy Data",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   testProxies,
				format: "csv",
			},
			want:    string(testProxiesToString) + "\n",
			wantErr: nil,
		},
		{
			name: "Encode CSV With Advanced Proxy Data",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   testAdvancedProxies,
				format: "csv",
			},
			want:    string(testAdvancedProxiesToText) + "\n",
			wantErr: nil,
		},
		{
			name: "Encode CSV With Error Data Type",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   []error{},
				format: "csv",
			},
			want:    "",
			wantErr: errors.New("invalid data type for CSV encoding"),
		},
		{
			name: "Encode XML",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   testProxies,
				format: "xml",
			},
			want:    string(testProxiesToString),
			wantErr: nil,
		},
		{
			name: "Encode XML Error",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &mockWriter{
						err: errors.New(testErrWriting),
					}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   testProxies,
				format: "xml",
			},
			want:    "",
			wantErr: fmt.Errorf(testErrEncode, "XML", testErrWriting),
		},
		{
			name: "Encode YAML",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   testProxies,
				format: "yaml",
			},
			want:    string(testProxiesToString),
			wantErr: nil,
		},
		{
			name: "Encode YAML Error",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &mockWriter{
						err: errors.New(testErrWriting),
					}, nil
				},
			},
			args: args{
				path:   testFilePath,
				data:   testProxies,
				format: "yaml",
			},
			want:    "",
			wantErr: fmt.Errorf(testErrEncode, "YAML", "yaml: write error: error writing"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &FileRepository{
				mkdirAll: tt.fields.mkdirAll,
				create:   tt.fields.create,
			}

			err := repo.SaveFile(tt.args.path, tt.args.data, tt.args.format)
			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("SaveFile() error = %v, want %v", err, tt.wantErr)
				}
			} else if (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("SaveFile() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if got, ok := tt.args.data.(string); ok {
					if got != tt.want {
						t.Errorf("SaveFile() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}
