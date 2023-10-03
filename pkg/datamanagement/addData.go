package datamanagement

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func AddLineIntoTargetTable(data DataContainer, table string) {
	var res sql.Result
	switch true {
	case table == "Users":
		res = AddDeleteUpdateDB("INSERT INTO Users (UserID, Username,Email,Password,Firstname,Lastname,Description,CreationDate,ProfilePicture,IsAdmin,ValidUser) VALUES (?,?,?,?,?,?,?,?,?,?,?);", data.Users.UserID, data.Users.Username, data.Users.Email, data.Users.Password, data.Users.Firstname, data.Users.Lastname, data.Users.Description, data.Users.CreationDate, data.Users.ProfilePicture, data.Users.IsAdmin, data.Users.ValidUser)
		break
	case table == "Posts":
		res = AddDeleteUpdateDB("INSERT INTO Posts (Content,AuthorID,TopicID,Likes,Dislikes,CreationDate) VALUES(?,?,?,?,?,?);", data.Posts.Content, data.Posts.AuthorID, data.Posts.TopicID, data.Posts.Likes, data.Posts.Dislikes, data.Posts.CreationDate)
		break
	case table == "Topics":
		res = AddDeleteUpdateDB("INSERT INTO Topics (Title,Description,Picture,CreationDate,CreatorID,Upvotes,Follows) VALUES(?,?,?,?,?,?,?);", data.Topics.Title, data.Topics.Description, data.Topics.Picture, data.Topics.CreationDate, data.Topics.CreatorID, data.Topics.Upvotes, data.Topics.Follows)
		break
	case table == "Tags":
		res = AddDeleteUpdateDB("INSERT INTO Tags (Title,CreatorID) VALUES(?,?);", data.Tags.Title, data.Tags.CreatorID)
		break
	case table == "WordsBlacklist":
		res = AddDeleteUpdateDB("INSERT INTO WordsBlacklist (word) VALUES(?);", data.WordsBlacklist.Word)
		break
	case table == "Reports":
		res = AddDeleteUpdateDB("INSERT INTO Reports (PostID,ReportUserID,Comment) VALUES (?,?,?);", data.Reports.PostID, data.Reports.ReportUserID, data.Reports.Comment)
		break
	case table == "Follows":
		res = AddDeleteUpdateDB("INSERT INTO Follows (TopicID,UserID) VALUES (?,?);", data.Follows.TopicID, data.Follows.UserID)
		break
	case table == "Likes":
		res = AddDeleteUpdateDB("INSERT INTO Likes (PostID,UserID) VALUES (?,?);", data.Likes.PostID, data.Likes.UserID)
		break
	case table == "Dislikes":
		res = AddDeleteUpdateDB("INSERT INTO Dislikes (PostID,UserID) VALUES (?,?);", data.Dislikes.PostID, data.Dislikes.UserID)
		break
	case table == "TopicsTags":
		res = AddDeleteUpdateDB("INSERT INTO TopicsTags (TopicID,TagID) VALUES (?,?);", data.TopicsTags.TopicID, data.TopicsTags.TagID)
		break
	case table == "Upvotes":
		res = AddDeleteUpdateDB("INSERT INTO Upvotes (TopicID,UserID) VALUES (?,?);", data.Upvotes.TopicID, data.Upvotes.UserID)
		break
	default:
		fmt.Println("Invalid Table")
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(affected, " ", table, " has been add to the database")
}

func UpvotesTopic(topicID int, idUser string) {
	AddDeleteUpdateDB("UPDATE Topics SET Upvotes = Upvotes + 1 WHERE TopicID = ?;", topicID)
	AddLineIntoTargetTable(DataContainer{Upvotes: Upvotes{TopicID: topicID, UserID: idUser}}, "Upvotes")
}

func UnUpvotesTopic(topicID int, idUser string) {
	AddDeleteUpdateDB("UPDATE Topics SET Upvotes = Upvotes - 1 WHERE TopicID = ?;", topicID)
	AddDeleteUpdateDB("DELETE FROM Upvotes WHERE TopicID = ? AND UserID = ?;", fmt.Sprint(topicID), fmt.Sprint(idUser))
}

/*
likOrdIS: 'Likes' - 'Dislikes';
*/
func LikePostManager(idPost int, idUser string, likOrdIS string) {
	if likOrdIS == "Likes" {
		AddDeleteUpdateDB("DELETE FROM Dislikes WHERE PostID = ? AND UserID = ?;", fmt.Sprint(idPost), fmt.Sprint(idUser))
		AddLineIntoTargetTable(DataContainer{Likes: Likes{PostID: idPost, UserID: idUser}}, "Likes")
	} else {
		AddDeleteUpdateDB("DELETE FROM Likes WHERE PostID = ? AND UserID = ?;", fmt.Sprint(idPost), fmt.Sprint(idUser))
		AddLineIntoTargetTable(DataContainer{Dislikes: Dislikes{PostID: idPost, UserID: idUser}}, "Dislikes")
	}
	updateDislikeFromPost(idPost)
	updateLikeFromPost(idPost)
}

func UnLikePostManager(idPost int, idUser string, likOrdIS string) {
	if likOrdIS == "unLike" {
		AddDeleteUpdateDB("DELETE FROM Likes WHERE PostID = ? AND UserID = ?;", fmt.Sprint(idPost), fmt.Sprint(idUser))
	} else {
		AddDeleteUpdateDB("DELETE FROM Dislikes WHERE PostID = ? AND UserID = ?;", fmt.Sprint(idPost), fmt.Sprint(idUser))
	}
	updateDislikeFromPost(idPost)
	updateLikeFromPost(idPost)
}

func updateDislikeFromPost(idPost int) {
	rows := SelectDB("SELECT * FROM Dislikes WHERE PostID=?;", idPost)
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
	}
	AddDeleteUpdateDB("UPDATE Posts SET Dislikes=? WHERE PostID=?", fmt.Sprint(count), idPost)
}

