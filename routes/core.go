package routes

import (
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/index.html", "web/base.layout.html")

	tmpl.Execute(w, nil)
}
