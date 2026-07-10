// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

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
