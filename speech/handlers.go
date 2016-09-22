package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/lentregu/voiceauth/oxford"
)

const (
	// Key for voice recognition API
	voiceKey = "af90809f8a0d4430ba2aabd44785ebc4"
)

type createProfileRequest struct {
	Locale  string `json:"locale"`
	Sample1 []byte `json:"sample1"`
	Sample2 []byte `json:"sample2"`
	Sample3 []byte `json:"sample3"`
}

type createProfileResponse struct {
	IdentificationProfileId string `json:"identificationProfileId"`
}

type recogniseHandlerResponse struct {
	IdentificationProfileId string `json:"identificationProfileId"`
	Audio                   string `json:audio`
}

var speak oxford.SpeakAPI

func init() {
	speak = oxford.NewSpeak("af90809f8a0d4430ba2aabd44785ebc4")
}

// index is the welcome handler
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
	fmt.Printf("key: %s\n", speak.GetKey())
}

func createProfileHandler(w http.ResponseWriter, r *http.Request) {
	parts := getParts(r)
	locale := string(parts["locale"])
	// audio1 := parts["audio1"]
	// audio2 := parts["audio2"]
	// audio3 := parts["audio3"]

	if locale == "en-US" {
		response := createProfileResponse{}
		response.IdentificationProfileId = "a34e82f4-5530-4fb9-8b7c-ebf86697865b"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	http.NotFound(w, r)
}

func recognizeHandler(w http.ResponseWriter, r *http.Request) {
	// frase + audio
	parts := getParts(r)
	passPhrase := string(parts["passPhrase"])
	audio := parts["audio"]
	response := recogniseHandlerResponse{}
	if passPhrase == "hola" {
		response.Audio = byteArrayToBase64(audio)
		response.IdentificationProfileId = "a34e82f4-5530-4fb9-8b7c-ebf86697865b"
		json.NewEncoder(w).Encode(response)
	}
	http.NotFound(w, r)
}

func byteArrayToBase64(binaryByteArray []byte) string {
	base64 := base64.StdEncoding.EncodeToString(binaryByteArray)
	return base64
}

func getParts(r *http.Request) map[string][]byte {
	//4= 1 text part is locale + 3 audio wav
	parts := make(map[string][]byte)
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, params["boundary"])
		i := 0
		for {
			i++
			p, err := mr.NextPart()
			if err == io.EOF {
				return parts
			}
			fmt.Printf("FormName: %s", p.FormName())
			if err != nil {
				log.Fatal(err)
			}
			part, err := ioutil.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}
			parts[p.FormName()] = part
		}
	}
	return parts
}
