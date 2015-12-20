package client

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
)

type Client struct {
	Login    string
	Password string
}

func New(login, password string) *Client {
	return &Client{
		Login:    login,
		Password: password,
	}
}

func (c *Client) GetCobSessionToken() (string, []error) {
	req := gorequest.New()
	_, body, errs := req.Post("https://rest.developer.yodlee.com/services/srest/restserver/v1.0/authenticate/coblogin").
		Type("form").
		Send("cobrandLogin=" + c.Login).
		Send("cobrandPassword=" + c.Password).
		End()
	if errs != nil {
		return "", errs
	}
	var j struct {
		CobrandConversationCredentials struct {
			SessionToken string
		}
	}
	if err := json.Unmarshal([]byte(body), &j); err != nil {
		return "", []error{err}
	}
	return j.CobrandConversationCredentials.SessionToken, nil
}
