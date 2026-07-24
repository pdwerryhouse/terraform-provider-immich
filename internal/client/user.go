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

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
type UserDelete struct {
	Force bool `json:"bool"`
}

type User struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	IsAdmin              bool   `json:"isAdmin"`
	Status               string `json:"status"`
	Notify               bool   `json:"notify"`
	QuotaSizeInBytes     int64  `json:"quotaSizeInBytes"`
	ShouldChangePassword bool   `json:"shouldChangePassword"`
	StorageLabel         string `json:"storageLabel"`
	PinCode              string `json:"pinCode"`
	AvatarColor          string `json:"avatarColor"`
	OauthId              string `json:"oauthId"`
	ProfileChangedAt     string `json:"profileChangedAt"`
	ProfileImagePath     string `json:"profileImagePath"`
	QuotaUsageInBytes    int64  `json:"quotaUsageInBytes"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt            string `json:"updatedAt"`
	DeletedAt            string `json:"deletedAt"`
}
type UserUpdate struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	IsAdmin              bool   `json:"isAdmin,omitempty"`
	Notify               bool   `json:"notify,omitempty"`
	QuotaSizeInBytes     int64  `json:"quotaSizeInBytes,omitempty"`
	ShouldChangePassword bool   `json:"shouldChangePassword,omitempty"`
	StorageLabel         string `json:"storageLabel,omitempty"`
	PinCode              string `json:"pinCode,omitempty"`
	AvatarColor          string `json:"avatarColor,omitempty"`
}
*/

type UserResponseDto struct {
	AvatarColor      string `json:"avatarColor"`
	Email            string `json:"email"`
	Id               string `json:"id"`
	Name             string `json:"name"`
	ProfileChangedAt string `json:"profileChangedAt"`
	ProfileImagePath string `json:"profileImagePath"`
}

type UserAdminResponseDto struct {
	AvatarColor          string      `json:"avatarColor"`
	CreatedAt            string      `json:"createdAt"`
	DeletedAt            string      `json:"deletedAt"`
	Email                string      `json:"email"`
	Id                   string      `json:"id"`
	IsAdmin              bool        `json:"isAdmin"`
	License              UserLicense `json:"license"`
	Name                 string      `json:"name"`
	OauthId              string      `json:"oauthId"`
	ProfileChangedAt     string      `json:"profileChangedAt"`
	ProfileImagePath     string      `json:"profileImagePath"`
	QuotaSizeInBytes     int64       `json:"quotaSizeInBytes"`
	QuotaUsageInBytes    int64       `json:"quotaUsageInBytes"`
	ShouldChangePassword bool        `json:"shouldChangePassword"`
	Status               string      `json:"status"`
	StorageLabel         string      `json:"storageLabel"`
	UpdatedAt            string      `json:"updatedAt"`
}

type UserAdminCreateDto struct {
	AvatarColor          string `json:"avatarColor"`
	Email                string `json:"email"`
	IsAdmin              bool   `json:"isAdmin"`
	Name                 string `json:"name"`
	Notify               bool   `json:"notify"`
	Password             string `json:"password"`
	PinCode              string `json:"pinCode"`
	QuotaSizeInBytes     int64  `json:"quotaSizeInBytes"`
	ShouldChangePassword bool   `json:"shouldChangePassword"`
	StorageLabel         string `json:"storageLabel"`
}

type UserAdminUpdateDto struct {
	AvatarColor          string `json:"avatarColor,omitempty"`
	Email                string `json:"email"`
	IsAdmin              bool   `json:"isAdmin,omitempty"`
	Name                 string `json:"name"`
	Password             string `json:"password"`
	PinCode              string `json:"pinCode,omitempty"`
	QuotaSizeInBytes     int64  `json:"quotaSizeInBytes,omitempty"`
	ShouldChangePassword bool   `json:"shouldChangePassword,omitempty"`
	StorageLabel         string `json:"storageLabel,omitempty"`
}

type UserAdminDeleteDto struct {
	Force bool `json:"force"`
}

func (c *Client) GetUsers() ([]UserAdminResponseDto, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/users", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	users := []UserAdminResponseDto{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *Client) GetUser(userId string) (*UserAdminResponseDto, error) {
	user, err := get_by_id[UserAdminResponseDto](c, userId, "admin/users")

	return user, err
}

func (c *Client) CreateUser(user UserAdminCreateDto) (*UserAdminResponseDto, error) {
	newUser, err := post[UserAdminResponseDto](c, "admin/users", user)

	return newUser, err
}

func (c *Client) UpdateUser(userId string, user UserAdminUpdateDto) (*UserAdminResponseDto, error) {
	newUser, err := patch[UserAdminResponseDto](c, userId, "admin/users", user)

	return newUser, err
}

func (c *Client) DeleteUser(ID string) error {

	delete := UserAdminDeleteDto{
		Force: false,
	}

	rb, err := json.Marshal(delete)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/admin/users/%s", c.Endpoint, ID), strings.NewReader(string(rb)))

	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	updatedUser := UserAdminResponseDto{}
	err = json.Unmarshal(body, &updatedUser)
	if err != nil {
		return err
	}

	return nil
}
