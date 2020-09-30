# Package speech-sync

> Suitable for TUPU short speech recognition interface, providing access


## Example

> The processing entry is `SpeechHandler` struct, using its three methods(`Perform()`, `PerformWithURL()`, `PerformWithPath()`) to get the recognition results


   1. By methods `func (spHdler *SpeechHandler) Perform(secretID string, spSlice []*Speech) (string, int, error)` to get recognition result, [click to see](##)

   2. By methods `func (spHdler *SpeechHandler) PerformWithURL(secretID string, URLs []string) (string, int, error)` to get recognition result, [click to see](##)

   3. By methods `func (spHdler *SpeechHandler) PerformWithPath(secretID string, speechPaths []string) (string, int, error)` to get recognition result, [click to see](##)