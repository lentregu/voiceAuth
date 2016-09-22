package oxford

import (
	"net/http"

	"encoding/json"

	"github.com/TDAF/gologops"
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
