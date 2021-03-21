package main

import (
	"fmt"

	"github.com/futuretea/go-harvester/pkg/client"
)

func main() {
	c, err := client.New("https://10.5.6.11:30443", nil)
	if err != nil {
		panic(err)
	}
	if err = c.Auth.Login("admin", "password"); err != nil {
		panic(err)
	}
	users, err := c.User.List()
	if err != nil {
		return
	}
	for _, user := range users.Data {
		fmt.Println(user)
	}

	nodes, err := c.Node.List()
	if err != nil {
		return
	}
	for _, node := range nodes.Data {
		fmt.Println(c.Node.Get(node.Name))
	}

	services, err := c.Service.List()
	if err != nil {
		return
	}
	for _, service := range services.Data {
		fmt.Println(service)
	}
}
