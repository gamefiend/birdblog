package main

import (
	"birdblog"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	TweetIDFromArgs(os.Args)
	NewTwitterClientFromEnv()
	NewGhostClientFromEnv()
	NewTweetID()
	NewAuthorIDFromEnv()

}
func oldmain() {
	token := os.Getenv("TWITTER_BEARER_TOKEN")
	api := os.Getenv("GHOST_API_TOKEN")
	now := time.Now()
	ghostJWT, err := birdblog.MakeGhostJWT(api, now)
	if err != nil {
		log.Fatal(err)
	}
	ID := "1430623231535423492"
	tc := birdblog.NewTwitterClient(token)
	conversation, err := tc.GetConversation(ID)
	filter := conversation.FilterAuthor("223680024")
	ghostReq, err := birdblog.GhostPostRequest(filter, ghostJWT, "thoughtcrime-games.ghost.io")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v", filter.String())

	ghostResp, err := http.DefaultClient.Do(ghostReq)
	io.Copy(os.Stdout, ghostResp.Body)
	if err != nil {
		log.Fatal("Can't make the ghost post", err)
	}
	ghostURL, err := birdblog.RetrieveGhostURL(ghostResp.Body)
	if err != nil {
		log.Fatalf("Error retrieving Ghost URL: %s", err.Error())
	}
	defer ghostResp.Body.Close()

	fmt.Printf("Ghost URL: %s", ghostURL)
	if ghostResp.StatusCode != http.StatusCreated {
		r, _ := io.ReadAll(ghostResp.Body)
		log.Fatal("\nCan't make a Ghost draft: ", string(r))
	}

}
