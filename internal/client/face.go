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

type AssetFaceResponseDto struct {
	BoundingBoxX1 int64             `json:"boundingBoxX1"`
	BoundingBoxX2 int64             `json:"boundingBoxX2"`
	BoundingBoxY1 int64             `json:"boundingBoxY1"`
	BoundingBoxY2 int64             `json:"boundingBoxY2"`
	Id            string            `json:"id"`
	ImageHeight   int64             `json:"imageHeight"`
	ImageWidth    int64             `json:"imageWidth"`
	Person        PersonResponseDto `json:"person"`
	SourceType    string            `json:"sourceType"`
}

type AssetFaceCreateDto struct {
	AssetId     string `json:"assetId"`
	Height      int64  `json:"height"`
	ImageHeight int64  `json:"imageHeight"`
	ImageWidth  int64  `json:"imageWidth"`
	PersonId    string `json:"personId"`
	Width       int64  `json:"width"`
	X           int64  `json:"x"`
	Y           int64  `json:"y"`
}

type FaceDto struct {
	Id string `json:"id"`
}

type AssetFaceDeleteDto struct {
	Force bool `json:"force"`
}

func (c *Client) GetFaces(id string) ([]AssetFaceResponseDto, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/faces", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("id", id)
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	faces := []AssetFaceResponseDto{}
	err = json.Unmarshal(body, &faces)
	if err != nil {
		return nil, err
	}

	return faces, nil
}

func (c *Client) GetFace(faceId string) (*AssetFaceResponseDto, error) {
	face, err := get_by_id[AssetFaceResponseDto](c, faceId, "faces")

	return face, err
}

func (c *Client) CreateFace(face AssetFaceCreateDto) (*AssetFaceResponseDto, error) {
	newFace, err := post[AssetFaceResponseDto](c, "faces", face)

	return newFace, err
}

func (c *Client) UpdateFace(faceId string, face FaceDto) (*AssetFaceResponseDto, error) {
	newFace, err := put[AssetFaceResponseDto](c, faceId, "faces", face)

	return newFace, err
}

func (c *Client) DeleteFace(faceId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/faces/%s", c.Endpoint, faceId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
