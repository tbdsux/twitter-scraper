package twitterscraper

import "net/url"

type ThreadCursor struct {
	FocalTweetID string
	ThreadID     string
	Cursor       string
	CursorType   string
}

func (s *Scraper) GetTweetReplies(id string, cursor string) ([]*Tweet, []*ThreadCursor, error) {
	req, err := s.newRequest("GET", "https://x.com/i/api/graphql/ldqoq5MmFHN1FhMGvzC9Jg/TweetDetail")
	if err != nil {
		return nil, nil, err
	}

	variables := map[string]interface{}{
		"focalTweetId":                           id,
		"referrer":                               "tweet",
		"with_rux_injections":                    false,
		"rankingMode":                            "Relevance",
		"includePromotedContent":                 true,
		"withCommunity":                          true,
		"withQuickPromoteEligibilityTweetFields": true,
		"withBirdwatchNotes":                     true,
		"withVoice":                              true,
	}

	if cursor != "" {
		variables["cursor"] = cursor
	}

	features := map[string]interface{}{
		"rweb_tipjar_consumption_enabled":                                         true,
		"responsive_web_graphql_exclude_directive_enabled":                        true,
		"verified_phone_label_enabled":                                            false,
		"creator_subscriptions_tweet_preview_api_enabled":                         true,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"communities_web_enable_tweet_community_results_fetch":                    true,
		"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
		"articles_preview_enabled":                                                true,
		"tweetypie_unmention_optimization_enabled":                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                true,
		"tweet_awards_web_tipping_enabled":                                        false,
		"creator_subscriptions_quote_tweet_preview_enabled":                       false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
		"rweb_video_timestamps_enabled":                                           true,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                true,
		"responsive_web_enhance_cards_enabled":                                    false,
	}

	fieldToggles := map[string]interface{}{
		"withArticleRichContentState": true,
		"withArticlePlainText":        false,
		"withGrokAnalyze":             false,
		"withDisallowedReplyControls": false,
	}

	query := url.Values{}
	query.Set("variables", mapToJSONString(variables))
	query.Set("features", mapToJSONString(features))
	query.Set("fieldToggles", mapToJSONString(fieldToggles))
	req.URL.RawQuery = query.Encode()

	var threads threadedConversation

	err = s.RequestAPI(req, &threads)
	if err != nil {
		return nil, nil, err
	}

	tweets, cursors := threads.parse(id)

	return tweets, cursors, nil
}
