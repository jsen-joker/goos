package env

import (
	"github.com/jsen-joker/goos/core/utils"
	"os"
	"strings"
)

// goos global env

var GoosHome = os.Getenv("HOME") + "/goos"
var GoosSecurityIgnoreUrls = []string{}
//var GoosDatabase = "user:password@tcp(127.0.0.1:3306)/test"
var GoosDatabase = "root:1234@tcp(127.0.0.1:3306)/goos"
var GoosVersion = "0.0.1"

func InitEnv()  {
	GoosHome = GetEnv(KeyGoosHome, GoosHome)
	GoosSecurityIgnoreUrls = strings.Split(GetEnv(KeyGoosSecurityIgnoreUrls, "/**"), ",")
	GoosDatabase = GetEnv(KeyGoosDatabase, GoosDatabase)
	GoosVersion = GetEnv(KeyGoosVersion, GoosVersion)
}

func GetEnv(key string, defaultValue string) string  {
	env := os.Getenv(key)
	if utils.NotEmpty(&env) {
		return env
	}
	return defaultValue
}
