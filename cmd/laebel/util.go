package main

import (
	"log"
	"os"
)

// fatal logs an error and exits the program with a non-zero status code.
//
// The message is always logged. If err is not nil, it is logged as the cause.
// If hints are provided, they are logged as well.
func fatal(err error, message string, hints ...string) {
	log.Println("Fatal: " + message)
	if err != nil {
		log.Printf("Cause: %s\n", err)
	}
	for _, hint := range hints {
		log.Printf("Hint: %s\n", hint)
	}
	os.Exit(1)
}
