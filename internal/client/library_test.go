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

func TestCreateLibrary(t *testing.T) {

	data := LibraryResponseDto{
		Id:                "b097dd49-8328-4331-8402-1f6144fdb291",
		OwnerId:           "ebd580a9-f019-4bbc-b557-900377becfd0",
		Name:              "New External Library",
		CreatedAt:         "2026-07-14T08:16:20.132Z",
		UpdatedAt:         "2026-07-14T08:16:20.132Z",
		RefreshedAt:       "",
		AssetCount:        0,
		ImportPaths:       nil,
		ExclusionPatterns: nil,
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	library, err := c.CreateLibrary(CreateLibraryDto{
		Name:              "blah",
		ImportPaths:       nil,
		ExclusionPatterns: nil,
		OwnerId:           "ebd580a9-f019-4bbc-b557-900377becfd0",
	})

	if err != nil {
		t.Error("Got an error when trying to run CreateLibrary.\n")
	}

	assert.Equal(t, library.Id, "b097dd49-8328-4331-8402-1f6144fdb291", "library.Id is incorrect")
}

func TestGetLibraries(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetLibrarys should use the GET method")

		librarys := [1]LibraryResponseDto{{
			Id:                id,
			OwnerId:           "ebd580a9-f019-4bbc-b557-900377becfd0",
			Name:              "New External Library",
			CreatedAt:         "2026-07-14T08:16:20.132Z",
			UpdatedAt:         "2026-07-14T08:16:20.132Z",
			RefreshedAt:       "",
			AssetCount:        0,
			ImportPaths:       nil,
			ExclusionPatterns: nil,
		}}

		data, err := json.Marshal(librarys)
		if err != nil {
			t.Error("Got an error when unmarshalling library.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	librarys, err := c.GetLibraries()
	if err != nil {
		t.Error("Got an error when trying to run GetLibrary.\n")
	}

	// XXX fix this
	for _, library := range librarys {
		assert.Equal(t, id, library.Id, "library.Id is incorrect")
	}
}

func TestGetLibrary(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetLibrary should use the GET method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		library := LibraryResponseDto{
			Id:                id,
			OwnerId:           "ebd580a9-f019-4bbc-b557-900377becfd0",
			Name:              "New External Library",
			CreatedAt:         "2026-07-14T08:16:20.132Z",
			UpdatedAt:         "2026-07-14T08:16:20.132Z",
			RefreshedAt:       "",
			AssetCount:        0,
			ImportPaths:       nil,
			ExclusionPatterns: nil,
		}

		data, err := json.Marshal(library)
		if err != nil {
			t.Error("Got an error when unmarshalling library.\n")
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

	library, err := c.GetLibrary(id)
	if err != nil {
		t.Error("Got an error when trying to run GetLibrary.\n")
	}

	assert.Equal(t, library.Id, id, "library.Id is incorrect")
}

func TestUpdateLibrary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "PATCH", "UpdateLibrary should use the PUT method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		library := LibraryResponseDto{
			Id:                id,
			OwnerId:           "ebd580a9-f019-4bbc-b557-900377becfd0",
			Name:              "New External Library",
			CreatedAt:         "2026-07-14T08:16:20.132Z",
			UpdatedAt:         "2026-07-14T08:16:20.132Z",
			RefreshedAt:       "",
			AssetCount:        0,
			ImportPaths:       nil,
			ExclusionPatterns: nil,
		}

		data, err := json.Marshal(library)
		if err != nil {
			t.Error("Got an error when unmarshalling library.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	library := UpdateLibraryDto{
		Name:              "blah",
		ImportPaths:       nil,
		ExclusionPatterns: nil,
	}

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	newLibrary, err := c.UpdateLibrary(id, library)

	if err != nil {
		t.Error("Got an error when trying to run UpdateLibrary.\n")
	}

	assert.Equal(t, id, newLibrary.Id, "Incorrect newLibrary.Id")
}

func TestDeleteLibrary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "DELETE", "DeleteLibrary should use the DELETE method")

		library := LibraryResponseDto{
			Id:                "b097dd49-8328-4331-8402-1f6144fdb291",
			OwnerId:           "ebd580a9-f019-4bbc-b557-900377becfd0",
			Name:              "New External Library",
			CreatedAt:         "2026-07-14T08:16:20.132Z",
			UpdatedAt:         "2026-07-14T08:16:20.132Z",
			RefreshedAt:       "",
			AssetCount:        0,
			ImportPaths:       nil,
			ExclusionPatterns: nil,
		}

		data, err := json.Marshal(library)
		if err != nil {
			t.Error("Got an error when unmarshalling library.\n")
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

	err = c.DeleteLibrary(id)

	assert.Equal(t, err, nil, "DeleteLibrary returned an error")
}
