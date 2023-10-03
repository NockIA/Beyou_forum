package handler

import (
	"crypto/sha256"
	"fmt"
	"forum/pkg/datamanagement"
	"net/http"
	"os/exec"
	"text/template"
	"time"
)

type register struct {
	EmailAndUsernameIsNotValid bool
	EmailIsNotValid            bool
	UsernameIsNotValid         bool
	IsConnected                bool
	IsAdmin                    bool
}

func CreateUser(UserId string, userName string, userFirstName string, userLastName string, userEmail string, stringPasswordInSha256 string) datamanagement.Users {
	nUser := datamanagement.Users{}
	nUser.UserID = UserId
	nUser.Username = userName
	nUser.Email = userEmail
	nUser.Password = stringPasswordInSha256
	nUser.Firstname = userFirstName
	nUser.Lastname = userLastName
	nUser.Description = ""
	nUser.CreationDate = time.Now()
	nUser.ProfilePicture = "../img/PP_wb.png"
	nUser.IsAdmin = false
	nUser.ValidUser = true
	return nUser
}

func deleteLastByte(tab []byte) []byte {
	returnedTab := []byte{}
	for index, element := range tab {
		if index != len(tab)-1 {
			returnedTab = append(returnedTab, element)
		}
	}
	return returnedTab
}

func Register(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./static/html/register.html", "./static/html/navBar.html"))
	userEmail := r.FormValue("email")
	userName := r.FormValue("username")
	userPassword := r.FormValue("password")
	rememberMe := r.FormValue("rememberMe")
	registerDisplay := register{false, false, false, false, false}
	if userEmail != "" && userName != "" && userPassword != "" {
		isUsernameAlreadyExist := datamanagement.IsUsernameAlreadyExist(userName)
		isEmailAlreadyExist := datamanagement.IsEmailAlreadyExist(userEmail)
		if !isUsernameAlreadyExist && !isEmailAlreadyExist {
			newUUID, err := exec.Command("uuidgen").Output()
			if err != nil {
				fmt.Println("user creation error ", err)
			}
			passwordByte := []byte(userPassword)
			passwordInSha256 := sha256.Sum256(passwordByte)
			stringPasswordInSha256 := fmt.Sprintf("%x", passwordInSha256[:])
			newUUID = deleteLastByte(newUUID)
			nUser := CreateUser(string(newUUID), userName, userName, userName, userEmail, stringPasswordInSha256)
			nDataContainer := datamanagement.DataContainer{}
			nDataContainer.Users = nUser
			datamanagement.AddLineIntoTargetTable(nDataContainer, "Users")
			cookieIdUser := http.Cookie{Name: "idUser", Value: string(newUUID)}
			http.SetCookie(w, &cookieIdUser)
			if rememberMe == "true" {
				cookieRememberMe := http.Cookie{Name: "Remember", Value: "true"}
				cookieIsConnected := http.Cookie{Name: "isConnected", Value: "true"}
				http.SetCookie(w, &cookieRememberMe)
				http.SetCookie(w, &cookieIsConnected)
			} else {
				cookieRememberMe := http.Cookie{Name: "Remember", Value: "false"}
				cookieIsConnected := http.Cookie{Name: "isConnected", Value: "true", Expires: time.Now().Add(6 * time.Hour)}
				http.SetCookie(w, &cookieRememberMe)
				http.SetCookie(w, &cookieIsConnected)
			}
			http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
		} else {
			if isEmailAlreadyExist && isUsernameAlreadyExist {
				registerDisplay.EmailAndUsernameIsNotValid = true
			} else if isEmailAlreadyExist {
				registerDisplay.EmailIsNotValid = true
			} else if isUsernameAlreadyExist {
				registerDisplay.UsernameIsNotValid = true
			}
		}
	}
	t.ExecuteTemplate(w, "register", registerDisplay)
}
