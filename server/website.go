// Filename:    website.go
// Author:      Joseph DeVictoria
// Date:        June_16_2018
// Purpose:     Embedded website used for account registration and login for Oldentide.

package main

import (
	"Oldentide/shared"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
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
		"<label for=\"registration_username\">Username (min 3)</label><br>"+
		"<input type=\"text\" id=\"registration_username\" name=\"registration_username\" "+
		"placeholder=\"User\" required=\"required\" maxlength=\"30\" pattern=\"[A-Za-z0-9]{3,}\"><br>"+
		"<label for=\"registration_email\">Email</label><br>"+
		"<input type=\"email\" id=\"registration_email\" name=\"registration_email\" "+
		"placeholder=\"user@domain.com\" required=\"required\"><br>"+
		"<label for=\"registration_password_first\">Password (min 6)</label><br>"+
		"<input type=\"password\" id=\"registration_password_first\" name=\"registration_password_first\" "+
		"placeholder=\"******\" required=\"required\" maxlength=\"30\" pattern=\".{6,}\"><br>"+
		"<label for=\"registration_password_second\">Confirm (min 6)</label><br>"+
		"<input type=\"password\" id=\"registration_password_second\" name=\"registration_password_second\" "+
		"placeholder=\"******\" required=\"required\" maxlength=\"30\" pattern=\".{6,}\" oninput=\"check(this)\"><br><br>"+
		"<script language='javascript' type='text/javascript'> function check(input) { "+
		"    if (input.value != document.getElementById('registration_password_first').value) "+
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
		registration_username := r.Form["registration_username"][0]
		registration_email := r.Form["registration_email"][0]
		registration_password_first := r.Form["registration_password_first"][0]
		registration_password_second := r.Form["registration_password_second"][0]
		if accountExists(registration_username) {
			fmt.Fprintf(w, "An account with the username "+registration_username+
				" already exists.\n\nPlease go back and choose another.")
		} else if emailExists(registration_email) {
			fmt.Fprintf(w, "An account is already associated with the email "+
				registration_email+"\n\nPlease use a unique email.")
		} else if registration_password_first == registration_password_second {
			// Generate a unique random verification key.
			//findKey := true

			verify_key := generateUniqueVerify(20)
			salt_key := generateUniqueSalt(40)
			hashed_key := shared.SaltAndHash(registration_password_first, salt_key)

			if everify {
				// Create email message to send to user.
				text_wport := ""
				if wport != 80 {
					text_wport = ":" + strconv.Itoa(wport)
				}
				msg := []byte("Hello " + registration_username +
					"\n\nPlease verify your Oldentide account by clicking the following link: " +
					"http://" + webadd + text_wport + "/verify/" + verify_key +
					" \n\nRegards,\nOldentide Server Admin")
				to := []string{registration_email}

				// Store user account information in the database!
				if !createAccount(registration_username, registration_email, verify_key, hashed_key, salt_key) {
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
					registration_username, registration_email)
			} else {
				// Store user account information in the database!
				if !createAccount(registration_username, registration_email, verify_key, hashed_key, salt_key) {
					fmt.Fprint(w, "Account could not be created, it caused a database error")
					return
				}
				fmt.Fprintf(w, "<html>You posted data to the register page.<br><br>"+
					"Email verification is disabled, your account has been created.<br>"+
					"Username: %s<br>Email: %s<br><br>"+
					registration_username, registration_email)
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
	verify_key := r.URL.Path[8:]
	accountname := getAccountnameFromVerifyKey(verify_key)
	if accountname == "" {
		fmt.Fprintf(w, "No account found for verify key "+verify_key+".")
	} else {
		if activateAccount(accountname) {
			fmt.Fprintf(w, "Successfully verified account "+accountname+".")
		} else {
			fmt.Println(w, "Found account "+accountname+" but couldn't activate for some reason.")
		}
	}
}

var session_id = 0;
// Web handler that handles logging in a player and returning them a sessionID.
func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Can only login via the Oldentide Client", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	username := r.Form["username"][0]
	password := r.Form["password"][0]
	fmt.Println("username=", username)
	fmt.Println("password=", password)
	if !accountExists(username) {
		fmt.Println("username dne")
		http.Error(w, "Username does not exist", http.StatusBadRequest)
		return;
	}

	// Get the salt of the user (the salt is public)
	salt := getSaltFromAccount(username)
	// Create a hash from the password and salt
	hash := getHashFromAccount(username)
	// Check the hash against what's in the DB
	supplied_hash := shared.SaltAndHash(password, salt)

	if supplied_hash != hash {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Return session_id and write session_id to db
	// TODO: Always assume new session ID? Check for old session id first?
	session_id++
	if !setSessionId(username, session_id) {
		http.Error(w, "Could not save the session ID to the DB", http.StatusInternalServerError)
		return
	}

	fmt.Println("Login successful! Sending user back session ID =", session_id)

	// Success! Return the session id
	fmt.Fprintf(w, strconv.Itoa(session_id))
}
