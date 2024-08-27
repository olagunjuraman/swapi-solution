package client

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	defaultHTTPTimeout = 10 * time.Second
	BaseURL            = "https://swapi.dev/api"
)

var client = &http.Client{Timeout: defaultHTTPTimeout}
var ErrNotFound = errors.New("not found")

func Call(path string, out interface{}) error {
	url := path
	if path[:4] != "http" {
		url = BaseURL + path
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "busha-go")

	res, err := client.Do(req)
	log.Println("Making SWAPI request")

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.StatusCode == 404 {
		return ErrNotFound
	}

	if err = json.NewDecoder(res.Body).Decode(out); err != nil {
		return err
	}

	return nil
}
