package twitterscraper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type newNoteTweet struct {
	Data struct {
		CreateTweet struct {
			TweetResults struct {
				Result tweetNote `json:"result"`
			} `json:"tweet_results"`
		} `json:"notetweet_create"`
	} `json:"data"`
}

type tweetNote struct {
	Core struct {
		UserResults struct {
			Result struct {
				IsBlueVerified bool       `json:"is_blue_verified"`
				Legacy         legacyUser `json:"legacy"`
			} `json:"result"`
		} `json:"user_results"`
	} `json:"core"`
	Views struct {
		Count string `json:"count"`
	} `json:"views"`
	NoteTweet struct {
		IsExpandable     bool `json:"is_expandable"`
		NoteTweetResults struct {
			Result struct {
				Id        string `json:"id"`
				Text      string `json:"text"`
				EntitySet any    `json:"entity_set"`
				Media     any    `json:"media"`
				RichText  any    `json:"richtext"`
			} `json:"result"`
		} `json:"note_tweet_results"`
	} `json:"note_tweet"`
	QuotedStatusResult struct {
		Result *result `json:"result"`
	} `json:"quoted_status_result"`
	Legacy legacyTweet `json:"legacy"`
}

func (newTweet *newNoteTweet) parse() *Tweet {
	var tweet = &newTweet.Data.CreateTweet.TweetResults.Result

	if tweet.NoteTweet.NoteTweetResults.Result.Text != "" {
		tweet.Legacy.FullText = tweet.NoteTweet.NoteTweetResults.Result.Text
	}
	var legacy *legacyTweet = &tweet.Legacy
	var user *legacyUser = &tweet.Core.UserResults.Result.Legacy

	fmt.Println(legacy)
	fmt.Println(legacy.IDStr)

	tw := parseLegacyTweet(user, legacy)
	if tw == nil {
		return nil
	}
	if tw.Views == 0 && tweet.Views.Count != "" {
		tw.Views, _ = strconv.Atoi(tweet.Views.Count)
	}
	if tweet.QuotedStatusResult.Result != nil {
		tw.QuotedStatus = tweet.QuotedStatusResult.Result.parse()
	}
	return tw
}

func (s *Scraper) CreateNoteTweet(tweet NewTweet) (*Tweet, error) {
	req, err := s.newRequest("POST", "https://x.com/i/api/graphql/AYb_zzIpA0IGC2rxLzXffQ/CreateNoteTweet")
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")

	media_entities := []map[string]interface{}{}

	if len(tweet.Medias) > 0 {
		for _, media := range tweet.Medias {
			media_entities = append(media_entities, map[string]interface{}{
				"media_id":     strconv.Itoa(media.ID),
				"tagged_users": []string{},
			})
		}
	}

	post_medias := map[string]interface{}{
		"media_entities":     media_entities,
		"possibly_sensitive": false,
	}

	variables := map[string]interface{}{
		"dark_request":            false,
		"media":                   post_medias,
		"semantic_annotation_ids": []string{},
		"tweet_text":              tweet.Text,
	}

	// If reply is set
	if tweet.Reply != nil {
		reply := map[string]interface{}{}

		if tweet.Reply.ReplyToTweetId != "" {
			reply["in_reply_to_tweet_id"] = tweet.Reply.ReplyToTweetId
			reply["exclude_reply_user_ids"] = []string{}
		}

		if len(tweet.Reply.ExcludeReplyUserIds) > 0 {
			reply["exclude_reply_user_ids"] = tweet.Reply.ExcludeReplyUserIds
		}

		if len(reply) > 0 {
			variables["reply"] = reply
		}
	}

	body := map[string]interface{}{
		"dark_request":             false,
		"disallowed_reply_options": nil,
		"features":                 tweetFeatures,
		"variables":                variables,
		"semantic_annotation_ids":  []string{},
		"queryId":                  "AYb_zzIpA0IGC2rxLzXffQ",
	}

	b, _ := json.Marshal(body)
	req.Body = io.NopCloser(bytes.NewReader(b))

	var response newNoteTweet
	err = s.RequestAPI(req, &response)
	if err != nil {
		return nil, err
	}

	if result := response.parse(); result != nil {
		return result, nil
	}

	return nil, errors.New("tweet wasn't post")
}
