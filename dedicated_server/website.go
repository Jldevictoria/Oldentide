// Filename:    website.go
// Author:      Joseph DeVictoria
// Date:        June_16_2018
// Purpose:     Embedded website used for account registration and login for Oldentide.

package main

import (
	"Oldentide/shared"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var homereg = regexp.MustCompile("^/$")
var regireg = regexp.MustCompile("/register")
var verireg = regexp.MustCompile("/verify/[a-z,A-Z]{20}")
var loginreg = regexp.MustCompile("/login")

// Router for web traffic.
func routeWebTraffic(w http.ResponseWriter, r *http.Request) {
	switch {
	case homereg.MatchString(r.URL.Path):
		homePage(w, r)
	case regireg.MatchString(r.URL.Path):
		registerPage(w, r)
	case verireg.MatchString(r.URL.Path):
		verifyPage(w, r)
	case loginreg.MatchString(r.URL.Path):
		loginPage(w, r)
	default:
		fmt.Fprintf(w, "Unknown Path.")
	}
}

// Web handler that loads the home page.
func homePage(w http.ResponseWriter, r *http.Request) {
	hn, err := os.Hostname()
	shared.CheckErr(err)
	fmt.Fprintf(w, "<html>"+
		"<p>Welcome to the Oldentide dedicated server running on %s.</p>"+
		"<p>Account Registration:</p>"+
		"<form name=\"registration\" method=\"post\" action=\"/register\" role=\"form\">"+
		"<label for=\"registrationUsername\">Username (min 3)</label><br>"+
		"<input type=\"text\" id=\"registrationUsername\" name=\"registrationUsername\" "+
		"placeholder=\"User\" required=\"required\" maxlength=\"30\" pattern=\"[A-Za-z0-9]{3,}\"><br>"+
		"<label for=\"registrationEmail\">Email</label><br>"+
		"<input type=\"email\" id=\"registrationEmail\" name=\"registrationEmail\" "+
		"placeholder=\"user@domain.com\" required=\"required\"><br>"+
		"<label for=\"registrationPasswordFirst\">Password (min 6)</label><br>"+
		"<input type=\"password\" id=\"registrationPasswordFirst\" name=\"registrationPasswordFirst\" "+
		"placeholder=\"******\" required=\"required\" maxlength=\"30\" pattern=\".{6,}\"><br>"+
		"<label for=\"registrationPasswordSecond\">Confirm (min 6)</label><br>"+
		"<input type=\"password\" id=\"registrationPasswordSecond\" name=\"registrationPasswordSecond\" "+
		"placeholder=\"******\" required=\"required\" maxlength=\"30\" pattern=\".{6,}\" oninput=\"check(this)\"><br><br>"+
		"<script language='javascript' type='text/javascript'> function check(input) { "+
		"    if (input.value != document.getElementById('registrationPasswordFirst').value) "+
		"        { input.setCustomValidity('Password Must be Matching.'); "+
		"    } else { input.setCustomValidity(''); }	} </script>"+
		"<input class=\"button-primary\" type=\"submit\" value=\"Register\">"+
		"</form>"+
		"</html>", hn)
}

// Web handler that uses the results of the registration POST to initiate a new player account registration.
func registerPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "Hi this should have the registration page.")
	} else if r.Method == "POST" {
		r.ParseForm()
		registrationUsername := r.Form["registrationUsername"][0]
		registrationEmail := r.Form["registrationEmail"][0]
		registrationPasswordFirst := r.Form["registrationPasswordFirst"][0]
		registrationPasswordSecond := r.Form["registrationPasswordSecond"][0]
		if accountExists(registrationUsername) {
			fmt.Fprintf(w, "An account with the username "+registrationUsername+
				" already exists.\n\nPlease go back and choose another.")
		} else if emailExists(registrationEmail) {
			fmt.Fprintf(w, "An account is already associated with the email "+
				registrationEmail+"\n\nPlease use a unique email.")
		} else if registrationPasswordFirst == registrationPasswordSecond {
			// Generate a unique random verification key.
			//findKey := true

			verifyKey := generateUniqueVerify(20)
			saltKey := generateUniqueSalt(40)
			hashedKey := shared.SaltAndHash(registrationPasswordFirst, saltKey)

			if everify {
				// Create email message to send to user.
				textWPort := ""
				if wport != 80 {
					textWPort = ":" + strconv.Itoa(wport)
				}
				msg := []byte("Hello " + registrationUsername +
					"\n\nPlease verify your Oldentide account by clicking the following link: " +
					"http://" + webadd + textWPort + "/verify/" + verifyKey +
					" \n\nRegards,\nOldentide Server Admin")
				to := []string{registrationEmail}

				// Store user account information in the database!
				if !createAccount(registrationUsername, registrationEmail, verifyKey, hashedKey, saltKey) {
					fmt.Fprint(w, "Account could not be created, it caused a database error")
					return
				}

				// Send registration email using
				err = smtp.SendMail("smtp.gmail.com:587", eauth, email, to, msg)
				shared.CheckErr(err)
				fmt.Fprintf(w, "<html>You posted data to the register page.<br><br>"+
					"An email has been sent to verify this information:<br><br>"+
					"Username: %s<br>Email: %s<br><br>"+
					"<b>It may take up to 5 minutes for this email to arrive.</b></html>",
					registrationUsername, registrationEmail)
			} else {
				// Store user account information in the database!
				if !createAccount(registrationUsername, registrationEmail, verifyKey, hashedKey, saltKey) {
					fmt.Fprint(w, "Account could not be created, it caused a database error")
					return
				}
				fmt.Fprintf(w, "<html>You posted data to the register page.<br><br>"+
					"Email verification is disabled, your account has been created.<br>"+
					"Username: %s<br>Email: %s<br><br>"+
					registrationUsername, registrationEmail)
			}
		} else {
			fmt.Fprintf(w, "Your passwords did not match...")
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Web handler that uses the url key to validate email for an account.
func verifyPage(w http.ResponseWriter, r *http.Request) {
	verifyKey := r.URL.Path[8:]
	accountname := getAccountnameFromVerifyKey(verifyKey)
	if accountname == "" {
		fmt.Fprintf(w, "No account found for verify key "+verifyKey+".")
	} else {
		if activateAccount(accountname) {
			fmt.Fprintf(w, "Successfully verified account "+accountname+".")
		} else {
			fmt.Println(w, "Found account "+accountname+" but couldn't activate for some reason.")
		}
	}
}

// Web handler that handles logging in a player and returning them a sessionID.
func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Can only login via the Oldentide Client", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	username := r.Form["username"][0]
	password := r.Form["password"][0]
	if !accountExists(username) {
		fmt.Println("Username:", username, "Password:", password, " does not exist in databse.")
		http.Error(w, "Username does not exist", http.StatusBadRequest)
		return
	}

	// Get the salt of the user (the salt is public)
	salt := getSaltFromAccount(username)
	// Use the salt to get the hashed version of the given password.
	suppliedHash := shared.SaltAndHash(password, salt)
	// Get the hash key (password + salt) for the user from the database (the hash is private).
	hash := getHashFromAccount(username)

	if suppliedHash != hash {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		fmt.Println("Username:", username, " tried to log in with the wrong password.")
		return
	}

	// Return sessionID and write sessionID to db
	sessionID := generateUniqueSessionID() // sql_connector.go
	if !setSessionID(username, sessionID) {
		http.Error(w, "Could not save the session ID to the DB", http.StatusInternalServerError)
		return
	}

	fmt.Println("Username:", username, "Password:", password, "successfuly logged in and was awarded sessionID:", sessionID)

	// Success! Return the session id
	fmt.Fprintf(w, strconv.FormatInt(sessionID, 10))
}
