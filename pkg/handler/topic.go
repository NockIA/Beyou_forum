package handler

import (
	"forum/pkg/datamanagement"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type TopicsDate struct {
	TopicID      int
	Title        string
	Description  string
	Picture      string
	CreationDate string
	CreatorID    string
	Upvotes      int
	Follows      int
}

type DataTopicPage struct {
	Topic       TopicsDate
	Tags        []datamanagement.Tags
	Posts       []PostInTopicPage
	IsFollow    bool
	IsUpvote    bool
	IsConnected string
	IsAdmin     bool
	IsNotValid  bool
}

type PostInTopicPage struct {
	PostID                   int
	Content                  string
	Likes                    int
	Dislikes                 int
	StructuredDate           string
	IsLikeByConnectedUser    bool
	IsDislikeByConnectedUser bool
	ProfilePicture           string
	AuthorName               string
}

func transformPostInPostInTopicPage(posts []datamanagement.Posts, userID string) []PostInTopicPage {
	result := []PostInTopicPage{}
	for _, element := range posts {
		var post PostInTopicPage
		post.Content = element.Content
		post.PostID = element.PostID
		post.Likes = element.Likes
		post.Dislikes = element.Dislikes
		post.StructuredDate = datamanagement.TransformDateInPostFormat(element.CreationDate)
		user := datamanagement.GetUserById(element.AuthorID)
		post.ProfilePicture = user.ProfilePicture
		post.AuthorName = user.Username
		if datamanagement.IsLikeByUser(userID, post.PostID) {
			post.IsLikeByConnectedUser = true
			post.IsDislikeByConnectedUser = false
		} else if datamanagement.IsDislikeByUser(userID, post.PostID) {
			post.IsDislikeByConnectedUser = true
			post.IsLikeByConnectedUser = false
		}
		result = append(result, post)
	}
	return result
}

func isFollowTopic(topicName string, topicDisplayStruct DataTopicPage, idUser string) DataTopicPage {
	if idUser != "" {
		rows := datamanagement.SelectDB("SELECT * FROM Follows WHERE UserID LIKE ? AND TopicID LIKE ?;", idUser, strconv.Itoa(topicDisplayStruct.Topic.TopicID))
		defer rows.Close()
		if !rows.Next() {
			topicDisplayStruct.IsFollow = false
		} else {
			topicDisplayStruct.IsFollow = true
		}
	}
	return topicDisplayStruct
}

func isUpvoteTopic(topicID int, idUser string) bool {
	if idUser != "" {
		rows := datamanagement.SelectDB("SELECT TopicID FROM Upvotes WHERE UserID = ? AND TopicID = ?;", idUser, topicID)
		defer rows.Close()
		if rows.Next() {
			return true
		}
	}
	return false
}
func Topic(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.String(), "/")
	topicName := url[2]
	cookie, _ := r.Cookie("idUser")
	idUser := getCookieValue(cookie)
	newPost := r.FormValue("postContent")
	clickFollow := r.FormValue("follow")
	clickUpvote := r.FormValue("upvote")
	unLike := r.FormValue("unLike")
	unDislike := r.FormValue("unDislike")
	like := r.FormValue("like")
	dislike := r.FormValue("dislike")
	reportPostID := r.FormValue("reportPostID")
	reportTopicID := r.FormValue("reportTopicID")
	reportReason := r.FormValue("reportReason")
	cookieIsConnected, _ := r.Cookie("isConnected")
	isConnected := getCookieValue(cookieIsConnected)
	topicDisplayStruct := DataTopicPage{}
	topicDisplayStruct.IsNotValid = false
	topic := datamanagement.GetOneTopicByName(topicName)
	topicDisplayStruct.Topic = TopicsDate{
		TopicID:      topic.TopicID,
		Title:        topic.Title,
		Description:  topic.Description,
		Picture:      topic.Picture,
		CreationDate: datamanagement.TransformDateInPostFormat(topic.CreationDate),
		CreatorID:    topic.CreatorID,
		Upvotes:      topic.Upvotes,
		Follows:      topic.Follows,
	}
	topicDisplayStruct.Tags = datamanagement.GetTagsByTopic(topicDisplayStruct.Topic.TopicID)
	topicDisplayStruct = isFollowTopic(topicName, topicDisplayStruct, idUser)
	topicDisplayStruct.IsUpvote = isUpvoteTopic(topic.TopicID, idUser)
	topicDisplayStruct.IsAdmin = false
	if isConnected == "true" {
		cookieIdUser, _ := r.Cookie("idUser")
		currentUser := datamanagement.GetUserById(getCookieValue(cookieIdUser))
		topicDisplayStruct.IsAdmin = currentUser.IsAdmin
		topicDisplayStruct.IsConnected = isConnected
		switch true {
		case len(newPost) > 0 && len(newPost) <= 500:
			if datamanagement.CheckContentByBlackListWord(newPost) {
				post := datamanagement.DataContainer{Posts: datamanagement.Posts{Content: newPost, AuthorID: idUser, TopicID: topicDisplayStruct.Topic.TopicID, Likes: 0, Dislikes: 0, CreationDate: time.Now()}}
				datamanagement.AddLineIntoTargetTable(post, "Posts")
				http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
			} else {
				topicDisplayStruct.IsNotValid = true
			}
			break
		case clickFollow != "":
			if topicDisplayStruct.IsFollow {
				datamanagement.AddDeleteUpdateDB("UPDATE Topics SET Follows = Follows - 1 WHERE TopicID = ?;", topicDisplayStruct.Topic.TopicID)
				datamanagement.AddDeleteUpdateDB("DELETE FROM Follows WHERE UserID = ?;", idUser)
				topicDisplayStruct.IsFollow = false
			} else {
				line := datamanagement.DataContainer{Follows: datamanagement.Follows{TopicID: topicDisplayStruct.Topic.TopicID, UserID: idUser}}
				datamanagement.AddDeleteUpdateDB("UPDATE Topics SET Follows = Follows + 1 WHERE TopicID = ?;", topicDisplayStruct.Topic.TopicID)
				datamanagement.AddLineIntoTargetTable(line, "Follows")
				topicDisplayStruct.IsFollow = true
			}
			break
		case clickUpvote != "":
			if topicDisplayStruct.IsUpvote {
				topicDisplayStruct.IsUpvote = false
				datamanagement.UnUpvotesTopic(topicDisplayStruct.Topic.TopicID, idUser)
			} else {
				topicDisplayStruct.IsUpvote = true
				datamanagement.UpvotesTopic(topicDisplayStruct.Topic.TopicID, idUser)
			}
			break
		case like != "":
			idPost, _ := strconv.Atoi(like)
			datamanagement.LikePostManager(idPost, idUser, "Likes")
			break
		case dislike != "":
			idPost, _ := strconv.Atoi(dislike)
			datamanagement.LikePostManager(idPost, idUser, "Dislikes")
			break
		case unLike != "":
			idPost, _ := strconv.Atoi(unLike)
			datamanagement.UnLikePostManager(idPost, idUser, "unLike")
			break
		case unDislike != "":
			idPost, _ := strconv.Atoi(unDislike)
			datamanagement.UnLikePostManager(idPost, idUser, "unDislike")
			break
		case reportPostID != "" && reportReason != "":
			datamanagement.AddPostReport(reportPostID, reportReason)
			break
		case reportTopicID != "" && reportReason != "":
			datamanagement.AddTopicReport(reportTopicID, reportReason)
			break
		}
	}
	topicDisplayStruct.Posts = transformPostInPostInTopicPage(datamanagement.GetPostByTopic(datamanagement.GetTopicId(topicName)), idUser)
	t := template.Must(template.ParseFiles("./static/html/topic.html", "./static/html/navBar.html"))
	t.ExecuteTemplate(w, "topic", topicDisplayStruct)
}
