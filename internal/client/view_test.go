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

func TestGetAssetsByOriginalPath(t *testing.T) {

	id := "c906aacc-fae0-4291-af97-6aef5b26eb43"
	path := "/data/upload/ebd580a9-f019-4bbc-b557-900377becfd0/4c/68"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Method, "GET", "GetFaces should use the GET method")

		assets := [1]AssetResponseDto{{
			Id:       id,
			Checksum: "6sl5aqLbtsFkNMGl1TxYNqbCcNg=",
		}}

		data, err := json.Marshal(assets)
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

	assets, err := c.GetAssetsByOriginalPath(path)
	if err != nil {
		t.Error("Got an error when trying to run GetFace.\n")
	}

	// XXX fix this
	for _, asset := range assets {
		assert.Equal(t, id, asset.Id, "asset.Id is incorrect")
	}
}
