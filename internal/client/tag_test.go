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

func TestCreateTag(t *testing.T) {

	data := TagResponseDto{
		Id:        "bcfbcff7-d844-40d0-8e87-50afae70a628",
		Name:      "test",
		Color:     "#443322",
		CreatedAt: "2026-07-09T05:27:10.749066+00:00",
		ParentId:  "",
		Value:     "",
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	tag, err := c.CreateTag(TagCreateDto{
		Name:  "blah",
		Color: "#443322",
	})

	if err != nil {
		t.Error("Got an error when trying to run CreateTag.\n")
	}

	assert.Equal(t, tag.Id, "bcfbcff7-d844-40d0-8e87-50afae70a628", "tag.Id is incorrect")
}

func TestGetTags(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetTags should use the GET method")

		tags := [1]TagResponseDto{{
			Id:        id,
			Name:      "test",
			Color:     "#443322",
			CreatedAt: "2026-07-09T05:27:10.749066+00:00",
			ParentId:  "",
			Value:     "",
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

	tags, err := c.GetTags()
	if err != nil {
		t.Error("Got an error when trying to run GetTag.\n")
	}

	// XXX fix this
	for _, tag := range tags {
		assert.Equal(t, id, tag.Id, "tag.Id is incorrect")
		assert.Equal(t, "test", tag.Name, "tag.Name is incorrect")
	}
}

func TestGetTagsBadData(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetTags should use the GET method")

		data := "blah"

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	_, err = c.GetTags()
	assert.NotEqual(t, nil, err, "err should not be nil.")
}

func TestGetTagsServerFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "something went wrong", http.StatusInternalServerError)
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	_, err = c.GetTags()

	assert.NotEqual(t, nil, err, "err should not be nil.")
}

func TestGetTag(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetTag should use the GET method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		tag := TagResponseDto{
			Id:        id,
			Name:      "test",
			Color:     "#443322",
			CreatedAt: "2026-07-09T05:27:10.749066+00:00",
			ParentId:  "",
			Value:     "",
		}

		data, err := json.Marshal(tag)
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

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	tag, err := c.GetTag(id)
	if err != nil {
		t.Error("Got an error when trying to run GetTag.\n")
	}

	assert.Equal(t, tag.Id, id, "tag.Id is incorrect")
}

func TestUpdateTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "PUT", "UpdateTag should use the PUT method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		tag := TagResponseDto{
			Id:        id,
			Name:      "test",
			Color:     "#443322",
			CreatedAt: "2026-07-09T05:27:10.749066+00:00",
			ParentId:  "",
			Value:     "",
		}

		data, err := json.Marshal(tag)
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

	tag := TagUpdateDto{
		Color: "#443322",
	}

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	newTag, err := c.UpdateTag(id, tag)

	if err != nil {
		t.Error("Got an error when trying to run UpdateTag.\n")
	}

	if newTag.Id != id {
		t.Errorf("Expected %v. Got %v\n", id, newTag.Id)
	}
}

func TestDeleteTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "DELETE", "DeleteTag should use the DELETE method")

		tag := TagResponseDto{
			Id:        "bcfbcff7-d844-40d0-8e87-50afae70a628",
			Name:      "test",
			Color:     "#443322",
			CreatedAt: "2026-07-09T05:27:10.749066+00:00",
			ParentId:  "",
			Value:     "",
		}

		data, err := json.Marshal(tag)
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

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	err = c.DeleteTag(id)
	if err != nil {
		t.Error("Got an error when trying to run DeleteTag.\n")
	}
}
