// Copyright (C) 2026 Paul Dwerryhouse <paul@dwerryhouse.com.au>
//
// This file is part of terraform-provider-immich.
//
// terraform-provider-immich is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// terraform-provider-immich is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with terraform-provider-immich.  If not, see <https://www.gnu.org/licenses/>.

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func get[T any](c *Client, path string) (*T, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.Endpoint, path), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var item T
	err = json.Unmarshal(body, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil

}

func get_by_id[T any](c *Client, id string, path string) (*T, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.Endpoint, path, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var item T
	err = json.Unmarshal(body, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil

}

func post[T any](c *Client, path string, item any) (*T, error) {
	rb, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.Endpoint, path), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var newItem T
	err = json.Unmarshal(body, &newItem)
	if err != nil {
		return nil, err
	}

	return &newItem, nil
}

func patch[T any](c *Client, id string, path string, item any) (*T, error) {

	rb, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s/%s", c.Endpoint, path, id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var newItem T
	err = json.Unmarshal(body, &newItem)
	if err != nil {
		return nil, err
	}

	return &newItem, nil
}

func put[T any](c *Client, id string, path string, item any) (*T, error) {

	rb, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%s", c.Endpoint, path, id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var newItem T
	err = json.Unmarshal(body, &newItem)
	if err != nil {
		return nil, err
	}

	return &newItem, nil
}
