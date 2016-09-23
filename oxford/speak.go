package oxford

import (
	"fmt"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/TDAF/gologops"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type SpeakAPI struct {
	apiKey string
}

// NewSpeak creates a speakAPI client
func NewSpeak(key string) SpeakAPI {
	s := SpeakAPI{}
	s.apiKey = key
	return s
}

func (s SpeakAPI) GetKey() string {
	return s.apiKey
}

type speakCreateProfileResponse struct {
	IdentificationProfileId string `json:"identificationProfileId,omitempty"`
}

type UserEnrollmentResponse struct {
	UrlOp string `json:"urlop,omitempty"`
}

type verifyReponse struct {
	Result     string `json:"result,omitempty"`     // [Accept | Reject]
	Confidence string `json:"confidence,omitempty"` // [Low | Normal | High]
	Phrase     string `json:"urlop,omitempty"`      //"recognized phrase"
}

type phraseResponse struct {
	Phrase string
}

type PhrasesListResponse []phraseResponse

func (s SpeakAPI) CreateProfile(locale string) (profileID string, err error) {
	url := GetResource(SpeakerRecognition, V1, "identificationProfiles")
	resp, err := POST(url, nil, s.apiKey, nil, "application/json", M{"locale": locale})

	if err != nil {
		return "", err
	}

	var successResponse speakCreateProfileResponse
	switch resp.StatusCode {
	case http.StatusOK:
		json.NewDecoder(resp.Body).Decode(&successResponse)
		gologops.InfoC(gologops.C{"op": "speak:CreateProfile", "result": "OK"}, "%s", resp.Status)
		profileID = successResponse.IdentificationProfileId
	default:
		var errorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		err = errorResponse.Err
		gologops.ErrorE(err, gologops.C{"op": "speak:CreateProfile", "result": "NOK"}, "%s", resp.Status)
	}

	return profileID, err
}

func (s SpeakAPI) UserEnrollment(id string, audio []byte) (urlOP string, err error) {
	url := GetResource(SpeakerRecognition, V1, "identificationProfiles")
	url = url + "/" + id + "/enroll"
	resp, err := POST(url, M{"shortAudio": "true"}, s.apiKey, M{"Content-Length": strconv.Itoa(len(audio))},
		"application/octet-stream", audio)
	if err != nil {
		return "", err
	}

	var successResponse UserEnrollmentResponse
	switch resp.StatusCode {
	case http.StatusAccepted:
		json.NewDecoder(resp.Body).Decode(&successResponse)
		gologops.InfoC(gologops.C{"op": "speak:EnrollUser", "result": "OK"}, "%s", resp.Status)
		urlOP = resp.Header.Get("Operation-Location")
	default:
		var errorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		err = errorResponse.Err
		gologops.ErrorE(err, gologops.C{"op": "speak:EnrollUser", "result": "NOK"}, "%s", resp.Status)
	}

	return urlOP, err
}

func (s SpeakAPI) Verify(userID string, passPhrase string, audio []byte) (id string, err error) {
	url := GetResource(SpeakerRecognition, V1, "verify")
	url = url + "/" + id + "/enroll"
	resp, err := POST(url, M{"verificationProfileId": id}, s.apiKey, M{"Content-Length": strconv.Itoa(len(audio))},
		"application/octet-stream", audio)
	if err != nil {
		return "", err
	}

	var successResponse verifyReponse
	var distance int
	switch resp.StatusCode {
	case http.StatusAccepted:
		json.NewDecoder(resp.Body).Decode(&successResponse)
		gologops.InfoC(gologops.C{"op": "speak:Verify", "result": "OK"}, "%s", resp.Status)
		phrase := successResponse.Phrase
		distance = levenshtein.DistanceForStrings([]rune(passPhrase), []rune(phrase), levenshtein.DefaultOptions)
	default:
		var errorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		err = errorResponse.Err
		gologops.ErrorE(err, gologops.C{"op": "speak:Verify", "result": "NOK"}, "%s", resp.Status)
	}

	if distance <= 5 {
		err = errInvalidPhrase
	}

	return userID, err
}
func (s SpeakAPI) PassList(locale string) (passPhrases PhrasesListResponse, err error) {

	url := GetResource(SpeakerRecognition, V1, "verificationPhrases")
	resp, err := GET(url, s.apiKey, M{"locale": locale}, nil)

	if err != nil {
		return nil, err
	}

	var successResponse PhrasesListResponse
	switch resp.StatusCode {
	case http.StatusOK:
		json.NewDecoder(resp.Body).Decode(&successResponse)
		passPhrases = successResponse
		fmt.Print(successResponse)
		gologops.InfoC(gologops.C{"op": "speak:passList", "result": "OK"}, "%s", resp.Status)
	default:
		var errorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		err = errorResponse.Err
		gologops.ErrorE(err, gologops.C{"op": "speak:passList", "result": "NOK"}, "%s", resp.Status)
	}

	return passPhrases, err
}
