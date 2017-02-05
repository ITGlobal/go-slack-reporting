package slackreporting

import "net/url"

type message struct {
	reporter *reporter
	ts       string
	c        string
}

func (m *message) Update(text string) error {
	if m.ts != "" && m.c != "" {
		err := m.updateMessage(text)
		if err != nil {
			err = m.postMessage(text)
		}
		return err
	}

	err := m.postMessage(text)
	return err
}

func (m *message) Delete() error {
	if m.ts != "" && m.c != "" {
		err := m.deleteMessage()
		return err
	}
	return nil
}

func (m *message) postMessage(text string) error {
	args := make(url.Values)
	args["text"] = []string{text}
	args["channel"] = []string{m.reporter.opt.Channel}
	if m.reporter.opt.Username != "" {
		args["username"] = []string{m.reporter.opt.Username}
	}
	if m.reporter.opt.Icon != nil {
		m.reporter.opt.Icon.apply(args)
	}

	m.reporter.printf("Posting message to '%s'...\n", m.reporter.opt.Channel)
	var resp postResponse
	err := m.reporter.callMethod("chat.postMessage", args, &resp)
	if err != nil {
		return err
	}

	m.ts = resp.Ts
	m.c = resp.Channel
	m.reporter.printf("Message has been posted: ts=%s, c=%s\n", m.ts, m.c)
	return nil
}

type postResponse struct {
	response

	Ts      string `json:"ts"`
	Channel string `json:"channel"`
}

func (r postResponse) validate() (bool, string) {
	return r.OK, r.Error
}

func (m *message) updateMessage(text string) error {
	args := make(url.Values)
	args["text"] = []string{text}
	args["ts"] = []string{m.ts}
	args["channel"] = []string{m.c}

	m.reporter.printf("Updating message ts=%s, c=%s...\n", m.ts, m.c)

	var resp updateResponse
	err := m.reporter.callMethod("chat.update", args, &resp)
	if err != nil {
		return err
	}

	m.ts = resp.Ts
	m.c = resp.Channel
	m.reporter.printf("Message has been updated: ts=%s, c=%s\n", m.ts, m.c)
	return nil
}

type updateResponse struct {
	response

	Ts      string `json:"ts"`
	Channel string `json:"channel"`
}

func (r updateResponse) validate() (bool, string) {
	return r.OK, r.Error
}

func (m *message) deleteMessage() error {
	args := make(url.Values)
	args["ts"] = []string{m.ts}
	args["channel"] = []string{m.c}

	m.reporter.printf("Deleting message ts=%s, c=%s...\n", m.ts, m.c)

	var resp deleteResponse
	err := m.reporter.callMethod("chat.delete", args, &resp)
	if err != nil {
		return err
	}

	m.ts = ""
	m.c = ""
	m.reporter.printf("Message has been deleted: ts=%s, c=%s\n", resp.Ts, resp.Channel)
	return nil
}

type deleteResponse struct {
	response

	Ts      string `json:"ts"`
	Channel string `json:"channel"`
}

func (r deleteResponse) validate() (bool, string) {
	return r.OK, r.Error
}
