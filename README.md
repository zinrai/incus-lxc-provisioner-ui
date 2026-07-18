# incus-lxc-provisioner-ui

A web frontend for [incus-lxc-provisioner](https://github.com/zinrai/incus-lxc-provisioner).
The provisioner is a headless HTTP API by design, so its browser interface lives
here as a separate component. The UI is server-rendered and reads the identity a
front proxy hands it, then calls the provisioner scoped to the caller's project.

## What it does

It reads `X-Forwarded-Email` and `X-Forwarded-Groups` (the identity an auth proxy
asserts), lets the user pick one of their groups as the active Incus project, and
lists, creates, starts, stops, and deletes containers in that project through the
provisioner. Groups are projects: the UI offers only the caller's groups and
narrows `X-Auth-Groups` to the single selected project before each call, so the
provisioner keeps its one-group-to-one-project routing unchanged.

The UI holds no authorization policy of its own beyond that narrowing. Who may use
which project is decided by group membership in the directory, upstream of here.

## Run it for development

When no auth proxy has set the `X-Forwarded-*` headers and `-dev-groups` is given,
the UI uses that identity, so you can open it directly in a browser without any
proxy.

Run the provisioner in its memory backend (from a checkout of
incus-lxc-provisioner), then the UI against it:

```
# incus-lxc-provisioner: in-memory backend, no cluster needed
go run . -backend memory -listen 127.0.0.1:8081

# incus-lxc-provisioner-ui: point at it, with a dev identity
go run . -provisioner-url http://127.0.0.1:8081 \
  -dev-email alice@example.com -dev-groups web,data
```

Then open http://localhost:8080. The development identity is inert in production:
there the auth proxy always sets the headers and `-dev-groups` is not given.
Because the UI trusts the headers, it must never be reachable except through that
proxy.

## Run it for real

Behind an auth proxy, pointing at a live provisioner:

```
incus-lxc-provisioner-ui -listen :8080 -provisioner-url http://<provisioner-host>:8080
```

The [incus-lxc-provisioner-auth-lab](https://github.com/zinrai/incus-lxc-provisioner-auth-lab)
repository wires oauth2-proxy, Dex, and glauth in front of this UI and points it
at a real provisioner on an Incus cluster.

## License

This project is licensed under the [MIT License](LICENSE).
