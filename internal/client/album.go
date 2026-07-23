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

type AlbumUserResponseDto struct {
	Role string          `json:"role"`
	User UserResponseDto `json:"user"`
}

type ContributorCountResponseDto struct {
	AssetCount int64  `json:"assetCount"`
	UserId     string `json:"userId"`
}

type AlbumResponseDto struct {
	AlbumName                  string                        `json:"albumName"`
	AlbumThumbnailAssetId      string                        `json:"albumThumbnailAssetId"`
	AlbumUsers                 []AlbumUserResponseDto        `json:"albumUsers"`
	AssetCount                 int64                         `json:"assetCount"`
	ContributorCounts          []ContributorCountResponseDto `json:"contributorCounts"`
	CreatedAt                  string                        `json:"createdAt"`
	Description                string                        `json:"description"`
	EndDate                    string                        `json:"endDate"`
	HasSharedLink              bool                          `json:"hasSharedLink"`
	Id                         string                        `json:"id"`
	IsActivityEnabled          bool                          `json:"isActivityEnabled"`
	LastModifiedAssetTimestamp string                        `json:"lastModifiedAssetTimestamp"`
	Order                      string                        `json:"order"`
	Shared                     bool                          `json:"shared"`
	StartDate                  string                        `json:"startDate"`
	UpdatedAt                  string                        `json:"updatedAt"`
}

type AlbumUserCreateDto struct {
	Role   string `json:"role"`
	UserId string `json:"userId"`
}

type CreateAlbumDto struct {
	AlbumName   string               `json:"albumName"`
	AlbumUsers  []AlbumUserCreateDto `json:"albumUsers,omitempty"`
	AssetIds    []string             `json:"assetIds,omitempty"`
	Description string               `json:"description"`
}

type UpdateAlbumDto struct {
	AlbumName             string `json:"albumName"`
	AlbumThumbnailAssetId string `json:"albumThumbnailAssetId,omitempty"`
	Description           string `json:"description"`
	IsActivityEnabled     bool   `json:"isActivityEnabled"`
	Order                 string `json:"order"`
}

func (c *Client) GetAlbums() ([]AlbumResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/albums", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	albums := []AlbumResponseDto{}
	err = json.Unmarshal(body, &albums)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (c *Client) GetAlbum(albumId string) (*AlbumResponseDto, error) {
	album, err := get_by_id[AlbumResponseDto](c, albumId, "albums")

	return album, err
}

func (c *Client) CreateAlbum(album CreateAlbumDto) (*AlbumResponseDto, error) {
	newAlbum, err := post[AlbumResponseDto](c, "albums", album)

	return newAlbum, err
}

func (c *Client) UpdateAlbum(albumId string, album UpdateAlbumDto) (*AlbumResponseDto, error) {
	newAlbum, err := patch[AlbumResponseDto](c, albumId, "albums", album)

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
