package goharv

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"sigs.k8s.io/yaml"
)

type apiClient struct {
	Public     bool
	Debug      bool
	BaseURL    *url.URL
	APIVersion string
	PluralName string
	HTTPClient *http.Client
}

func (r *apiClient) BuildAPIURL() string {
	apiVersion := r.APIVersion
	if r.Public {
		apiVersion += "-public"
	}
	return fmt.Sprintf("%s/%s/%s", r.BaseURL, apiVersion, r.PluralName)
}

func (r *apiClient) BuildResourceURL(namespacedName string) string {
	if namespacedName == "" {
		return r.BuildAPIURL()
	}
	return fmt.Sprintf("%s/%s", r.BuildAPIURL(), namespacedName)
}

func (r *apiClient) NewRequest() *dataflow.Gout {
	return gout.New(r.HTTPClient)
}

func (r *apiClient) Create(object interface{}) (respCode int, respBody []byte, err error) {
	err = r.NewRequest().
		POST(r.BuildAPIURL()).
		SetJSON(object).
		SetHeader().
		BindBody(&respBody).
		Code(&respCode).
		Debug(r.Debug).
		Do()
	return
}

func (r *apiClient) CreateByYAML(object interface{}) (respCode int, respBody []byte, err error) {
	var yamlData []byte
	yamlData, err = yaml.Marshal(object)
	if err != nil {
		return
	}
	err = r.NewRequest().
		POST(r.BuildAPIURL()).
		SetBody(yamlData).
		SetCookies().
		SetHeader(gout.H{"content-type": "application/yaml"}).
		BindBody(&respBody).
		Code(&respCode).
		Debug(r.Debug).
		Do()
	return
}

func (r *apiClient) List() (respCode int, respBody []byte, err error) {
	err = r.NewRequest().
		GET(r.BuildAPIURL()).
		BindBody(&respBody).
		Code(&respCode).
		Debug(r.Debug).
		Do()
	if err != nil {
		return
	}
	if respCode != http.StatusOK {
		return
	}
	return
}

func (r *apiClient) Get(namespacedName string, obj ...interface{}) (respCode int, respBody []byte, err error) {
	err = r.NewRequest().
		GET(r.BuildResourceURL(namespacedName)).
		SetQuery(obj...).
		BindBody(&respBody).
		Code(&respCode).
		Debug(r.Debug).
		Do()
	return
}

func (r *apiClient) Update(namespacedName string, object interface{}) (respCode int, respBody []byte, err error) {
	err = r.NewRequest().
		PUT(r.BuildResourceURL(namespacedName)).
		SetJSON(object).
		BindBody(&respBody).
		Code(&respCode).
		Debug(r.Debug).
		Do()
	return
}

func (r *apiClient) Delete(namespacedName string, obj ...interface{}) (respCode int, respBody []byte, err error) {
	err = r.NewRequest().
		DELETE(r.BuildResourceURL(namespacedName)).
		SetQuery(obj...).
		BindBody(&respBody).
		Code(&respCode).
		Debug(r.Debug).
		Do()
	return
}

func (r *apiClient) Action(namespacedName string, action string, object interface{}) (respCode int, respBody []byte, err error) {
	dataFlow := r.NewRequest().
		POST(fmt.Sprintf("%s?action=%s", r.BuildResourceURL(namespacedName), action)).
		SetHeader().
		BindBody(&respBody).
		Code(&respCode).
		Debug(r.Debug)
	if object != nil {
		dataFlow = dataFlow.SetJSON(object)
	}
	err = dataFlow.Do()
	return
}
