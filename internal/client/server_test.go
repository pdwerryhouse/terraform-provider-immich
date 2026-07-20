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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerAbout(t *testing.T) {

	data := ServerAbout{
		Version:       "v3.0.1",
		VersionUrl:    "https://github.com/immich-app/immich/releases/tag/v3.0.1",
		Licensed:      false,
		Build:         "28622919952",
		BuildUrl:      "https://github.com/immich-app/immich/actions/runs/28622919952",
		BuildImage:    "v3.0.1",
		BuildImageUrl: "https://github.com/immich-app/immich/pkgs/container/immich-server",
		Repository:    "immich-app/immich",
		RepositoryUrl: "https://github.com/immich-app/immich",
		SourceRef:     "v3.0.1",
		SourceCommit:  "f77c8a4699aec043ec4a1b404559a7a291d592dc",
		SourceUrl:     "https://github.com/immich-app/immich/commit/f77c8a4699aec043ec4a1b404559a7a291d592dc",
		Nodejs:        "v24.14.1",
		Exiftool:      "13.59",
		Ffmpeg:        "7.1.4-3",
		Libvips:       "8.18.2",
		Imagemagick:   "7.1.2-21",
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	server_about, err := c.GetServerAbout()
	if err != nil {
		t.Error("Got an error when trying to run GetServer.\n")
	}

	assert.Equal(t, "v3.0.1", server_about.Version, "GetServerAbout not returning the correct values.")
}

func TestGetServerConfig(t *testing.T) {

	data := ServerConfig{
		LoginPageMessage: "",
		TrashDays:        30,
		UserDeleteDelay:  7,
		OauthButtonText:  "Login with OAuth",
		IsInitialized:    true,
		IsOnboarded:      true,
		ExternalDomain:   "",
		PublicUsers:      true,
		MapDarkStyleUrl:  "https://tiles.immich.cloud/v1/style/dark.json",
		MapLightStyleUrl: "https://tiles.immich.cloud/v1/style/light.json",
		MaintenanceMode:  false,
		MinFaces:         3,
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	server_config, err := c.GetServerConfig()
	if err != nil {
		t.Error("Got an error when trying to run GetServer.\n")
	}

	assert.Equal(t, "Login with OAuth", server_config.OauthButtonText, "GetServerConfig not returning the correct values.")
}

func TestGetServerFeatures(t *testing.T) {

	data := ServerFeatures{
		SmartSearch:         true,
		FacialRecognition:   true,
		DuplicateDetection:  true,
		Map:                 true,
		ReverseGeocoding:    true,
		ImportFaces:         false,
		Sidecar:             true,
		Search:              true,
		Trash:               true,
		Oauth:               false,
		OauthAutoLaunch:     false,
		Ocr:                 true,
		PasswordLogin:       true,
		ConfigFile:          false,
		Email:               false,
		RealtimeTranscoding: false,
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	server_features, err := c.GetServerFeatures()
	if err != nil {
		t.Error("Got an error when trying to run GetServer.\n")
	}

	assert.Equal(t, true, server_features.SmartSearch, "GetServerFeatures not returning the correct values.")
}

func TestGetServerStorage(t *testing.T) {

	data := ServerStorage{
		DiskAvailable:       "183.7 GiB",
		DiskSize:            "617.2 GiB",
		DiskUse:             "424.2 GiB",
		DiskAvailableRaw:    197298446336,
		DiskSizeRaw:         662684303360,
		DiskUseRaw:          455459086336,
		DiskUsagePercentage: 68.73,
	}

	server := createTestServer(t, data)

	defer server.Close()

	endpoint := server.URL + "/api"

	c, err := NewClient(&endpoint, &testApiKey)

	if err != nil {
		t.Error("Got an error when trying to create a client.\n")
	}

	server_storage, err := c.GetServerStorage()
	if err != nil {
		t.Error("Got an error when trying to run GetServer.\n")
	}

	assert.Equal(t, "183.7 GiB", server_storage.DiskAvailable, "GetServerFeatures not returning the correct values.")
}
