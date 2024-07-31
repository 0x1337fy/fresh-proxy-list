package utils

import (
	"fmt"
	"io"
	"net/http"
)

type FetcherUtil struct {
	client     http.RoundTripper
	newRequest func(method string, url string, body io.Reader) (*http.Request, error)
}

type FetcherUtilInterface interface {
	FetchData(url string) ([]byte, error)
	Do(req *http.Request) (*http.Response, error)
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
	SetClient(client http.RoundTripper)
}

func NewFetcher(client http.RoundTripper, newRequestFunc func(method, url string, body io.Reader) (*http.Request, error)) FetcherUtilInterface {
	return &FetcherUtil{
		client:     client,
		newRequest: newRequestFunc,
	}
}

func (u *FetcherUtil) FetchData(url string) ([]byte, error) {
	req, err := u.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return body, fmt.Errorf("failed to fetch data: %s", http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (u *FetcherUtil) Do(req *http.Request) (*http.Response, error) {
	return u.client.RoundTrip(req)
}

func (u *FetcherUtil) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func (u *FetcherUtil) SetClient(client http.RoundTripper) {
	u.client = client
}
