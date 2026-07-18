package main

import "testing"

// allowed is the UI's authorization gate. If it breaks, a user can reach another
// tenant's project, so it is the one piece of this package worth a unit test.
func TestAllowed(t *testing.T) {
	cases := []struct {
		name     string
		projects []string
		p        string
		want     bool
	}{
		{"member of the project", []string{"web", "data"}, "web", true},
		{"member of a second project", []string{"web", "data"}, "data", true},
		{"not a member", []string{"web", "data"}, "sandbox", false},
		{"no projects at all", nil, "web", false},
		{"empty target project", []string{"web"}, "", false},
	}
	for _, c := range cases {
		if got := allowed(c.projects, c.p); got != c.want {
			t.Errorf("%s: allowed(%v, %q) = %v, want %v", c.name, c.projects, c.p, got, c.want)
		}
	}
}
