package resultstruct

type Voiceprint struct {
	TaskID  string
	Speechs []SpeechVoiceprint `json:"speechs"`
}

type SpeechVoiceprint struct {
	Name    string             `json:"name"`
	Details []DetailVoiceprint `json:"details"`
}

type DetailVoiceprint struct {
	PersonName string  `json:"personName"`
	Similarity string  `json:"similarity"`
	TypeName   string  `json:"typeName"`
	StartTime  float32 `json:"startTime"`
	EndTime    float32 `json:"endTime"`
	Review     bool    `json:"review"`
}
