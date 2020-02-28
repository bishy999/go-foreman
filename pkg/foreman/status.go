package foreman

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"path"
)

const (
	statusapi = "api/status"
)

// ConnectionInfo represents data that is needed for a establishig connection to foreman
type ConnectionInfo struct {
	Username string
	Password string
	BaseURL  string
	Client   *http.Client
	Hostname string
	Size     string
	Group    int
	Profile  string
	Action   string
}

// CheckStatus check to see if successfully connected to api
func (cli *ConnectionInfo) CheckStatus(ctx context.Context) (bool, string, error) {

	var jsonData []byte
	var name string

	url, _ := url.Parse(cli.BaseURL)
	url.Path = path.Join(url.Path, statusapi)
	apiused := url.String()
	log.Printf("API Used: %s", apiused)

	exists, status, err := sendRequest(ctx, name, cli.Username, cli.Password, cli.Client, jsonData, apiused, http.MethodGet)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	return exists, status, err
}
