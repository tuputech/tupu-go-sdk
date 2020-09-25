module longasync

go 1.15

replace github.com/tuputech/tupu-go-sdk/recognition-api/general => ../recognition-api/general

replace github.com/tuputech/tupu-go-sdk/recognition-api/speech/longasync => ../recognition-api/speech/longasync

replace github.com/tuputech/tupu-go-sdk/recognition-api/speech/shortsync => ../recognition-api/speech/shortsync

require (
	github.com/tuputech/tupu-go-sdk/recognition-api/general v0.0.0-00010101000000-000000000000 // indirect
	github.com/tuputech/tupu-go-sdk/recognition-api/speech/longasync v0.0.0-00010101000000-000000000000
	github.com/tuputech/tupu-go-sdk/recognition-api/speech/shortsync v0.0.0-00010101000000-000000000000
)
