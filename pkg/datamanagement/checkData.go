package datamanagement

import (
	"crypto/sha256"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func IsAdmin(userID string) bool {
	rows := SelectDB("SELECT IsAdmin FROM Users WHERE UserID = ?;", userID)
	defer rows.Close()
	for rows.Next() {
		var isAdmin bool
		rows.Scan(&isAdmin)
		return isAdmin
	}
	return false
}

func IsPostDLikeByBYser(PostID int, UserID string, DisOrLike string) bool {
	rows := SelectDB("SELECT * FROM ? WHERE PostID = ? AND UserID = ?", DisOrLike, strconv.Itoa(PostID), UserID)
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}

func IsEmailAlreadyExist(userEmail string) bool {
	rows := SelectDB("SELECT * FROM Users WHERE (Email = ?);", string(userEmail))
	defer rows.Close()

	if !rows.Next() {
		return false
	} else {
		return true
	}
}

func IsUsernameAlreadyExist(userProfilName string) bool {
	rows := SelectDB("SELECT * FROM Users WHERE ( Username = ?);", string(userProfilName))
	defer rows.Close()

	if !rows.Next() {
		return false
	} else {
		return true
	}
}

func IsValidTopic(topic string) bool {
	rows := SelectDB("SELECT * FROM Topics WHERE (Title = ?);", string(topic))
	defer rows.Close()

	if !rows.Next() {
		return false
	} else {
		return true
	}
}

func IsUserExist(userEmail string, userProfilName string) bool {
	rows := SelectDB("SELECT * FROM Users WHERE ( Username = ? OR Email = ?);", string(userProfilName), string(userEmail))
	defer rows.Close()

	if !rows.Next() {
		return false
	} else {
		return true
	}
}

func IsRegister(userInput string, password string) (bool, string) {
	passwordByte := []byte(password)
	passwordInSha256 := sha256.Sum256(passwordByte)
	stringPasswordInSha256 := fmt.Sprintf("%x", passwordInSha256[:])
	rows := SelectDB("SELECT userID FROM Users WHERE ( Username = ? OR Email = ? ) AND Password = ?;", string(userInput), string(userInput), string(stringPasswordInSha256))
	defer rows.Close()

	for rows.Next() {
		var id string
		rows.Scan(&id)
		return true, id
	}

	return false, ""
}

func IsValidPassword(password string, idUser string) bool {
	dataUser := GetUserByID(idUser)
	passwordByte := []byte(password)
	passwordInSha256 := sha256.Sum256(passwordByte)
	stringPasswordInSha256 := fmt.Sprintf("%x", passwordInSha256[:])
	if dataUser.Password == stringPasswordInSha256 {
		return true
	} else {
		return false
	}
}

func IsLikeByUser(userID string, postID int) bool {
	rows := SelectDB("SELECT * FROM Likes WHERE (PostID = ? AND UserID = ?);", postID, userID)
	defer rows.Close()
	if !rows.Next() {
		return false
	} else {
		return true
	}
}

func IsDislikeByUser(userID string, postID int) bool {
	rows := SelectDB("SELECT * FROM Dislikes WHERE (PostID = ? AND UserID = ?);", postID, userID)
	defer rows.Close()
	if !rows.Next() {
		return false
	} else {
		return true
	}
}

func IsWordInBlacklist(word string) bool {
	rows := SelectDB("SELECT word FROM WordsBlacklist WHERE word = ?", word)
	defer rows.Close()

	if !rows.Next() {
		return false
	}
	return true
}
