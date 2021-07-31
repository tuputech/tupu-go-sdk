// Package videoasync provide interface of TUPU video recognition
package videoasync

// video extends recognition.DataInfo to descripton video file
type (
	TaskCallbackRule struct {
		Review     bool     `json:"review,omitempty"`
		Offset     bool     `json:"offset,omitempty"`
		Total      bool     `json:"total,omitempty"`
		Label      uint8    `json:"label,omitempty"`
		FaceId     []string `json:"faceId,omitempty"`
		TypeName   []string `json:"typeName,omitempty"`
		Similarity float64  `json:"similarity,omitempty"`
	}

	VideoAsync struct {
		RealTimeCallback bool                          `json:"realTimeCallback,omitempty"`
		Audio            bool                          `json:"audio,omitempty"`
		Interval         uint8                         `json:"interval,omitempty"`
		Task             []string                      `json:"task,omitempty"`
		Video            string                        `json:"video"`
		CallbackUrl      string                        `json:"callbackUrl"`
		CallbackRules    map[string][]TaskCallbackRule `json:"callbackRules,omitempty"`
		CustomInfo       map[string]interface{}        `json:"customInfo,omitempty"`
	}

	AsyncOptFunc func(*VideoAsync)
)

const (
	DefalutInterval  = 1
	DefaultMaxFrames = 200
)

func newVidoASync(optFuncs ...AsyncOptFunc) *VideoAsync {
	var (
		video = new(VideoAsync)
	)
	for _, setConf := range optFuncs {
		setConf(video)
	}
	return video
}

// ClearBuffer is an helper to clear video binary content
func (vda *VideoAsync) ClearData() {
	vda.Interval = DefalutInterval
	vda.Task = nil
	vda.CallbackRules = nil
	vda.Video = ""
	vda.Audio = false
	vda.CallbackUrl = ""
	vda.RealTimeCallback = false
	vda.CustomInfo = nil
}

func (vdSync *VideoAsync) InitOptionParams(optFuncs ...AsyncOptFunc) {
	for _, opt := range optFuncs {
		opt(vdSync)
	}
}

func WithCallbackRules(callbackRules map[string][]TaskCallbackRule) AsyncOptFunc {
	return func(vs *VideoAsync) {
		vs.CallbackRules = callbackRules
	}
}

func WithAudio(audio bool) AsyncOptFunc {
	return func(vs *VideoAsync) {
		vs.Audio = audio
	}
}

func WithRealTimeCallback(realTimeCallback bool) AsyncOptFunc {
	return func(vs *VideoAsync) {
		vs.RealTimeCallback = realTimeCallback
	}
}

func WithCustomInfo(customInfo map[string]interface{}) AsyncOptFunc {
	return func(vds *VideoAsync) {
		vds.CustomInfo = customInfo
	}
}

func WithInterval(interval uint8) AsyncOptFunc {
	return func(vs *VideoAsync) {
		vs.Interval = interval
	}
}

func WithTask(tasks ...string) AsyncOptFunc {
	return func(vs *VideoAsync) {
		vs.Task = tasks
	}
}
