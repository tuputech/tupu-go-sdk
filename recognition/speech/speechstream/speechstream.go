// Package speechstream provide interface of TUPU speech stream recognition
package speechstream

type (
	// SpeechStream  extends recognition.DataInfo to descripton speech file
	SpeechStream struct {
		// Url is your speech stream url
		URL string `json:"url,omitempty"`
		// Callback is your receive the recognition result url
		Callback string `json:"callback,omitempty"`
		// RoomId is customer params
		RoomID string `json:"roomId,omitempty"`
		// UserId is customer params
		UserID string `json:"userId,omitempty"`
		// ForumId is customer params
		ForumID string `json:"forumId,omitempty"`
		// CallbackRules mean callback recognition result, if this value equal 1
		CallbackRules uint8 `json:"callbackRules,omitempty"`
		// ReturnPreSpeech represents the link back 1min before the offending audio
		ReturnPreSpeech bool `json:"returnPreSpeech,omitempty"`
		tasks           []string
	}

	StreamOptFunc func(*SpeechStream)
)

const (
	SpeechVulgarTaskID      = "5c8213b9bc807806aab0a574"
	SpeechAnalysisTaskID    = "5caee6b2a76925c55a09a6d2"
	SpeechTranslationTaskID = "5ca1bd6b3872ecc9afb99132"
	SpeechGenderTaskID      = "5f59e4b71b29fa890e5472fb"
	CallbackNone            = 0
	CallbackAllRecognition  = 1
	CallbackEndStatus       = 2
)

func newSpeechStream(optFuncs ...StreamOptFunc) *SpeechStream {
	var (
		spstrm = new(SpeechStream)
	)
	for _, setConf := range optFuncs {
		setConf(spstrm)
	}
	return spstrm
}

// ClearBuffer is an helper to clear video binary content
func (spstrm *SpeechStream) ClearData() {
	spstrm.tasks = nil
	spstrm.URL = ""
	spstrm.Callback = ""
	spstrm.CallbackRules = CallbackNone
	spstrm.RoomID = ""
	spstrm.ForumID = ""
	spstrm.UserID = ""
	spstrm.ReturnPreSpeech = false
}

func (spstrm *SpeechStream) InitOptionParams(optFuncs ...StreamOptFunc) {
	for _, opt := range optFuncs {
		opt(spstrm)
	}
}

func WithURL(url string) StreamOptFunc {
	return func(spstrm *SpeechStream) {
		spstrm.URL = url
	}
}

func WithCallbackUrl(url string) StreamOptFunc {
	return func(spstrm *SpeechStream) {
		spstrm.Callback = url
	}
}

func WithRoomID(roomId string) StreamOptFunc {
	return func(spstrm *SpeechStream) {
		spstrm.RoomID = roomId
	}
}

func WithUserId(userId string) StreamOptFunc {
	return func(spstrm *SpeechStream) {
		spstrm.UserID = userId
	}
}

func WithFormID(forumId string) StreamOptFunc {
	return func(spstrm *SpeechStream) {
		spstrm.ForumID = forumId
	}
}

func WithCallbackRules(callbackRules uint8) StreamOptFunc {
	return func(spstrm *SpeechStream) {
		spstrm.CallbackRules = callbackRules
	}
}

func WithTask(tasks ...string) StreamOptFunc {
	return func(spstrm *SpeechStream) {
		spstrm.tasks = tasks
	}
}
