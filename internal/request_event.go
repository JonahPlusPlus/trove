package trove

import (
	"encoding/json"
)

type RequestEvent struct {
	Method string  `json:"method"`
	Host   string  `json:"host"`
	Path   string  `json:"path"`
	Remote string  `json:"remote"`
	Time   float64 `json:"time"`
	Error  error   `json:"error,omitempty"`

	encoded []byte
	err     error
}

func (re *RequestEvent) ensureEncoded() {
	if re.encoded == nil && re.err == nil {
		re.encoded, re.err = json.Marshal(re)
	}
}

func (re *RequestEvent) Length() int {
	re.ensureEncoded()
	return len(re.encoded)
}

func (re *RequestEvent) Encode() ([]byte, error) {
	re.ensureEncoded()
	return re.encoded, re.err
}
