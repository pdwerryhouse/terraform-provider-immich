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

func TestCreateAlbum(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		album := AlbumResponseDto{
			AlbumName:             "Test",
			Description:           "A Test Album",
			AlbumThumbnailAssetId: "",
			CreatedAt:             "2026-07-08T01:46:57.450Z",
			UpdatedAt:             "2026-07-08T01:46:57.450Z",
			Id:                    "fe2217e0-b071-4e32-8440-c499d79a469b",
			Shared:                false,
			HasSharedLink:         false,
			IsActivityEnabled:     true,
			Order:                 "desc",
		}

		data, err := json.Marshal(album)
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

	album, err := c.CreateAlbum(CreateAlbumDto{
		AlbumName:   "Test",
		Description: "A Test Album",
	})

	if err != nil {
		t.Error("Got an error when trying to run CreateAlbum.\n")
	}

	assert.Equal(t, album.Id, "fe2217e0-b071-4e32-8440-c499d79a469b", "album.ID is incorrect")
}

func TestGetAlbums(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetTags should use the GET method")

		albums := [1]AlbumResponseDto{{
			AlbumName:             "Test",
			Description:           "A Test Album",
			AlbumThumbnailAssetId: "",
			CreatedAt:             "2026-07-08T01:46:57.450Z",
			UpdatedAt:             "2026-07-08T01:46:57.450Z",
			Id:                    id,
			Shared:                false,
			HasSharedLink:         false,
			IsActivityEnabled:     true,
			Order:                 "desc",
		}}

		data, err := json.Marshal(albums)
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

	albums, err := c.GetAlbums()
	if err != nil {
		t.Error("Got an error when trying to run GetTag.\n")
	}

	// XXX fix this
	for _, album := range albums {
		assert.Equal(t, id, album.Id, "album.Id is incorrect")
	}
}

func TestGetAlbum(t *testing.T) {

	id := "fe2217e0-b071-4e32-8440-c499d79a469b"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//equals(t, req.URL.String(), "/albums/feb9ee33-b35a-4a20-952e-04d7969d7006")

		album := AlbumResponseDto{
			AlbumName:             "Test",
			Description:           "A Test Album",
			AlbumThumbnailAssetId: "",
			CreatedAt:             "2026-07-08T01:46:57.450Z",
			UpdatedAt:             "2026-07-08T01:46:57.450Z",
			Id:                    id,
			Shared:                false,
			HasSharedLink:         false,
			IsActivityEnabled:     true,
			Order:                 "desc",
		}

		data, err := json.Marshal(album)
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

	album, err := c.GetAlbum("feb9ee33-b35a-4a20-952e-04d7969d7006")
	if err != nil {
		t.Error("Got an error when trying to run GetAlbum.\n")
	}

	assert.Equal(t, id, album.Id, "album.ID is incorrect")
}

func TestUpdateAlbum(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "PATCH", "UpdateAlbum should use the PUT method")

		/*
			i := strings.LastIndex(req.URL.Path, "/")
			id := req.URL.Path[i+1:]
		*/

		album := AlbumResponseDto{
			AlbumName:             "Test",
			Description:           "A Test Album",
			AlbumThumbnailAssetId: "",
			CreatedAt:             "2026-07-08T01:46:57.450Z",
			UpdatedAt:             "2026-07-08T01:46:57.450Z",
			Id:                    "fe2217e0-b071-4e32-8440-c499d79a469b",
			Shared:                false,
			HasSharedLink:         false,
			IsActivityEnabled:     true,
			Order:                 "desc",
		}

		data, err := json.Marshal(album)
		if err != nil {
			t.Error("Got an error when unmarshalling album.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	album := UpdateAlbumDto{
		AlbumName:   "Test",
		Description: "A Test Album",
	}

	id := "fe2217e0-b071-4e32-8440-c499d79a469b"

	newAlbum, err := c.UpdateAlbum(id, album)

	if err != nil {
		t.Error("Got an error when trying to run UpdateAlbum.\n")
	}

	assert.Equal(t, id, newAlbum.Id, "Incorrect newAlbum.Id")
}

func TestDeleteAlbum(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "DELETE", "DeleteAlbum should use the DELETE method")

		album := AlbumResponseDto{
			AlbumName:             "Test",
			Description:           "A Test Album",
			AlbumThumbnailAssetId: "",
			CreatedAt:             "2026-07-08T01:46:57.450Z",
			UpdatedAt:             "2026-07-08T01:46:57.450Z",
			Id:                    "fe2217e0-b071-4e32-8440-c499d79a469b",
			Shared:                false,
			HasSharedLink:         false,
			IsActivityEnabled:     true,
			Order:                 "desc",
		}

		data, err := json.Marshal(album)
		if err != nil {
			t.Error("Got an error when unmarshalling album.\n")
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

	err = c.DeleteAlbum(id)

	assert.Equal(t, err, nil, "DeleteAlbum returned an error")
}
