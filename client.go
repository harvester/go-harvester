package goharv

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	defaultAPIVersion = "v1"
)

type Client struct {
	HTTPClient *http.Client
	APIVersion string
	BaseURL    *url.URL

	Auth                    *AuthClient
	Users                   *UsersClient
	Images                  *ImagesClient
	Settings                *SettingsClient
	KeyPairs                *KeyPairsClient
	DataVolumes             *DataVolumesClient
	VirtualMachines         *VirtualMachinesClient
	VirtualMachineInstances *VirtualMachineInstanceClient
	Nodes                   *NodesClient
	SVCs                    *ServicesClient
	Networks                *NetworksClient
}

func newAPIClient(c *Client, pluralName string, publicOptions ...bool) *apiClient {
	var public bool
	if len(publicOptions) > 0 {
		public = publicOptions[0]
	}
	return &apiClient{
		PluralName: pluralName,
		Public:     public,
		Debug:      false,
		BaseURL:    c.BaseURL,
		APIVersion: c.APIVersion,
		HTTPClient: c.HTTPClient,
	}
}

func New(harvesterURL string, httpClient *http.Client) *Client {
	jar, _ := cookiejar.New(nil)
	if httpClient == nil {
		httpClient = &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	baseURL, _ := url.Parse(harvesterURL)

	c := &Client{
		HTTPClient: httpClient,
		APIVersion: defaultAPIVersion,
		BaseURL:    baseURL,
	}

	c.Auth = newAuthClient(c)
	c.Users = newUsersClient(c)
	c.Images = newImagesClient(c)
	c.Settings = newSettingsClient(c)
	c.KeyPairs = newKeyPairsClient(c)
	c.DataVolumes = newDataVolumesClient(c)
	c.VirtualMachines = newVirtualMachinesClient(c)
	c.VirtualMachineInstances = newVirtualMachineInstanceClient(c)
	c.Nodes = newNodesClient(c)
	c.SVCs = newServicesClient(c)
	c.Networks = newNetworksClient(c)
	return c
}
