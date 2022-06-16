package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ghuvrons/fota-server-go/API"
	"github.com/ghuvrons/fota-server-go/models"
	"github.com/joho/godotenv"

	giotgo "github.com/ghuvrons/g-IoT-Go"
)

func main() {
	godotenv.Load()
	models.DBReadEnv()

	fmt.Println("Start")

	var server = giotgo.NewServer()

	setCmdHandlers(server)

	server.ClientAuth(func(username, password string) bool {
		fmt.Println(username, password)
		return true
	})

	defaultHost := "127.0.0.1"
	defaultPort := 2000

	if host := os.Getenv("IOT_HOST"); host != "" {
		defaultHost = host
	}

	if port := os.Getenv("IOT_PORT"); port != "" {
		nerPort, err := strconv.Atoi(port)
		if err == nil {
			defaultPort = nerPort
		}
	}

	addr := fmt.Sprintf("%s:%d", defaultHost, defaultPort)
	fmt.Println("IOT server started at", addr)

	go server.Serve(addr)

	API.Serve()
}
