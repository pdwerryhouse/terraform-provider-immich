// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAlbum(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{
  "albumName": "Test",
  "description": "A Test Album",
  "albumThumbnailAssetId": null,
  "createdAt": "2026-07-08T01:46:57.450Z",
  "updatedAt": "2026-07-08T01:46:57.450Z",
  "id": "fe2217e0-b071-4e32-8440-c499d79a469b",
  "ownerId": "ebd580a9-f019-4bbc-b557-900377becfd0",
  "owner": {
    "id": "ebd580a9-f019-4bbc-b557-900377becfd0",
    "email": "paul@dwerryhouse.com.au",
    "name": "Paul Dwerryhouse",
    "profileImagePath": "",
    "avatarColor": "primary",
    "profileChangedAt": "2026-07-02T06:11:20.062077+00:00"
  },
  "albumUsers": [],
  "shared": false,
  "hasSharedLink": false,
  "assets": [],
  "assetCount": 0,
  "isActivityEnabled": true,
  "order": "desc"
}
`))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"
	apikey := "xxx"

	c, err := NewClient(&endpoint, &apikey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	album, err := c.CreateAlbum(AlbumUpdate{
		AlbumName:   "Test",
		Description: "A Test Album",
	})

	if err != nil {
		t.Error("Got an error when trying to run CreateAlbum.\n")
	}

	if album.ID != "fe2217e0-b071-4e32-8440-c499d79a469b" {
		t.Errorf("Expected %v. Got %v\n", "fe2217e0-b071-4e32-8440-c499d79a469b", album.ID)
	}
}

func TestGetAlbums(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//equals(t, req.URL.String(), "/albums/feb9ee33-b35a-4a20-952e-04d7969d7006")
		rw.Write([]byte(`{
  "albumName": "Test1",
  "description": "Test1 Album",
  "albumThumbnailAssetId": null,
  "createdAt": "2026-07-02T06:17:05.708Z",
  "updatedAt": "2026-07-02T06:17:15.613Z",
  "id": "feb9ee33-b35a-4a20-952e-04d7969d7006",
  "ownerId": "ebd580a9-f019-4bbc-b557-900377becfd0",
  "owner": {
    "id": "ebd580a9-f019-4bbc-b557-900377becfd0",
    "email": "test@example.com",
    "name": "Test User",
    "profileImagePath": "",
    "avatarColor": "primary",
    "profileChangedAt": "2026-07-02T06:11:20.062077+00:00"
  },
  "albumUsers": [],
  "shared": false,
  "hasSharedLink": false,
  "assets": [],
  "assetCount": 0,
  "isActivityEnabled": true,
  "order": "desc"
}
`))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"
	apikey := "xxx"

	c, err := NewClient(&endpoint, &apikey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	album, err := c.GetAlbum("feb9ee33-b35a-4a20-952e-04d7969d7006")
	if err != nil {
		t.Error("Got an error when trying to run GetAlbum.\n")
	}

	if album.ID != "feb9ee33-b35a-4a20-952e-04d7969d7006" {
		t.Errorf("Expected %v. Got %v\n", "feb9ee33-b35a-4a20-952e-04d7969d7006", album.ID)
	}
}
