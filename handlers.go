package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type pageData struct {
	Email      string
	Projects   []string
	Project    string
	Containers []Container
	Error      string
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	email, projects := identity(r)

	project := r.URL.Query().Get("project")
	if project == "" && len(projects) > 0 {
		project = projects[0]
	}

	data := pageData{Email: email, Projects: projects, Project: project}

	if project != "" {
		if !allowed(projects, project) {
			http.Error(w, "forbidden project", http.StatusForbidden)
			return
		}
		cs, err := listContainers(email, project)
		if err != nil {
			data.Error = err.Error()
		}
		data.Containers = cs
	}

	render(w, data)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	email, projects := identity(r)
	project := r.FormValue("project")
	if !allowed(projects, project) {
		http.Error(w, "forbidden project", http.StatusForbidden)
		return
	}
	// The create body is name+image only; the provisioner rejects anything else.
	body, _ := json.Marshal(map[string]string{"name": r.FormValue("name"), "image": r.FormValue("image")})
	resp, err := call(http.MethodPost, "/containers", email, project, bytes.NewReader(body))
	if err == nil {
		resp.Body.Close()
	}
	redirect(w, r, project)
}

func handleAction(w http.ResponseWriter, r *http.Request) {
	email, projects := identity(r)
	project := r.FormValue("project")
	if !allowed(projects, project) {
		http.Error(w, "forbidden project", http.StatusForbidden)
		return
	}
	name := r.FormValue("name")
	var method, path string
	switch r.FormValue("action") {
	case "start":
		method, path = http.MethodPost, "/containers/"+url.PathEscape(name)+"/start"
	case "stop":
		method, path = http.MethodPost, "/containers/"+url.PathEscape(name)+"/stop"
	case "delete":
		method, path = http.MethodDelete, "/containers/"+url.PathEscape(name)
	default:
		http.Error(w, "bad action", http.StatusBadRequest)
		return
	}
	resp, err := call(method, path, email, project, nil)
	if err == nil {
		resp.Body.Close()
	}
	redirect(w, r, project)
}

func redirect(w http.ResponseWriter, r *http.Request, project string) {
	http.Redirect(w, r, "/?project="+url.QueryEscape(project), http.StatusSeeOther)
}

func render(w http.ResponseWriter, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("template: %v", err)
	}
}
