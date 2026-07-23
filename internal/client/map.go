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

type MapMarkerResponseDto struct {
	City    string  `json:"city"`
	Country string  `json:"country"`
	Id      string  `json:"id"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	State   string  `json:"state"`
}

type MapReverseGeocodeResponseDto struct {
	City    string `json:"city"`
	Country string `json:"country"`
	State   string `json:"state"`
}

func (c *Client) GetMapMarkers() ([]MapMarkerResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/map/markers", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	markers := []MapMarkerResponseDto{}
	err = json.Unmarshal(body, &markers)
	if err != nil {
		return nil, err
	}

	return markers, nil
}

func (c *Client) GetReverseGeocode(lat float64, lon float64) ([]MapReverseGeocodeResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/map/reverse-geocode", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("lat", fmt.Sprintf("%f", lat))
	q.Add("lon", fmt.Sprintf("%f", lon))
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	responses := []MapReverseGeocodeResponseDto{}
	err = json.Unmarshal(body, &responses)
	if err != nil {
		return nil, err
	}

	return responses, nil
}
