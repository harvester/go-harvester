package apis

import (
	"net/http"

	"github.com/guonaihong/gout"
)

type AuthClient struct {
	*Resource
}

func (s *AuthClient) Login(username, password string) error {
	respCode, respBody, err := s.Resource.Action("", "login", gout.H{
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
