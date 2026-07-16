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

type ServerAbout struct {
	Build                      string `json:"build"`
	BuildImage                 string `json:"buildImage"`
	BuildImageUrl              string `json:"buildImageUrl"`
	BuildUrl                   string `json:"buildUrl"`
	Exiftool                   string `json:"exiftool"`
	Ffmpeg                     string `json:"ffmpeg"`
	Imagemagick                string `json:"imagemagick"`
	Libvips                    string `json:"libvips"`
	Licensed                   bool   `json:"licensed"`
	Nodejs                     string `json:"nodejs"`
	Repository                 string `json:"repository"`
	RepositoryUrl              string `json:"repositoryUrl"`
	SourceCommit               string `json:"sourceCommit"`
	SourceRef                  string `json:"sourceRef"`
	SourceUrl                  string `json:"sourceUrl"`
	ThirdPartyBugFeatureUrl    string `json:"thirdPartyBugFeatureUrl"`
	ThirdPartyDocumentationUrl string `json:"thirdPartyDocumentationUrl"`
	ThirdPartySourceUrl        string `json:"thirdPartySourceUrl"`
	ThirdPartySupportUrl       string `json:"thirdPartySupportUrl"`
	Version                    string `json:"version"`
	VersionUrl                 string `json:"versionUrl"`
}

type ServerConfig struct {
	ExternalDomain   string `json:"externalDomain"`
	IsInitialized    bool   `json:"isInitialized"`
	IsOnboarded      bool   `json:"isOnboarded"`
	LoginPageMessage string `json:"loginPageMessage"`
	MaintenanceMode  bool   `json:"maintenanceMode"`
	MapDarkStyleUrl  string `json:"mapDarkStyleUrl"`
	MapLightStyleUrl string `json:"mapLightStyleUrl"`
	MinFaces         int64  `json:"minFaces"`
	OauthButtonText  string `json:"oauthButtonText"`
	PublicUsers      bool   `json:"publicUsers"`
	TrashDays        int64  `json:"trashDays"`
	UserDeleteDelay  int64  `json:"userDeleteDelay"`
}

type ServerFeatures struct {
	ConfigFile          bool `json:"configFile"`
	DuplicateDetection  bool `json:"duplicateDetection"`
	Email               bool `json:"email"`
	FacialRecognition   bool `json:"facialRecognition"`
	ImportFaces         bool `json:"importFaces"`
	Map                 bool `json:"map"`
	Oauth               bool `json:"oauth"`
	OauthAutoLaunch     bool `json:"oauthAutoLaunch"`
	Ocr                 bool `json:"ocr"`
	PasswordLogin       bool `json:"passwordLogin"`
	RealtimeTranscoding bool `json:"realtimeTranscoding"`
	ReverseGeocoding    bool `json:"reverseGeocoding"`
	Search              bool `json:"search"`
	Sidecar             bool `json:"sidecar"`
	SmartSearch         bool `json:"smartSearch"`
	Trash               bool `json:"trash"`
}

type ServerStorage struct {
	DiskAvailable       string  `json:"diskAvailable"`
	DiskAvailableRaw    int64   `json:"diskAvailableRaw"`
	DiskSize            string  `json:"diskSize"`
	DiskSizeRaw         int64   `json:"diskSizeRaw"`
	DiskUsagePercentage float64 `json:"diskUsagePercentage"`
	DiskUse             string  `json:"diskUse"`
	DiskUseRaw          int64   `json:"diskUseRaw"`
}

func (c *Client) GetServerAbout() (*ServerAbout, error) {
	ServerAbout, err := get[ServerAbout](c, "server/about")

	return ServerAbout, err
}

func (c *Client) GetServerConfig() (*ServerConfig, error) {
	ServerConfig, err := get[ServerConfig](c, "server/config")

	return ServerConfig, err
}

func (c *Client) GetServerFeatures() (*ServerFeatures, error) {
	ServerFeatures, err := get[ServerFeatures](c, "server/features")

	return ServerFeatures, err
}

func (c *Client) GetServerStorage() (*ServerStorage, error) {
	ServerStorage, err := get[ServerStorage](c, "server/storage")

	return ServerStorage, err
}
