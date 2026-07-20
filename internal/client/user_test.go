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

func TestCreateUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		user := User{
			ID:                   "c906aacc-fae0-4291-af97-6aef5b26eb43",
			Email:                "test4@example.com",
			Name:                 "Test4",
			ProfileImagePath:     "",
			AvatarColor:          "orange",
			ProfileChangedAt:     "2026-07-08T02:03:42.039Z",
			StorageLabel:         "",
			ShouldChangePassword: true,
			IsAdmin:              false,
			CreatedAt:            "2026-07-08T02:03:42.039Z",
			DeletedAt:            "",
			UpdatedAt:            "2026-07-08T02:03:42.039Z",
			OauthId:              "",
			QuotaSizeInBytes:     0,
			QuotaUsageInBytes:    0,
			Status:               "active",
		}

		data, err := json.Marshal(user)
		if err != nil {
			t.Error("Got an error when unmarshalling user.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

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
		t.Error("Got an error when trying to run CreateUser.\n")
	}

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	if newUser.ID != id {
		t.Errorf("Expected %v. Got %v\n", id, newUser.ID)
	}
}

func TestGetUsers(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetUsers should use the GET method")

		users := [1]User{{
			ID:                   id,
			Email:                "test4@example.com",
			Name:                 "Test4",
			ProfileImagePath:     "",
			AvatarColor:          "orange",
			ProfileChangedAt:     "2026-07-08T02:03:42.039Z",
			StorageLabel:         "",
			ShouldChangePassword: true,
			IsAdmin:              false,
			CreatedAt:            "2026-07-08T02:03:42.039Z",
			DeletedAt:            "",
			UpdatedAt:            "2026-07-08T02:03:42.039Z",
			OauthId:              "",
			QuotaSizeInBytes:     0,
			QuotaUsageInBytes:    0,
			Status:               "active",
		}}

		data, err := json.Marshal(users)
		if err != nil {
			t.Error("Got an error when unmarshalling user.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	users, err := c.GetUsers()
	if err != nil {
		t.Error("Got an error when trying to run GetUser.\n")
	}

	// XXX fix this
	for _, user := range users {
		assert.Equal(t, id, user.ID, "user.Id is incorrect")
	}
}

func TestGetUser(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetUser should use the GET method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		user := User{
			ID:                   id,
			Email:                "test4@example.com",
			Name:                 "Test4",
			ProfileImagePath:     "",
			AvatarColor:          "orange",
			ProfileChangedAt:     "2026-07-08T02:03:42.039Z",
			StorageLabel:         "",
			ShouldChangePassword: true,
			IsAdmin:              false,
			CreatedAt:            "2026-07-08T02:03:42.039Z",
			DeletedAt:            "",
			UpdatedAt:            "2026-07-08T02:03:42.039Z",
			OauthId:              "",
			QuotaSizeInBytes:     0,
			QuotaUsageInBytes:    0,
			Status:               "active",
		}

		data, err := json.Marshal(user)
		if err != nil {
			t.Error("Got an error when unmarshalling user.\n")
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

	user, err := c.GetUser(id)
	if err != nil {
		t.Error("Got an error when trying to run GetUser.\n")
	}

	assert.Equal(t, user.ID, id, "user.Id is incorrect")
}

func TestUpdateUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "PATCH", "UpdateUser should use the PATCH method")

		user := User{
			ID:                   "c906aacc-fae0-4291-af97-6aef5b26eb43",
			Email:                "test4@example.com",
			Name:                 "Test4",
			ProfileImagePath:     "",
			AvatarColor:          "orange",
			ProfileChangedAt:     "2026-07-08T02:03:42.039Z",
			StorageLabel:         "",
			ShouldChangePassword: true,
			IsAdmin:              false,
			CreatedAt:            "2026-07-08T02:03:42.039Z",
			DeletedAt:            "",
			UpdatedAt:            "2026-07-08T02:03:42.039Z",
			OauthId:              "",
			QuotaSizeInBytes:     0,
			QuotaUsageInBytes:    0,
			Status:               "active",
		}

		data, err := json.Marshal(user)
		if err != nil {
			t.Error("Got an error when unmarshalling user.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	user := UserUpdate{
		Name:     "Test4",
		Email:    "test4@example.com",
		Password: "xxxxxxxxxxxxxxxxxx",
	}

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	newUser, err := c.UpdateUser(id, user)
	if err != nil {
		t.Error("Got an error when trying to run UpdateUser.\n")
	}

	if newUser.ID != id {
		t.Errorf("Expected %v. Got %v\n", id, newUser.ID)
	}
}

func TestDeleteUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "DELETE", "DeleteUser should use the DELETE method")

		user := User{
			ID:                   "c906aacc-fae0-4291-af97-6aef5b26eb43",
			Email:                "test4@example.com",
			Name:                 "Test4",
			ProfileImagePath:     "",
			AvatarColor:          "orange",
			ProfileChangedAt:     "2026-07-08T02:03:42.039Z",
			StorageLabel:         "",
			ShouldChangePassword: true,
			IsAdmin:              false,
			CreatedAt:            "2026-07-08T02:03:42.039Z",
			DeletedAt:            "2026-07-08T02:03:42.039Z",
			UpdatedAt:            "2026-07-08T02:03:42.039Z",
			OauthId:              "",
			QuotaSizeInBytes:     0,
			QuotaUsageInBytes:    0,
			Status:               "removing",
		}

		data, err := json.Marshal(user)
		if err != nil {
			t.Error("Got an error when unmarshalling user.\n")
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

	err = c.DeleteUser(id)
	if err != nil {
		t.Error("Got an error when trying to run DeleteUser.\n")
	}
}
