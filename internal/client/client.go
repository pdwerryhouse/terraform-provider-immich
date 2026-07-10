// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Endpoint   string
	HTTPClient *http.Client
	ApiKey     string
}

func NewClient(endpoint *string, apikey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Endpoint:   *endpoint,
		ApiKey:     *apikey,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	apikey := c.ApiKey

	req.Header.Set("x-api-key", apikey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

type Owner struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	ProfileImagePath string `json:"profileImagePath"`
	AvatarColor      string `json:"avatarColor"`
	ProfileChangedAt string `json:"profileChangedAt"`
}
