package models

import "os"

var goBlobBaseUrl string
var goBlobAuthorization string

func InitializeGoBlobBaseUrl() {
	goBlobBaseUrl = os.Getenv("GO_BLOB_BASE_URL")
}
func InitializeGoBlobAuthorization() {
	goBlobAuthorization = os.Getenv("GO_BLOB_AUTHORIZATION")
}
