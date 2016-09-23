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
	"time"

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
	audio1 := parts["audio1"]
	audio2 := parts["audio2"]
	audio3 := parts["audio3"]

	if locale != "en-US" {
		http.NotFound(w, r)
		return
	}

	id, err := speak.CreateProfile(locale)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	id, err = enrollUser(id, audio1, audio2, audio3)

	if err == nil {
		response := createProfileResponse{}
		response.IdentificationProfileId = id
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
	userID := string(parts["userID"])

	id, err := speak.Verify(userID, passPhrase, audio)

	response := recogniseHandlerResponse{}
	if err == nil {
		response.Audio = byteArrayToBase64(audio)
		response.IdentificationProfileId = id
		json.NewEncoder(w).Encode(response)
		return
	}
	http.NotFound(w, r)
}

func passPhrases(w http.ResponseWriter, r *http.Request) {
	locale := "en-US"
	response, err := speak.PassList(locale)
	if err == nil {
		json.NewEncoder(w).Encode(response)
		return
	}
	http.NotFound(w, r)

}

//------------------------------------------------------------------------------

func enrollUser(userID string, audios ...[]byte) (id string, err error) {
	duration := time.Duration(5) * time.Second
	for _, audio := range audios {
		urlOP, _ := speak.UserEnrollment(userID, audio)
		fmt.Print(urlOP)
		time.Sleep(duration)
	}
	return userID, err
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
