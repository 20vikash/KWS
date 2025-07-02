package main

import (
	"html/template"
	"kws/kws/consts/config"
	"kws/kws/consts/services"
	"kws/kws/models"
	"kws/kws/models/web"
	"log"
	"net/http"
	"strconv"
)

type Device struct {
	PublicKey string
	IP        string
	Active    bool
}

type Data struct {
	Username string
	Devices  []Device
}

type ServicesData struct {
	Username string
	Services []services.WebService
}

var templates = template.Must(template.ParseGlob("../web/*.html"))

func (app *Application) RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	isAuthorized := app.SessionManager.GetBool(r.Context(), "isAuthorized")
	if isAuthorized {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := templates.ExecuteTemplate(w, "register", nil)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func (app *Application) RenderSignInPage(w http.ResponseWriter, r *http.Request) {
	isAuthorized := app.SessionManager.GetBool(r.Context(), "isAuthorized")
	if isAuthorized {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := templates.ExecuteTemplate(w, "signin", nil)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	data := struct {
		Username string
	}{
		Username: userName,
	}

	err := templates.ExecuteTemplate(w, "home", data)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func (app *Application) RenderDevicesPage(w http.ResponseWriter, r *http.Request) {
	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	// Fetch peers information
	peers, err := app.Store.Wireguard.GetDevices(r.Context(), uid)
	if err != nil {
		http.Error(w, "cannot load this page right now", http.StatusInternalServerError)
	}

	var data = new(Data)

	for _, peer := range peers {
		isActive, err := app.Wg.IsOnline(peer.PublicKey)
		if err != nil {
			log.Println("Cannot find online status")
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		ipAddress := app.IpAlloc.GenerateIP(peer.IpAddress)
		data.Devices = append(data.Devices, Device{PublicKey: peer.PublicKey, IP: ipAddress, Active: isActive})
	}

	data.Username = userName

	err = templates.ExecuteTemplate(w, "devices", data)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (app *Application) RenderServicesPage(w http.ResponseWriter, r *http.Request) {
	// Get services list
	services := services.GetServiceList()

	// Get the username
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	webServices := ServicesData{
		Username: userName,
		Services: services,
	}

	err := templates.ExecuteTemplate(w, "services", webServices)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (app *Application) RenderPgUsersPage(w http.ResponseWriter, r *http.Request) {
	// Get uid and username
	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	// Get all the users list
	users, err := app.Store.PgService.GetUsers(r.Context(), uid)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	pg := services.GetPgServiceData()

	pgData := web.PGUserPageData{
		HostName:     pg.Hostname,
		ServiceIP:    pg.IP,
		Port:         "5432",
		LoggedInUser: userName,
		UserLimit:    config.MAX_SERVICE_DB_USERS,
		Users:        users,
	}

	err = templates.ExecuteTemplate(w, "pgusers", pgData)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (app *Application) RenderPgDatabasesPage(w http.ResponseWriter, r *http.Request) {
	// Get the pid from the query
	pidStr := r.URL.Query().Get("pid")
	owner := r.URL.Query().Get("owner")
	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	// Convert to int
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Invalid pid", http.StatusBadRequest)
		return
	}

	count, dbs, err := app.Store.PgService.GetDatabases(r.Context(), pid, uid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	pg := services.GetPgServiceData()

	password, err := app.Store.PgService.GetPassword(r.Context(), pid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	pgDB := web.PgDatabase{
		Password:       password,
		HostName:       pg.Hostname,
		Username:       userName,
		Owner:          owner,
		TotalDatabases: count,
		Limit:          config.MAX_SERVICE_DB_DB,
		AvailableSlots: config.MAX_SERVICE_DB_DB - count,
		Databases:      dbs,
		UsagePercent:   int(float64(count) / float64(config.MAX_SERVICE_DB_DB) * 100),
	}

	err = templates.ExecuteTemplate(w, "db_management", pgDB)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (app *Application) RenderInstancePage(w http.ResponseWriter, r *http.Request) {
	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	data, err := app.Store.Instance.GetData(r.Context(), uid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	data.Username = userName

	ip, err := app.Docker.FindContainerIP(r.Context(), data.ContainerName)
	if err != nil {
		data.Instance.IP = ""
	} else {
		data.Instance.IP = ip
	}

	containerID, err := app.Docker.GetContainerIDByName(r.Context(), data.ContainerName)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if containerID != "" {
		containerID = containerID[:15]
		data.ContainerName = containerID
	}

	err = templates.ExecuteTemplate(w, "instance_management", data)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func (app *Application) RenderPublishPage(w http.ResponseWriter, r *http.Request) {
	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	// Get domains data
	domains, err := app.Store.Domains.GetUserDomains(r.Context(), &models.Domain{Uid: uid})
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	hasDomains := len(*domains) > 0

	pubinsData := web.PublishInstancePageData{
		LoggedInUser: userName,
		Domains:      *domains,
		HasDomains:   hasDomains,
	}

	err = templates.ExecuteTemplate(w, "publish_instance", pubinsData)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}
