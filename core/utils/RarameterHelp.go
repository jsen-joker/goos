package utils

import (
	"encoding/json"
	"errors"
	"github.com/jsen-joker/gend/exporter/rest/utils"
	"net/http"
	"strconv"
)

func getParameter(name string, r *http.Request) []string  {
	v := r.PostForm[name]

	if len(v) == 0 {
		v = r.URL.Query()[name]
	}
	return v
}

func GetParameter(name string, nullAble bool, defaultValue string, w http.ResponseWriter, r *http.Request) (value string, err error) {
	v := getParameter(name, r)

	if len(v) == 0 {
		if !nullAble {
			resp := utils.Failed("goos required parameter " + name)
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				return "", err
			} else {
				return "", errors.New("goos required parameter " + name)
			}
		} else {
			return "", nil
		}
	} else {
		return v[0], nil
	}
}

func GetInt64Parameter(name string, nullAble bool, defaultValue int64, w http.ResponseWriter, r *http.Request) (value int64, err error) {
	if v, e := GetParameter(name, nullAble, strconv.FormatInt(defaultValue, 10), w, r); e != nil {
		return 0, e
	} else {
		if v, e := strconv.ParseInt(v, 10, 64); e == nil {
			return v, e
		} else {
			if !nullAble {
				return v, e
			} else {
				return defaultValue, nil
			}
		}

	}
}
