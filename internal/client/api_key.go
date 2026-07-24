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
	"net/http"
)

type ApiKeyResponseDto struct {
	CreatedAt   string   `json:"createdAt"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	UpdatedAt   string   `json:"updatedAt"`
}

type ApiKeyCreateDto struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type ApiKeyCreateResponseDto struct {
	ApiKey ApiKeyResponseDto `json:"apiKey"`
	Secret string            `json:"secret"`
}

type ApiKeyUpdateDto struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func (c *Client) GetApiKeys() ([]ApiKeyResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api-keys", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiKeys := []ApiKeyResponseDto{}
	err = json.Unmarshal(body, &apiKeys)
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (c *Client) GetApiKey(apiKeyId string) (*ApiKeyResponseDto, error) {
	apiKey, err := get_by_id[ApiKeyResponseDto](c, apiKeyId, "api-keys")

	return apiKey, err
}

func (c *Client) GetMyApiKey() (*ApiKeyResponseDto, error) {
	apiKey, err := get[ApiKeyResponseDto](c, "api-keys/me")

	return apiKey, err
}

func (c *Client) CreateApiKey(apiKey ApiKeyCreateDto) (*ApiKeyCreateResponseDto, error) {
	newApiKey, err := post[ApiKeyCreateResponseDto](c, "api-keys", apiKey)

	return newApiKey, err
}

func (c *Client) UpdateApiKey(apiKeyId string, apiKey ApiKeyUpdateDto) (*ApiKeyResponseDto, error) {
	newApiKey, err := put[ApiKeyResponseDto](c, apiKeyId, "api-keys", apiKey)

	return newApiKey, err
}

func (c *Client) DeleteApiKey(apiKeyId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api-keys/%s", c.Endpoint, apiKeyId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
