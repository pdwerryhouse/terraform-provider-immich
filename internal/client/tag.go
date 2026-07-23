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

type Tag struct {
	ID        string `json:"id"`
	Color     string `json:"color"`
	CreatedAt string `json:"createdAt"`
	Name      string `json:"name"`
	ParentId  string `json:"parentId"`
	UpdatedAt string `json:"updatedAt"`
	Value     string `json:"value"`
}

type TagResponseDto struct {
	Color     string `json:"color"`
	CreatedAt string `json:"createdAt"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	ParentId  string `json:"parentId"`
	UpdatedAt string `json:"updatedAt"`
	Value     string `json:"value"`
}

type TagUpdate struct {
	Name     string `json:"name"`
	Color    string `json:"color,omitempty"`
	ParentId string `json:"parentId,omitempty"`
}

func (c *Client) GetTags() ([]Tag, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tags", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	tags := []Tag{}
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (c *Client) GetTag(tagId string) (*Tag, error) {
	tag, err := get_by_id[Tag](c, tagId, "tags")

	return tag, err
}

func (c *Client) CreateTag(tag TagUpdate) (*Tag, error) {
	newTag, err := post[Tag](c, "tags", tag)

	return newTag, err
}

func (c *Client) UpdateTag(tagId string, tag TagUpdate) (*Tag, error) {
	newTag, err := put[Tag](c, tagId, "tags", tag)

	return newTag, err
}

func (c *Client) DeleteTag(tagId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/tags/%s", c.Endpoint, tagId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
