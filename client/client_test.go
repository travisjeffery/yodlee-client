package yodlee

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var c *Client

func init() {
	c = &Client{
		Login:    os.Getenv("YODLEE_COB_LOGIN"),
		Password: os.Getenv("YODLEE_COB_PASSWORD"),
	}
}

func TestGetCobSessionToken(t *testing.T) {
	token, err := c.GetCobSessionToken()
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestGetUserSessionToken(t *testing.T) {
	c.Authenticate()
	login := os.Getenv("YODLEE_USER_LOGIN")
	pass := os.Getenv("YODLEE_USER_PASSWORD")
	token, err := c.GetUserSessionToken(login, pass)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestGetAccounts(t *testing.T) {
	c.Authenticate()
	login := os.Getenv("YODLEE_USER_LOGIN")
	pass := os.Getenv("YODLEE_USER_PASSWORD")
	token, err := c.GetUserSessionToken(login, pass)
	assert.Nil(t, err)
	accounts, errs := c.GetAccounts(token)
	assert.Nil(t, errs)
	assert.NotZero(t, len(accounts))
}

func TestGetTransactions(t *testing.T) {
	c.Authenticate()
	login := os.Getenv("YODLEE_USER_LOGIN")
	pass := os.Getenv("YODLEE_USER_PASSWORD")
	token, err := c.GetUserSessionToken(login, pass)
	assert.Nil(t, err)
	transactions, errs := c.GetTransactions(token, NewGetTransactionInput())
	assert.Nil(t, errs)
	assert.NotNil(t, transactions)
}
