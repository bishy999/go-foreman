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
	hostsTimeout = 180
)

func ExampleConnectionInfo_CheckHost() {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		check(rw.Write([]byte("ok")))
	}))
	defer server.Close()

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult bool
		hostname       string
	}{
		{name: "Host Exists", username: "test", password: "test", url: server.URL, client: server.Client(), expectedresult: true, hostname: "test"},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, hostsTimeout*time.Second)
	defer cancel()

	for _, tc := range tt {

		api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client, Hostname: tc.hostname}
		exists, status, err := api.CheckHost(ctx)

		fmt.Printf("%t %s %v", exists, status, err)

		// Output: true ok <nil>

	}

}

func TestHostDoesNotExist(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		check(rw.Write([]byte("Resource host not found by i")))
	}))
	defer server.Close()

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult bool
		hostname       string
	}{
		{name: "Host Doesn't Exist", username: "test", password: "test", url: server.URL, client: server.Client(), expectedresult: false, hostname: "test"},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, hostsTimeout*time.Second)
	defer cancel()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client}
			exists, status, err := api.CheckHost(ctx)
			if err != nil {
				t.Fatalf("Could not read response %v correctly", err)
			}
			if tc.expectedresult != exists {
				t.Errorf("Test %v result should be %v, got  `%v`", tc.name, tc.expectedresult, status)
			}
			log.Printf("Response: %s", status)
		})
	}

}

func ExampleConnectionInfo_CreateHost() {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		check(rw.Write([]byte("ok")))
	}))
	defer server.Close()

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult bool
		hostname       string
	}{
		{name: "Host Exists", username: "test", password: "test", url: server.URL, client: server.Client(), expectedresult: true, hostname: "test"},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, hostsTimeout*time.Second)
	defer cancel()

	for _, tc := range tt {
		api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client, Hostname: tc.hostname}
		exists, _, err := api.CreateHost(ctx)

		fmt.Printf("%t %v", exists, err)

		// Output:
		// true <nil>
	}
}

func BenchmarkConnectionInfo_CreateHost(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		check(rw.Write([]byte("ok")))
	}))
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		tt := []struct {
			name           string
			username       string
			password       string
			url            string
			client         *http.Client
			expectedresult bool
			hostname       string
		}{
			{name: "Host Exists", username: "test", password: "test", url: server.URL, client: server.Client(), expectedresult: true, hostname: "test"},
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, hostsTimeout*time.Second)
		defer cancel()

		for _, tc := range tt {
			api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client, Hostname: tc.hostname}
			exists, _, err := api.CreateHost(ctx)

			fmt.Printf("%t %v", exists, err)

			// Output:
			// true <nil>
		}
	}
}

func ExampleConnectionInfo_DeleteHost() {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		check(rw.Write([]byte("ok")))

	}))
	defer server.Close()

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult string
		hostname       string
	}{
		{name: "Host Delete", username: "test", password: "test", url: server.URL, client: server.Client(), expectedresult: "ok", hostname: "test"},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, hostsTimeout*time.Second)
	defer cancel()

	for _, tc := range tt {
		api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client, Hostname: tc.hostname}
		_, status, err := api.DeleteHost(ctx)

		fmt.Printf(" %s %v", status, err)

		// Output:  ok <nil>
	}

}

func check(n int, err error) {
    if err != nil {
        log.Printf("Write failed: %v", err)
    }
}
