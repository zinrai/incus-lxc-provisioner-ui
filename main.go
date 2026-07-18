// Command incus-lxc-provisioner-ui is the web application in front of
// incus-lxc-provisioner. It is server-rendered, reads the identity an auth proxy
// passes, lets the user pick one of their projects, and calls the provisioner
// scoped to that project.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// version, commit and date are set by goreleaser via -ldflags -X.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:8080", "address to listen on")
	provisioner := flag.String("provisioner-url", "http://127.0.0.1:8080", "base URL of the incus-lxc-provisioner")
	devEmail := flag.String("dev-email", "", "development identity email, used only when no proxy set the X-Forwarded-* headers")
	devGroups := flag.String("dev-groups", "", "development identity groups, comma-separated, used only when no proxy set the X-Forwarded-* headers")
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("incus-lxc-provisioner-ui version %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	provisionerBase = *provisioner
	devIdentityEmail = *devEmail
	devIdentityGroups = *devGroups

	http.HandleFunc("/create", handleCreate)
	http.HandleFunc("/action", handleAction)
	http.HandleFunc("/", handleIndex)

	log.Printf("incus-lxc-provisioner-ui listening on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
