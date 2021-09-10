package birdblog

import (
	"encoding/json"
	"io"
)

type TwitterResponse struct {
	Data []struct {
		Text     string
		AuthorID string `json:"author_id"`
	}
	Meta struct {
		NewestID    string
		OldestID    string
		ResultCount int
		NextToken   string
	}
}

type Tweet struct {
	Content  string
	AuthorID string
}


func (t Tweet) String() string {
	return t.Content
}

// DecodeJSONIntoStruct returns Conversation
func DecodeJSONIntoStruct(r io.Reader) (Conversation, error) {
	var tr TwitterResponse

	err := json.NewDecoder(r).Decode(&tr)
	if err != nil {
		return Conversation{}, err
	}
	// grab the data in reverse order and put it into conversation
	convo := make(Conversation, 0, len(tr.Data))
	for i := len(tr.Data) - 1; i >= 0; i-- {
		convo = append(convo, Tweet{
			Content:  tr.Data[i].Text,
			AuthorID: tr.Data[i].AuthorID,
		})
	}
	return convo, nil
}

