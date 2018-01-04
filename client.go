package dgo2poc

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Client for the Discord REST API.
type Client interface {
	// Make an arbitrary request.
	Request(ctx context.Context, method, urlStr string, body []byte, opts ...ReqOption) ([]byte, error)

	// Make an arbitrary request, which returns a JSON object.
	RequestJSON(ctx context.Context, method, urlStr string, body []byte, out interface{}, opts ...ReqOption) error

	// Returns a user object for a given user ID.
	// The special ID "@me" returns the authenticating user.
	User(ctx context.Context, id string, opts ...ReqOption) (*discordgo.User, error)

	// Returns a gateway for a websocket connection.
	// Depending on the type of token used, this will call either /gateway or /gateway/bot;
	// the two are identical, except the latter will also provide a suggested shard count.
	Gateway(ctx context.Context, opts ...ReqOption) (*Gateway, error)

	// Returns the token used by this client.
	Token() *oauth2.Token
}

type client struct {
	Tok        *oauth2.Token
	HTTPClient *http.Client
	BaseURL    string
	Opts       []ReqOption
}

// Create a new client. Use UserToken() or BotToken() to wrap a token.
func NewClient(t *oauth2.Token, opts ...ReqOption) Client {
	return &client{
		Tok: t,
		HTTPClient: oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(t),
		),
		BaseURL: BaseURL,
		Opts:    opts,
	}
}

func (c *client) Request(ctx context.Context, method, urlStr string, body []byte, opts ...ReqOption) ([]byte, error) {
	// Create a request, set defaults.
	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json")

	// Apply options.
	reqOpts := ReqOptions{Request: req}
	for _, opt := range c.Opts {
		opt(&reqOpts)
	}
	for _, opt := range opts {
		opt(&reqOpts)
	}

	// Send the request...
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Always read and close the body, else connections can't be reused.
	data, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Handle status codes.
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		var apiErr APIError
		if err := json.Unmarshal(data, &apiErr); err != nil {
			return data, errors.Errorf("%d: %s", resp.StatusCode, string(data))
		}
		return data, errors.Errorf("%d: %s", resp.StatusCode, apiErr)
	}

	return data, nil
}

func (c *client) RequestJSON(ctx context.Context, method, urlStr string, body []byte, out interface{}, opts ...ReqOption) error {
	data, err := c.Request(ctx, method, urlStr, body, opts...)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func (c *client) User(ctx context.Context, id string, opts ...ReqOption) (*discordgo.User, error) {
	var user discordgo.User
	return &user, c.RequestJSON(ctx, "GET", c.BaseURL+EndpointUser(id), nil, &user, opts...)
}

func (c *client) Gateway(ctx context.Context, opts ...ReqOption) (*Gateway, error) {
	ep := EndpointGateway
	if c.Tok.TokenType == "Bot" {
		ep = EndpointGatewayBot
	}
	var gw Gateway
	return &gw, c.RequestJSON(ctx, "GET", c.BaseURL+ep, nil, &gw, opts...)
}

func (c *client) Token() *oauth2.Token {
	return c.Tok
}
