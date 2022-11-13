package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)
	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	return string(hashedPasswordBytes), err

}

func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil

}

func readPasswordFile() string {
	content, err := ioutil.ReadFile(".password")
	if err != nil {
		fmt.Println("Err")
	}
	st := string(content)
	return st
}

func savePassword(hash string) {
	s := []byte(hash)
	ioutil.WriteFile(".password", s, 0600)
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func PasswordPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
		b, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(b)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}

func lockedBin() bool {
	passwordEntry := PasswordPrompt("Password: ")
	newHashed, _ := hashPassword(passwordEntry)
	var oldHashed string
	if checkFileExists(".password") {
		oldHashed = readPasswordFile()
		match := doPasswordsMatch(oldHashed, passwordEntry)
		if match {
			fmt.Println("Match")
			return true
		} else {
			fmt.Printf("Old hash: %s\n", oldHashed)
			fmt.Printf("New hash: %s\n", newHashed)
			fmt.Println("No Match")
			return false
		}

	} else {
		savePassword(newHashed)
		fmt.Println("Saving password")
		return lockedBin()
	}

}

// I need to change something
func main() {
	lockedBin()
}
