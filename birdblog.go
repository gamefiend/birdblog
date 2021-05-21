package birdblog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func NewTweetRequest(token, ID string) (*http.Request, error) {
	URL := fmt.Sprintf("https://api.twitter.com/2/tweets?ids=%s&tweet.fields=author_id,conversation_id,created_at,in_reply_to_user_id,referenced_tweets&expansions=author_id,in_reply_to_user_id,referenced_tweets.id&user.fields=name,username", ID)
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return req, nil
}