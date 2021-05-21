package birdblog_test

import (
	"birdblog"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGivenJSONDataShouldDecodeIntoStruct(t *testing.T) {
	file, err := os.Open("testdata/response.json")
	if err != nil {
		t.Fatal(err)
	}
	response, err := birdblog.DecodeJSONIntoStruct(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Data) != 10 {
		t.Fatalf("response is the wrong length. Wanted 10, got %d", len(response.Data))
	}
	want := "@qh_murphy 100%. The Attack of Opportunity (and ESPECIALLY the Full Attack) led to just standing still in melee and wailing on each other."
	got := response.Data[0].Text

	if !(cmp.Equal(want, got)) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGivenIDGenerateTweetRequest(t *testing.T) {
	var ID string = "1390300179098701826"
	var token string = "dummyToken"
	wantURL, err := url.Parse("https://api.twitter.com/2/tweets?ids=1390300179098701826&tweet.fields=author_id,conversation_id,created_at,in_reply_to_user_id,referenced_tweets&expansions=author_id,in_reply_to_user_id,referenced_tweets.id&user.fields=name,username")
	if err != nil {
		t.Fatal(err)
	}
	got, err := birdblog.NewTweetRequest(token, ID)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(http.MethodGet, got.Method) {
		t.Error(cmp.Diff(http.MethodGet, got.Method))
	}
	if !cmp.Equal(wantURL, got.URL) {
		t.Error(cmp.Diff(wantURL, got.URL))
	}
	wantHeader := "Bearer dummyToken"
	gotHeader := got.Header.Get("Authorization")
	if !cmp.Equal(wantHeader, gotHeader) {
		t.Errorf("bad Authorization header: %s", cmp.Diff(wantHeader, gotHeader))
	}
}