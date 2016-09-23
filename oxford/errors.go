package oxford

import "errors"

type SpeakError error

var (
	errUnknownLocale  = SpeakError(errors.New("Unknown locale"))
	errInvalidSpeaker = SpeakError(errors.New("Invalid Speaker"))
	errInvalidPhrase  = SpeakError(errors.New("Invalid Phrase"))
)

func IsSpeakError(err error) bool {
	_, ok := err.(SpeakError)
	return ok
}
