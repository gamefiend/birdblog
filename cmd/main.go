package main

import (
	"birdblog"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	token := os.Getenv("TWITTER_BEARER_TOKEN")
	ID := "1390300179098701826"
	req, err := birdblog.NewTweetRequest(token, ID)
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
	fmt.Printf("%+v", tweet)
}