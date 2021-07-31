package speechasync

// SpeechAsync is a structure that encapsulates long speech messages
type (
	SpeechAsync struct {
		// FileRemoteURL represents the address of the long voice, can't be empty
		FileRemoteURL string `json:"url"`
		// CallbackUrl represents the address of the callback result, cant' be empty
		CallbackURL string `json:"callbackUrl"`
		// CallbackRule represents the Rule of the callback, empty is using default rule, `all` is callback all result
		CallbackRule string `json:"callbackRule,omitempty"`
		// RoomID represents the room id
		RoomID string `json:"roomId,omitempty"`
		// UserID represents the user id
		UserID string `json:"userId,omitempty"`
		// ForumID represents the forum id
		ForumID    string            `json:"forumId,omitempty"`
		CustomInfo map[string]string `json:"customInfo,omitempty"`
	}
	SPAsyncOptFunc func(*SpeechAsync)
)

func (sa *SpeechAsync) ClearData() {
	sa.FileRemoteURL = ""
	sa.CallbackURL = ""
	sa.CallbackRule = ""
	sa.RoomID = ""
	sa.UserID = ""
	sa.ForumID = ""
	sa.CustomInfo = nil
}

const (
	CallbackRuleALL = "all"
)

func newSpeechASync(optFuncs ...SPAsyncOptFunc) *SpeechAsync {
	speech := new(SpeechAsync)
	for _, setConf := range optFuncs {
		setConf(speech)
	}
	return speech
}

func WithCallbackRule(callbackRule string) SPAsyncOptFunc {
	return func(sa *SpeechAsync) {
		sa.CallbackRule = callbackRule
	}
}

func WithRoomID(roomId string) SPAsyncOptFunc {
	return func(sa *SpeechAsync) {
		sa.RoomID = roomId
	}
}

func WithUserId(userId string) SPAsyncOptFunc {
	return func(sa *SpeechAsync) {
		sa.UserID = userId
	}
}

func WithFormID(forumId string) SPAsyncOptFunc {
	return func(sa *SpeechAsync) {
		sa.ForumID = forumId
	}
}

func WithCustomInfo(customInfo map[string]string) SPAsyncOptFunc {
	return func(sa *SpeechAsync) {
		sa.CustomInfo = customInfo
	}
}
