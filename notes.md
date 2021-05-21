# Notes+

*Get the original tweet*

`curl -X GET -H "Authorization: Bearer $TWITTER_BEARER_TOKEN" "https://api.twitter.com/2/tweets?ids=1390300179098701826&tweet.fields=author_id,conversation_id,created_at,in_reply_to_user_id,referenced_tweets&expansions=author_id,in_reply_to_user_id,referenced_tweets.id&user.fields=name,username"`

*Use the tweet ID as the conversation ID and get the conversation*

  curl -X GET -H "Authorization: Bearer $TWITTER_BEARER_TOKEN" "https://api.twitter.com/2/tweets/search/recent?query=conversation_id:1395721291072692231&max_results=100&tweet.fields=in_reply_to_user_id,author_id,created_at,conversation_id"

*Filter conversation by author*

TODO:
    NewConversationRequest(token, conversationID)
    FormatThread(tweet, conversation)
    PostToBlog(thread)