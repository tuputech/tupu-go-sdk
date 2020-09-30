package resultstruct

type TranslationEvent struct {
	TaskID  string
	Speechs []SpeechEvent `json:"speechs"`
}

type SpeechEvent struct {
	Name     string            `json:"name"`
	Metadata MetadataEvent     `json:"metadata"`
	Duration map[string]string `json:"duration"`
	Prob     map[string]string `json:"prob"`
}

type MetadataEvent struct {
	Threshold  float32           `json:"threshold"`
	Labels     map[string]string `json:"labels"`
	MsPerFrame float32           `json:"msPerFrame"`
	Smoothing  float32           `json:"smoothing"`
}
