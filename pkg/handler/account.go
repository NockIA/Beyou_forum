package handler

import (
	"crypto/sha256"
	"fmt"
	"forum/pkg/datamanagement"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type AccountPage struct {
	Username                  string
	Email                     string
	Profile_picture           string
	Description               string
	FirstName                 string
	LastName                  string
	IsNotValidchangedPwd      bool
	IsNotValidchangedBIO      bool
	IsNotValidEditMail        bool
	IsNotValidchangedUsername bool
	IsConnected               bool
	IsAdmin                   bool
}

func setDisplayStructAccount(displayStructAccountPage AccountPage, currentUser datamanagement.Users) AccountPage {
	displayStructAccountPage.Username = currentUser.Username
	displayStructAccountPage.Email = currentUser.Email
	displayStructAccountPage.Profile_picture = currentUser.ProfilePicture
	displayStructAccountPage.Description = currentUser.Description
	displayStructAccountPage.FirstName = currentUser.Firstname
	displayStructAccountPage.LastName = currentUser.Lastname
	return displayStructAccountPage
}

func setDefaultValue(displayStructAccountPage AccountPage) AccountPage {
	displayStructAccountPage.Profile_picture = "../img/PP_wb.png"
	displayStructAccountPage.Username = "Guest"
	displayStructAccountPage.Email = "No email for guest"
	displayStructAccountPage.Description = "No bio added"
	displayStructAccountPage.FirstName = "No First Name"
	displayStructAccountPage.LastName = "No Last Name"
	return displayStructAccountPage
}

func Account(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("isConnected")
	if err != nil {
		fmt.Println("Cookie does not exist so redirect.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	isConnected := getCookieValue(cookie)
	if isConnected == "false" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	t := template.Must(template.ParseFiles("./static/html/account.html", "./static/html/navBar.html"))
	delAccount := r.FormValue("delAccount")
	file, handler, err := r.FormFile("photo")
	editMail := r.FormValue("editMail")
	changedPwd1 := r.FormValue("changedPwd1")
	changedPwd2 := r.FormValue("changedPwd2")
	currentPwd := r.FormValue("currentPwd")
	changedBIO := r.FormValue("changedBIO")
	changedFirstname := r.FormValue("changedFirstname")
	changedLastname := r.FormValue("changedLastname")
	changedUsername := r.FormValue("changedUsername")
	disconnect := r.FormValue("disconnect")
	cookieIdUser, _ := r.Cookie("idUser")
	displayStructAccountPage := AccountPage{}
	idUser := getCookieValue(cookieIdUser)
	switch true {
	case delAccount != "":
		datamanagement.AddDeleteUpdateDB("DELETE FROM Users WHERE UserID = ?;", idUser)
		cookieIsConnected := http.Cookie{Name: "isConnected", Value: "false"}
		http.SetCookie(w, &cookieIsConnected)
		cookieIdUser := http.Cookie{Name: "idUser", Value: ""}
		http.SetCookie(w, &cookieIdUser)
		cookieRemenberMe := http.Cookie{Name: "Remember", Value: "false"}
		http.SetCookie(w, &cookieRemenberMe)
		break
	case disconnect != "":
		cookieIsConnected := http.Cookie{Name: "isConnected", Value: "false"}
		http.SetCookie(w, &cookieIsConnected)
		cookieIdUser := http.Cookie{Name: "idUser", Value: ""}
		http.SetCookie(w, &cookieIdUser)
		cookieRemenberMe := http.Cookie{Name: "Remember", Value: "false"}
		http.SetCookie(w, &cookieRemenberMe)
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
	case editMail != "":
		if !datamanagement.IsEmailAlreadyExist(editMail) && datamanagement.CheckContentByBlackListWord(editMail) {
			datamanagement.AddDeleteUpdateDB("UPDATE Users SET Email = ? WHERE UserID = ?;", editMail, idUser)
		} else {
			displayStructAccountPage.IsNotValidEditMail = true
		}
		break
	case changedPwd1 != "" && changedPwd2 != "":
		if changedPwd1 == changedPwd2 && datamanagement.IsValidPassword(currentPwd, idUser) {
			passwordByte := []byte(changedPwd1)
			passwordInSha256 := sha256.Sum256(passwordByte)
			stringPasswordInSha256 := fmt.Sprintf("%x", passwordInSha256[:])
			datamanagement.AddDeleteUpdateDB("UPDATE Users SET Password = ? WHERE UserID = ?;", stringPasswordInSha256, idUser)
		} else {
			displayStructAccountPage.IsNotValidchangedPwd = true
		}
		break
	case changedBIO != "" && datamanagement.CheckContentByBlackListWord(changedBIO):
		datamanagement.AddDeleteUpdateDB("UPDATE Users SET Description = ? WHERE UserID = ?;", changedBIO, idUser)
		break
	case changedFirstname != "" && datamanagement.CheckContentByBlackListWord(changedFirstname):
		datamanagement.AddDeleteUpdateDB("UPDATE Users SET Firstname = ? WHERE UserID = ?;", changedFirstname, idUser)
		break
	case changedLastname != "" && datamanagement.CheckContentByBlackListWord(changedLastname):
		datamanagement.AddDeleteUpdateDB("UPDATE Users SET Lastname = ? WHERE UserID = ?;", changedLastname, idUser)
		break
	case changedUsername != "" && datamanagement.CheckContentByBlackListWord(changedUsername):
		if !datamanagement.IsUsernameAlreadyExist(changedUsername) {
			datamanagement.AddDeleteUpdateDB("UPDATE Users SET Username = ? WHERE UserID = ?;", changedUsername, idUser)
		} else {
			displayStructAccountPage.IsNotValidchangedUsername = true
		}
		break
	}
	currentUser := datamanagement.GetUserById(idUser)
	if file != nil && err == nil {
		defer file.Close()
		destinationPath := "./static/img/" + idUser + "." + strings.Split(handler.Filename, ".")[1]
		destinationFile, err := os.Create(destinationPath)
		if err != nil {
			fmt.Println("Failed to create destination file")
		}
		defer destinationFile.Close()
		_, err = io.Copy(destinationFile, file)
		if err != nil {
			fmt.Println("Failed to save photo on server")
		}
		fileName := "../img/" + idUser + "." + strings.Split(handler.Filename, ".")[1]
		datamanagement.AddDeleteUpdateDB("UPDATE Users SET ProfilePicture = ? WHERE UserID = ?;", fileName, idUser)
		http.Redirect(w, r, "http://localhost:8080/account", http.StatusSeeOther)
	}
	displayStructAccountPage.IsAdmin = currentUser.IsAdmin
	displayStructAccountPage = setDisplayStructAccount(displayStructAccountPage, currentUser)
	displayStructAccountPage.Profile_picture = currentUser.ProfilePicture
	cookieIsConnected, _ := r.Cookie("isConnected")
	isConnected = getCookieValue(cookieIsConnected)
	if isConnected != "true" {
		displayStructAccountPage = setDefaultValue(displayStructAccountPage)
		displayStructAccountPage.IsConnected = false
	} else {
		displayStructAccountPage.IsConnected = true
	}
	t.ExecuteTemplate(w, "account", displayStructAccountPage)
}
