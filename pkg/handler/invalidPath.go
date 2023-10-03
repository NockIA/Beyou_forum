package handler

import (
	"forum/pkg/datamanagement"
	"net/http"
	"text/template"
)

type invalidPath struct {
	IsConnected bool
	IsAdmin     bool
}

func InvalidPath(w http.ResponseWriter, r *http.Request) {
	cookieConnected, _ := r.Cookie("isConnected")
	IsConnected := getCookieValue(cookieConnected)
	invalidPath := invalidPath{}
	invalidPath.IsAdmin = false
	if IsConnected == "true" {
		invalidPath.IsConnected = true
		cookieIdUser, _ := r.Cookie("idUser")
		currentUser := datamanagement.GetUserById(getCookieValue(cookieIdUser))
		invalidPath.IsAdmin = currentUser.IsAdmin
	}
	t := template.Must(template.ParseFiles("./static/html/invalidPath.html", "./static/html/navBar.html"))
	t.ExecuteTemplate(w, "errorPath", invalidPath)
}
