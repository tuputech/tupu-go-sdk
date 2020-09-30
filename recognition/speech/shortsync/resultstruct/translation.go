package resultstruct

type Translation struct {
	TaskID   string
	FileList []FileTranslation `json:"fileList"`
}

type FileTranslation struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
