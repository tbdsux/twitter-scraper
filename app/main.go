package main

import (
	"fmt"

	twitterscraper "github.com/tbdsux/twitter-scraper"
)



func main() {
	X_AUTH_TOKEN := "da58e595eeb19651787bbad6c5663f56a28f6a50"
	X_CSRF_TOKEN := "1e3ac5c3f95b49d2dbd8f580997a61f006de954e4688cafa57521162ecaa07c94176557f6c5c1cf76d27a806e019d53ebca25fe35e4f1f29610c70433ff2d02ed26bf6bc0f96bee02cb4d030965f5e7d"

	x := twitterscraper.New()
	x.SetAuthToken(twitterscraper.AuthToken{
		Token: X_AUTH_TOKEN,
		CSRFToken: X_CSRF_TOKEN,
	})

	if (!x.IsLoggedIn()) {
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
			ReplyToTweetId: "1888966650000724182",
		},
		Medias: []*twitterscraper.Media{
			media,
		},
	})


	fmt.Println(res, err)
}