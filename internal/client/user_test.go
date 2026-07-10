// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{
  "id": "c906aacc-fae0-4291-af97-6aef5b26eb43",
  "email": "test4@example.com",
  "name": "Test4",
  "profileImagePath": "",
  "avatarColor": "orange",
  "profileChangedAt": "2026-07-08T02:03:42.039Z",
  "storageLabel": null,
  "shouldChangePassword": true,
  "isAdmin": false,
  "createdAt": "2026-07-08T02:03:42.039Z",
  "deletedAt": null,
  "updatedAt": "2026-07-08T02:03:42.039Z",
  "oauthId": "",
  "quotaSizeInBytes": null,
  "quotaUsageInBytes": 0,
  "status": "active",
  "license": null
}`))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"
	apikey := "xxx"

	c, err := NewClient(&endpoint, &apikey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	user := UserUpdate{
		Name:     "Test4",
		Email:    "test4@example.com",
		Password: "xxxxxxxxxxxxxxxxxx",
	}

	newUser, err := c.CreateUser(user)
	if err != nil {
		t.Error("Got an error when trying to run CreateAlbum.\n")
	}

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	if newUser.ID != id {
		t.Errorf("Expected %v. Got %v\n", id, newUser.ID)
	}
}

func TestGetUser(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{
  "id": "3c8157f6-c60d-4be3-9b81-3d2a2d7037ce",
  "email": "test@dwerryhouse.com.au",
  "name": "Test",
  "profileImagePath": "",
  "avatarColor": "blue",
  "profileChangedAt": "2026-07-04T05:58:00.660Z",
  "storageLabel": null,
  "shouldChangePassword": false,
  "isAdmin": false,
  "createdAt": "2026-07-04T05:58:00.660Z",
  "deletedAt": null,
  "updatedAt": "2026-07-06T07:07:32.396Z",
  "oauthId": "",
  "quotaSizeInBytes": null,
  "quotaUsageInBytes": 0,
  "status": "active",
  "license": null
}`))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"
	apikey := "xxx"

	c, err := NewClient(&endpoint, &apikey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	id := "3c8157f6-c60d-4be3-9b81-3d2a2d7037ce"

	user, err := c.GetUser(id)
	if err != nil {
		t.Error("Got an error when trying to run GetAlbum.\n")
	}

	if user.ID != id {
		t.Errorf("Expected %v. Got %v\n", id, user.ID)
	}
}
