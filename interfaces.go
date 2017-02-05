package slackreporting

import (
	"log"
	"net/url"
)

// Config defines slack reporter options
type Config interface {
	// Set Slack API URL (e.g. "https://slack.com/")
	SetURL(url string) Config

	// Set Slack channel or username (e.g. "#channel" or "@SetUsername")
	SetChannel(channel string) Config

	// Set bot's username
	SetUsername(username string) Config

	// Set bot's icon
	SetIcon(icon Icon) Config

	// Set diagnostics logger
	SetLogger(logger *log.Logger) Config

	// Create new reporter
	CreateReporter() (Reporter, error)
}

// Icon defines either an icon URL or emoji
type Icon interface {
	apply(args url.Values)
}

// Reporter defines methods to create new updatable messages
type Reporter interface {
	// Start new updatable message
	BeginMessage(text string) (Message, error)
}

// Message represents a single updatable Slack message
type Message interface {
	// Update message text
	Update(text string) error

	// Delete message
	Delete() error
}
