package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bishy999/go-foreman/pkg/foreman"
)

var (
	version    string
	buildstamp string
)

// timeout for context
const (
	timeout = 240
)

func main() {

	log.Printf("Version    : %s\n", version)
	log.Printf("Build Time : %s\n", buildstamp)

	user := os.Getenv("FOREMAN_USER")
	password := os.Getenv("FOREMAN_PASSWORD")
	url := os.Getenv("FOREMAN_URL")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{
		Transport: tr,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()

	connection := foreman.ConnectionInfo{
		Username: user,
		Password: password,
		BaseURL:  url,
		Client:   client,
	}

	_, err := connection.CheckUserInput()
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	_, status, err := connection.CheckStatus(ctx)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	log.Printf("Response: %s", status)

	if strings.Contains(status, "\"status\":422") {
		log.Fatalf("Error: Status code 422 indicates this hostname already exists in a terminated state. Try again with a different hostname")
	}

	exists, _, err := connection.CheckHost(ctx)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	if exists {
		log.Printf("Response: %s already exists ⚠️", connection.Hostname)

		if connection.Action == "delete" {
			_, deleted, err := connection.DeleteHost(ctx)
			if err != nil {
				log.Fatalf("Error: %s", err.Error())
			}
			log.Printf("Response: %s", deleted)
		} else {
			log.Fatalf("Cannot create a host that already exists. Please try a different host name.")
		}
	} else {
		if connection.Action == "delete" {
			log.Printf("Response: [%s] doesn't exist so let's not do any delete action", connection.Hostname)
		} else {
			log.Printf("Response: [%s] doesn't exist so let's create the host via foreman", connection.Hostname)
			_, created, err := connection.CreateHost(ctx)
			if err != nil {
				log.Fatalf("Error: %s", err.Error())
			}
			log.Printf("Response: %s", created)
		}

	}
}
