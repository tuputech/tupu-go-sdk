module example

go 1.15

replace github.com/tuputech/tupu-go-sdk/base-recognition => ../base-recognition

replace github.com/tuputech/tupu-go-sdk/speech-sync => ../speech-sync

require (
	github.com/tuputech/tupu-go-sdk/base-recognition v0.0.0-00010101000000-000000000000
	github.com/tuputech/tupu-go-sdk/speech-sync v0.0.0-00010101000000-000000000000
)
