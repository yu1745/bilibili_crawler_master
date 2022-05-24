package util

import (
	"bytes"
	"encoding/json"
)

func EncodeTask(t any) ([]byte, error) {
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	e.SetEscapeHTML(false)
	err := e.Encode(t)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
