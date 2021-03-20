package goharv

import (
	"encoding/json"
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	harv1 "github.com/rancher/harvester/pkg/apis/harvester.cattle.io/v1alpha1"
)

type User harv1.User

type UserList struct {
	types.Collection
	Data []*User `json:"data"`
}

type UsersClient struct {
	*apiClient
}

func newUsersClient(c *Client) *UsersClient {
	return &UsersClient{
		apiClient: newAPIClient(c, "harvester.cattle.io.users"),
	}
}

func (s *UsersClient) List() (*UserList, error) {
	var collection UserList
	respCode, respBody, err := s.apiClient.List()
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, NewResponseError(respCode, respBody)
	}
	err = json.Unmarshal(respBody, &collection)
	return &collection, err
}

func (s *UsersClient) Create(obj *User) (*User, error) {
	var created *User
	respCode, respBody, err := s.apiClient.Create(obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusCreated {
		return nil, NewResponseError(respCode, respBody)
	}
	if err = json.Unmarshal(respBody, &created); err != nil {
		return nil, err
	}
	return created, nil
}
