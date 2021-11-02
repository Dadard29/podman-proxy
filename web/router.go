package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Dadard29/podman-proxy/models"
	"github.com/gorilla/mux"
)

const basePath = "web/templates/"

type mainData struct {
	NotLogged bool
}

func getTemplateFile(name string) string {
	return fmt.Sprintf("%s%s.html", basePath, name)
}

func renderTemplates(templateList []string, w http.ResponseWriter, data interface{}) {
	tmpl := template.Must(template.ParseFiles(templateList...))
	tmpl.Execute(w, data)
}

func isLogged(r *http.Request) bool {
	_, err := models.NewAccessTokenFromCookie(r)
	if err != nil {
		return false
	}

	// fixme
	return true

}

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		getTemplateFile("index"),
		getTemplateFile("main"),
		getTemplateFile("login"),
	}

	renderTemplates(templates, w, nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// username := r.Form["username"]
	// password := r.Form["password"]

}

func handleContent(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		getTemplateFile("index"),
		getTemplateFile("main"),
		getTemplateFile("content"),
	}

	renderTemplates(templates, w, nil)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	if isLogged(r) {
		handleContent(w, r)
	} else {
		handleLoginPage(w, r)
	}
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", handleMain).Methods(http.MethodGet)
	r.HandleFunc("/login", handleLogin).Methods(http.MethodPost)
}
