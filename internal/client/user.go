// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type UserDelete struct {
	Force bool `json:"bool"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"isAdmin"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"deletedAt"`
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
}

func (c *Client) GetUsers() ([]User, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/users", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	users := []User{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *Client) GetUser(userId string) (*User, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/users/%s", c.Endpoint, userId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) CreateUser(user UserUpdate) (*User, error) {

	rb, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/admin/users", c.Endpoint), strings.NewReader(string(rb)))

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newUser := User{}
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (c *Client) UpdateUser(ID string, user UserUpdate) (*User, error) {

	rb, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/admin/users/%s", c.Endpoint, ID), strings.NewReader(string(rb)))

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedUser := User{}
	err = json.Unmarshal(body, &updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (c *Client) DeleteUser(ID string) error {

	delete := UserDelete{
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

	updatedUser := User{}
	err = json.Unmarshal(body, &updatedUser)
	if err != nil {
		return err
	}

	return nil
}
