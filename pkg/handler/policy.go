package handler

import (
	"forum/pkg/datamanagement"
	"net/http"
	"text/template"
)

type informationConnection struct {
	IsAdmin     bool
	IsConnected string
}

func Policy(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./static/html/policy.html", "./static/html/navBar.html"))
	dataPolicy := informationConnection{}
	cookieConnected, _ := r.Cookie("isConnected")
	IsConnected := getCookieValue(cookieConnected)
	dataPolicy.IsConnected = IsConnected
	cookieIdUser, _ := r.Cookie("idUser")
	dataPolicy.IsAdmin = datamanagement.IsAdmin(getCookieValue(cookieIdUser))
	t.ExecuteTemplate(w, "policy", dataPolicy)
}
