package birdblog

import (
	"encoding/json"
	"io"
)

type TwitterResponse struct {
	Data []struct {
		Text string
	}
	Meta struct {
		NewestID    string
		OldestID    string
		ResultCount int
		NextToken   string
	}
}

// DecodeJSONIntoStruct returns TwitterResponse
func DecodeJSONIntoStruct(r io.Reader) (TwitterResponse, error) {
	var tr TwitterResponse
	err := json.NewDecoder(r).Decode(&tr)
	if err != nil {
		return TwitterResponse{}, err
	}
	return tr, nil
}
