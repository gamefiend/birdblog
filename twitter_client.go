package birdblog

import (
	"fmt"
	"net/http"
)

type TwitterClient struct {
	TwitterBearerToken      string
	BaseURL    string
	HTTPClient *http.Client
}

func NewTwitterClient(token string) *TwitterClient {
	return &TwitterClient{
		TwitterBearerToken:      token,
		BaseURL:    "https://api.twitter.com/2",
		HTTPClient: http.DefaultClient,
	}
}

func (tc *TwitterClient) GetConversation(ID string) (Conversation, error) {
	req, err := tc.NewConversationRequest(ID)
	if err != nil {
		return nil, err
	}

	resp, err := tc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error retrieving tweets: %v", resp.Status)
	}
	conversation, err := DecodeJSONIntoStruct(resp.Body)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func (tc TwitterClient) NewConversationRequest(ID string) (*http.Request, error) {
	URL := fmt.Sprintf("%s/tweets/search/recent?query=conversation_id:%s&max_results=100&tweet.fields=in_reply_to_user_id,author_id,created_at,conversation_id", tc.BaseURL, ID)
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.TwitterBearerToken))
	return req, nil
}
