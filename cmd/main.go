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
	token := os.Getenv("TWITTER_BEARER_TOKEN")
	api := os.Getenv("GHOST_API_TOKEN")
	now := time.Now()
	ghostJWT, err := birdblog.MakeGhostJWT(api, now)
	if err != nil {
		log.Fatal(err)
	}
	ID := "1405513134245388292"
	req, err := birdblog.NewConversationRequest(token, ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req.URL)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("no tweets for you: %v", resp.Status)
	}
	tweet, err := birdblog.DecodeJSONIntoStruct(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	ghostReq, err := birdblog.GhostPostRequest(tweet, ghostJWT, "thoughtcrime-games.ghost.io")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v", tweet)
	b, _ := io.ReadAll(ghostReq.Body)
	fmt.Printf("%s\n", string(b))
	ghostResp, err := http.DefaultClient.Do(ghostReq)
	if err != nil {
		log.Fatal(err)
	}
	if ghostResp.StatusCode != http.StatusOK {
		r, _ := io.ReadAll(ghostResp.Body)
		log.Fatal("\nCan't make a Ghost draft: ", string(r))
	}

}
