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

type ActivityResponseDto struct {
	AssetId   string          `json:"assetId"`
	Comment   string          `json:"comment"`
	CreatedAt string          `json:"createdAt"`
	Id        string          `json:"id"`
	Type      string          `json:"type"`
	User      UserResponseDto `json:"user"`
}

type ActivityCreateDto struct {
	AlbumId string `json:"albumId"`
	AssetId string `json:"assetId,omitempty"`
	Comment string `json:"comment"`
	Type    string `json:"type"`
}

type ActivityStatisticsResponseDto struct {
	Comments int64 `json:"comments"`
	Likes    int64 `json:"likes"`
}

func (c *Client) GetActivities(albumId string) ([]ActivityResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/activities", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("albumId", albumId)
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	activities := []ActivityResponseDto{}
	err = json.Unmarshal(body, &activities)
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (c *Client) CreateActivity(activity ActivityCreateDto) (*ActivityResponseDto, error) {
	newActivity, err := post[ActivityResponseDto](c, "activities", activity)

	return newActivity, err
}

func (c *Client) GetActivityStatistics(albumId string, assetId string) (*ActivityStatisticsResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/activities/statistics", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("albumId", albumId)
	q.Add("assetId", assetId)
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := &ActivityStatisticsResponseDto{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DeleteActivity(activityId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/activities/%s", c.Endpoint, activityId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
