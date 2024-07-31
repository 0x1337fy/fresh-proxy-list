package entity

import (
	"encoding/json"
	"testing"
)

const (
	expectedNoErrorMsg = "expected no error, got %v"
)

func TestUnmarshalJSONWithIsChecked(t *testing.T) {
	data := []byte(`{
		"method": "GET",
		"category": "api",
		"url": "http://example.com",
		"is_checked": false
	}`)

	var s Source
	err := json.Unmarshal(data, &s)
	if err != nil {
		t.Fatalf(expectedNoErrorMsg, err)
	}

	if s.Method != "GET" {
		t.Fatalf("expected method GET, got %s", s.Method)
	}
	if s.Category != "api" {
		t.Fatalf("expected category api, got %s", s.Category)
	}
	if s.URL != "http://example.com" {
		t.Fatalf("expected url http://example.com, got %s", s.URL)
	}
	if s.IsChecked != false {
		t.Fatalf("expected is_checked false, got %v", s.IsChecked)
	}
}

func TestUnmarshalJSONWithoutIsChecked(t *testing.T) {
	data := []byte(`{
		"method": "POST",
		"category": "web",
		"url": "http://example.org"
	}`)

	var s Source
	err := json.Unmarshal(data, &s)
	if err != nil {
		t.Fatalf(expectedNoErrorMsg, err)
	}

	if s.Method != "POST" {
		t.Fatalf("expected method POST, got %s", s.Method)
	}
	if s.Category != "web" {
		t.Fatalf("expected category web, got %s", s.Category)
	}
	if s.URL != "http://example.org" {
		t.Fatalf("expected url http://example.org, got %s", s.URL)
	}
	if s.IsChecked != true {
		t.Fatalf("expected is_checked true, got %v", s.IsChecked)
	}
}

func TestUnmarshalJSONWithInvalidData(t *testing.T) {
	data := []byte(`{
		"method": "GET",
		"category": "api",
		"url": "http://example.com",
		"is_checked": "string_instead_of_bool"
	}`)

	var s Source
	err := json.Unmarshal(data, &s)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUnmarshalJSONWithEmptyData(t *testing.T) {
	data := []byte(`{}`)

	var s Source
	err := json.Unmarshal(data, &s)
	if err != nil {
		t.Fatalf(expectedNoErrorMsg, err)
	}

	if s.Method != "" {
		t.Fatalf("expected method empty, got %s", s.Method)
	}
	if s.Category != "" {
		t.Fatalf("expected category empty, got %s", s.Category)
	}
	if s.URL != "" {
		t.Fatalf("expected url empty, got %s", s.URL)
	}
	if s.IsChecked != true {
		t.Fatalf("expected is_checked true, got %v", s.IsChecked)
	}
}
