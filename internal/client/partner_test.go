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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePartner(t *testing.T) {

	data := PartnerResponseDto{
		Id:               "bcfbcff7-d844-40d0-8e87-50afae70a628",
		Email:            "test447a@example.com",
		Name:             "test447a",
		ProfileImagePath: "",
		AvatarColor:      "red",
		ProfileChangedAt: "2026-07-09T05:27:10.749066+00:00",
		InTimeline:       false,
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	partner, err := c.CreatePartner(PartnerCreateDto{
		SharedWithId: "bcfbcff7-d844-40d0-8e87-50afae70a628",
	})

	if err != nil {
		t.Error("Got an error when trying to run CreatePartner.\n")
	}

	assert.Equal(t, partner.Id, "bcfbcff7-d844-40d0-8e87-50afae70a628", "partner.Id is incorrect")
}

func TestGetPartners(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetPartners should use the GET method")

		partners := [1]PartnerResponseDto{{
			Id:               id,
			Email:            "test447a@example.com",
			Name:             "test447a",
			ProfileImagePath: "",
			AvatarColor:      "red",
			ProfileChangedAt: "2026-07-09T05:27:10.749066+00:00",
			InTimeline:       false,
		}}

		data, err := json.Marshal(partners)
		if err != nil {
			t.Error("Got an error when unmarshalling partner.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	partners, err := c.GetPartners()
	if err != nil {
		t.Error("Got an error when trying to run GetPartner.\n")
	}

	// XXX fix this
	for _, partner := range partners {
		assert.Equal(t, id, partner.Id, "partner.Id is incorrect")
		assert.Equal(t, "test447a", partner.Name, "partner.Name is incorrect")
	}
}

func TestGetPartner(t *testing.T) {

	data := PartnerResponseDto{
		Id:               "bcfbcff7-d844-40d0-8e87-50afae70a628",
		Email:            "test447a@example.com",
		Name:             "test447a",
		ProfileImagePath: "",
		AvatarColor:      "red",
		ProfileChangedAt: "2026-07-09T05:27:10.749066+00:00",
		InTimeline:       false,
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	id := "bcfbcff7-d844-40d0-8e87-50afae70a628"

	partner, err := c.GetPartner(id)

	if err != nil {
		t.Error("Got an error when trying to run CreatePartner.\n")
	}

	assert.Equal(t, partner.Id, "bcfbcff7-d844-40d0-8e87-50afae70a628", "partner.Id is incorrect")
}

func TestUpdatePartner(t *testing.T) {

	data := PartnerResponseDto{
		Id:               "bcfbcff7-d844-40d0-8e87-50afae70a628",
		Email:            "test447a@example.com",
		Name:             "test447a",
		ProfileImagePath: "",
		AvatarColor:      "red",
		ProfileChangedAt: "2026-07-09T05:27:10.749066+00:00",
		InTimeline:       true,
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	id := "bcfbcff7-d844-40d0-8e87-50afae70a628"

	partner, err := c.UpdatePartner(id, PartnerUpdateDto{
		InTimeline: true,
	})

	if err != nil {
		t.Error("Got an error when trying to run CreatePartner.\n")
	}

	assert.Equal(t, partner.Id, "bcfbcff7-d844-40d0-8e87-50afae70a628", "partner.Id is incorrect")
}

func TestDeletePartner(t *testing.T) {

	data := PartnerResponseDto{
		Id:               "bcfbcff7-d844-40d0-8e87-50afae70a628",
		Email:            "test447a@example.com",
		Name:             "test447a",
		ProfileImagePath: "",
		AvatarColor:      "red",
		ProfileChangedAt: "2026-07-09T05:27:10.749066+00:00",
		InTimeline:       false,
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	id := "bcfbcff7-d844-40d0-8e87-50afae70a628"

	err = c.DeletePartner(id)

	if err != nil {
		t.Error("Got an error when trying to run CreatePartner.\n")
	}
}
