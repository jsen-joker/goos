package controller

import (
	"fmt"
	"github.com/jsen-joker/goos/plugin-goos-ui/constants"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func HandleStaticResource(w http.ResponseWriter, r *http.Request)  {
	path := r.URL.Path
	if strings.Contains(path, ".") {
		requestType := path[strings.LastIndex(path, ".") : ]

		switch requestType {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "text/javascript")
		default:
			w.Header().Set("Content-Type", "text/html")
		}
	}

	fin, err := os.Open(constants.GOOS_HOME + path)
	fmt.Println(constants.GOOS_HOME + path)
	defer fin.Close()
	if err != nil {
		log.Fatal("static resource:", err)
	}

	fd, _ := ioutil.ReadAll(fin)
	w.Write(fd)
}
