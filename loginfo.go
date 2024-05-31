// Logger

package main

import (
	"fmt"
	"log"
)

const (
	Err_bucket_maked_previous = "Your previous request to create the named bucket succeeded and you already own it"
	Err_name_bucket_empty     = "Bucket name cannot be empty"
)

var logger = log.Default()

func LogErrorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func LogInfo(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}
