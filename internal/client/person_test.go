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

func TestCreatePerson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		person := PersonResponseDto{
			Id:            "bcfbcff7-d844-40d0-8e87-50afae70a628",
			Name:          "test",
			Color:         "#443322",
			BirthDate:     "2026-07-09",
			IsFavorite:    false,
			IsHidden:      false,
			ThumbnailPath: "",
			UpdatedAt:     "2026-07-09T05:27:10.749066+00:00",
		}

		data, err := json.Marshal(person)
		if err != nil {
			t.Error("Got an error when unmarshalling person.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	person, err := c.CreatePerson(PersonCreateDto{
		Name:  "blah",
		Color: "#443322",
	})

	if err != nil {
		t.Error("Got an error when trying to run CreatePerson.\n")
	}

	assert.Equal(t, person.Id, "bcfbcff7-d844-40d0-8e87-50afae70a628", "person.Id is incorrect")
}

func TestGetPeople(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetPeople should use the GET method")

		person := PersonResponseDto{
			Id:            id,
			Name:          "test",
			Color:         "#443322",
			BirthDate:     "2026-07-09",
			IsFavorite:    false,
			IsHidden:      false,
			ThumbnailPath: "",
			UpdatedAt:     "2026-07-09T05:27:10.749066+00:00",
		}

		people_response := PeopleResponseDto{
			HasNextPage: false,
			Total:       1,
			Hidden:      0,
			People:      []PersonResponseDto{person},
		}

		data, err := json.Marshal(people_response)
		if err != nil {
			t.Error("Got an error when unmarshalling person.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Errorf("Got an error when trying to create a client: %v.\n", err)
	}

	people, err := c.GetPeople()
	if err != nil {
		t.Errorf("Got an error when trying to run GetPeople: %v.\n", err)
	}

	// XXX fix this
	for _, person := range people {
		assert.Equal(t, id, person.Id, "person.Id is incorrect")
		assert.Equal(t, "test", person.Name, "person.Name is incorrect")
	}
}

func TestGetPeopleBadData(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetPeople should use the GET method")

		data := "blah"

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	_, err = c.GetPeople()
	assert.NotEqual(t, nil, err, "err should not be nil.")
}

func TestGetPeopleServerFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "something went wrong", http.StatusInternalServerError)
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	_, err = c.GetPeople()

	assert.NotEqual(t, nil, err, "err should not be nil.")
}

func TestGetPerson(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetPerson should use the GET method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		person := PersonResponseDto{
			Id:            id,
			Name:          "test",
			Color:         "#443322",
			BirthDate:     "2026-07-09",
			IsFavorite:    false,
			IsHidden:      false,
			ThumbnailPath: "",
			UpdatedAt:     "2026-07-09T05:27:10.749066+00:00",
		}

		data, err := json.Marshal(person)
		if err != nil {
			t.Error("Got an error when unmarshalling person.\n")
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

	person, err := c.GetPerson(id)
	if err != nil {
		t.Error("Got an error when trying to run GetPerson.\n")
	}

	assert.Equal(t, person.Id, id, "person.Id is incorrect")
}

func TestUpdatePerson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "PATCH", "UpdatePerson should use the PATCH method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		person := PersonResponseDto{
			Id:            id,
			Name:          "test",
			Color:         "#443322",
			BirthDate:     "2026-07-09",
			IsFavorite:    false,
			IsHidden:      false,
			ThumbnailPath: "",
			UpdatedAt:     "2026-07-09T05:27:10.749066+00:00",
		}

		data, err := json.Marshal(person)
		if err != nil {
			t.Error("Got an error when unmarshalling person.\n")
		}

		rw.Write([]byte(data))
	}))

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	person := PersonUpdateDto{
		Name:  "test",
		Color: "#443322",
	}

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"

	newPerson, err := c.UpdatePerson(id, person)

	if err != nil {
		t.Error("Got an error when trying to run UpdatePerson.\n")
	}

	if newPerson.Id != id {
		t.Errorf("Expected %v. Got %v\n", id, newPerson.Id)
	}
}

func TestDeletePerson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "DELETE", "DeletePerson should use the DELETE method")

		i := strings.LastIndex(req.URL.Path, "/")
		id := req.URL.Path[i+1:]

		person := PersonResponseDto{
			Id:            id,
			Name:          "test",
			Color:         "#443322",
			BirthDate:     "2026-07-09",
			IsFavorite:    false,
			IsHidden:      false,
			ThumbnailPath: "",
			UpdatedAt:     "2026-07-09T05:27:10.749066+00:00",
		}

		data, err := json.Marshal(person)
		if err != nil {
			t.Error("Got an error when unmarshalling person.\n")
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

	err = c.DeletePerson(id)
	if err != nil {
		t.Error("Got an error when trying to run DeletePerson.\n")
	}
}
