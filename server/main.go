package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	serverAddress = "192.168.1.80"
)

var (
	consulClient *consul.Client
	serverID     string
)

func init() {
	var err error
	config := consul.DefaultConfig()
	consulClient, err = consul.NewClient(config)
	if err != nil {
		log.Fatalf("Cannot initialize consul client: %v\n", err)
	}

	serverID = generateServerID()
}

func generateServerID() string {
	rand.Seed(time.Now().UnixNano())
	sid := rand.Intn(65534)
	return strconv.Itoa(sid)
}

func registerService(name, address string, port int) (string, error) {
	serviceID := fmt.Sprintf("%s-%s", name, serverID)
	consulService := consul.AgentServiceRegistration{
		ID: serviceID,
		Tags: []string{
			"urlprefix-/api",
		},
		Check: &consul.AgentServiceCheck{
			Interval: "10s",
			Timeout:  "8s",
			HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
			Status:   "passing",
		},
		Name:    name,
		Address: address,
		Port:    port,
	}
	if err := consulClient.Agent().ServiceRegister(&consulService); err != nil {
		return "", err
	}
	return serviceID, nil
}

func deregisterService(id string) error {
	err := consulClient.Agent().ServiceDeregister(id)
	log.Printf("Deregistered service with id %s\n", id)
	return err
}

func main() {
	// Setup server port
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	// Register service
	serviceID, err := registerService("api", serverAddress, port)
	if err != nil {
		log.Fatalf("Cannot register consul service: %v\n", err)
	}
	log.Printf("Registered service with id %s\n", serviceID)
	defer deregisterService(serviceID)

	// Run api
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/api", func(c echo.Context) error {
		msg := fmt.Sprintf("Hello from %s", serverID)
		return c.String(http.StatusOK, msg)
	})
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	go func() {
		terminate := make(chan os.Signal, 1)
		signal.Notify(terminate, os.Interrupt, syscall.SIGINT)
		<-terminate
		e.Shutdown(context.Background())
	}()

	e.Start(fmt.Sprintf(":%d", port))
}
