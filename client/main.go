package main

import (
	"fmt"
	"log"

	consul "github.com/hashicorp/consul/api"
)

var (
	consulClient *consul.Client
)

func init() {
	var err error
	consulClient, err = consul.NewClient(consul.DefaultConfig())
	if err != nil {
		log.Fatalf("Cannot initialize consul client: %v\n", err)
	}
}

func resolveService(name string) (string, error) {
	services, _, err := consulClient.Catalog().Service(name, "", nil)
	if err != nil {
		return "", err
	}
	for _, s := range services {
		addr := fmt.Sprintf("%s:%d", s.Address, s.ServicePort)
		fmt.Println("Address", addr)
	}

	return "", nil
}

func doRequest() error {
	return nil
}

func main() {
	_, err := resolveService("api")
	if err != nil {
		log.Fatalf("Failed to resolve service: %v", err)
	}
}
