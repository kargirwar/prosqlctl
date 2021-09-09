package utils

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const CURRENT_RELEASE = "https://raw.githubusercontent.com/kargirwar/prosql-agent/master/current-release.json"
const STATUS_URL = "http://localhost:23890/about"

type Release struct {
	Version string
	Mac     string
	Linux   string
	Windows string
}

type Response struct {
	Status    string      `json:"status"`
	Msg       string      `json:"msg"`
	ErrorCode string      `json:"error-code"`
	Data      interface{} `json:"data"`
	Eof       bool        `json:"eof"`
}

func DownloadFile(fileName string, url string) (err error) {
	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	return nil
}

func GetStatus() *Response {
	response := &Response{}
	err := GetJson(STATUS_URL, response)

	if err != nil {
		return &Response{
			Status: "error",
			Msg:    "prosql-agent not installed",
		}
	}

	return response
}

//https://gist.github.com/paulerickson/6d8650947ee4e3f3dbcc28fde10eaae7
func Unzip(source, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()
	for _, file := range archive.Reader.File {
		reader, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()
		path := filepath.Join(destination, file.Name)
		// Remove file if it already exists; no problem if it doesn't; other cases can error out below
		_ = os.Remove(path)
		// Create a directory at path, including parents
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		// If file is _supposed_ to be a directory, we're done
		if file.FileInfo().IsDir() {
			continue
		}
		// otherwise, remove that directory (_not_ including parents)
		err = os.Remove(path)
		if err != nil {
			log.Fatal(err)
		}
		// and create the actual file.  This ensures that the parent directories exist!
		// An archive may have a single file with a nested path, rather than a file for each parent dir
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			log.Fatal(err)
		}
		defer writer.Close()
		_, err = io.Copy(writer, reader)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func GetCwd() string {
	//get current working dir
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func GetLatestRelease() *Release {
	release := &Release{}
	err := GetJson(CURRENT_RELEASE, release)
	if err != nil {
		log.Fatal(err)
	}

	return release
}

func GetJson(url string, target interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(r)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(target)
}

func GetValue(r *Response, k string) string {
	m, ok := r.Data.(map[string]interface{})
	if !ok {
		log.Fatal("Unable to parse JSON")
	}
	return m[k].(string)
}
