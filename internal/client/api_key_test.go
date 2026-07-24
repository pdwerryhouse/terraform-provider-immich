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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateApiKey(t *testing.T) {

	data := ApiKeyCreateResponseDto{
		ApiKey: ApiKeyResponseDto{
			Id:          "bcfbcff7-d844-40d0-8e87-50afae70a628",
			Name:        "test",
			Permissions: []string{"all"},
			CreatedAt:   "2026-07-09T05:27:10.749066+00:00",
			UpdatedAt:   "2026-07-09T05:27:10.749066+00:00",
		},
		Secret: "fdsmsdffdsksfdkfdskldfskmdfsmkdfsmklfdskl",
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	apiKey, err := c.CreateApiKey(ApiKeyCreateDto{
		Name:        "test",
		Permissions: []string{"all"},
	})

	if err != nil {
		t.Error("Got an error when trying to run CreateApiKey.\n")
	}

	assert.Equal(t, apiKey.ApiKey.Id, "bcfbcff7-d844-40d0-8e87-50afae70a628", "apiKey.ApiKey.Id is incorrect")
}

func TestGetApiKeys(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetApiKeys should use the GET method")

		apiKeys := [1]ApiKeyResponseDto{{
			Id:          id,
			Name:        "test",
			Permissions: []string{"all"},
			CreatedAt:   "2026-07-09T05:27:10.749066+00:00",
			UpdatedAt:   "2026-07-09T05:27:10.749066+00:00",
		}}

		data, err := json.Marshal(apiKeys)
		if err != nil {
			t.Error("Got an error when unmarshalling apiKey.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	apiKeys, err := c.GetApiKeys()
	if err != nil {
		t.Error("Got an error when trying to run GetApiKey.\n")
	}

	// XXX fix this
	for _, apiKey := range apiKeys {
		assert.Equal(t, id, apiKey.Id, "apiKey.Id is incorrect")
		assert.Equal(t, "test", apiKey.Name, "apiKey.Name is incorrect")
	}
}

func TestGetApiKeysBadData(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetApiKeys should use the GET method")

		data := "blah"

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	_, err = c.GetApiKeys()
	assert.NotEqual(t, nil, err, "err should not be nil.")
}

func TestGetApiKeysServerFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "something went wrong", http.StatusInternalServerError)
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	_, err = c.GetApiKeys()

	assert.NotEqual(t, nil, err, "err should not be nil.")
}

func TestGetApiKey(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetApiKey should use the GET method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		apiKey := ApiKeyResponseDto{
			Id:          id,
			Name:        "test",
			Permissions: []string{"all"},
			CreatedAt:   "2026-07-09T05:27:10.749066+00:00",
			UpdatedAt:   "2026-07-09T05:27:10.749066+00:00",
		}

		data, err := json.Marshal(apiKey)
		if err != nil {
			t.Error("Got an error when unmarshalling apiKey.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	apiKey, err := c.GetApiKey(id)
	if err != nil {
		t.Error("Got an error when trying to run GetApiKey.\n")
	}

	assert.Equal(t, apiKey.Id, id, "apiKey.Id is incorrect")
}

func TestUpdateApiKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "PUT", "UpdateApiKey should use the PUT method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		apiKey := ApiKeyResponseDto{
			Id:          id,
			Name:        "test",
			Permissions: []string{"all"},
			CreatedAt:   "2026-07-09T05:27:10.749066+00:00",
			UpdatedAt:   "2026-07-09T05:27:10.749066+00:00",
		}

		data, err := json.Marshal(apiKey)
		if err != nil {
			t.Error("Got an error when unmarshalling apiKey.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	apiKey := ApiKeyUpdateDto{
		Name:        "test",
		Permissions: []string{"all"},
	}

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	newApiKey, err := c.UpdateApiKey(id, apiKey)

	if err != nil {
		t.Error("Got an error when trying to run UpdateApiKey.\n")
	}

	if newApiKey.Id != id {
		t.Errorf("Expected %v. Got %v\n", id, newApiKey.Id)
	}
}

func TestDeleteApiKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "DELETE", "DeleteApiKey should use the DELETE method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		apiKey := ApiKeyResponseDto{
			Id:          id,
			Name:        "test",
			Permissions: []string{"all"},
			CreatedAt:   "2026-07-09T05:27:10.749066+00:00",
			UpdatedAt:   "2026-07-09T05:27:10.749066+00:00",
		}

		data, err := json.Marshal(apiKey)
		if err != nil {
			t.Error("Got an error when unmarshalling apiKey.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	err = c.DeleteApiKey(id)
	if err != nil {
		t.Error("Got an error when trying to run DeleteApiKey.\n")
	}
}
