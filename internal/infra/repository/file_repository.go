package repository

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"fresh-proxy-list/internal/entity"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileRepository struct {
	mkdirAll  func(path string, perm os.FileMode) error
	create    func(name string) (io.Writer, error)
	modePerm  fs.FileMode
	csvWriter *csv.Writer
}

type FileRepositoryInterface interface {
	SaveFile(filePath string, data interface{}, format string) error
}

type MkdirAllFunc func(path string, perm os.FileMode) error
type CreateFunc func(name string) (io.Writer, error)

func NewFileRepository(mkdirAll MkdirAllFunc, create CreateFunc) FileRepositoryInterface {
	return &FileRepository{
		mkdirAll:  mkdirAll,
		create:    create,
		modePerm:  fs.ModePerm,
		csvWriter: csv.NewWriter(&bytes.Buffer{}),
	}
}

func (r *FileRepository) SaveFile(filePath string, data interface{}, format string) error {
	if err := r.createDirectory(filePath); err != nil {
		return err
	}

	file, err := r.create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer func() {
		if f, ok := file.(io.Closer); ok {
			f.Close()
		}
	}()

	switch format {
	case "txt":
		return r.writeTxt(file, data)
	case "json":
		return r.encodeJSON(file, data)
	case "csv":
		return r.encodeCSV(file, data)
	case "xml":
		return r.encodeXML(file, data)
	case "yaml":
		return r.encodeYAML(file, data)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func (r *FileRepository) createDirectory(filePath string) error {
	err := r.mkdirAll(filepath.Dir(filePath), r.modePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %v", filePath, err)
	}
	return nil
}

func (r *FileRepository) writeTxt(file io.Writer, data interface{}) error {
	var dataString string
	if stringData, ok := data.([]string); ok {
		dataString = strings.Join(stringData, "\n")
	}

	_, err := file.Write([]byte(dataString))
	if err != nil {
		return fmt.Errorf("error writing TXT: %v", err)
	}
	return nil
}

func (r *FileRepository) encodeCSV(w io.Writer, data interface{}) error {
	switch proxyData := data.(type) {
	case []string:
		rows := make([][]string, len(proxyData))
		for i, rowElem := range proxyData {
			rows[i] = []string{rowElem}
		}
		return r.writeCSV(w, nil, rows)
	case []entity.Proxy:
		header := []string{"Proxy", "IP", "Port", "TimeTaken", "CheckedAt"}
		rows := make([][]string, len(proxyData))
		for i, proxy := range proxyData {
			rows[i] = []string{proxy.Proxy, proxy.IP, proxy.Port, fmt.Sprintf("%v", proxy.TimeTaken), proxy.CheckedAt}
		}
		return r.writeCSV(w, header, rows)
	case []entity.AdvancedProxy:
		header := []string{"Proxy", "IP", "Port", "Categories", "TimeTaken", "CheckedAt"}
		rows := make([][]string, len(proxyData))
		for i, proxy := range proxyData {
			rows[i] = []string{proxy.Proxy, proxy.IP, proxy.Port, strings.Join(proxy.Categories, ","), fmt.Sprintf("%v", proxy.TimeTaken), proxy.CheckedAt}
		}
		return r.writeCSV(w, header, rows)
	default:
		return fmt.Errorf("invalid data type for CSV encoding")
	}
}

func (r *FileRepository) writeCSV(w io.Writer, header []string, rows [][]string) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	if header != nil {
		if err := csvWriter.Write(header); err != nil {
			return fmt.Errorf("error writing CSV header: %v", err)
		}
	}

	for _, row := range rows {
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("error writing CSV row: %v", err)
		}
	}

	return nil
}

func (r *FileRepository) encodeJSON(w io.Writer, data interface{}) error {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil
}

func (r *FileRepository) encodeXML(w io.Writer, data interface{}) error {
	err := xml.NewEncoder(w).Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding XML: %v", err)
	}
	return nil
}

func (r *FileRepository) encodeYAML(w io.Writer, data interface{}) error {
	err := yaml.NewEncoder(w).Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding YAML: %v", err)
	}
	return nil
}
