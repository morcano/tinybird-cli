package tokensapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TokensResponse struct {
	Tokens []Token `json:"tokens"`
}

type Token struct {
	ID          string  `json:"id"`
	Token       string  `json:"token"`
	Scopes      []Scope `json:"scopes"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Origin      Origin  `json:"origin"`
	Host        string  `json:"host"`
}

type Scope struct {
	Type     string `json:"type"`
	Resource string `json:"resource"`
	Filter   string `json:"filter"`
}

type Origin struct {
	Type string `json:"type"`
}

type Client struct {
	Token      string
	BaseURL    string
	HTTPClient *http.Client
}

type PutTokenParams struct {
	Token       string   `json:"token"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Scope       []string `json:"scope"`
}

func (token TokensResponse) FilterByName(filterName string) TokensResponse {
	var filteredTokens []Token

	for _, token := range token.Tokens {
		if strings.Contains(token.Name, filterName) {
			filteredTokens = append(filteredTokens, token)
		}
	}

	return TokensResponse{
		Tokens: filteredTokens,
	}
}

func (token Token) GetAccountIdFromName() string {
	parts := strings.Split(token.Name, "_")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		BaseURL:    "https://api.tinybird.co",
		HTTPClient: &http.Client{},
	}
}

func (c *Client) Get() (*TokensResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v0/tokens", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokensResponse TokensResponse
	err = json.NewDecoder(resp.Body).Decode(&tokensResponse)

	return &tokensResponse, err
}

func (c *Client) Put(params PutTokenParams) error {
	baseURL, err := url.Parse(fmt.Sprintf("%s/v0/tokens/%s", c.BaseURL, params.Token))
	if err != nil {
		return err
	}

	query := baseURL.Query()
	if params.Name != "" {
		query.Set("name", params.Name)
	}
	if params.Description != "" {
		query.Set("description", params.Description)
	}
	for _, scope := range params.Scope {
		query.Add("scope", scope)
	}

	baseURL.RawQuery = query.Encode()

	req, err := http.NewRequest("PUT", baseURL.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("bad status")
	}

	return nil
}
