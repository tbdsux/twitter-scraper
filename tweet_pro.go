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

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
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

	features := map[string]interface{}{
		"communities_web_enable_tweet_community_results_fetch":                    true,
		"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
		"tweetypie_unmention_optimization_enabled":                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                true,
		"tweet_awards_web_tipping_enabled":                                        false,
		"creator_subscriptions_quote_tweet_preview_enabled":                       false,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                true,
		"articles_preview_enabled":                                                true,
		"rweb_video_timestamps_enabled":                                           true,
		"rweb_tipjar_consumption_enabled":                                         true,
		"responsive_web_graphql_exclude_directive_enabled":                        true,
		"verified_phone_label_enabled":                                            false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_enhance_cards_enabled":                                    false,
		"responsive_web_grok_analyze_button_fetch_trends_enabled":                 false,
		"responsive_web_grok_analyze_post_followups_enabled":                      true,
		"responsive_web_grok_analysis_button_from_backend":                        true,
		"responsive_web_grok_image_annotation_enabled":                            true,
		"responsive_web_grok_share_attachment_enabled":                            true,
		"responsive_web_jetfuel_frame":                                            false,
		"profile_label_improvements_pcf_label_in_post_enabled":                    true,
		"premium_content_api_read_enabled":                                        false,
	}

	body := map[string]interface{}{
		"dark_request":             false,
		"disallowed_reply_options": nil,
		"features":                 features,
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
