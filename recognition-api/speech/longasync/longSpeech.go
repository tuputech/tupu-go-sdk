package longasync

// LongSpeech is a structure that encapsulates long speech messages
type LongSpeech struct {
	// Url represents the address of the long voice, can't be empty
	URL string `json:"url"`
	// CallbackUrl represents the address of the callback result, cant' be empty
	CallbackURL string `json:"callbackUrl"`
	// CallbackRule represents the Rule of the callback, empty is using default rule, `all` is callback all result
	CallbackRule string `json:"callbackRule,omitempty"`
	// RoomID represents the room id
	RoomID string `json:"roomId,omitempty"`
	// UserID represents the user id
	UserID string `json:"userId,omitempty"`
	// ForumID represents the forum id
	ForumID string `json:"forumId,omitempty"`
}
