package resultstruct

type TranslationReview struct {
	TaskID   string
	FileList []FileReview `json:"fileList"`
}

type FileReview struct {
	File_name string       `json:"file_name"`
	Result    ResultReview `json:"result"`
}

type ResultReview struct {
	Content  string         `json:"content"`
	Action   string         `json:"action"`
	Label    string         `json:"label"`
	Review   bool           `json:"review"`
	rate     string         `json:"rate"`
	HasVoice bool           `json:"hasVoice"`
	Details  []DetailReview `json:"details"`
}

type DetailReview struct {
	Keyword   string `json:"key"`
	Hint      string `json:"hint"`
	MainLabel string `json:"mainLabel"`
	SubLabel  string `json:"subLabel"`
}
