package goharv

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/futuretea/go-harvester/apis"
)

const (
	defaultAPIVersion = "v1"
)

type Client struct {
	HTTPClient *http.Client
	APIVersion string
	BaseURL    *url.URL

	Auth                    *apis.AuthClient
	Users                   *apis.UsersClient
	Images                  *apis.ImagesClient
	Settings                *apis.SettingsClient
	KeyPairs                *apis.KeyPairsClient
	DataVolumes             *apis.DataVolumesClient
	VirtualMachines         *apis.VirtualMachinesClient
	VirtualMachineInstances *apis.VirtualMachineInstanceClient
	Nodes                   *apis.NodesClient
	SVCs                    *apis.ServicesClient
	Networks                *apis.NetworksClient
}

func NewService(c *Client, pluralName string, publicOptions ...bool) *apis.Resource {
	var public bool
	if len(publicOptions) > 0 {
		public = publicOptions[0]
	}
	return &apis.Resource{
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
	c.Auth = &apis.AuthClient{
		Resource: NewService(c, "auth", true),
	}
	c.Users = &apis.UsersClient{
		Resource: NewService(c, "harvester.cattle.io.users"),
	}
	c.Images = &apis.ImagesClient{
		Resource: NewService(c, "harvester.cattle.io.virtualmachineimages"),
	}
	c.Settings = &apis.SettingsClient{
		Resource: NewService(c, "harvester.cattle.io.settings"),
	}
	c.KeyPairs = &apis.KeyPairsClient{
		Resource: NewService(c, "harvester.cattle.io.keypairs"),
	}
	c.DataVolumes = &apis.DataVolumesClient{
		Resource: NewService(c, "cdi.kubevirt.io.datavolumes"),
	}
	c.VirtualMachines = &apis.VirtualMachinesClient{
		Resource: NewService(c, "kubevirt.io.virtualmachines"),
	}
	c.VirtualMachineInstances = &apis.VirtualMachineInstanceClient{
		Resource: NewService(c, "kubevirt.io.virtualmachineinstance"),
	}
	c.Nodes = &apis.NodesClient{
		Resource: NewService(c, "nodes"),
	}
	c.SVCs = &apis.ServicesClient{
		Resource: NewService(c, "services"),
	}
	c.Networks = &apis.NetworksClient{
		Resource: NewService(c, "k8s.cni.cncf.io.network-attachment-definitions"),
	}
	return c
}
