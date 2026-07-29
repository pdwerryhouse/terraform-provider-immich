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

func (c *Client) GetAssetsByOriginalPath(path string) ([]AssetResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/view/folder", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("path", fmt.Sprintf("%s", path))
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	responses := []AssetResponseDto{}
	err = json.Unmarshal(body, &responses)
	if err != nil {
		return nil, err
	}

	return responses, nil
}

func (c *Client) GetUniqueOriginalPaths() ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/view/folder/unique-paths", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	responses := []string{}
	err = json.Unmarshal(body, &responses)
	if err != nil {
		return nil, err
	}

	return responses, nil
}
