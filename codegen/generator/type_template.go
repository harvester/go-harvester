package generator

var (
	typeTemplate = `package client

import (
	"net/http"
{{- if .schema | hasGet }}
	"encoding/json"
{{- end}}

	"github.com/harvester/go-harvester/pkg/clientbase"
	"github.com/harvester/go-harvester/pkg/errors"
{{- if .schema | hasGet }}
	"github.com/rancher/apiserver/pkg/types"
	{{.schema.Attributes.importAlias}} "{{.schema.Attributes.importPackage}}"
{{- end}}
)

{{ if .schema | hasGet }}
type {{.schema.CodeName}} {{.schema.Attributes.importAlias}}.{{.schema.Attributes.importType}}

type {{.schema.CodeName}}List struct {
	types.Collection
	Data []*{{.schema.CodeName}} %BACK%json:"data"%BACK%
}
{{end}}

type {{.schema.CodeName}}Client struct {
	*clientbase.APIClient
}

func new{{.schema.CodeName}}Client(c *Client) *{{.schema.CodeName}}Client {
	return &{{.schema.CodeName}}Client{
		APIClient: clientbase.NewAPIClient(c.BaseURL, c.HTTPClient, "{{.schema.Attributes.version}}", "{{.schema.PluralName}}"),
	}
}

{{- $namespaced :=.schema | namespaced}}

{{ if .schema | hasGet }}
func (c *{{.schema.CodeName}}Client) List() (*{{.schema.CodeName}}List, error) {
	var collection {{.schema.CodeName}}List
	respCode, respBody, err := c.APIClient.List()
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	err = json.Unmarshal(respBody, &collection)
	return &collection, err
}

func (c *{{.schema.CodeName}}Client) Create(obj *{{.schema.CodeName}}) (*{{.schema.CodeName}}, error) {
	var created *{{.schema.CodeName}}
	respCode, respBody, err := c.APIClient.Create(obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusCreated {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	err = json.Unmarshal(respBody, &created)
	return created, nil
}

{{if $namespaced }}
func (c *{{.schema.CodeName}}Client) Update(namespace, name string, obj *{{.schema.CodeName}}) (*{{.schema.CodeName}}, error) {
	resourceName := namespace + "/" + name
{{- else}}
func (c *{{.schema.CodeName}}Client) Update(name string, obj *{{.schema.CodeName}}) (*{{.schema.CodeName}}, error) {
	resourceName := name
{{- end}}
	respCode, respBody, err := c.APIClient.Update(resourceName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var updated *{{.schema.CodeName}}
	if err = json.Unmarshal(respBody, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

{{if $namespaced }}
func (c *{{.schema.CodeName}}Client) Get(namespace, name string, opts ...interface{}) (*{{.schema.CodeName}}, error) {
	resourceName := namespace + "/" + name
{{- else}}
func (c *{{.schema.CodeName}}Client) Get(name string, opts ...interface{}) (*{{.schema.CodeName}}, error) {
	resourceName := name
{{- end}}
	respCode, respBody, err := c.APIClient.Get(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *{{.schema.CodeName}}
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}

{{if $namespaced }}
func (c *{{.schema.CodeName}}Client) Delete(namespace, name string, opts ...interface{}) (*{{.schema.CodeName}}, error) {
	resourceName := namespace + "/" + name
{{- else}}
func (c *{{.schema.CodeName}}Client) Delete(name string, opts ...interface{}) (*{{.schema.CodeName}}, error) {
	resourceName := name
{{- end}}
	respCode, respBody, err := c.APIClient.Delete(resourceName, opts...)
	if err != nil {
		return nil, err
	}
	if respCode == http.StatusNoContent {
		return nil, nil
	}
	if respCode != http.StatusOK {
		return nil, errors.NewResponseError(respCode, respBody)
	}
	var obj *{{.schema.CodeName}}
	err = json.Unmarshal(respBody, &obj)
	return obj, nil
}
{{- end}}

{{range $key, $value := .resourceActions}}
    {{if (and (eq $value.Input "") (eq $value.Output ""))}}
{{- if $namespaced }}
        func (c *{{$.schema.CodeName}}Client) {{$key | capitalize}} (namespace, name string) (error) {
			resourceName := namespace + "/" + name
{{- else}}
		func (c *{{$.schema.CodeName}}Client) Action{{$key | capitalize}} (name string) (error) {
			resourceName := name
{{- end}}
			respCode, respBody, err := c.APIClient.Action(resourceName, "{{$key}}", nil)
			if err != nil {
				return err
			}
			if respCode != http.StatusNoContent {
				return errors.NewResponseError(respCode, respBody)
			}
			return nil
    {{else if (and (eq $value.Input "") (ne $value.Output ""))}}
    {{else if (and (ne $value.Input "") (eq $value.Output ""))}}
{{- if $namespaced }}
        func (c *{{$.schema.CodeName}}Client) {{$key | capitalize}} (namespace, name string, {{$value.Input}} interface{}) (error) {
			resourceName := namespace + "/" + name
{{- else}}
		func (c *{{$.schema.CodeName}}Client) Action{{$key | capitalize}} (name string, {{$value.Input}} interface{}) (error) {
			resourceName := name
{{- end}}
			respCode, respBody, err := c.APIClient.Action(resourceName, "{{$key}}", {{$value.Input}})
			if err != nil {
				return err
			}
			if respCode != http.StatusNoContent {
				return errors.NewResponseError(respCode, respBody)
			}
			return nil
    {{else}}
    {{- end -}}
    }
{{end}}

`
)
