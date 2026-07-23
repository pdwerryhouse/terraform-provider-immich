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

//import (
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestUpdateAlbumActivity(t *testing.T) {
//	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
//
//		assert.Equal(t, req.Method, "PATCH", "UpdateAlbumActivity should use the PATCH method")
//
//		/*
//			i := strings.LastIndex(req.URL.Path, "/")
//			id := req.URL.Path[i+1:]
//		*/
//
//		album := Album{
//			AlbumName:             "Test",
//			Description:           "A Test Album",
//			AlbumThumbnailAssetId: "",
//			CreatedAt:             "2026-07-08T01:46:57.450Z",
//			UpdatedAt:             "2026-07-08T01:46:57.450Z",
//			ID:                    "fe2217e0-b071-4e32-8440-c499d79a469b",
//			OwnerId:               "ebd580a9-f019-4bbc-b557-900377becfd0",
//			Owner: Owner{
//				ID:               "ebd580a9-f019-4bbc-b557-900377becfd0",
//				Email:            "test001@example.com",
//				Name:             "Test 001",
//				ProfileImagePath: "",
//				AvatarColor:      "primary",
//				ProfileChangedAt: "2026-07-02T06:11:20.062077+00:00",
//			},
//			Shared:            false,
//			HasSharedLink:     false,
//			IsActivityEnabled: true,
//			Order:             "desc",
//		}
//
//		data, err := json.Marshal(album)
//		if err != nil {
//			t.Error("Got an error when unmarshalling album.\n")
//		}
//
//		rw.Write([]byte(data))
//	}))
//
//	defer server.Close()
//
//	endpoint := server.URL + "/api"
//
//	c, err := NewClient(&endpoint, &testApiKey)
//
//	if err != nil {
//		t.Error("Got an error when trying to create a client.\n")
//	}
//
//	album := AlbumActivityUpdate{
//		IsActivityEnabled: true,
//	}
//
//	id := "fe2217e0-b071-4e32-8440-c499d79a469b"
//
//	newAlbum, err := c.UpdateAlbumActivity(id, album)
//
//	if err != nil {
//		t.Error("Got an error when trying to run UpdateAlbumActivity.\n")
//	}
//
//	assert.Equal(t, id, newAlbum.ID, "Incorrect newAlbum.Id")
//}
//
