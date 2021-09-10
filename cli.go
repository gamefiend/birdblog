package birdblog

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func RunCLI() {
	c, err := NewTwitterClientFromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ID, err := TweetIDFromArgs(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	conv, err := c.GetConversation(ID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	filter := conv.FilterAuthor()
	ghostToken := os.Getenv("GHOST_API_TOKEN")
	if ghostToken == "" {
		fmt.Fprintln(os.Stderr, "please set GHOST_API_TOKEN")
		os.Exit(1)
	}
	ghostJWT, err := MakeGhostJWT(ghostToken, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ghostReq, err := GhostPostRequest(filter, ghostJWT, "thoughtcrime-games.ghost.io")
	ghostResp, err := http.DefaultClient.Do(ghostReq)
	io.Copy(os.Stdout, ghostResp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't make the ghost post:", err)
		os.Exit(1)
	}
	ghostURL, err := RetrieveGhostURL(ghostResp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving Ghost URL: %v\n", err)
		os.Exit(1)
	}
	defer ghostResp.Body.Close()

	fmt.Printf("Ghost URL: %s", ghostURL)
	if ghostResp.StatusCode != http.StatusCreated {
		r, _ := io.ReadAll(ghostResp.Body)
		fmt.Fprintf(os.Stderr, "unexpected Ghost status code %q, response: %q\n", ghostResp.Status, string(r))
		os.Exit(1)
	}
}

func NewTwitterClientFromEnv() (*TwitterClient, error) {
	token := os.Getenv("TWITTER_BEARER_TOKEN")
	if token == "" {
		return nil, errors.New("please set TWITTER_BEARER_TOKEN")
	}
	return NewTwitterClient(token), nil
}

func TweetIDFromArgs(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("expected tweet ID")
	}
	return args[1], nil
}