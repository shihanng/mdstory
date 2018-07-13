package medium

import (
	"github.com/pkg/errors"
	medium "mod/github.com/Medium/medium-sdk-go@v0.0.0-20171230201202-4daca056cf6a"
)

// Client is a wrapper for medium.Medium with user information.
type Client interface {
	Username() string
}

// New gets the user id from Medium with the given accessToken and creates a new client.
func New(accessToken string) (Client, error) {
	m := medium.NewClientWithAccessToken(accessToken)

	u, err := m.GetUser("")
	if err != nil {
		return nil, errors.Wrap(err, "get Medium's user")
	}

	return &client{
		medium:   m,
		username: u.Username,
		userID:   u.ID,
	}, nil
}

type client struct {
	medium   *medium.Medium
	username string
	userID   string
}

// Username of the user identified by the access token.
func (c *client) Username() string {
	return c.username
}
