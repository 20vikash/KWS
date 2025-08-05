package main

import (
	"encoding/json"
	"kws/kws/consts/config"
	"kws/kws/internal/nginx"
	"kws/kws/models"
	"log"
	"net/http"
	"strconv"
)

type DomainResponse struct {
	Domain string
	Port   string
	Status string
}

func (app *Application) AddUserDomain(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form")
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}

	domain := r.FormValue("domain_name")
	portStr := r.FormValue("port")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		http.Error(w, "Invalid port number", http.StatusBadRequest)
		return
	}

	uid := app.SessionManager.GetInt(r.Context(), "id")

	err = app.Store.Domains.AddUserDomain(r.Context(), &models.Domain{Domain: domain, Port: port, Uid: uid})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	ipInt, err := app.Store.Instance.GetIPFromUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	nginxTemplate := nginx.Template{
		Domain: domain,
		IP:     app.IpAlloc.GenerateIPLXC(ipInt),
		Port:   portStr,
	}

	err = nginxTemplate.AddNewConf()
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	domainResponse := DomainResponse{
		Domain: domain,
		Port:   portStr,
		Status: "Active",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domainResponse)
}

func (app *Application) RemoveUserDomain(w http.ResponseWriter, r *http.Request) {
	uid := app.SessionManager.GetInt(r.Context(), "id")

	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form")
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}

	domain := r.FormValue("domain_name")

	err = app.Store.Domains.RemoveDomain(r.Context(), &models.Domain{Uid: uid, Domain: domain})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	nginxTemplate := nginx.Template{
		Domain: domain,
	}

	err = nginxTemplate.RemoveConf()
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
