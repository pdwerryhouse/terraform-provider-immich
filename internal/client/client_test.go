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

import "testing"

func TestNewClient(t *testing.T) {
	endpoint := "http://localhost:2283"
	apikey := "xxx"

	c, err := NewClient(&endpoint, &apikey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	if c.Endpoint != endpoint {
		t.Errorf("Expected %v. Got %v\n", endpoint, c.Endpoint)
	}

	if c.ApiKey != apikey {
		t.Errorf("Expected %v. Got %v\n", apikey, c.ApiKey)
	}
}
