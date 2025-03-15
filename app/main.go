package main

import (
	"fmt"

	twitterscraper "github.com/tbdsux/twitter-scraper"
)

func main() {
	X_AUTH_TOKEN := "auth_token"
	X_CSRF_TOKEN := "ct0"

	x := twitterscraper.New()
	x.SetAuthToken(twitterscraper.AuthToken{
		Token:     X_AUTH_TOKEN,
		CSRFToken: X_CSRF_TOKEN,
	})

	if !x.IsLoggedIn() {
		panic("not logged in")
	}

	// res, err := x.GetTweet("1888971932655304948")

	media, err := x.UploadMedia("./421426081932771394.png")
	if err != nil {
		fmt.Printf("error uploading media: %v", err)
		return
	}

	res, err := x.CreateTweet(twitterscraper.NewTweet{
		Text: "this is a reply!!!",
		Reply: &twitterscraper.TweetReply{
			ReplyToTweetId: "1894809694239367413",
		},
		Medias: []*twitterscraper.Media{
			media,
		},
	})

	fmt.Println(res, err)
}
