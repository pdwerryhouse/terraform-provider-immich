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

func TestCreateFace(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	data := AssetFaceResponseDto{
		Id:            id,
		BoundingBoxX1: 10,
		BoundingBoxX2: 10,
		BoundingBoxY1: 10,
		BoundingBoxY2: 10,
		ImageHeight:   100,
		ImageWidth:    100,
		SourceType:    "manual",
		Person: PersonResponseDto{
			Id:            "bcfbcff7-d844-40d0-8e87-50afae70a628",
			Name:          "test",
			Color:         "#443322",
			BirthDate:     "2026-07-09",
			IsFavorite:    false,
			IsHidden:      false,
			ThumbnailPath: "",
			UpdatedAt:     "2026-07-09T05:27:10.749066+00:00",
		},
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	face, err := c.CreateFace(AssetFaceCreateDto{
		Height:      100,
		Width:       100,
		ImageHeight: 100,
		ImageWidth:  100,
		X:           100,
		Y:           100,
		PersonId:    "bcfbcff7-d844-40d0-8e87-50afae70a628",
	})

	if err != nil {
		t.Error("Got an error when trying to run CreateFace.\n")
	}

	assert.Equal(t, id, face.Id, "face.Id is incorrect")
}

func TestGetFaces(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetFaces should use the GET method")

		faces := [1]AssetFaceResponseDto{{
			Id:            id,
			BoundingBoxX1: 10,
			BoundingBoxX2: 10,
			BoundingBoxY1: 10,
			BoundingBoxY2: 10,
			ImageHeight:   100,
			ImageWidth:    100,
			SourceType:    "manual",
			Person: PersonResponseDto{
				Id:            "bcfbcff7-d844-40d0-8e87-50afae70a628",
				Name:          "test",
				Color:         "#443322",
				BirthDate:     "2026-07-09",
				IsFavorite:    false,
				IsHidden:      false,
				ThumbnailPath: "",
				UpdatedAt:     "2026-07-09T05:27:10.749066+00:00",
			},
		}}

		data, err := json.Marshal(faces)
		if err != nil {
			t.Error("Got an error when unmarshalling face.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	faces, err := c.GetFaces(id)
	if err != nil {
		t.Error("Got an error when trying to run GetFace.\n")
	}

	// XXX fix this
	for _, face := range faces {
		assert.Equal(t, id, face.Id, "face.Id is incorrect")
	}
}

func TestGetFace(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetFace should use the GET method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		face := AssetFaceResponseDto{
			Id:            id,
			BoundingBoxX1: 10,
			BoundingBoxX2: 10,
			BoundingBoxY1: 10,
			BoundingBoxY2: 10,
			ImageHeight:   100,
			ImageWidth:    100,
			SourceType:    "manual",
			Person: PersonResponseDto{
				Id:            "bcfbcff7-d844-40d0-8e87-50afae70a628",
				Name:          "test",
				Color:         "#443322",
				BirthDate:     "2026-07-09",
				IsFavorite:    false,
				IsHidden:      false,
				ThumbnailPath: "",
				UpdatedAt:     "2026-07-09T05:27:10.749066+00:00",
			},
		}

		data, err := json.Marshal(face)
		if err != nil {
			t.Error("Got an error when unmarshalling face.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	face, err := c.GetFace(id)
	if err != nil {
		t.Error("Got an error when trying to run GetFace.\n")
	}

	assert.Equal(t, face.Id, id, "face.Id is incorrect")
}

func TestUpdateFace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "PUT", "UpdateFace should use the PUT method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		face := FaceDto{
			Id: id,
		}

		data, err := json.Marshal(face)
		if err != nil {
			t.Error("Got an error when unmarshalling face.\n")
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

	face := FaceDto{
		Id: id,
	}

	newFace, err := c.UpdateFace(id, face)

	if err != nil {
		t.Error("Got an error when trying to run UpdateFace.\n")
	}

	assert.Equal(t, id, newFace.Id, "Incorrect newFace.Id")
}

func TestDeleteFace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "DELETE", "DeleteFace should use the DELETE method")

		response := ""

		data, err := json.Marshal(response)
		if err != nil {
			t.Error("Got an error when unmarshalling face.\n")
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

	err = c.DeleteFace(id)

	assert.Equal(t, err, nil, "DeleteFace returned an error")
}
