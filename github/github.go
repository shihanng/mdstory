package github

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Client is a wrapper for github.Client with user information.
type Client interface {
	Username() string
}

// New gets the user id from GitHub with the given accessToken and creates a new client.
func New(accessToken string) (Client, error) {
	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	oauth2Client := oauth2.NewClient(ctx, tokenSource)
	c := github.NewClient(oauth2Client)

	u, _, err := c.Users.Get(ctx, "")
	if err != nil {
		return nil, errors.Wrap(err, "get GitHub's user")
	}

	return &client{
		github: c,
		userID: u.GetLogin(),
	}, nil
}

type client struct {
	github   *github.Client
	username string
	userID   string
}

// Username of the user identified by the access token.
func (c *client) Username() string {
	return c.userID
}
