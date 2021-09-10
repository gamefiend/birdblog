package birdblog

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
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

type TwitterClient struct {
	token      string
	BaseURL    string
	HTTPClient *http.Client
}

type Tweet struct {
	Content  string
	AuthorID string
}

type Conversation []Tweet

func (c Conversation) String() string {
	var content strings.Builder
	for _, t := range c {
		content.WriteString(t.String())
	}
	return content.String()
}

func (c Conversation) FilterAuthor(ID string) Conversation {
	var filter Conversation
	for _, tweet := range c {
		if tweet.AuthorID == ID {
			filter = append(filter, tweet)
		}
	}
	return filter
}

func (c Conversation) FormatForGhost() string {
	var format strings.Builder
	for _, tweet := range c {
		format.WriteString("<p>")
		format.WriteString(strings.ReplaceAll(tweet.Content, "\n\n", "</p>\n<p>"))
		format.WriteString("</p>\n")
	}
	return format.String()
}

func NewTwitterClient(token string) TwitterClient {
	return TwitterClient{
		token:      token,
		BaseURL:    "https://api.twitter.com/2",
		HTTPClient: http.DefaultClient,
	}
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
func (tc TwitterClient) GetConversation(ID string) (Conversation, error) {
	req, err := tc.NewConversationRequest(ID)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := tc.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.token))
	return req, nil
}

func GhostPostRequest(content Conversation, token, site string) (*http.Request, error) {
	URL := fmt.Sprintf("https://%s/ghost/api/v4/admin/posts/?source=html", site)
	req, err := http.NewRequest(http.MethodPost, URL, nil)
	if err != nil {
		return nil, err
	}

	posts := html.EscapeString(content.String())
	req.Header.Set("Authorization", fmt.Sprintf("Ghost %s", token))
	req.Header.Set("Content-Type", "application/json")
	body := `
{
	"posts": [{
		"title": "Twitter Draft",
		"html": %q,
		"status": "draft"
	}]
}`
	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(body, posts)))
	return req, nil

}

func MakeGhostJWT(api string, start time.Time) (string, error) {
	info := strings.Split(api, ":")
	id := info[0]
	secret := info[1]
	expiration := 300 * time.Second

	key := make([]byte, hex.DecodedLen(len(secret)))
	_, err := hex.Decode(key, []byte(secret))
	if err != nil {
		return "", err
	}
	hs256 := jwt.NewHS256(key)
	// h := jwt.Header{KeyID: id}
	p := jwt.Payload{
		Audience:       jwt.Audience{"/v4/admin/"},
		ExpirationTime: jwt.NumericDate(start.Add(expiration)),
		IssuedAt:       jwt.NumericDate(start),
	}
	token, err := jwt.Sign(p, hs256, jwt.ContentType("JWT"), jwt.KeyID(id))
	if err != nil {
		return "", err
	}
	return string(token), nil

}

func RetrieveGhostURL(r io.Reader) (string, error) {
	type GhostResponse struct {
		Posts []struct {
			Url string
		}
	}
	var gr GhostResponse
	err := json.NewDecoder(r).Decode(&gr)
	if err != nil {
		return "", err
	}

	if len(gr.Posts) < 1 {
		return "", errors.New("Invalid Ghost Response")
	}

	return gr.Posts[0].Url, nil
}
