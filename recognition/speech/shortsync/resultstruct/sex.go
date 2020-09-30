package resultstruct

type Sex struct {
	TaskID  string
	Speechs []Speech `json:"speechs"`
}
