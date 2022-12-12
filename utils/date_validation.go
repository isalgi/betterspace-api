package utils

import "time"

func DateValidation(date string) error {
	var err error

	_, err = time.Parse("02/01/2006", date)

	if err != nil {
		return err
	}

	return err
}