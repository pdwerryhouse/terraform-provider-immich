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

type PartnerCreateDto struct {
	SharedWithId string `json:"sharedWithId"`
}

type PartnerResponseDto struct {
	AvatarColor      string `json:"avatarColor"`
	Email            string `json:"email"`
	Id               string `json:"id"`
	InTimeline       bool   `json:"inTimeline"`
	Name             string `json:"name"`
	ProfileChangedAt string `json:"profileChangedAt"`
	ProfileImagePath string `json:"profileImagePath"`
}

type PartnerUpdateDto struct {
	InTimeline bool `json:"inTimeline"`
}

func (c *Client) GetPartners() ([]PartnerResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/partners", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("direction", "shared-by")
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	partners := []PartnerResponseDto{}
	err = json.Unmarshal(body, &partners)
	if err != nil {
		return nil, err
	}

	return partners, nil
}

func (c *Client) GetPartner(partnerId string) (*PartnerResponseDto, error) {
	partner, err := get_by_id[PartnerResponseDto](c, partnerId, "partners")

	return partner, err
}

func (c *Client) CreatePartner(partner PartnerCreateDto) (*PartnerResponseDto, error) {
	newPartner, err := post[PartnerResponseDto](c, "partners", partner)

	return newPartner, err
}

func (c *Client) UpdatePartner(partnerId string, partner PartnerUpdateDto) (*PartnerResponseDto, error) {
	newPartner, err := patch[PartnerResponseDto](c, partnerId, "partners", partner)

	return newPartner, err
}

func (c *Client) DeletePartner(partnerId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/partners/%s", c.Endpoint, partnerId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
