package main

import (
	"net/http"
	"strings"
)

// devIdentityEmail and devIdentityGroups are a development identity, set from the
// -dev-email and -dev-groups flags. When no auth proxy has set the X-Forwarded-*
// headers and the dev groups are non-empty, identity uses them, so the UI can be
// opened directly in a browser. They are empty in production.
var (
	devIdentityEmail  string
	devIdentityGroups string
)

// identity returns the caller's email and the projects they may use, as asserted
// by the auth proxy via X-Forwarded-* headers. The UI is only reachable through
// the proxy, which overwrites client-supplied values, so these cannot be spoofed.
// With no proxy and a dev identity configured, that identity is used instead.
func identity(r *http.Request) (email string, projects []string) {
	email = r.Header.Get("X-Forwarded-Email")
	projects = splitGroups(r.Header.Get("X-Forwarded-Groups"))
	if email == "" && len(projects) == 0 && devIdentityGroups != "" {
		email = devIdentityEmail
		projects = splitGroups(devIdentityGroups)
	}
	return email, projects
}

func splitGroups(s string) []string {
	var out []string
	for _, g := range strings.Split(s, ",") {
		if g = strings.TrimSpace(g); g != "" {
			out = append(out, g)
		}
	}
	return out
}
