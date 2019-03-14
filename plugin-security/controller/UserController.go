package controller

import (
	"encoding/json"
	"github.com/jsen-joker/goos/core/utils"
	"github.com/jsen-joker/goos/plugin-security/service"
	utils2 "github.com/jsen-joker/goos/plugin-security/utils"
	"net/http"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {

	auth := r.Header.Get("Authorization")

	if utils.NotEmpty(&auth) {
		id, _, _, err := utils2.GetSubject(auth[7:])
		if err != nil {
			resp := utils.Failed(err.Error())
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
			return
		}

		user, err := service.GetUserByAccountId(id)
		if err != nil {
			resp := utils.Failed(err.Error())
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
			return
		}

		resp := utils.Succeed(nil).PutAll(user)
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	} else {
		resp := utils.Failed("auth error")
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	}

}

func ChangeSelectNamespace(w http.ResponseWriter, r *http.Request)  {
	namespace, err := utils.GetInt64Parameter("namespace", true, 1, w, r)

	if err != nil {
		return
	}

	auth := r.Header.Get("Authorization")

	if utils.NotEmpty(&auth) {
		id, _, _, err := utils2.GetSubject(auth[7:])
		if err != nil {
			resp := utils.Failed(err.Error())
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
			return
		}

		eff, err := service.ChangeSelectedNamespace(id, namespace)
		if err != nil {
			resp := utils.Failed(err.Error())
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
			return
		}

		resp := utils.Succeed(eff)
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	} else {
		resp := utils.Failed("auth error")
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	}
}