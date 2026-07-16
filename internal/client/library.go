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

type LibraryUpdate struct {
	ExclusionPatterns []string `json:"exclusionPatterns"`
	ImportPaths       []string `json:"importPaths"`
	Name              string   `json:"name"`
	OwnerId           string   `json:"ownerId"`
}

type Library struct {
	AssetCount        int64    `json:"assetCount"`
	CreatedAt         string   `json:"createdAt"`
	ExclusionPatterns []string `json:"exclusionPatterns"`
	Id                string   `json:"id"`
	ImportPaths       []string `json:"importPaths"`
	Name              string   `json:"name"`
	OwnerId           string   `json:"ownerId"`
	RefreshedAt       string   `json:"refreshedAt"`
	UpdatedAt         string   `json:"updatedAt"`
}

func (c *Client) GetLibraries() ([]Library, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/libraries", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	libraries := []Library{}
	err = json.Unmarshal(body, &libraries)
	if err != nil {
		return nil, err
	}

	return libraries, nil
}

func (c *Client) GetLibrary(libraryId string) (*Library, error) {
	library, err := get_by_id[Library](c, libraryId, "libraries")

	return library, err
}

func (c *Client) CreateLibrary(library LibraryUpdate) (*Library, error) {
	newLibrary, err := post[Library](c, "libraries", library)

	return newLibrary, err
}

func (c *Client) UpdateLibrary(libraryId string, library LibraryUpdate) (*Library, error) {
	newLibrary, err := patch[Library](c, libraryId, "libraries", library)

	return newLibrary, err
}

func (c *Client) DeleteLibrary(libraryId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/libraries/%s", c.Endpoint, libraryId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
