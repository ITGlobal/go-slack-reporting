package slackreporting

import "net/url"

type urlIcon struct {
	URL string
}

func (i urlIcon) apply(args url.Values) {
	args["icon_url"] = []string{i.URL}
}

// NewIconFromURL creates new icon from an image URL
func NewIconFromURL(url string) Icon {
	return urlIcon{url}
}

type emojiIcon struct {
	Emoji string
}

func (i emojiIcon) apply(args url.Values) {
	args["icon_emoji"] = []string{i.Emoji}
}

// NewIconFromEmoji creates new icon from an emoji
func NewIconFromEmoji(emoji string) Icon {
	return emojiIcon{emoji}
}
