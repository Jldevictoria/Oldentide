// Filename:    util.go
// Author:      Joseph DeVictoria
// Date:        June_16_2018
// Purpose:     Utility functions used by Oldentide dedicated server.

package shared

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

// CheckErr checks the error status of an operation.
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// IfErrPrintErr checks the error status of an operation and prints the error if it exists.
func IfErrPrintErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

// SaltAndHash takes in a password and salt and hashes them together into a db safe key.
func SaltAndHash(password string, salt string) string {
	hasher := sha256.New()
	saltpass := password + salt
	hasher.Write([]byte(saltpass))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hash
}

// GenerateRandomLetters is used to generate a string of lower and upper case letters.
func GenerateRandomLetters(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	key := make([]rune, n)
	for i := range key {
		key[i] = letters[rand.Intn(len(letters))]
	}
	return string(key)
}

// GenerateRandomAlnums is used to generate a string of lowe and upper case letters and numbers.
func GenerateRandomAlnums(n int) string {
	var alnums = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	key := make([]rune, n)
	for i := range key {
		key[i] = alnums[rand.Intn(len(alnums))]
	}
	return string(key)
}

// ValidateName will check that a character name fits within the rules; 3-20 letters, alphabetic.
func ValidateName(name string) bool {
	var regexName = regexp.MustCompile("^[a-zA-Z]{3,20}$")
	return regexName.MatchString(name)
}

// ValidateUsername will check that an account Username fits within the rules; 3-30 characters, alphanumeric.
func ValidateUsername(name string) bool {
	var regexUsername = regexp.MustCompile("^[a-zA-Z0-9]{3,20}$")
	return regexUsername.MatchString(name)
}

// DefaultGOPATH will get the default GOPATH variable that is set up when you install Go.
func DefaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}

// Use is a go hack that lets you null use a variable or any interface to suppress errors when commenting code.
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

// PrettyPrint is a very simple helper function that converts struct data to JSON for pretty printing to the command line.
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
