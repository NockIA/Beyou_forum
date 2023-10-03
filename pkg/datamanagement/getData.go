package datamanagement

import (
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func GetPostFromUser(idUser string) []Posts {
	rows := SelectDB("SELECT * FROM Posts WHERE AuthorID = ?;", idUser)
	defer rows.Close()
	result := []Posts{}
	for rows.Next() {
		var post Posts
		rows.Scan(&post.PostID, &post.Content, &post.AuthorID, &post.TopicID, &post.Likes, &post.Dislikes, &post.CreationDate)
		result = append(result, post)
	}
	return result
}

func GetPostById(idPost int) Posts {
	row := SelectDB("SELECT * FROM Posts LEFT JOIN Users ON AuthorID = UserID LEFT JOIN Topics ON Posts.TopicID = Topics.TopicID WHERE PostID = ?;", strconv.Itoa(idPost))
	defer row.Close()
	var post Posts
	for row.Next() {
		row.Scan(&post.PostID, &post.Content, &post.AuthorID, &post.TopicID, &post.Likes, &post.Dislikes, &post.CreationDate)
	}
	return post
}

func GetSortPost() []Posts {
	rows := SelectDB("SELECT * FROM Posts ORDER BY Likes - Dislikes DESC;")
	defer rows.Close()

	result := []Posts{}
	for rows.Next() {
		var post Posts
		rows.Scan(&post.PostID, &post.Content, &post.AuthorID, &post.TopicID, &post.Likes, &post.Dislikes, &post.CreationDate)
		result = append(result, post)
	}

	return result
}

func SearchUserByName(search string) []Users {
	rows := SelectDB("SELECT * FROM Users WHERE Username LIKE '%" + search + "%';")
	defer rows.Close()

	result := []Users{}
	for rows.Next() {
		var user Users
		rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Firstname, &user.Lastname, &user.Description, &user.CreationDate, &user.ProfilePicture, &user.IsAdmin, &user.ValidUser)
		result = append(result, user)
	}
	return result
}

func GetUserByName(userName string) Users {
	rows := SelectDB("SELECT * FROM Users WHERE Username LIKE ?;", userName)
	defer rows.Close()

	var user Users
	for rows.Next() {
		rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Firstname, &user.Lastname, &user.Description, &user.CreationDate, &user.ProfilePicture, &user.IsAdmin, &user.ValidUser)
	}
	return user
}

func GetUserByID(userId string) Users {
	rows := SelectDB("SELECT * FROM Users WHERE UserID LIKE ?;", userId)
	defer rows.Close()

	var user Users
	for rows.Next() {
		rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Firstname, &user.Lastname, &user.Description, &user.CreationDate, &user.ProfilePicture, &user.IsAdmin, &user.ValidUser)
	}
	return user
}

func GetTopicId(topicName string) int {
	rows := SelectDB("SELECT TopicID FROM Topics WHERE Title LIKE '" + topicName + "';")
	defer rows.Close()
	var id int
	for rows.Next() {
		rows.Scan(&id)
	}
	return id
}

func GetPostByTopic(topic int) []Posts {
	rows := SelectDB("SELECT * FROM Posts WHERE TopicID LIKE " + strconv.Itoa(topic) + ";")
	defer rows.Close()

	result := []Posts{}
	for rows.Next() {
		var post Posts
		rows.Scan(&post.PostID, &post.Content, &post.AuthorID, &post.TopicID, &post.Likes, &post.Dislikes, &post.CreationDate)
		result = append(result, post)
	}

	return result
}

