package birdblog_test

import (
	"birdblog"
	"io"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateGhostPostRequest(t *testing.T) {
	var content birdblog.Conversation
	content = append(content, birdblog.Tweet{
		Content:  "hello blog!",
		AuthorID: "12355233",
	})
	credentials := "someJWTToken"
	site := "myblog.ghost.io"
	gpr, err := birdblog.GhostPostRequest(content, credentials, site)
	wantBody := `
{
	"posts": [{
		"title": "Twitter Draft",
		"html": "hello blog!",
		"status": "draft"
	}]
}`
	gotBody, err := io.ReadAll(gpr.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantBody, string(gotBody)) {
		t.Errorf(cmp.Diff(wantBody, string(gotBody)))
	}

}

func TestCanCreateValidGhostJWT(t *testing.T) {
	now := time.Unix(162401290, 0)
	// dummy values to test that we get consistent results
	wantJWT := "eyJhbGciOiJIUzI1NiIsImN0eSI6IkpXVCIsImtpZCI6ImI0OWMxYzRmNjBlMDAzYjhhYTZjNmYwYyIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIvdjQvYWRtaW4vIiwiZXhwIjoxNjI0MDE1OTAsImlhdCI6MTYyNDAxMjkwfQ._ZJAG7d5_69HnGqvNWTbXPu_XsQtPQ8bwKVjCRua-oA"
	api := "b49c1c4f60e003b8aa6c6f0c:d396b0a338dc588295efe7d71e381cf6a82729d4f2783e47123803cbc0cca621"
	gotJWT, err := birdblog.MakeGhostJWT(api, now)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantJWT, gotJWT) {
		t.Errorf(cmp.Diff(wantJWT, gotJWT))
	}
}
