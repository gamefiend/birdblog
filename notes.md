# Notes+

*Get tweets from the page*
`curl -X GET -H "Authorization: Bearer $TWITTER_BEARER_TOKEN" "https://api.twitter.com/2/tweets/tweets?id=1390300179098701826&tweet.fields=author_id,conversation_id,created_at,in_reply_to_user_id,referenced_tweets&expansions=author_id,in_reply_to_user_id,referenced_tweets.id&user.fields=name,username"`
