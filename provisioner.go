package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// provisionerBase is the base URL of the incus-lxc-provisioner, set from
// PROVISIONER_URL in main.
var provisionerBase string

// Container mirrors the real incus-lxc-provisioner list element. It reports a
// live status and ipv4 read back from Incus, and deliberately carries no image
// (the provisioner does not read config back).
type Container struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	IPv4      string `json:"ipv4"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
}

func listContainers(email, project string) ([]Container, error) {
	resp, err := call(http.MethodGet, "/containers", email, project, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provisioner returned %s", resp.Status)
	}
	var cs []Container
	if err := json.NewDecoder(resp.Body).Decode(&cs); err != nil {
		return nil, err
	}
	return cs, nil
}

// call is the trust hand-off: the UI sets X-Auth-Email and X-Auth-Groups from
// the authenticated identity, scoping the request to the single chosen project.
func call(method, path, email, project string, body *bytes.Reader) (*http.Response, error) {
	var r *http.Request
	var err error
	if body != nil {
		r, err = http.NewRequest(method, provisionerBase+path, body)
	} else {
		r, err = http.NewRequest(method, provisionerBase+path, nil)
	}
	if err != nil {
		return nil, err
	}
	r.Header.Set("X-Auth-Email", email)
	r.Header.Set("X-Auth-Groups", project)
	if body != nil {
		r.Header.Set("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(r)
}
