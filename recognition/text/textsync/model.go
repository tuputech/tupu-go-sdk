package textsync

type TextAsyncItem struct {
	Content   string `json:"content"`
	ContentID string `json:"contentId,omitempty"`
	UserID    string `json:"userId,omitempty"`
	ForumID   string `json:"forumId,omitempty"`
}
