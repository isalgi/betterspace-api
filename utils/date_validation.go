package utils

import (
	"strings"
	"time"
)

func DateValidation(date string) error {
	var err error

	_, err = time.Parse("02/01/2006", date)

	if err != nil {
		return err
	}

	return err
}

func StatusValidation(str string) error {
	var err error
	var list = []string{"pending","on process","accepted","cancelled"}

	for _, v := range list {
		if strings.Contains(str, v) {
			return err
		}
	}

	return &statusError{}
}