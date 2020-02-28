package foreman_test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/bishy999/go-foreman/pkg/foreman"
)

func ExampleConnectionInfo_CheckUserInput() {

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult string
	}{
		{name: "standard input", username: "test", password: "test", url: "http://mytest.com", expectedresult: "testdev"},
	}

	os.Args = []string{"/fake/loc/main", "create", "-name=testdev", "-size=i3.2xlarge", "-group=1", "-profile=2"}
	for _, tc := range tt {
		conn := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client}
		ok, err := conn.CheckUserInput()

		fmt.Printf("%t %v", ok, err)

		// Output:
		// true <nil>

	}
}
func TestValidation(t *testing.T) {

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult string
	}{
		{name: "standard input", username: "test", password: "test", url: "http://mytest.com", expectedresult: "testdev"},
	}

	os.Args = []string{"/fake/loc/main", "create", "-name=testdev", "-size=i3.2xlarge"}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client}
			_, err := api.CheckUserInput()
			if tc.expectedresult != api.Hostname {
				t.Errorf("Test %v result should be %v, got  `%v`", tc.name, tc.expectedresult, api.Hostname)
			}
			log.Printf("Response: %s", err)
		})

	}
}

func TestInvalidCredentials(t *testing.T) {

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult string
	}{
		{name: "no username|password set", username: "", password: "", url: "http://mytest.com", expectedresult: "username, password & url should be part of the api call"},
		{name: "no url set", username: "", password: "", url: "", expectedresult: "username, password & url should be part of the api call"},
	}

	os.Args = []string{"/fake/loc/main", "create", "-name=testdev", "-size=i3.2xlarge"}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client}
			_, err := api.CheckUserInput()
			if tc.expectedresult != err.Error() {
				t.Errorf("Test %v result should be %v, got  `%v`", tc.name, tc.expectedresult, err)
			}
			log.Printf("Response: %s", err)
		})

	}
}

func TestInvalidFlags(t *testing.T) {

	tt := []struct {
		name           string
		username       string
		password       string
		url            string
		client         *http.Client
		expectedresult string
		hostname       string
	}{
		{name: "size not set", username: "test", password: "test", url: "http://test.com", expectedresult: "size needs to be provided"},
	}

	os.Args = []string{"/fake/loc/main", "create", "-name=testdev", "-size="}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			api := foreman.ConnectionInfo{Username: tc.username, Password: tc.password, BaseURL: tc.url, Client: tc.client}
			_, err := api.CheckUserInput()
			if tc.expectedresult != err.Error() {
				t.Errorf("Test %v result should be %v, got  `%v`", tc.name, tc.expectedresult, err)
			}
			log.Printf("Response: %s", err)
		})

	}
}
