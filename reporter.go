package slackreporting

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type reporter struct {
	opt options
}

type response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type responseValidate interface {
	validate() (bool, string)
}

type authTestResponse struct {
	response
	URL    string `json:"url"`
	Team   string `json:"team"`
	User   string `json:"user"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}

func (r authTestResponse) validate() (bool, string) {
	return r.OK, r.Error
}

func (r *reporter) initialize() error {
	var resp authTestResponse
	err := r.callMethod("auth.test", make(url.Values), &resp)
	if err != nil {
		return err
	}

	r.printf("Logged in as [%s] '%s' from [%s] '%s'\n", resp.UserID, resp.User, resp.TeamID, resp.Team)
	return nil
}

// Start new updatable message
func (r *reporter) BeginMessage(text string) (Message, error) {
	msg := &message{reporter: r}
	err := msg.Update(text)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *reporter) callMethod(method string, args url.Values, out interface{}) error {
	url := r.opt.URL + method
	args["token"] = []string{r.opt.AccessToken}

	tracetxt := "POST " + url + "\n"
	for k := range args {
		tracetxt += k + ": "
		for _, s := range args[k] {
			tracetxt += s + ";"
		}
		tracetxt += "\n"
	}
	r.printf("%s\n", tracetxt)

	resp, err := http.DefaultClient.PostForm(url, args)

	if err != nil {
		r.printf("POST %s -> %s\n", url, err)
		return err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.printf("POST %s: Unable to read response (%s)\n", url, err)
		return err
	}

	err = json.Unmarshal(respBody, out)
	if err != nil {
		r.printf("POST %s: Unable to parse response (%s)\n", url, err)
		return err
	}

	result, ok := out.(responseValidate)
	if !ok {
		r.printf("POST %s: Result type is not an API Response\n", url)
		return errors.New("TypeAssertError")
	}

	ok, errcode := result.validate()
	if !ok {
		r.printf("POST %s -> %s\n", url, errcode)
		return errors.New(errcode)
	}

	r.printf("POST %s -> ok\n", url)
	return nil
}

func (r *reporter) printf(format string, v ...interface{}) {
	r.opt.Logger.Printf(format, v...)
}
