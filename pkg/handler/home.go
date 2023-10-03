package handler

import (
	"fmt"
	"forum/pkg/datamanagement"
	"net/http"
	"text/template"
	"time"
)

func getCookieValue(cookie *http.Cookie) string {
	var valueReturned string
	test := false
	value := fmt.Sprint(cookie)
	for _, element := range value {
		if test {
			valueReturned += string(element)
		}
		if element == 61 {
			test = true
		}
	}
	return valueReturned
}

type TopicsWithUserInfo struct {
	TopicID      int
	Title        string
	Description  string
	Picture      string
	CreatorID    string
	Upvotes      int
	Follows      int
	CreationDate string
	CreatorName  string
}

type structDisplayHome struct {
	AllTopics         []TopicsWithUserInfo
	HasAcceptedCookie string
	IsConnected       bool
	IsAdmin           bool
}

func updateTopicsInTopicsWithUserInfo(topics []datamanagement.Topics) []TopicsWithUserInfo {
	result := []TopicsWithUserInfo{}
	for _, element := range topics {
		var topic TopicsWithUserInfo
		topic.TopicID = element.TopicID
		topic.Title = element.Title
		topic.Description = element.Description
		topic.Picture = element.Picture
		topic.CreatorID = element.CreatorID
		topic.Upvotes = element.Upvotes
		topic.Follows = element.Follows
		topic.CreationDate = fmt.Sprint(element.CreationDate.Day()) + " " + fmt.Sprint(element.CreationDate.Month()) + " " + fmt.Sprint(element.CreationDate.Year())
		user := datamanagement.GetUserById(topic.CreatorID)
		topic.CreatorName = user.Username
		result = append(result, topic)
	}
	return result
}

func Home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./static/html/home.html", "./static/html/navBar.html", "./static/html/cookiesPopup.html"))
	allTop := datamanagement.SortTopics("DESC-Upvote-Home")
	structDisplayHome := structDisplayHome{}
	structDisplayHome.AllTopics = updateTopicsInTopicsWithUserInfo(allTop)
	cookieConnected, _ := r.Cookie("isConnected")
	IsConnected := getCookieValue(cookieConnected)
	if IsConnected == "true" {
		structDisplayHome.IsConnected = true
	} else {
		structDisplayHome.IsConnected = false
	}
	structDisplayHome.IsAdmin = false
	getAcceptCookies := r.FormValue("AcceptCookies")
	if getAcceptCookies == "OK" {
		cookieHasAccepted := http.Cookie{Name: "hasAcceptedCookies", Value: getAcceptCookies, Expires: time.Now().Add(24 * time.Hour)}
		http.SetCookie(w, &cookieHasAccepted)
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
	}
	cookieHasAcceptedCookies, _ := r.Cookie("hasAcceptedCookies")
	CAcceptCookies := getCookieValue(cookieHasAcceptedCookies)
	structDisplayHome.HasAcceptedCookie = CAcceptCookies
	if IsConnected == "true" {
		cookieIdUser, _ := r.Cookie("idUser")
		currentUser := datamanagement.GetUserById(getCookieValue(cookieIdUser))
		structDisplayHome.IsAdmin = currentUser.IsAdmin
	}
	t.ExecuteTemplate(w, "home", structDisplayHome)
}
