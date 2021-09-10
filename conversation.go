package birdblog

import "strings"

type Conversation []Tweet

func (c Conversation) String() string {
	var content strings.Builder
	for _, t := range c {
		content.WriteString(t.String())
	}
	return content.String()
}

func (c Conversation) FilterAuthor() Conversation {
	var filter Conversation
	if len(c) < 1 {
		return c
	}
	ID := c[0].AuthorID
	for _, tweet := range c {
		if tweet.AuthorID == ID {
			filter = append(filter, tweet)
		}
	}
	return filter
}

func (c Conversation) FormatForGhost() string {
	var format strings.Builder
	for _, tweet := range c {
		format.WriteString("<p>")
		format.WriteString(strings.ReplaceAll(tweet.Content, "\n\n", "</p>\n<p>"))
		format.WriteString("</p>\n")
	}
	return format.String()
}
