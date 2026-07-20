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

type PersonUpdate struct {
	BirthDate  string `json:"birthDate,omitempty"`
	Color      string `json:"color,omitzero"`
	IsFavorite bool   `json:"isFavorite,omitempty"`
	IsHidden   bool   `json:"isHidden,omitempty"`
	Name       string `json:"name,omitempty"`
}

type Person struct {
	ID            string `json:"id"`
	BirthDate     string `json:"birthDate,omitempty"`
	Color         string `json:"color,omitempty"`
	IsFavorite    bool   `json:"isFavorite,omitempty"`
	IsHidden      bool   `json:"isHidden,omitempty"`
	Name          string `json:"name,omitempty"`
	ThumbNailPath string `json:"thumbNailPath"`
	UpdatedAt     string `json:"updatedAt"`
}

type PeopleResponse struct {
	HasNextPage bool     `json:"hasNextPage"`
	Total       int64    `json:"total"`
	Hidden      int64    `json:"hidden"`
	People      []Person `json:"people"`
}

// XXX Update this to handle pages
func (c *Client) GetPeople() ([]Person, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/people", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	people_response := PeopleResponse{}
	err = json.Unmarshal(body, &people_response)
	if err != nil {
		return nil, err
	}

	return people_response.People, nil
}

func (c *Client) GetPerson(personId string) (*Person, error) {
	person, err := get_by_id[Person](c, personId, "people")

	return person, err
}

func (c *Client) CreatePerson(person PersonUpdate) (*Person, error) {
	newPerson, err := post[Person](c, "people", person)

	return newPerson, err
}

func (c *Client) UpdatePerson(personId string, person PersonUpdate) (*Person, error) {
	newPerson, err := patch[Person](c, personId, "people", person)

	return newPerson, err
}

func (c *Client) DeletePerson(personId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/people/%s", c.Endpoint, personId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
