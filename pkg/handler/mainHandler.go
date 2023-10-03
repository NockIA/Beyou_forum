package handler

import (
	"forum/pkg/datamanagement"
	"net/http"
	"strings"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.String(), "/")
	switch true {
	case url[1] == "" && len(url) == 2:
		Home(w, r)
	case url[1] == "account" && len(url) == 2:
		Account(w, r)
	case url[1] == "moderation" && len(url) == 2:
		Moderation(w, r)
	case url[1] == "explore" && len(url) == 2:
		Explore(w, r)
	case url[1] == "login" && len(url) == 2:
		Login(w, r)
	case url[1] == "policy" && len(url) == 2:
		Policy(w, r)
	case url[1] == "register" && len(url) == 2:
		Register(w, r)
	case url[1] == "profile" && len(url) == 2:
		Profile(w, r, true)
	case url[1] == "profile" && len(url) == 3 && datamanagement.IsUsernameAlreadyExist(url[2]):
		Profile(w, r, false)
	case url[1] == "topic" && len(url) == 3 && datamanagement.IsValidTopic(url[2]):
		Topic(w, r)
	default:
		if len(url) >= 3 {
			http.Redirect(w, r, "http://localhost:8080/InvalidPath", http.StatusSeeOther)
		}
		InvalidPath(w, r)
	}
}
