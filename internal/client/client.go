package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	Endpoint   string
	HTTPClient *http.Client
	ApiKey     string
}

func NewClient(endpoint *string, apikey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Endpoint:   *endpoint,
		ApiKey:     *apikey,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	apikey := c.ApiKey

	req.Header.Set("x-api-key", apikey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

type Owner struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	ProfileImagePath string `json:"profileImagePath"`
	AvatarColor      string `json:"avatarColor"`
	ProfileChangedAt string `json:"profileChangedAt"`
}

type AlbumUpdate struct {
	ID          string `json:"id,omitempty"`
	AlbumName   string `json:"albumName"`
	Description string `json:"description,omitempty"`
}

type Album struct {
	ID                    string `json:"id"`
	AlbumName             string `json:"albumName"`
	AlbumThumbnailAssetId string `json:"albumThumbnailAssetId"`
	Description           string `json:"description"`
	Shared                bool   `json:"shared"`
	HasSharedLink         bool   `json:"hasSharedLink"`
	Order                 string `json:"order"`
	IsActivityEnabled     bool   `json:"isActivityEnabled"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	OwnerId               string `json:"ownerId"`
	Owner                 Owner  `json:"owner"`
}

func (c *Client) GetAlbums() ([]Album, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/albums", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	albums := []Album{}
	err = json.Unmarshal(body, &albums)
	if err != nil {
		return nil, err
	}

	return albums, nil
}
func (c *Client) GetAlbum(albumId string) (*Album, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/albums/%s", c.Endpoint, albumId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	album := Album{}
	err = json.Unmarshal(body, &album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

func (c *Client) CreateAlbum(albumName string, description string) (*Album, error) {

	album := AlbumUpdate{
		AlbumName:   albumName,
		Description: description,
	}

	rb, err := json.Marshal(album)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile("/tmp/x.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(string(rb))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/albums", c.Endpoint), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newAlbum := Album{}
	err = json.Unmarshal(body, &newAlbum)
	if err != nil {
		return nil, err
	}

	return &newAlbum, nil
}
