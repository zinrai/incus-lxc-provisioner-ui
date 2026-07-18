package main

// allowed is the UI's authorization gate: it reports whether the caller may act
// in project p. It is load-bearing because the provisioner trusts the group the
// UI sends, so a wrong answer here would let a user reach another tenant's
// project. The caller may name any project, but only ones in their group set are
// honored.
func allowed(projects []string, p string) bool {
	for _, x := range projects {
		if x == p {
			return true
		}
	}
	return false
}
