package repository

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"fresh-proxy-list/internal/entity"
	"fresh-proxy-list/pkg/utils"
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
	csvWriter utils.CSVWriterUtilInterface
}

type FileRepositoryInterface interface {
	SaveFile(filePath string, data interface{}, format string) error
}

type MkdirAllFunc func(path string, perm os.FileMode) error
type CreateFunc func(name string) (io.Writer, error)

func NewFileRepository(mkdirAll MkdirAllFunc, create CreateFunc, csvWriter utils.CSVWriterUtilInterface) FileRepositoryInterface {
	return &FileRepository{
		mkdirAll:  mkdirAll,
		create:    create,
		modePerm:  fs.ModePerm,
		csvWriter: csvWriter,
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
	if err := r.csvWriter.Open(w, nil); err != nil {
		return fmt.Errorf("failed to open CSV writer: %w", err)
	}
	defer r.csvWriter.Close()

	if header != nil {
		if err := r.csvWriter.Write(header); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	for _, row := range rows {
		if err := r.csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	r.csvWriter.Flush()
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
	var err error
	switch proxyData := data.(type) {
	case []string:
		view := entity.ProxyXMLClassicView{
			XMLName: xml.Name{Local: "proxies"},
			Proxies: make([]string, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = xml.NewEncoder(w).Encode(view)
	case []entity.Proxy:
		view := entity.ProxyXMLAdvancedView{
			XMLName: xml.Name{Local: "proxies"},
			Proxies: make([]entity.Proxy, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = xml.NewEncoder(w).Encode(view)
	case []entity.AdvancedProxy:
		view := entity.ProxyXMLAllAdvancedView{
			XMLName: xml.Name{Local: "Proxies"},
			Proxies: make([]entity.AdvancedProxy, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = xml.NewEncoder(w).Encode(view)
	}

	if err != nil {
		return fmt.Errorf("error encoding XML: %v", err)
	}
	return nil
}

func (r *FileRepository) encodeYAML(w io.Writer, data interface{}) error {
	var err error
	switch proxyData := data.(type) {
	case []string:
		view := struct {
			Proxies []string `yaml:"proxies"`
		}{
			Proxies: make([]string, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = yaml.NewEncoder(w).Encode(view)
	case []entity.Proxy:
		view := struct {
			Proxies []entity.Proxy `yaml:"proxies"`
		}{
			Proxies: make([]entity.Proxy, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = yaml.NewEncoder(w).Encode(view)
	case []entity.AdvancedProxy:
		view := struct {
			Proxies []entity.AdvancedProxy `yaml:"proxies"`
		}{
			Proxies: make([]entity.AdvancedProxy, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = yaml.NewEncoder(w).Encode(view)
	}

	if err != nil {
		return fmt.Errorf("error encoding YAML: %v", err)
	}
	return nil
}
