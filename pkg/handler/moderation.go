package handler

import (
	"fmt"
	"forum/pkg/datamanagement"
	"net/http"
	"text/template"
)

type ReportPostInformations struct {
	ReportID       int
	Reason         string
	UserID         string
	Username       string
	ProfilePicture string
	PostID         int
	Message        string
}

type ReportTopicInformations struct {
	ReportID    int
	Reason      string
	UserID      string
	Username    string
	TopicID     int
	Title       string
	Description string
	Picture     string
}

type ModerationInformations struct {
	ReportsPostInformations  []ReportPostInformations
	ReportsTopicInformations []ReportTopicInformations
	BannedUsers              []datamanagement.Users
	WordsBlacklist           []datamanagement.WordsBlacklist
	IsConnected              bool
	IsAdmin                  bool
}

func GetReportsPostInformations() []ReportPostInformations {
	rows := datamanagement.SelectDB("SELECT Reports.ReportID, Posts.PostID, Reports.Comment, Users.UserID, Users.Username, Users.ProfilePicture, Posts.Content FROM Reports JOIN Users ON Reports.ReportUserID = Users.UserID JOIN Posts ON Reports.PostID = Posts.PostID;")
	defer rows.Close()

	var reportsPosts []ReportPostInformations
	for rows.Next() {
		var reportPost ReportPostInformations
		rows.Scan(&reportPost.ReportID, &reportPost.PostID, &reportPost.Reason, &reportPost.UserID, &reportPost.Username, &reportPost.ProfilePicture, &reportPost.Message)
		reportsPosts = append(reportsPosts, reportPost)
	}

	return reportsPosts
}

func GetReportsTopicInformations() []ReportTopicInformations {
	rows := datamanagement.SelectDB("SELECT Reports.ReportID, Reports.Comment, Users.UserID, Users.Username, Topics.TopicID, Topics.Title, Topics.Description, Topics.Picture FROM Reports JOIN Users ON Users.UserID = Reports.ReportUserID JOIN Topics ON Topics.TopicID = Reports.TopicID;")
	defer rows.Close()

	var reportsTopics []ReportTopicInformations
	for rows.Next() {
		var reportTopic ReportTopicInformations
		rows.Scan(&reportTopic.ReportID, &reportTopic.Reason, &reportTopic.UserID, &reportTopic.Username, &reportTopic.TopicID, &reportTopic.Title, &reportTopic.Description, &reportTopic.Picture)
		reportsTopics = append(reportsTopics, reportTopic)
	}

	return reportsTopics
}

func Moderation(w http.ResponseWriter, r *http.Request) {
	cookieIdUser, _ := r.Cookie("idUser")
	if !datamanagement.IsAdmin(getCookieValue(cookieIdUser)) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	t := template.Must(template.ParseFiles("./static/html/moderation.html", "./static/html/navBar.html"))
	banPost := r.FormValue("banUser")
	deletePost := r.FormValue("deletePost")
	deleteTopic := r.FormValue("deleteTopic")
	removeReport := r.FormValue("removeReport")
	unbanUser := r.FormValue("unbanUser")
	addWord := r.FormValue("addWord")
	deleteWord := r.FormValue("deleteWord")
	if banPost != "" {
		datamanagement.BanUser(banPost)
	} else if deletePost != "" {
		datamanagement.DeletePost(deletePost)
		datamanagement.DeleteReportsFromPost(deletePost)
	} else if deleteTopic != "" {
		datamanagement.DeletePostsFromTopic(deleteTopic)
		datamanagement.DeleteTopic(deleteTopic)
		datamanagement.DeleteReportsFromTopic(deleteTopic)
	} else if removeReport != "" {
		datamanagement.DeleteReport(removeReport)
	} else if unbanUser != "" {
		datamanagement.SetUserStatus(unbanUser, true)
	} else if addWord != "" {
		datamanagement.AddWordIntoBlacklist(addWord)
	} else if deleteWord != "" {
		res := datamanagement.AddDeleteUpdateDB("DELETE FROM WordsBlacklist WHERE Word = ?", deleteWord)
		_, err := res.RowsAffected()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(deleteWord, "deleted!")
	}

	t.ExecuteTemplate(w, "moderation", ModerationInformations{
		ReportsPostInformations:  GetReportsPostInformations(),
		ReportsTopicInformations: GetReportsTopicInformations(),
		BannedUsers:              datamanagement.GetAllBannedUsers(),
		WordsBlacklist:           datamanagement.GetAllBlacklistWords(),
		IsConnected:              true,
		IsAdmin:                  true,
	})
}
