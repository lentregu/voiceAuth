package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for i := 1; i < 3; i++ {
		part, err := writer.CreateFormFile(paramName+strconv.Itoa(i), filepath.Base(path))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func main() {
	fmt.Println("--------------------------")
	extraParams := map[string]string{
		"locale": "en-US",
	}
	//request, err := newfileUploadRequest("http://localhost:8080/profiles", extraParams, "audio", "/Users/gfr/Dropbox/Personal/RingsABell/audios/left.wav")
	request, err := newfileUploadRequest("http://localhost:8080/profiles", extraParams, "audio", "/Users/gfr/Develop/Go/src/github.com/lentregu/voiceAuth/samples/left1.wav")
	fmt.Println("--------------------------")
	//fmt.Print(request)
	fmt.Println("--------------------------")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
	}
}
