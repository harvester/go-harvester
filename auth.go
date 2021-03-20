package goharv

import (
	"net/http"

	"github.com/guonaihong/gout"
)

type AuthClient struct {
	*apiClient
}

func newAuthClient(c *Client) *AuthClient {
	return &AuthClient{
		apiClient: newAPIClient(c, "auth", true),
	}
}

func (s *AuthClient) Login(username, password string) error {
	respCode, respBody, err := s.apiClient.Action("", "login", gout.H{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}
	if respCode != http.StatusOK {
		return NewResponseError(respCode, respBody)
	}
	return nil
}
