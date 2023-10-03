package datamanagement

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	username = ""
	password = ""
	hostname = ""
	dbname   = ""
)

type DataFilter struct {
	Number int
}

type Display struct {
	Logout bool
}

type Dislikes struct {
	PostID int
	UserID string
}

type Follows struct {
	FollowID int
	TopicID  int
	UserID   string
}

type Likes struct {
	PostID int
	UserID string
}

type Posts struct {
	PostID       int
	Content      string
	AuthorID     string
	TopicID      int
	Likes        int
	Dislikes     int
	CreationDate time.Time
}

type Reports struct {
	ReportID     int
	PostID       int
	ReportUserID string
	Comment      string
	TopicID      int
}

type Tags struct {
	TagID     int
	Title     string
	CreatorID string
}

type UserConnected struct {
	IsUserConnected bool
	IdUser          string
}

type Topics struct {
	TopicID      int
	Title        string
	Description  string
	Picture      string
	CreationDate time.Time
	CreatorID    string
	Upvotes      int
	Follows      int
}

type TopicsTags struct {
	TopicID int
	TagID   int
}

type Upvotes struct {
	TopicID int
	UserID  string
}

type Users struct {
	UserID         string
	Username       string
	Email          string
	Password       string
	Firstname      string
	Lastname       string
	Description    string
	CreationDate   time.Time
	ProfilePicture string
	IsAdmin        bool
	ValidUser      bool
}

type WordsBlacklist struct {
	WordID int
	Word   string
}

type DataContainer struct {
	Dislikes   Dislikes
	Follows    Follows
	Likes      Likes
	Posts      Posts
	Reports    Reports
	Tags       Tags
	Topics     Topics
	TopicsTags TopicsTags
	Upvotes    Upvotes
	Users
	WordsBlacklist
}

type DataExplorePage struct {
	Topics []Topics
	Users  []string
}

func CheckContentByBlackListWord(content string) bool {
	blackListWords := GetAllFromTable("WordsBlacklist")
	contentArray := strings.Split(content, " ")
	for _, w := range blackListWords {
		if arrayContains(contentArray, w.WordsBlacklist.Word) {
			return false
		}
	}
	contentArray = strings.Split(content, "")
	blackListCaracter := []string{"<", ">", "(", ")", "{", "}", "|", "[", "]", "/", "\\", "#", "$", "£", "¤", "@"}
	for _, c := range content {
		for _, cb := range blackListCaracter {
			if string(c) == cb {
				return false
			}
		}
	}
	if IsAnyWordToLong(content) {
		return false
	}
	return true
}

func arrayContains(array []string, word string) bool {
	for _, val := range array {
		if word == val {
			return true
		}
	}
	return false
}

func IsAnyWordToLong(s string) bool {
	as := strings.Split(s, " ")
	for _, w := range as {
		if len(w) > 20 {
			return true
		}
	}
	return false
}

func SelectDB(query string, args ...interface{}) *sql.Rows {
	db, err := sql.Open("sqlite3", "./DB-Forum.db")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return rows
}

func AddDeleteUpdateDB(query string, args ...interface{}) sql.Result {
	db, err := sql.Open("sqlite3", "./DB-Forum.db")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer db.Close()

	res, err := db.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func TransformDateInPostFormat(CreationDate time.Time) string {
	pastTime := math.Trunc(CreationDate.Sub(time.Now()).Minutes() * -1)
	if pastTime < 60 {
		return fmt.Sprintf("%v", pastTime) + " min"
	} else {
		pastTime = math.Trunc(CreationDate.Sub(time.Now()).Hours() * -1)
		if pastTime < 24 {
			return fmt.Sprintf("%v", pastTime) + " h"
		} else {
			pastTime = math.Trunc(pastTime / 24)
			if pastTime < 30 {
				if pastTime <= 1 {
					return fmt.Sprintf("%v", pastTime) + " day"
				} else {
					return fmt.Sprintf("%v", pastTime) + " days"
				}
			} else {
				pastTime = math.Trunc(pastTime / 30)
				if pastTime < 12 {
					if pastTime <= 1 {
						return fmt.Sprintf("%v", pastTime) + " month"
					} else {
						return fmt.Sprintf("%v", pastTime) + " months"
					}
				} else {
					pastTime = math.Trunc(pastTime / 12)
					if pastTime <= 1 {
						return fmt.Sprintf("%v", pastTime) + " year"
					} else {
						return fmt.Sprintf("%v", pastTime) + " years"
					}
				}
			}
		}
	}
}

func TransformTags(topicId int) []string {
	tags := GetTagsByTopic(topicId)
	result := []string{}
	for i, tag := range tags {
		if i > 4 {
			break
		}
		result = append(result, tag.Title)
	}
	return result
}
