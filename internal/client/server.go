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

type ServerAboutResponseDto struct {
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

type ServerConfigDto struct {
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

type ServerFeaturesDto struct {
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

type ServerStorageResponseDto struct {
	DiskAvailable       string `json:"diskAvailable"`
	DiskAvailableRaw    int64  `json:"diskAvailableRaw"`
	DiskSize            string `json:"diskSize"`
	DiskSizeRaw         int64  `json:"diskSizeRaw"`
	DiskUsagePercentage int64  `json:"diskUsagePercentage"`
	DiskUse             string `json:"diskUse"`
	DiskUseRaw          int64  `json:"diskUseRaw"`
}

type UserLicense struct {
	ActivatedAt   string `json:"activatedAt"`
	ActivationKey string `json:"activationKey"`
	LicenseKey    string `json:"licenseKey"`
}

func (c *Client) GetServerAbout() (*ServerAboutResponseDto, error) {
	ServerAbout, err := get[ServerAboutResponseDto](c, "server/about")

	return ServerAbout, err
}

func (c *Client) GetServerConfig() (*ServerConfigDto, error) {
	ServerConfig, err := get[ServerConfigDto](c, "server/config")

	return ServerConfig, err
}

func (c *Client) GetServerFeatures() (*ServerFeaturesDto, error) {
	ServerFeatures, err := get[ServerFeaturesDto](c, "server/features")

	return ServerFeatures, err
}

func (c *Client) GetServerStorage() (*ServerStorageResponseDto, error) {
	ServerStorage, err := get[ServerStorageResponseDto](c, "server/storage")

	return ServerStorage, err
}

func (c *Client) GetServerLicense() (*ServerStorageResponseDto, error) {
	ServerStorage, err := get[ServerStorageResponseDto](c, "server/storage")

	return ServerStorage, err
}
