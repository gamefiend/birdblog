package birdblog_test

import (
	"birdblog"
	"testing"
)

func TestNewTwitterClientFromEnv(t *testing.T) {
	want := "dummy_twitter_token"
	t.Setenv("TWITTER_BEARER_TOKEN", want)
	c, err := birdblog.NewTwitterClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	got := c.TwitterBearerToken
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}

}

func TestNewTwitterClientFromEnvNoTokenError(t *testing.T) {
	_, err := birdblog.NewTwitterClientFromEnv()
	if err == nil {
		t.Fatal("want error with missing Twitter token, got nil")
	}
}

func TestTweetIDFromArgs(t *testing.T) {
	want := "1430623231535423492"
	got, err := birdblog.TweetIDFromArgs([]string{"/usr/bin/birdblog", want})
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}