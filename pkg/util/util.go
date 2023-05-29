package util

import "log"

func HandleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}

	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
