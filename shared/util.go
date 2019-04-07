// Filename:    util.go
// Author:      Joseph DeVictoria
// Date:        June_16_2018
// Purpose:     Utility functions used by Oldentide dedicated server.

package shared

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"math/rand"
	"regexp"
)

// Simple function to check the error status of an operation.
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IfErrPrintErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Function that takes in a password and salt and hashes them together into a db safe key.
func SaltAndHash(password string, salt string) string {
	hasher := sha256.New()
	saltpass := password + salt
	hasher.Write([]byte(saltpass))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hash
}

// Util used to generate a string of lower and upper case letters.
func GenerateRandomLetters(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	key := make([]rune, n)
	for i := range key {
		key[i] = letters[rand.Intn(len(letters))]
	}
	return string(key)
}

// Util used to generate a string of lowe and upper case letters and numbers.
func GenerateRandomAlnums(n int) string {
	var alnums = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	key := make([]rune, n)
	for i := range key {
		key[i] = alnums[rand.Intn(len(alnums))]
	}
	return string(key)
}


// Validate a character name. 3-20 letters, alphabetic
var regex_name = regexp.MustCompile("^[a-zA-Z]{3,20}$")
func ValidateName(name string) bool {
	return regex_name.MatchString(name)
}

// Validate a username. 3-30 characters, alphanumeric
var regex_username = regexp.MustCompile("^[a-zA-Z0-9]{3,20}$")
func ValidateUsername(name string) bool {
	return regex_username.MatchString(name)
}
