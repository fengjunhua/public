package utils

import (
	"os/user"
)

func GetUserName() (string, error) {

	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	userName := currentUser.Username
	return userName, nil

}
