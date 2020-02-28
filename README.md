
# go-foreman

go-foreman is a Go client library for accessing the Forenman API.

You can view the Foreman API docs here: [https://theforeman.org/api/](https://theforeman.org/api/)

You can view the client API docs by serving the docs from this repository : [http://localhost:6060/pkg/](http://localhost:6060/pkg/)
```go
 godoc -http :6060
```

## Status
[![Build Status](https://travis-ci.com/bishy999/go-foreman.svg?branch=master)](https://travis-ci.com/bishy999/go-foreman)
[![Go Report Card](https://goreportcard.com/badge/github.com/bishy999/go-foreman)](https://goreportcard.com/report/github.com/bishy999/go-foreman)
[![GoDoc](https://godoc.org/github.com/bishy999/go-foreman/pkg/foreman?status.svg)](https://godoc.org/github.com/bishy999/go-foremanpkg/foreman)
[![GolangCI](https://golangci.com/badges/github.com/bishy999/go-foreman.svg)](https://golangci.com)
![GitHub Repo size](https://img.shields.io/github/repo-size/bishy999/go-foreman)
[![GitHub Tag](https://img.shields.io/github/tag/bishy999/go-foreman.svg)](https://github.com/bishy999/go-foreman/releases/latest)
[![GitHub Activity](https://img.shields.io/github/commit-activity/m/bishy999/go-foreman)](https://github.com/bishy999/go-foreman)
[![GitHub Contributors](https://img.shields.io/github/contributors/bishy999/go-foreman)](https://github.com/bishy999/go-foreman)


## Usage (package)

### Download package
```go
 go get github.com/bishy999/go-foreman/foreman
 ```

### Use package
```go
import 
(
	"github.com/bishy999/go-foreman/pkg/foreman"
)
```

### Authentication
You will need a user and password with sufficent priviliges to perform actions against the formeman api

You can then use these credentials to create a new client. An example of a client is stored under the cmd directory in this repository

```go


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



```

## Usage (binary)

Download the client binary from the repository and compile it with version 

Go get will download from the master, as such when we download it give it the tag verison from the master

```go
go get -ldflags "-X main.version=v1.0.2 -X main.buildstamp=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'`)" github.com/bishy999/go-foreman/foreman/cmd/foreman-client

foreman-client create -name=mytestenv.com -size=i3.4xlarge -group=1 -profile=2

```


## Contributing

We love pull requests! Please see the [contribution guidelines](CONTRIBUTING.md).