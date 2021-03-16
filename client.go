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

	Auth                    *apis.AuthAPI
	Users                   *apis.UsersAPI
	Images                  *apis.ImagesAPI
	Settings                *apis.SettingsAPI
	KeyPairs                *apis.KeyPairsAPI
	DataVolumes             *apis.DataVolumesAPI
	VirtualMachines         *apis.VirtualMachinesAPI
	VirtualMachineInstances *apis.VirtualMachineInstanceAPI
	Nodes                   *apis.NodesAPI
	SVCs                    *apis.ServicesAPI
	Networks                *apis.NetworksAPI
}

func NewService(c *Client, pluralName string, public bool) *apis.Resource {
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
	c.Auth = &apis.AuthAPI{
		Resource: NewService(c, "auth", true),
	}
	c.Users = &apis.UsersAPI{
		Resource: NewService(c, "harvester.cattle.io.users", false),
	}
	c.Images = &apis.ImagesAPI{
		Resource: NewService(c, "harvester.cattle.io.virtualmachineimages", false),
	}
	c.Settings = &apis.SettingsAPI{
		Resource: NewService(c, "harvester.cattle.io.settings", false),
	}
	c.KeyPairs = &apis.KeyPairsAPI{
		Resource: NewService(c, "harvester.cattle.io.keypairs", false),
	}
	c.DataVolumes = &apis.DataVolumesAPI{
		Resource: NewService(c, "cdi.kubevirt.io.datavolumes", false),
	}
	c.VirtualMachines = &apis.VirtualMachinesAPI{
		Resource: NewService(c, "kubevirt.io.virtualmachines", false),
	}
	c.VirtualMachineInstances = &apis.VirtualMachineInstanceAPI{
		Resource: NewService(c, "kubevirt.io.virtualmachineinstance", false),
	}
	c.Nodes = &apis.NodesAPI{
		Resource: NewService(c, "nodes", false),
	}
	c.SVCs = &apis.ServicesAPI{
		Resource: NewService(c, "services", false),
	}
	c.Networks = &apis.NetworksAPI{
		Resource: NewService(c, "k8s.cni.cncf.io.network-attachment-definitions", false),
	}
	return c
}
