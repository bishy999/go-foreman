package foreman_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bishy999/go-foreman/pkg/foreman"
)

const (
	expected      = "ok"
	statusTimeout = 180
)

func ExampleConnectionInfo_CheckStatus() {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		check(rw.Write([]byte(expected)))
	}))
	defer server.Close()

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult string
	}{
		{name: "standard input", username: "test", password: "test", url: server.URL, client: server.Client(), expectedresult: "ok"},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, statusTimeout*time.Second)
	defer cancel()

	for _, tc := range tt {
		api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client}
		exists, status, err := api.CheckStatus(ctx)

		fmt.Printf("%t %s %v", exists, status, err)

		// Output:
		// true ok <nil>

	}
}

func TestBadAPI(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadGateway)
		check(rw.Write([]byte("502")))
	}))
	defer server.Close()

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult int
	}{
		{name: "bad api response", username: "test", password: "test", url: server.URL, client: server.Client(), expectedresult: http.StatusBadGateway},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, statusTimeout*time.Second)
	defer cancel()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client}
			exists, status, err := api.CheckStatus(ctx)
			if err != nil {
				t.Fatalf("Could not read response %v correctly", err)
			}
			if tc.expectedresult != http.StatusBadGateway {
				t.Errorf("Test %v result should be %v, got `%v`", tc.name, tc.expectedresult, exists)
			}
			log.Printf("Response: %s", status)
		})

	}

}
