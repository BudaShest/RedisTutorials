package helpers

import "log"

func FailOnError(err error) {
	if err != nil {
		log.Printf("Critiacal error: %v\n", err)
	}
}
