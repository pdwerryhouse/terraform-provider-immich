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
	"fmt"
	"net/http"
)

type ExifResponseDto struct {
	City             string `json:"city"`
	Country          string `json:"country"`
	DateTimeOriginal string `json:"dateTimeOriginal"`
	Description      string `json:"description"`
	ExifImageHeight  int64  `json:"exifImageHeight"`
	ExifImageWidth   int64  `json:"exifImageWidth"`
	ExposureTime     string `json:"exposureTime"`
	FNumber          int64  `json:"fNumber"`
	FileSizeInByte   int64  `json:"fileSizeInByte"`
	FocalLength      int64  `json:"focalLength"`
	Iso              int64  `json:"iso"`
	Latitude         int64  `json:"latitude"`
	LensModel        string `json:"lensModel"`
	Longitude        int64  `json:"longitude"`
	Make             string `json:"make"`
	Model            string `json:"model"`
	ModifyDate       string `json:"modifyDate"`
	Orientation      string `json:"orientation"`
	ProjectionType   string `json:"projectionType"`
	Rating           int64  `json:"rating"`
	State            string `json:"state"`
	TimeZone         string `json:"timeZone"`
}

type AssetStackResponseDto struct {
	AssetCount     int64  `json:"assetCount"`
	Id             string `json:"id"`
	PrimaryAssetId string `json:"primaryAssetId"`
}

type AssetResponseDto struct {
	Checksum         string                `json:"checksum"`
	CreatedAt        string                `json:"createdAt"`
	DuplicateId      string                `json:"duplicateId"`
	Duration         int64                 `json:"duration"`
	ExifInfo         ExifResponseDto       `json:"exifInfo"`
	FileCreatedAt    string                `json:"fileCreatedAt"`
	FileModifiedAt   string                `json:"fileModifiedAt"`
	HasMetadata      bool                  `json:"hasMetadata"`
	Height           int64                 `json:"height"`
	Id               string                `json:"id"`
	IsArchived       bool                  `json:"isArchived"`
	IsEdited         bool                  `json:"isEdited"`
	IsFavorite       bool                  `json:"isFavorite"`
	IsOffline        bool                  `json:"isOffline"`
	IsTrashed        bool                  `json:"isTrashed"`
	LibraryId        string                `json:"libraryId"`
	LivePhotoVideoId string                `json:"livePhotoVideoId"`
	LocalDateTime    string                `json:"localDateTime"`
	OriginalFileName string                `json:"originalFileName"`
	OriginalMimeType string                `json:"originalMimeType"`
	OriginalPath     string                `json:"originalPath"`
	Owner            UserResponseDto       `json:"owner"`
	OwnerId          string                `json:"ownerId"`
	People           []PersonResponseDto   `json:"people"`
	Resized          bool                  `json:"resized"`
	Stack            AssetStackResponseDto `json:"stack"`
	Tags             []TagResponseDto      `json:"tags"`
	Thumbhash        string                `json:"thumbhash"`
	Type             string                `json:"type"`
	UpdatedAt        string                `json:"updatedAt"`
	Visibility       string                `json:"visibility"`
	Width            int64                 `json:"width"`
}

// XXX Value is not a string
type AssetMetadataUpsertItemDto struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AssetMediaCreateDto struct {
	AssetData        string                       `json:"assetData"`
	Duration         int64                        `json:"duration"`
	FileCreatedAt    string                       `json:"fileCreatedAt"`
	FileModifiedAt   string                       `json:"fileModifiedAt"`
	Filename         string                       `json:"filename"`
	IsFavorite       bool                         `json:"isFavorite"`
	LivePhotoVideoId string                       `json:"livePhotoVideoId"`
	Metadata         []AssetMetadataUpsertItemDto `json:"metadata"`
	SidecarData      string                       `json:"sidecarData"`
	Visibility       string                       `json:"visibility"`
}

type AssetMediaResponseDto struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type UpdateAssetDto struct {
	DateTimeOriginal string `json:"dateTimeOriginal"`
	Description      string `json:"description"`
	IsFavorite       bool   `json:"isFavorite"`
	Latitude         int64  `json:"latitude"`
	LivePhotoVideoId string `json:"livePhotoVideoId"`
	Longitude        int64  `json:"longitude"`
	Rating           int64  `json:"rating"`
	Visibility       string `json:"visibility"`
}

type AssetStatsResponseDto struct {
	Images int64 `json:"images"`
	Total  int64 `json:"total"`
	Videos int64 `json:"videos"`
}

func (c *Client) GetAssetInfo(assetId string, key string, slug string) (*AssetResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/assets/%s", c.Endpoint, assetId), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if key != "" {
		q.Add("key", key)
	}
	if slug != "" {
		q.Add("slug", slug)
	}
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var asset *AssetResponseDto
	err = json.Unmarshal(body, &asset)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (c *Client) CreateAsset(asset AssetMediaCreateDto) (*AssetMediaResponseDto, error) {
	newAsset, err := post[AssetMediaResponseDto](c, "assets", asset)

	return newAsset, err
}

func (c *Client) UpdateAsset(assetId string, asset UpdateAssetDto) (*AssetResponseDto, error) {
	newAsset, err := put[AssetResponseDto](c, assetId, "assets", asset)

	return newAsset, err
}

func (c *Client) DeleteAsset(assetId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/assets/%s", c.Endpoint, assetId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *Client) GetAssetStatistics(isFavorite *bool, isTrashed *bool, visibility string) (*AssetStatsResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/assets/statistics", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if isFavorite != nil {
		q.Add("isFavorite", fmt.Sprintf("%s", *isFavorite))
	}
	if isTrashed != nil {
		q.Add("isTrashed", fmt.Sprintf("%s", *isTrashed))
	}
	if visibility != "" {
		q.Add("visibility", visibility)
	}
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response *AssetStatsResponseDto
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
