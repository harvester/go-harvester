package main

import (
	"fmt"

	"github.com/harvester/go-harvester/pkg/client"
)

func main() {
	c, err := client.New("https://10.5.6.11:30443", nil)
	if err != nil {
		panic(err)
	}
	c.Auth.V1AuthMode.Debug = true
	if err = c.Auth.Login("admin", "admin"); err != nil {
		panic(err)
	}
	users, err := c.Users.List()
	if err != nil {
		return
	}
	for _, user := range users.Data {
		fmt.Println(user)
	}

	nodes, err := c.Nodes.List()
	if err != nil {
		return
	}
	for _, node := range nodes.Data {
		fmt.Println(c.Nodes.Get(node.Name))
	}

	services, err := c.Services.List()
	if err != nil {
		return
	}
	for _, service := range services.Data {
		fmt.Println(service)
	}
}