func GetAllFromTable(table string) []DataContainer {
	QUERY := "SELECT * FROM '" + table + "';"
	rows := SelectDB(QUERY)
	defer rows.Close()

	var result []DataContainer
	for rows.Next() {
		var line DataContainer
		switch true {
		case table == "Users":
			rows.Scan(&line.Users.UserID, &line.Users.Username, &line.Users.Email, &line.Users.Password, &line.Users.Firstname, &line.Users.Lastname, &line.Users.Description, &line.Users.CreationDate, &line.Users.ProfilePicture, &line.Users.IsAdmin, &line.Users.ValidUser)
			break
		case table == "Posts":
			rows.Scan(&line.Posts.PostID, &line.Posts.Content, &line.Posts.AuthorID, &line.Posts.TopicID, &line.Posts.Likes, &line.Posts.Dislikes, &line.Posts.CreationDate)
			break
		case table == "Topics":
			rows.Scan(&line.Topics.TopicID, &line.Topics.Title, &line.Topics.Description, &line.Topics.Picture, &line.CreationDate, &line.Topics.CreatorID, &line.Topics.Upvotes, &line.Topics.Follows)
			break
		case table == "Tags":
			rows.Scan(&line.Tags.TagID, &line.Tags.Title, &line.Tags.CreatorID)
			break
		case table == "Reports":
			rows.Scan(&line.Reports.ReportID, &line.Reports.PostID, &line.Reports.ReportUserID, &line.Reports.Comment)
			break
		case table == "Dislikes":
			rows.Scan(&line.Dislikes.PostID, &line.Dislikes.UserID)
			break
		case table == "Likes":
			rows.Scan(&line.Dislikes.PostID, &line.Dislikes.UserID)
			break
		case table == "Follows":
			rows.Scan(&line.Follows.FollowID, &line.Follows.TopicID, &line.Follows.UserID)
			break
		case table == "TopicsTags":
			rows.Scan(&line.TopicsTags.TopicID, &line.TopicsTags.TagID)
			break
		case table == "Upvotes":
			rows.Scan(&line.Upvotes.TopicID, &line.Upvotes.UserID)
			break
		case table == "WordsBlacklist":
			rows.Scan(&line.WordsBlacklist.WordID, &line.WordsBlacklist.Word)
			break
		}
		result = append(result, line)
	}
	return result
}

/*
typofsort: 'a-z' - 'z-a' - 'DESC-Upvote' - 'ASC-Upvote' - 'creator'
*/

func SortTopics(typOfSort string) []Topics {
	var query string
	switch typOfSort {
	case "a-z":
		query = "SELECT * FROM Topics ORDER BY Title ASC;"
		break
	case "z-a":
		query = "SELECT * FROM Topics ORDER BY Title DESC;"
		break
	case "DESC-Upvote":
		query = "SELECT * FROM Topics ORDER BY Upvotes DESC;"
		break
	case "DESC-Upvote-Home":
		query = "SELECT * FROM Topics ORDER BY Upvotes DESC LIMIT 3;"
		break
	case "ASC-Upvote":
		query = "SELECT * FROM Topics ORDER BY Upvotes ASC;"
		break
	case "creator":
		query = "SELECT * FROM Topics ORDER BY CreatorID DESC;"
		break
	case "default":
		query = "SELECT * FROM Topics;"
		break
	default:
		fmt.Println("invalid type of sort")
		return nil
	}

	rows := SelectDB(query)
	defer rows.Close()

	var result []Topics
	for rows.Next() {
		var line Topics
		rows.Scan(&line.TopicID, &line.Title, &line.Description, &line.Picture, &line.CreationDate, &line.CreatorID, &line.Upvotes, &line.Follows)
		result = append(result, line)
	}

	return result
}

func GetAllBlacklistWords() []WordsBlacklist {
	rows := SelectDB("SELECT * FROM WordsBlacklist;")
	defer rows.Close()

	var blackListWords []WordsBlacklist
	for rows.Next() {
		var word WordsBlacklist
		rows.Scan(&word.WordID, &word.Word)
		blackListWords = append(blackListWords, word)
	}

	return blackListWords
}

func GetUserById(id string) Users {
	row := SelectDB("SELECT * FROM Users WHERE UserID = ?;", id)
	defer row.Close()
	var user Users
	for row.Next() {
		row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Firstname, &user.Lastname, &user.Description, &user.CreationDate, &user.ProfilePicture, &user.IsAdmin, &user.ValidUser)
	}
	return user
}

func GetTopicsById(creatorID string) []Topics {
	rows := SelectDB("SELECT * FROM Topics WHERE CreatorID = ?;", creatorID)
	defer rows.Close()

	topics := []Topics{}
	for rows.Next() {
		var topic Topics
		rows.Scan(&topic.TopicID, &topic.Title, &topic.Description, &topic.Picture, &topic.CreationDate, &topic.CreatorID, &topic.Upvotes, &topic.Follows)
		topics = append(topics, topic)
	}
	return topics
}

func GetTopicsByName(search string) []Topics {
	rows := SelectDB("SELECT * FROM Topics WHERE Title LIKE '%" + search + "%';")
	defer rows.Close()
	result := []Topics{}
	for rows.Next() {
		var topic Topics
		rows.Scan(&topic.TopicID, &topic.Title, &topic.Description, &topic.Picture, &topic.CreationDate, &topic.CreatorID, &topic.Upvotes, &topic.Follows)
		result = append(result, topic)
	}
	return result
}

