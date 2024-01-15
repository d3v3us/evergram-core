package env

import (
	"os"
)

func IsDev() bool {
	return isDev
}
func IsProduction() bool {
	return isDev
}
func IsQa() bool {
	return isDev
}

var (
	isDev        = os.Getenv("APP_ENV") == "development"
	isProduction = os.Getenv("APP_ENV") == "production"
	isQa         = os.Getenv("APP_ENV") == "qa"
)
