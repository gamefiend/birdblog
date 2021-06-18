package birdblog

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

type TwitterResponse struct {
	Data []struct {
		Text     string
		AuthorID string
	}
	Meta struct {
		NewestID    string
		OldestID    string
		ResultCount int
		NextToken   string
	}
}

type Conversation []string

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
		convo = append(convo, tr.Data[i].Text)
	}
	return convo, nil
}

func NewConversationRequest(token, ID string) (*http.Request, error) {
	URL := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=conversation_id:%s&max_results=100&tweet.fields=in_reply_to_user_id,author_id,created_at,conversation_id", ID)
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return req, nil
}

func GhostPostRequest(content Conversation, token, site string) (*http.Request, error) {
	URL := fmt.Sprintf("https://%s/ghost/api/v4/admin/posts/?source=html", site)
	req, err := http.NewRequest(http.MethodPost, URL, nil)
	if err != nil {
		return nil, err
	}

	posts := html.EscapeString(strings.Join(content, "\n"))
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
