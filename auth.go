package goharv

import (
	"encoding/json"
	"net/http"

	"github.com/guonaihong/gout"
	"github.com/rancher/wrangler/pkg/slice"
)

func UnmarshalAuthModes(data []byte) (AuthModes, error) {
	var r AuthModes
	err := json.Unmarshal(data, &r)
	return r, err
}

type AuthModes struct {
	Modes []string `json:"modes"`
}

type AuthClient struct {
	v1AuthMode       *apiClient
	v1Auth           *apiClient
	v3localProviders *apiClient
}

func newAuthClient(c *Client) *AuthClient {
	return &AuthClient{
		v1AuthMode:       newAPIClient(c, "auth-modes", true),
		v1Auth:           newAPIClient(c, "auth", true),
		v3localProviders: newAPIClient(c, "localProviders", true),
	}
}

func (s *AuthClient) Login(username, password string) error {
	respCode, respBody, err := s.v1AuthMode.List()
	if err != nil {
		return err
	}
	if respCode != http.StatusOK {
		return NewResponseError(respCode, respBody)
	}
	authModes, err := UnmarshalAuthModes(respBody)
	if err != nil {
		return err
	}
	if slice.ContainsString(authModes.Modes, "rancher") {
		respCode, respBody, err = s.v3localProviders.Action("local", "login", gout.H{
			"username":     username,
			"password":     password,
			"ttl":          57600000,
			"description":  "UI Session",
			"responseType": "cookie",
		})
	} else {
		respCode, respBody, err = s.v1Auth.Action("", "login", gout.H{
			"username": username,
			"password": password,
		})
	}
	if err != nil {
		return err
	}
	if respCode != http.StatusOK {
		return NewResponseError(respCode, respBody)
	}
	return nil
}
