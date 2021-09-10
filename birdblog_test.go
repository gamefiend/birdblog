package birdblog_test

import (
	"birdblog"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGivenJSONDataShouldDecodeIntoStruct(t *testing.T) {
	file, err := os.Open("testdata/conversation.json")
	if err != nil {
		t.Fatal(err)
	}
	var got birdblog.Conversation
	want := birdblog.Conversation{
		birdblog.Tweet{
			Content:  "This is super spoiler free, as IMO the way you design games based on other media is to extract patterns and not to translate content.\n\nAlso: I've got thoughts on what I'd use for system, but I'd REALLY rather talk about what a system should be doing to emulate the feel.",
			AuthorID: "223680024",
		},
		birdblog.Tweet{
			Content:  "...feel free to apply whatever we are talking about to your system of choice, but I'm not shopping or asking for system recommendations. I  probably won't interact much with \"this system does X\" even though I'm sure it does! \n\nThis is just not that conversation.",
			AuthorID: "223680024",
		},
	}

	got, err = birdblog.DecodeJSONIntoStruct(file)
	// fmt.Printf("%#v", got)
	if err != nil {
		t.Fatal(err)
	}

	if !(cmp.Equal(want, got)) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGivenIDGenerateTweetRequest(t *testing.T) {
	var ID string = "1400096718788743180"
	var token string = "dummyToken"
	wantURL, err := url.Parse("https://api.twitter.com/2/tweets/search/recent?query=conversation_id:1400096718788743180&max_results=100&tweet.fields=in_reply_to_user_id,author_id,created_at,conversation_id")
	if err != nil {
		t.Fatal(err)
	}
	tc := birdblog.NewTwitterClient(token)
	got, err := tc.NewConversationRequest(ID)
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



func TestGetConversation(t *testing.T) {
	token := "asodkfasndf"
	tweetID := "12353"
	tc := birdblog.NewTwitterClient(token)

	s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/conversation.json")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		io.Copy(w, file)
	}))

	tc.BaseURL = s.URL
	tc.HTTPClient = s.Client()
	// fake server to get conversation request from
	got, err := tc.GetConversation(tweetID)
	if err != nil {
		t.Fatal(err)
	}

	want := birdblog.Conversation{
		birdblog.Tweet{
			Content:  "This is super spoiler free, as IMO the way you design games based on other media is to extract patterns and not to translate content.\n\nAlso: I've got thoughts on what I'd use for system, but I'd REALLY rather talk about what a system should be doing to emulate the feel.",
			AuthorID: "223680024",
		},
		birdblog.Tweet{
			Content:  "...feel free to apply whatever we are talking about to your system of choice, but I'm not shopping or asking for system recommendations. I  probably won't interact much with \"this system does X\" even though I'm sure it does! \n\nThis is just not that conversation.",
			AuthorID: "223680024",
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFilterAuthorReturnsOnlyTweetsFromChosenAuthor(t *testing.T) {
	c := birdblog.Conversation{
		birdblog.Tweet{
			Content:  "This is super spoiler free, as IMO the way you design games based on other media is to extract patterns and not to translate content.\n\nAlso: I've got thoughts on what I'd use for system, but I'd REALLY rather talk about what a system should be doing to emulate the feel.",
			AuthorID: "223680024",
		},
		birdblog.Tweet{
			Content:  "...feel free to apply whatever we are talking about to your system of choice, but I'm not shopping or asking for system recommendations. I  probably won't interact much with \"this system does X\" even though I'm sure it does! \n\nThis is just not that conversation.",
			AuthorID: "223680023",
		},
	}

	want := birdblog.Conversation{
		birdblog.Tweet{
			Content:  "This is super spoiler free, as IMO the way you design games based on other media is to extract patterns and not to translate content.\n\nAlso: I've got thoughts on what I'd use for system, but I'd REALLY rather talk about what a system should be doing to emulate the feel.",
			AuthorID: "223680024",
		},
	}
	got := c.FilterAuthor()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFormatForGhostProperlyFormatsConversation(t *testing.T) {
	c := birdblog.Conversation{
		birdblog.Tweet{
			Content:  "This is super spoiler free, as IMO the way you design games based on other media is to extract patterns and not to translate content.\n\nAlso: I've got thoughts on what I'd use for system, but I'd REALLY rather talk about what a system should be doing to emulate the feel.",
			AuthorID: "223680024",
		},
		birdblog.Tweet{
			Content:  "...feel free to apply whatever we are talking about to your system of choice, but I'm not shopping or asking for system recommendations. I  probably won't interact much with \"this system does X\" even though I'm sure it does! \n\nThis is just not that conversation.",
			AuthorID: "223680024",
		},
	}

	want :=
		`<p>This is super spoiler free, as IMO the way you design games based on other media is to extract patterns and not to translate content.</p>
<p>Also: I've got thoughts on what I'd use for system, but I'd REALLY rather talk about what a system should be doing to emulate the feel.</p>
<p>...feel free to apply whatever we are talking about to your system of choice, but I'm not shopping or asking for system recommendations. I  probably won't interact much with "this system does X" even though I'm sure it does! </p>
<p>This is just not that conversation.</p>
`
	got := c.FormatForGhost()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRetrieveGhostUrlReturnsProperURL(t *testing.T) {
	wantURL := "https://thoughtcrime-games.ghost.io/p/0d725064-4953-43da-a2da-6bded281b27b/"
	file, err := os.Open("testdata/ghostresponse.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	gotURL, err := birdblog.RetrieveGhostURL(file)
	if err != nil {
		t.Fatal(err)
	}

	if wantURL != gotURL {
		t.Errorf("\nWant: %s\nGot:  %s\n", wantURL, gotURL)
	}
}
