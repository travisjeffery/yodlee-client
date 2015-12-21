package client

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
)

type Client struct {
	Login        string
	Password     string
	SessionToken string
}

func New(login, password string) *Client {
	return &Client{
		Login:    login,
		Password: password,
	}
}

func (c *Client) Authenticate() []error {
	token, errs := c.GetCobSessionToken()
	if errs != nil {
		return errs
	}
	c.SessionToken = token
	return nil
}

func (c *Client) GetCobSessionToken() (string, []error) {
	var j struct {
		CobrandConversationCredentials struct {
			SessionToken string
		}
	}
	errs := request("https://rest.developer.yodlee.com/services/srest/restserver/v1.0/authenticate/coblogin", struct {
		CobrandLogin    string `json:"cobrandLogin"`
		CobrandPassword string `json:"cobrandPassword"`
	}{
		c.Login,
		c.Password,
	}, &j)
	if errs != nil {
		return "", errs
	}
	return j.CobrandConversationCredentials.SessionToken, nil
}

func (c *Client) GetUserSessionToken(login, password string) (string, []error) {
	var j struct {
		UserContext struct {
			ConversationCredentials struct {
				SessionToken string
			}
		}
	}
	errs := request("https://rest.developer.yodlee.com/services/srest/restserver/v1.0/authenticate/login", struct {
		Login           string `json:"login"`
		Password        string `json:"password"`
		CobSessionToken string `json:"cobSessionToken"`
	}{
		login,
		password,
		c.SessionToken,
	}, &j)
	if errs != nil {
		return "", errs
	}
	return j.UserContext.ConversationCredentials.SessionToken, nil
}

func (c *Client) GetAccounts(token string) ([]map[string]interface{}, []error) {
	var j []map[string]interface{}
	errs := request("https://rest.developer.yodlee.com/services/srest/restserver/v1.0/jsonsdk/SiteAccountManagement/getSiteAccounts", struct {
		CobSessiontoken  string `json:"cobSessionToken"`
		UserSessionToken string `json:"userSessionToken"`
	}{
		c.SessionToken,
		token,
	}, &j)
	return j, errs
}

type GetTransactionInput struct {
	ContainerType    string `json:"transactionSearchRequest.containerType"`
	HigherFetchLimit string `json:"transactionSearchRequest.higherFetchLimit"`
	LowerFetchLimit  string `json:"transactionSearchRequest.lowerFetchLimit"`
	IgnoreUserInput  string `json:"transactionSearchRequest.ignoreUserInput"`
	EndNumber        int    `json:"transactionSearchRequest.resultRange.endNumber"`
	StartNumber      int    `json:"transactionSearchRequest.resultRange.startNumber"`
	CurrencyCode     string `json:"transactionSearchRequest.searchFilter.currencyCode"`
}

func NewGetTransactionInput() *GetTransactionInput {
	return &GetTransactionInput{
		ContainerType:    "All",
		HigherFetchLimit: "500",
		LowerFetchLimit:  "1",
		IgnoreUserInput:  "true",
		EndNumber:        5,
		StartNumber:      1,
		CurrencyCode:     "USD",
	}
}

func (c *Client) GetTransactions(token string, input *GetTransactionInput) (map[string]interface{}, []error) {
	var j map[string]interface{}
	errs := request("https://rest.developer.yodlee.com/services/srest/restserver/v1.0/jsonsdk/TransactionSearchService/executeUserSearchRequest", struct {
		*GetTransactionInput
		CobSessionToken  string `json:"cobSessionToken"`
		UserSessionToken string `json:"userSessionToken"`
	}{
		input,
		c.SessionToken,
		token,
	}, &j)
	if errs != nil {
		return nil, errs
	}
	return j, nil
}

func request(url string, content interface{}, data interface{}) []error {
	req := gorequest.New()
	_, body, errs := req.Post(url).
		Type("form").
		Send(content).
		End()
	if errs != nil {
		return errs
	}
	if err := json.Unmarshal([]byte(body), data); err != nil {
		return []error{err}
	}
	return nil
}
