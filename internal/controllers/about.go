package controllers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/gorilla/csrf"
)

func About(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/layouts/main.html", "web/templates/components/core.html", "web/templates/pages/about.html")

	if err != nil {
		log.Fatal(err)
	}

	data := models.Site{
		AppName:  "HTMX + GO",
		Title:    "About",
		CSRF:     csrf.Token(r),
		Metatags: models.SEO{Description: "Basic boilerplate for go web apps", Keywords: "go,htmx,web"},
		Year:     time.Now().Year(),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
