package birdblog

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

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
