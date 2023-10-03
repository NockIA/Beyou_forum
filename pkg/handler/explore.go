package handler

import (
	"fmt"
	"forum/pkg/datamanagement"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type DataExplorePage struct {
	Topics       []datamanagement.Topics
	Users        []string
	Upvotes      []string
	CanPrevious  bool
	CanNext      bool
	InvalidTopic bool
	IsConnected  bool
	IsAdmin      bool
	Tags         [][]string
}

// return true if some content of the new topic is forbiden
func createTopic(w http.ResponseWriter, r *http.Request, creatorID string) bool {
	title := r.FormValue("topicTitle")
	description := r.FormValue("description")
	tags := r.FormValue("tags")
	file, handler, err := r.FormFile("photo")
	if title != "" && (datamanagement.GetOneTopicByName(title) == datamanagement.Topics{}) {
		if !datamanagement.CheckContentByBlackListWord(title) && !datamanagement.CheckContentByBlackListWord(description) && !datamanagement.CheckContentByBlackListWord(tags) && len(strings.Split(title, " ")) == 1 {
			return true
		}
		title = strings.Title(title)
		fileName := "../img/PP_wb.png"
		if file != nil && err == nil {
			defer file.Close()
			row := datamanagement.SelectDB("SELECT TopicID FROM Topics WHERE title =?", title)
			defer row.Close()
			var id int
			for row.Next() {
				row.Scan(&id)
			}
			destinationPath := "./static/img/" + title + "." + strings.Split(handler.Filename, ".")[1]
			destinationFile, err := os.Create(destinationPath)
			if err != nil {
				fmt.Println("Failed to create destination file")
			}
			defer destinationFile.Close()
			_, err = io.Copy(destinationFile, file)
			if err != nil {
				fmt.Println("Failed to save photo on server")
			}
			fileName = "../img/" + title + "." + strings.Split(handler.Filename, ".")[1]
		}
		if string(tags[len(tags)-1]) == " " {
			tags = tags[:len(tags)-1]
		}
		topic := datamanagement.DataContainer{Topics: datamanagement.Topics{Title: title, Description: description, Picture: fileName, CreationDate: time.Now(), CreatorID: creatorID, Upvotes: 0, Follows: 0}}
		datamanagement.AddLineIntoTargetTable(topic, "Topics")
		datamanagement.AddTagsToTopic(tags, creatorID, datamanagement.GetOneTopicByName(title).TopicID)
		http.Redirect(w, r, "http://localhost:8080/topic/"+title, http.StatusSeeOther)
	}
	return false
}

func changeUpvote(upvote int) string {
	if upvote < 1000 {
		return fmt.Sprintf("%v", upvote)
	}
	return fmt.Sprintf("%v", upvote/1000.0) + "k"
}

func Explore(w http.ResponseWriter, r *http.Request) {
	cookieUserID, _ := r.Cookie("idUser")
	userId := getCookieValue(cookieUserID)
	dataToSend := DataExplorePage{}
	dataToSend.IsAdmin = false
	dataToSend.InvalidTopic = false
	dataToSend.IsConnected = false
	if userId != "" {
		dataToSend.IsConnected = true
		currentUser := datamanagement.GetUserById(userId)
		dataToSend.IsAdmin = currentUser.IsAdmin
		if createTopic(w, r, userId) {
			dataToSend.InvalidTopic = true
		}
	}
	cookieFilter, _ := r.Cookie("filter")
	cookiePaging, _ := r.Cookie("paging")
	cookieSearch, _ := r.Cookie("search")
	sort := r.FormValue("sort")
	if getCookieValue(cookieFilter) == "" || getCookieValue(cookiePaging) == "" {
		if sort == "" {
			sort = "default"
		}
		cookieFilter = &http.Cookie{Name: "filter", Value: sort}
		cookiePaging = &http.Cookie{Name: "paging", Value: "1"}
		http.SetCookie(w, cookieFilter)
		http.SetCookie(w, cookiePaging)
	}
	if sort != "" {
		cookieFilter = &http.Cookie{Name: "filter", Value: sort}
		http.SetCookie(w, cookieFilter)
		cookieSearch = &http.Cookie{Name: "search", Value: ""}
		http.SetCookie(w, cookieSearch)
		cookiePaging = &http.Cookie{Name: "paging", Value: "1"}
		http.SetCookie(w, cookiePaging)
	}
	topic := r.FormValue("topicSearch")
	if topic != "" {
		if datamanagement.CheckContentByBlackListWord(topic) {
			cookiePaging = &http.Cookie{Name: "paging", Value: "1"}
			cookieSearch = &http.Cookie{Name: "search", Value: topic}
			http.SetCookie(w, cookiePaging)
			http.SetCookie(w, cookieSearch)
			dataToSend.Topics = datamanagement.GetTopicByTagAndTitle(topic)
		} else {
			dataToSend.Topics = nil
		}
	} else if getCookieValue(cookieSearch) != "" {
		dataToSend.Topics = datamanagement.GetTopicByTagAndTitle(getCookieValue(cookieSearch))
	} else if sort == "Follows" && userId != "" {
		dataToSend.Topics = datamanagement.GetTopicsByUser(userId)
	} else {
		if sort == "Follows" {
			sort = "default"
			cookieFilter = &http.Cookie{Name: "filter", Value: "default"}
			http.SetCookie(w, cookieFilter)
		}
		if getCookieValue(cookieFilter) == "Follows" {
			dataToSend.Topics = datamanagement.GetTopicsByUser(userId)
		} else {
			dataToSend.Topics = datamanagement.SortTopics(getCookieValue(cookieFilter))
		}
	}
	prev := r.FormValue("previous")
	next := r.FormValue("next")
	pagingInt, _ := strconv.Atoi(getCookieValue(cookiePaging))
	dataToSend.CanNext = true
	dataToSend.CanPrevious = true
	dataToSend.InvalidTopic = false
	if next != "" && pagingInt*2 < len(dataToSend.Topics) {
		cookiePaging = &http.Cookie{Name: "paging", Value: strconv.Itoa(pagingInt + 1)}
		pagingInt++
		http.SetCookie(w, cookiePaging)
	} else if prev != "" && pagingInt > 1 {
		cookiePaging = &http.Cookie{Name: "paging", Value: strconv.Itoa(pagingInt - 1)}
		pagingInt--
		http.SetCookie(w, cookiePaging)
	}
	if pagingInt == 1 {
		dataToSend.CanPrevious = false
	}
	t := template.Must(template.ParseFiles("./static/html/explore.html", "./static/html/navBar.html"))
	if len(dataToSend.Topics) == 0 {
		t.Execute(w, dataToSend)
	} else if pagingInt*2 >= len(dataToSend.Topics) {
		dataToSend.CanNext = false
		dataToSend.Upvotes = append(dataToSend.Upvotes, changeUpvote(dataToSend.Topics[(pagingInt-1)*2].Upvotes))
		dataToSend.Users = []string{datamanagement.GetUserById(dataToSend.Topics[(pagingInt-1)*2].CreatorID).Username}
		dataToSend.Tags = [][]string{datamanagement.TransformTags(dataToSend.Topics[(pagingInt-1)*2].TopicID)}
		if pagingInt*2 == len(dataToSend.Topics) {
			dataToSend.Tags = append(dataToSend.Tags, datamanagement.TransformTags(dataToSend.Topics[(pagingInt)*2-1].TopicID))
			dataToSend.Upvotes = append(dataToSend.Upvotes, changeUpvote(dataToSend.Topics[(pagingInt)*2-1].Upvotes))
			dataToSend.Users = append(dataToSend.Users, datamanagement.GetUserById(dataToSend.Topics[(pagingInt)*2-1].CreatorID).Username)
		}
		dataToSend.Topics = dataToSend.Topics[(pagingInt-1)*2:]
		t.Execute(w, dataToSend)
	} else {
		dataToSend.Tags = [][]string{datamanagement.TransformTags(dataToSend.Topics[(pagingInt-1)*2].TopicID)}
		dataToSend.Tags = append(dataToSend.Tags, datamanagement.TransformTags(dataToSend.Topics[(pagingInt)*2-1].TopicID))
		dataToSend.Upvotes = append(dataToSend.Upvotes, changeUpvote(dataToSend.Topics[(pagingInt-1)*2].Upvotes))
		dataToSend.Upvotes = append(dataToSend.Upvotes, changeUpvote(dataToSend.Topics[(pagingInt)*2-1].Upvotes))
		dataToSend.Users = []string{datamanagement.GetUserById(dataToSend.Topics[(pagingInt-1)*2].CreatorID).Username}
		dataToSend.Users = append(dataToSend.Users, datamanagement.GetUserById(dataToSend.Topics[(pagingInt)*2-1].CreatorID).Username)
		dataToSend.Topics = dataToSend.Topics[(pagingInt-1)*2 : pagingInt*2]
		t.Execute(w, dataToSend)
	}
}
