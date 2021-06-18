# Notes+

*Get the original tweet*

`curl -X GET -H "Authorization: Bearer $TWITTER_BEARER_TOKEN" "https://api.twitter.com/2/tweets?ids=1400096718788743180&tweet.fields=author_id,conversation_id,created_at,in_reply_to_user_id,referenced_tweets&expansions=author_id,in_reply_to_user_id,referenced_tweets.id&user.fields=name,username"`

*Use the tweet ID as the conversation ID and get the conversation*

  curl -X GET -H "Authorization: Bearer $TWITTER_BEARER_TOKEN" "https://api.twitter.com/2/tweets/search/recent?query=conversation_id:1390299295744659458,&max_results=100&tweet.fields=in_reply_to_user_id,author_id,created_at,conversation_id"

  curl -X GET -H "Authorization: Bearer $TWITTER_BEARER_TOKEN" "https://api.twitter.com/2/tweets/search/recent?query=conversation_id:1400096718788743180%20from:223680024&max_results=100&tweet.fields=in_reply_to_user_id,author_id,created_at,conversation_id"



*Filter conversation by author*
author id 223680024

TODO:
    NewConversationRequest(token, conversationID)
    FormatThread(tweet, conversation)
    PostToBlog(thread)