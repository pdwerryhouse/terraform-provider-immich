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

func TestGetMapMarkers(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetMapMarkers should use the GET method")

		tags := [1]MapMarkerResponseDto{{
			Id:      id,
			Lat:     55.932804,
			Lon:     -4.248018,
			City:    "Torrance",
			State:   "Scotland",
			Country: "United Kingdom",
		}}

		data, err := json.Marshal(tags)
		if err != nil {
			t.Error("Got an error when unmarshalling tag.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	mapMarkers, err := c.GetMapMarkers()
	if err != nil {
		t.Error("Got an error when trying to run GetMapMarker.\n")
	}

	// XXX fix this
	for _, mapMarker := range mapMarkers {
		assert.Equal(t, id, mapMarker.Id, "mapMarker.Id is incorrect")
		assert.Equal(t, "Scotland", mapMarker.State, "mapMarker.County is incorrect")
	}
}

func TestGetReverseGeocode(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetMapMarkers should use the GET method")

		tags := [1]MapReverseGeocodeResponseDto{{
			Country: "Australia",
			State:   "Victoria",
			City:    "Brunswick",
		}}

		data, err := json.Marshal(tags)
		if err != nil {
			t.Error("Got an error when unmarshalling tag.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	responses, err := c.GetReverseGeocode(-37.76667, 144.9628)
	if err != nil {
		t.Error("Got an error when trying to run GetReverseGeocode.\n")
	}

	// XXX fix this
	for _, response := range responses {
		assert.Equal(t, "Brunswick", response.City, "response.City is incorrect")
	}
}
