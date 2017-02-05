package slackreporting

import (
	"errors"
	"io/ioutil"
	"log"
)

type options struct {
	URL         string
	AccessToken string
	Channel     string
	Username    string
	Icon        Icon
	Logger      *log.Logger
}

const (
	// DefaultURL is a default value for Config.URL
	DefaultURL = "https://slack.com/"

	// DefaultChannel is a default value for Config.Channel
	DefaultChannel = "#general"
)

// NewConfig creates new config object
func NewConfig(accessToken string) Config {
	return &options{
		URL:         DefaultURL,
		AccessToken: accessToken,
		Channel:     DefaultChannel,
		Logger:      log.New(ioutil.Discard, "", 0),
	}
}

// Set Slack API URL (e.g. "https://slack.com/")
func (c *options) SetURL(url string) Config {
	c.URL = url
	return c
}

// Set Slack channel or username (e.g. "#channel" or "@SetUsername")
func (c *options) SetChannel(channel string) Config {
	c.Channel = channel
	return c
}

// Set bot's username
func (c *options) SetUsername(username string) Config {
	c.Username = username
	return c
}

// Set bot's icon
func (c *options) SetIcon(icon Icon) Config {
	c.Icon = icon
	return c
}

// Set diagnostics logger
func (c *options) SetLogger(logger *log.Logger) Config {
	c.Logger = logger
	return c
}

// Create new reporter
func (c *options) CreateReporter() (Reporter, error) {
	if c.URL == "" {
		return nil, errors.New("URL is not set")
	}
	if c.AccessToken == "" {
		return nil, errors.New("Access token is not set")
	}
	if c.Channel == "" {
		return nil, errors.New("Channel is not set")
	}

	url := c.URL
	if url[len(url)-1] != '/' {
		url = url + "/"
	}

	url = url + "api/"
	c.URL = url

	r := &reporter{opt: *c}
	e := r.initialize()
	if e != nil {
		return nil, e
	}

	return r, nil
}
