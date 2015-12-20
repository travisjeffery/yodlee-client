package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var c *Client

func init() {
	c = &Client{
		Login:    os.Getenv("YODLEE_LOGIN"),
		Password: os.Getenv("YODLEE_PASSWORD"),
	}
}

func TestGetCobSessionToken(t *testing.T) {
	token, err := c.GetCobSessionToken()
	assert.Nil(t, err)
	assert.NotEmpty(t, token, token)
}
