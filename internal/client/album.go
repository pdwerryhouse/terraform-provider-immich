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
	"strings"
)

type AlbumUpdate struct {
	AlbumName   string `json:"albumName,omitempty"`
	Description string `json:"description,omitempty"`
}

type Album struct {
	ID                    string `json:"id"`
	AlbumName             string `json:"albumName"`
	AlbumThumbnailAssetId string `json:"albumThumbnailAssetId"`
	Description           string `json:"description"`
	Shared                bool   `json:"shared"`
	HasSharedLink         bool   `json:"hasSharedLink"`
	Order                 string `json:"order"`
	IsActivityEnabled     bool   `json:"isActivityEnabled"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	OwnerId               string `json:"ownerId"`
	Owner                 Owner  `json:"owner"`
}

func (c *Client) GetAlbums() ([]Album, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/albums", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	albums := []Album{}
	err = json.Unmarshal(body, &albums)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (c *Client) GetAlbum(albumId string) (*Album, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/albums/%s", c.Endpoint, albumId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	album := Album{}
	err = json.Unmarshal(body, &album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

func (c *Client) CreateAlbum(album AlbumUpdate) (*Album, error) {

	rb, err := json.Marshal(album)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/albums", c.Endpoint), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newAlbum := Album{}
	err = json.Unmarshal(body, &newAlbum)
	if err != nil {
		return nil, err
	}

	return &newAlbum, nil
}

func (c *Client) UpdateAlbum(albumId string, album AlbumUpdate) (*Album, error) {

	newAlbum, err := patch[Album](c, albumId, "albums", album)

	return newAlbum, err
}

func (c *Client) DeleteAlbum(albumId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/albums/%s", c.Endpoint, albumId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