func updateLikeFromPost(idPost int) {
	rows := SelectDB("SELECT * FROM Likes WHERE PostID=?;", idPost)
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
	}
	AddDeleteUpdateDB("UPDATE Posts SET Likes=? WHERE PostID=?", fmt.Sprint(count), idPost)
}

func AddTagsToTopic(tags, creatorId string, TopicID int) {
	tagsArray := strings.Split(tags, " ")
	for _, tag := range tagsArray {
		if (GetTagByName(tag) == Tags{}) {
			AddLineIntoTargetTable(DataContainer{Tags: Tags{Title: tag, CreatorID: creatorId}}, "Tags")
		}
		AddLineIntoTargetTable(DataContainer{TopicsTags: TopicsTags{TopicID: TopicID, TagID: GetTagByName(tag).TagID}}, "TopicsTags")
	}
}

/*--------------------*/
/* MODERATION SYSTEM: */
/*--------------------*/
func BanUser(userID string) {
	AddDeleteUpdateDB("DELETE FROM Dislikes WHERE UserID = ?;", userID)
	AddDeleteUpdateDB("DELETE FROM Follows WHERE UserID = ?;", userID)
	AddDeleteUpdateDB("DELETE FROM Likes WHERE UserID = ?;", userID)
	AddDeleteUpdateDB("DELETE FROM Upvotes WHERE UserID = ?;", userID)
	AddDeleteUpdateDB("DELETE FROM Topics WHERE CreatorID = ?", userID)
	AddDeleteUpdateDB("DELETE FROM Tags WHERE CreatorID = ?", userID)
	AddDeleteUpdateDB("DELETE FROM Posts WHERE AuthorID = ?;", userID)
	AddDeleteUpdateDB("DELETE FROM Reports WHERE ReportUserID = ?;", userID)
	SetUserStatus(userID, false)
}

func AddWordIntoBlacklist(word string) {
	if IsWordInBlacklist(word) {
		fmt.Println(word, "is already in the blacklist.")
		return
	}
	res := AddDeleteUpdateDB("INSERT INTO WordsBlacklist (WordID, Word) VALUES (?,?);", nil, word)
	_, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(word, "added in the blacklist.")
}

func SetUserStatus(userID string, status bool) {
	res := AddDeleteUpdateDB("UPDATE Users SET ValidUser = ? WHERE UserID = ?;", status, userID)
	_, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(userID, "has been ban!")
}

func DeleteReport(reportID string) {
	res := AddDeleteUpdateDB("DELETE FROM Reports WHERE ReportID = ?", reportID)
	_, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ReportID:", reportID, "deleted!")
}

func DeleteReportsFromPost(postID string) {
	res := AddDeleteUpdateDB("DELETE FROM Reports WHERE PostID = ?", postID)
	_, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("All reports concerned by PostID:", postID, "were delete!")
}

func DeletePost(postID string) {
	res := AddDeleteUpdateDB("DELETE FROM Posts WHERE Posts.PostID = ?", postID)
	_, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("PostID:", postID, "deleted!")
}

func DeletePostsFromTopic(topicID string) {
	res := AddDeleteUpdateDB("DELETE FROM Posts WHERE Posts.TopicID = ?", topicID)
	_, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Posts of TopicID:", topicID, "were delete!")
}

func DeleteTopic(topicID string) {
	res := AddDeleteUpdateDB("DELETE FROM Topics WHERE Topics.TopicID = ?", topicID)
	_, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("TopicID:", topicID, "has been delete!")
}

func DeleteReportsFromTopic(topicID string) {
	res := AddDeleteUpdateDB("DELETE FROM Reports WHERE TopicID = ?", topicID)
	_, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("All reports concerned by PostID:", topicID, "were delete!")
}

func AddPostReport(postID string, reason string) {
	db, err := sql.Open("sqlite3", "./DB-Forum.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT AuthorID FROM Posts WHERE PostID = ?;", postID)

	var userID string
	if err := row.Scan(&userID); err != nil {
		fmt.Println(err)
		return
	}

	res := AddDeleteUpdateDB("INSERT INTO Reports (ReportID, PostID, ReportUserID, Comment, TopicID) VALUES(?,?,?,?,?);", nil, postID, userID, reason, nil)
	_, err = res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("New post report! PostID:", postID, "Reason:", reason)
}

func AddTopicReport(topicID string, reason string) {
	db, err := sql.Open("sqlite3", "./DB-Forum.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT CreatorID FROM Topics WHERE TopicID = ?;", topicID)

	var userID string
	if err := row.Scan(&userID); err != nil {
		fmt.Println(err)
		return
	}

	res := AddDeleteUpdateDB("INSERT INTO Reports (ReportID, PostID, ReportUserID, Comment, TopicID) VALUES(?,?,?,?,?);", nil, nil, userID, reason, topicID)
	_, err = res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("New topic report! TopicID:", topicID, "Reason:", reason)
}
