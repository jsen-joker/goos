package controller

import (
	"encoding/json"
	"github.com/jsen-joker/goos/core/utils"
	"net/http"
)

func UIEcho(w http.ResponseWriter, r *http.Request)  {
	resp := utils.Succeed(nil)
	if err := json.NewEncoder(w).Encode(resp); err != nil{
		panic(err)
	}
}
