// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AlbumActivityUpdate struct {
	IsActivityEnabled bool `json:"isActivityEnabled"`
}

func (c *Client) UpdateAlbumActivity(albumId string, album AlbumActivityUpdate) (*Album, error) {

	rb, err := json.Marshal(album)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/albums/%s", c.Endpoint, albumId), strings.NewReader(string(rb)))
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