func GetOneTopicByName(search string) Topics {
	rows := SelectDB("SELECT * FROM Topics WHERE Title = ?;", search)
	defer rows.Close()

	result := Topics{}
	for rows.Next() {
		rows.Scan(&result.TopicID, &result.Title, &result.Description, &result.Picture, &result.CreationDate, &result.CreatorID, &result.Upvotes, &result.Follows)
	}
	return result
}

func GetTagByName(search string) Tags {
	rows := SelectDB("SELECT * FROM Tags WHERE Title = ?;", search)
	defer rows.Close()

	result := Tags{}
	for rows.Next() {
		rows.Scan(&result.TagID, &result.Title, &result.CreatorID)
	}
	return result
}

/*--------------------*/
/* MODERATION SYSTEM: */
/*--------------------*/

func GetAllBannedUsers() []Users {
	rows := SelectDB("SELECT * FROM Users WHERE ValidUser = false;")
	defer rows.Close()

	var bannedUsers []Users
	for rows.Next() {
		var user Users
		rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Firstname, &user.Lastname, &user.Description, &user.CreationDate, &user.ProfilePicture, &user.IsAdmin, &user.ValidUser)
		bannedUsers = append(bannedUsers, user)
	}

	return bannedUsers
}

func GetTopicsByUser(userId string) []Topics {
	topicsRows := SelectDB("SELECT DISTINCT t.TopicID, t.Title, t.Description, t.Picture, t.CreationDate, t.CreatorID, t.Upvotes, t.Follows FROM Topics AS t INNER JOIN Follows AS f ON f.TopicID = t.TopicID WHERE UserID=?;", userId)
	defer topicsRows.Close()
	result := []Topics{}
	for topicsRows.Next() {
		var topic Topics
		topicsRows.Scan(&topic.TopicID, &topic.Title, &topic.Description, &topic.Picture, &topic.CreationDate, &topic.CreatorID, &topic.Upvotes, &topic.Follows)
		result = append(result, topic)
	}
	return result
}

func GetTopicID(topicID int) Topics {
	topicRow := SelectDB("SELECT * FROM Topics WHERE TopicID=?;", topicID)
	defer topicRow.Close()
	result := Topics{}
	for topicRow.Next() {
		topicRow.Scan(&result.TopicID, &result.Title, &result.Description, &result.Picture, &result.CreationDate, &result.CreatorID, &result.Upvotes, &result.Follows)
	}
	return result
}

func GetTagByID(tagID int) Tags {
	topicRow := SelectDB("SELECT * FROM Tags WHERE TagID=?;", tagID)
	defer topicRow.Close()
	result := Tags{}
	for topicRow.Next() {
		topicRow.Scan(&result.TagID, &result.Title, &result.CreatorID)
	}
	return result
}

func GetTagsByTopic(topicID int) []Tags {
	tagsRow := SelectDB("SELECT DISTINCT tg.TagID,tg.Title,tg.CreatorID FROM Tags AS tg LEFT JOIN TopicsTags AS tt ON tt.TagID = tg.TagID WHERE tt.TopicID=?;", topicID)
	defer tagsRow.Close()
	result := []Tags{}
	for tagsRow.Next() {
		var tag Tags
		tagsRow.Scan(&tag.TagID, &tag.Title, &tag.CreatorID)
		result = append(result, tag)
	}
	return result
}

func GetTopicByTagAndTitle(search string) []Topics {
	topicsRow := SelectDB("SELECT DISTINCT t.TopicID, t.Title, t.Description, t.Picture, t.CreationDate, t.CreatorID, t.Upvotes, t.Follows FROM Topics AS t LEFT JOIN TopicsTags AS tt ON tt.TopicID = t.TopicID WHERE t.Title LIKE '%"+search+"%' OR tt.TagID IN (SELECT TagID FROM Tags WHERE Title = ?)", search)
	defer topicsRow.Close()
	result := []Topics{}
	for topicsRow.Next() {
		var topic Topics
		topicsRow.Scan(&topic.TopicID, &topic.Title, &topic.Description, &topic.Picture, &topic.CreationDate, &topic.CreatorID, &topic.Upvotes, &topic.Follows)
		result = append(result, topic)
	}

	return result
}
