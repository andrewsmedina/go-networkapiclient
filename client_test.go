package network

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	endpoint := "http://localhost:4243"
	client, err := NewClient(endpoint, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if client.endpoint.String() != endpoint {
		t.Errorf("Expected endpoint %s. Got %s.", endpoint, client.endpoint)
	}
}

func TestListVlans(t *testing.T) {
	xmlVlans := `<?xml version="1.0" encoding="UTF-8"?>
  <networkapi versao="1.0">
    <vlan>
      <redeipv4>
        <network>10.10.10.0/24</network>
      </redeipv4>
    </vlan>
</networkapi>`
	var expected ListVlanResult
	err := xml.Unmarshal([]byte(xmlVlans), &expected)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, xmlVlans)
	}))
	defer ts.Close()
	client, err := NewClient(ts.URL, "login", "password")
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.ListVlans(ListVlansOptions{Name: "vlan_name"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected.Vlans) {
		t.Errorf("ListVlans: Expected %#v. Got %#v.", expected, result)
	}
}
