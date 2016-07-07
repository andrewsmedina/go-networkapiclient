package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	endpoint := "http://localhost:4243"
	client, err := NewClient(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint.String() != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
}

func TestListVlans(t *testing.T) {
	jsonVlans := `[
     {
             "Id": "8dfafdbc3a40",
             "Image": "base:latest",
             "Command": "echo 1",
             "Created": 1367854155,
             "Ports":[{"PrivatePort": 2222, "PublicPort": 3333, "Type": "tcp"}],
             "Status": "Exit 0"
     }
]`
	var expected []Vlan
	err := json.Unmarshal([]byte(jsonVlans), &expected)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, jsonVlans)
	}))
	defer ts.Close()
	client, err := NewClient(ts.URL)
	vlans, err := client.ListVlans(ListVlansOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(vlans, expected) {
		t.Errorf("ListVlans: Expected %#v. Got %#v.", expected, vlans)
	}
}
